package race

type Hub struct {
	Clients    map[string]*Client
	Boardcast  chan []byte
	Register   chan *Client
	Unregister chan *Client
}

// NewHub creates a new hub instrance.
func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[string]*Client),
		Boardcast:  make(chan []byte),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

// Run working the tasks, such as new device register event and boardcasting.
func (hub *Hub) Run() {
	for {
		select {
		case client := <-hub.Register:
			hub.registerEvent(client)
		case client := <-hub.Unregister:
			hub.unregisterEvent(client)
			// case message := <-hub.Boardcast:
			// 	hub.boardcastEvent(message)
		}
	}
}

// registerEvent adds new device.
func (hub *Hub) registerEvent(client *Client) {
	hub.Clients[client.Id] = client
}

// unregisterEvent removes this client.
func (hub *Hub) unregisterEvent(client *Client) {
	if _, exist := hub.Clients[client.Id]; exist {
		close(client.Send)
		delete(hub.Clients, client.Id)
	}
}

// // boardcastEvent boardcasts to all client.
// func (hub *Hub) boardcastEvent(message []byte) {
// 	for _, client := range hub.Clients {
// 		select {
// 		case client.Send <- message:
// 		default:
// 			close(client.Send)
// 			delete(hub.Clients, client.Id)
// 		}
// 	}
// }
