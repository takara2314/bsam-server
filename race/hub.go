package race

import (
	"fmt"
	"log"
	"sailing-assist-mie-api/bsamdb"
	"sailing-assist-mie-api/utils"
	"strings"
	"time"

	"github.com/lib/pq"
)

type Hub struct {
	RaceId     string
	Clients    map[string]*Client
	Managecast chan *ManageInfo
	Livecast   chan *LiveInfo
	Register   chan *Client
	Unregister chan *Client
	PointA     PointDevice
	PointB     PointDevice
	PointC     PointDevice
	Begin      bool
}

// NewHub creates a new hub instrance.
func NewHub(raceId string) *Hub {
	return &Hub{
		RaceId:     raceId,
		Clients:    make(map[string]*Client),
		Managecast: make(chan *ManageInfo),
		Livecast:   make(chan *LiveInfo),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

// Run working the tasks, such as new device register event and boardcasting.
// In addition, it updates mark positions every 2 seconds.
func (hub *Hub) Run() {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case client := <-hub.Register:
			hub.registerEvent(client)
		case client := <-hub.Unregister:
			hub.unregisterEvent(client)
		case message := <-hub.Managecast:
			hub.managecastEvent(message)
		case message := <-hub.Livecast:
			// fmt.Println("livecasting:", message)
			hub.livecastEvent(message)

		case <-ticker.C:
			hub.updateMarkPositions()
		}
	}
}

// registerEvent adds new device.
func (hub *Hub) registerEvent(client *Client) {
	log.Println(client.UserId, "joined.")

	hub.Clients[client.UserId] = client
	err := hub.addAthlete(client.UserId)
	if err != nil {
		panic(err)
	}
}

// unregisterEvent removes this client.
func (hub *Hub) unregisterEvent(client *Client) {
	log.Println(client.UserId, "left.")

	if hub.isExistUser(client.UserId) {
		close(client.Send)
		close(client.SendManage)
		close(client.SendLive)

		if !strings.HasPrefix(client.UserId, "NPC") {
			err := hub.removeAthlete(client.UserId)
			if err != nil {
				panic(err)
			}
		}

		delete(hub.Clients, client.UserId)
	}
}

// managecastEvent boardcasts to manage and admin client.
func (hub *Hub) managecastEvent(message *ManageInfo) {
	for _, client := range hub.Clients {
		if !(client.Role == "manage" || client.Role == "admin") {
			continue
		}

		if IsClosedSendManageChan(client.SendManage) {
			client.SendManage <- message
		} else {
			continue
		}
	}
}

// livecastEvent boardcasts live infomation.
func (hub *Hub) livecastEvent(message *LiveInfo) {
	for _, client := range hub.Clients {
		if IsClosedSendLiveChan(client.SendLive) {
			fmt.Println(client.UserId, "send!!!")
			client.SendLive <- message
		} else {
			continue
		}
	}
}

func (hub *Hub) isExistUser(userId string) bool {
	_, exist := hub.Clients[userId]
	return exist
}

func (hub *Hub) updateMarkPositions() {
	if hub.PointA.DeviceId != "" && hub.isExistUser(hub.PointA.DeviceId) {
		hub.PointA.Latitude = hub.Clients[hub.PointA.DeviceId].Position.Latitude
		hub.PointA.Longitude = hub.Clients[hub.PointA.DeviceId].Position.Longitude
	}
	if hub.PointB.DeviceId != "" && hub.isExistUser(hub.PointB.DeviceId) {
		hub.PointB.Latitude = hub.Clients[hub.PointB.DeviceId].Position.Latitude
		hub.PointB.Longitude = hub.Clients[hub.PointB.DeviceId].Position.Longitude
	}
	if hub.PointC.DeviceId != "" && hub.isExistUser(hub.PointC.DeviceId) {
		hub.PointC.Latitude = hub.Clients[hub.PointC.DeviceId].Position.Latitude
		hub.PointC.Longitude = hub.Clients[hub.PointC.DeviceId].Position.Longitude
	}

	// Livecast for all device
	go func() {
		hub.Livecast <- &LiveInfo{
			Begin:  hub.Begin,
			PointA: hub.PointA,
			PointB: hub.PointB,
			PointC: hub.PointC,
		}
	}()
}

// addAthlete adds a athlete in this race.
func (hub *Hub) addAthlete(userId string) error {
	// Connect to the database and insert such data.
	db, err := bsamdb.Open()
	if err != nil {
		panic(err)
	}
	defer db.DB.Close()

	// Obtain athletes info from database.
	rows, err := db.SelectSpecified(
		"races",
		[]bsamdb.Field{
			{Column: "id", Value: hub.RaceId},
		},
		[]string{"athlete"},
	)
	if err != nil {
		return err
	}

	rows.Next()
	var athletes []string
	rows.Scan(pq.Array(&athletes))

	// After append, update the database.
	_, err = db.Update(
		"races",
		"id",
		hub.RaceId,
		[]bsamdb.Field{{
			Column: "athlete",
			Value2d: utils.StrSliceToAnySlice(
				utils.StrSliceAdd(athletes, userId),
			),
		}},
	)

	return err
}

// removeAthlete removes a athlete from this race.
func (hub *Hub) removeAthlete(userId string) error {
	// Connect to the database and insert such data.
	db, err := bsamdb.Open()
	if err != nil {
		panic(err)
	}
	defer db.DB.Close()

	// Obtain athletes info from database.
	rows, err := db.SelectSpecified(
		"races",
		[]bsamdb.Field{
			{Column: "id", Value: hub.RaceId},
		},
		[]string{"athlete"},
	)
	if err != nil {
		return err
	}

	rows.Next()
	var athletes []string
	rows.Scan(pq.Array(&athletes))

	// After append, update the database.
	_, err = db.Update(
		"races",
		"id",
		hub.RaceId,
		[]bsamdb.Field{{
			Column: "athlete",
			Value2d: utils.StrSliceToAnySlice(
				utils.StrSliceRemove(athletes, userId),
			),
		}},
	)

	return err
}
