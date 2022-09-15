package racing

import "log"

type Hub struct {
	RaceID     string
	Clients    map[string]*Client
	Athletes   map[string]*Client
	Marks      map[string]*Client
	MarkNum    int
	IsStarted  bool
	Register   chan *Client
	Unregister chan *Client
}

func NewHub(raceID string) *Hub {
	return &Hub{
		RaceID:     raceID,
		Clients:    make(map[string]*Client),
		Athletes:   make(map[string]*Client),
		Marks:      make(map[string]*Client),
		MarkNum:    3,
		IsStarted:  false,
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
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

func (h *Hub) registerEvent(c *Client) {
	log.Println("Joined:", c.ID)
	h.Clients[c.ID] = c
}

func (h *Hub) unregisterEvent(c *Client) {
	if _, ok := h.Clients[c.ID]; !ok {
		return
	}

	log.Println("Left:", c.ID)
	c.Conn.Close()
	delete(h.Clients, c.ID)
	delete(h.Athletes, c.ID)
	delete(h.Marks, c.ID)
}

func (h Hub) getMarkPositions() []Position {
	positions := make([]Position, h.MarkNum)

	for _, c := range h.Marks {
		if c.MarkNo > h.MarkNum {
			panic("invalid mark no")
		}
		positions[c.MarkNo-1] = c.Position
	}

	return positions
}

func (h Hub) startRace(isStarted bool) {
	h.IsStarted = isStarted

	for _, c := range h.Athletes {
		c.sendStartRaceMsg()
	}
}

func (h Hub) setMarkNo(info *SetMarkNoInfo) {
	id := h.findClientID(info.UserID)
	if id == "" {
		return
	}

	log.Printf("Force Mark Changed: [%d] -> %s -> [%d]\n", info.MarkNo, info.UserID, info.NextMarkNo)

	h.Clients[id].MarkNo = info.MarkNo
	h.Clients[id].NextMarkNo = info.NextMarkNo

	h.Clients[id].sendSetMarkNoEvent(&SetMarkNoMsg{
		MarkNo:     info.MarkNo,
		NextMarkNo: info.NextMarkNo,
	})
}

func (h Hub) findClientID(userID string) string {
	for _, c := range h.Clients {
		if c.UserID == userID {
			return c.ID
		}
	}

	return ""
}
