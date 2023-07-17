package racing

import (
	"log"
	"sort"
	"time"
)

type Hub struct {
	AssociationID string
	Clients       map[string]*Client
	Athletes      map[string]*Client
	Marks         map[string]*Client
	Managers      map[string]*Client
	Disconnectors map[string]*Client
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
		Disconnectors: make(map[string]*Client),
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

// disconnectEvent disconnects the client.
func (h *Hub) disconnectEvent(c *Client) {
	if _, ok := h.Clients[c.ID]; !ok {
		return
	}

	log.Println("Disconnected:", c.ID)
	c.Conn.Close()

	// Register the client to the disconnector group
	h.Disconnectors[c.ID] = c

	// Unregister the client from the role group
	delete(c.Hub.Clients, c.ID)
	delete(c.Hub.Athletes, c.ID)
	deleteMarks(c.Hub.Marks, c.ID)
	delete(c.Hub.Managers, c.ID)
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
	deleteMarks(c.Hub.Marks, c.ID)
	delete(c.Hub.Managers, c.ID)
	delete(c.Hub.Disconnectors, c.ID)
}

// deleteMarks deletes the marks of the client.
func deleteMarks(marks map[string]*Client, id string) {
	for _, c := range marks {
		if c.UserID == id {
			c.Conn = nil
			c.UserID = ""
		}
	}
}

// getAthleteInfos returns the athlete infos.
func (h *Hub) getAthleteInfos() []Athlete {
	athletes := make([]Athlete, len(h.Athletes))

	cnt := 0
	for _, c := range h.Athletes {
		athletes[cnt] = Athlete{
			UserID:       c.UserID,
			NextMarkNo:   c.NextMarkNo,
			CourseLimit:  c.CourseLimit,
			BatteryLevel: c.BatteryLevel,
			CompassDeg:   c.CompassDeg,
			Location:     c.Location,
		}
		cnt++
	}

	// Sort by user id asc
	sort.Slice(athletes, func(i int, j int) bool {
		return athletes[i].UserID > athletes[j].UserID
	})

	return athletes
}

// getMarkInfos returns the mark infos.
func (h *Hub) getMarkInfos() []Mark {
	marks := make([]Mark, h.MarkNum)

	for _, c := range h.Marks {
		if c.MarkNo > h.MarkNum {
			panic("invalid mark no")
		}
		marks[c.MarkNo-1] = Mark{
			UserID:       c.UserID,
			MarkNo:       c.MarkNo,
			BatteryLevel: c.BatteryLevel,
			Position: Position{
				Lat: c.Location.Lat,
				Lng: c.Location.Lng,
				Acc: c.Location.Acc,
			},
		}
	}

	for i := range marks {
		marks[i].MarkNo = i + 1
	}

	return marks
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
		go c.sendStartRaceMsg()
	}
}

// setNextMarkNoForce force sets the client's next mark no.
func (h *Hub) setNextMarkNoForce(info *SetMarkNoInfo) {
	id := h.findClientID(info.UserID)
	if id == "" {
		log.Printf("Force Next Mark Failed: Not Found (%s)\n", info.UserID)
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

// findDisconnectedID returns the disconnected client id by user id.
func (h *Hub) findDisconnectedID(userID string) string {
	for _, c := range h.Disconnectors {
		if c.UserID == userID {
			return c.ID
		}
	}

	return ""
}
