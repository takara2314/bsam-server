package racing

import (
	"log"
	"time"
)

type Hub struct {
	AssociationID string
	Clients       map[string]*Client
	Athletes      map[string]*Client
	Marks         map[string]*Client
	Managers      map[string]*Client
	MarkNum       int
	IsStarted     bool
	StartAt       time.Time
	EndAt         time.Time
	Register      chan *Client
	Disconnect    chan *Client
	Unregister    chan *Client
}

func NewHub(assocID string) *Hub {
	return &Hub{
		AssociationID: assocID,
		Clients:       make(map[string]*Client),
		Athletes:      make(map[string]*Client),
		Marks:         make(map[string]*Client),
		Managers:      make(map[string]*Client),
		MarkNum:       3,
		IsStarted:     false,
		StartAt:       time.Unix(0, 0),
		EndAt:         time.Unix(0, 0),
		Register:      make(chan *Client),
		Disconnect:    make(chan *Client),
		Unregister:    make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.registerEvent(client)
		case client := <-h.Disconnect:
			h.disconnectEvent(client)
		case client := <-h.Unregister:
			h.unregisterEvent(client)
		}
	}
}

// registerEvent registers the client.
func (h *Hub) registerEvent(c *Client) {
	log.Println("Joined:", c.ID)
	h.Clients[c.ID] = c
}

// unregisterEvent unregisters the client.
func (h *Hub) disconnectEvent(c *Client) {
	if _, ok := h.Clients[c.ID]; !ok {
		return
	}

	log.Println("Disconnected:", c.ID)
	c.Conn.Close()
}

// unregisterEvent unregisters the client.
func (h *Hub) unregisterEvent(c *Client) {
	if _, ok := h.Clients[c.ID]; !ok {
		return
	}

	log.Println("Unregistered:", c.ID)
	c.Conn.Close()

	delete(c.Hub.Clients, c.ID)
	delete(c.Hub.Athletes, c.ID)
	delete(c.Hub.Marks, c.ID)
	delete(c.Hub.Managers, c.ID)
}

// getMarkPositions returns the mark positions.
func (h *Hub) getMarkPositions() []Position {
	positions := make([]Position, h.MarkNum)

	for _, c := range h.Marks {
		if c.MarkNo > h.MarkNum {
			panic("invalid mark no")
		}
		positions[c.MarkNo-1] = Position{
			Lat: c.Location.Lat,
			Lng: c.Location.Lng,
		}
	}

	return positions
}

// startRace sends start message to all clients.
func (h *Hub) startRace(isStarted bool) {
	h.IsStarted = isStarted

	if h.IsStarted {
		h.StartAt = time.Now()
	} else {
		h.EndAt = time.Now()
	}

	for _, c := range h.Clients {
		c.sendStartRaceMsg()
	}
}

// setNextMarkNoForce force sets the client's next mark no.
func (h *Hub) setNextMarkNoForce(info *SetMarkNoInfo) {
	id := h.findClientID(info.UserID)
	if id == "" {
		return
	}

	log.Printf("Force Next Mark Changed: %s -> [%d]\n", info.UserID, info.NextMarkNo)

	h.Clients[id].NextMarkNo = info.NextMarkNo

	h.Clients[id].sendSetNextMarkNoEvent(&SetNextMarkNoMsg{
		NextMarkNo: info.NextMarkNo,
	})
}

// findClientID returns the client id by user id.
func (h *Hub) findClientID(userID string) string {
	for _, c := range h.Clients {
		if c.UserID == userID {
			return c.ID
		}
	}

	return ""
}
