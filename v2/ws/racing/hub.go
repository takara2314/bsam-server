package racing

import "log"

type Hub struct {
	RaceID     string
	Clients    map[string]*Client
	Athletes   map[string]*Client
	Marks      map[string]*Client
	MarkNum    int
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
	// delete(h.Marks, c.ID)
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
