package racing

import (
	"log"
)

type Hub struct {
	AssociationID string
	Clients       map[string]*Client
	Athletes      map[string]*Client
	Marks         map[string]*Client
	Managers      map[string]*Client
	MarkNum       int
	IsStarted     bool
	Register      chan *Client
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
		Register:      make(chan *Client),
		Unregister:    make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.registerEvent(client)
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
func (h *Hub) unregisterEvent(c *Client) {
	if _, ok := h.Clients[c.ID]; !ok {
		return
	}

	log.Println("Disconnected:", c.ID)
	c.Conn.Close()

	c.Connecting = false
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

	for _, c := range h.Athletes {
		c.sendStartRaceMsg()
	}
	for _, c := range h.Managers {
		c.sendStartRaceMsg()
	}
}

// setMarkNoForce force sets the client's mark no.
func (h *Hub) setMarkNoForce(info *SetMarkNoInfo) {
	id := h.findClientID(info.UserID)
	if id == "" {
		return
	}

	log.Printf("Force Mark Changed: %s -> [%d]\n", info.UserID, info.NextMarkNo)

	h.Clients[id].NextMarkNo = info.NextMarkNo

	h.Clients[id].sendSetMarkNoEvent(&SetMarkNoMsg{
		MarkNo:     info.MarkNo,
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
