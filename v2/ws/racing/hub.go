package racing

type Hub struct {
	RaceID     string
	Clients    map[string]*Client
	Register   chan *Client
	Unregister chan *Client
}

func NewHub(raceID string) *Hub {
	return &Hub{
		RaceID:     raceID,
		Clients:    make(map[string]*Client),
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
	h.Clients[c.ID] = c
}

func (h *Hub) unregisterEvent(c *Client) {
	if _, ok := h.Clients[c.ID]; !ok {
		return
	}

	c.Conn.Close()
	close(c.Send)
	delete(h.Clients, c.ID)
}
