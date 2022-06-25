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
	RaceID     string
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
func NewHub(raceID string) *Hub {
	return &Hub{
		RaceID:     raceID,
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
			hub.livecastEvent(message)

		case <-ticker.C:
			hub.updateMarkPositions()
		}
	}
}

// registerEvent adds new device.
func (hub *Hub) registerEvent(client *Client) {
	log.Println(client.UserID, "joined.")

	if _, ok := hub.Clients[client.UserID]; ok {
		fmt.Println("!? 重複しています !?", client.UserID)
		hub.unregisterEvent(hub.Clients[client.UserID])
		fmt.Println("古いのを削除しました。")
	}

	hub.Clients[client.UserID] = client

	fmt.Println(len(hub.Clients))
	fmt.Println(hub.Clients)

	if client.Role == "athlete" {
		err := hub.addAthlete(client.UserID)
		if err != nil {
			panic(err)
		}
	}
}

// unregisterEvent removes this client.
func (hub *Hub) unregisterEvent(client *Client) {
	log.Println(client.UserID, "left.")

	if hub.isExistUser(client.UserID) {
		close(client.Send)
		close(client.SendManage)
		close(client.SendLive)

		if !strings.HasPrefix(client.UserID, "NPC") && client.Role == "athlete" {
			err := hub.removeAthlete(client.UserID)
			if err != nil {
				panic(err)
			}
		}

		delete(hub.Clients, client.UserID)
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
			client.SendLive <- message
		} else {
			continue
		}
	}
}

func (hub *Hub) isExistUser(userID string) bool {
	_, exist := hub.Clients[userID]
	return exist
}

func (hub *Hub) updateMarkPositions() {
	if hub.PointA.DeviceID != "" && hub.isExistUser(hub.PointA.DeviceID) {
		hub.PointA.Latitude = hub.Clients[hub.PointA.DeviceID].Position.Latitude
		hub.PointA.Longitude = hub.Clients[hub.PointA.DeviceID].Position.Longitude
	}
	if hub.PointB.DeviceID != "" && hub.isExistUser(hub.PointB.DeviceID) {
		hub.PointB.Latitude = hub.Clients[hub.PointB.DeviceID].Position.Latitude
		hub.PointB.Longitude = hub.Clients[hub.PointB.DeviceID].Position.Longitude
	}
	if hub.PointC.DeviceID != "" && hub.isExistUser(hub.PointC.DeviceID) {
		hub.PointC.Latitude = hub.Clients[hub.PointC.DeviceID].Position.Latitude
		hub.PointC.Longitude = hub.Clients[hub.PointC.DeviceID].Position.Longitude
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
func (hub *Hub) addAthlete(userID string) error {
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
			{Column: "id", Value: hub.RaceID},
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
		hub.RaceID,
		[]bsamdb.Field{{
			Column: "athlete",
			Value2d: utils.StrSliceToAnySlice(
				utils.StrSliceAdd(athletes, userID),
			),
		}},
	)

	return err
}

// removeAthlete removes a athlete from this race.
func (hub *Hub) removeAthlete(userID string) error {
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
			{Column: "id", Value: hub.RaceID},
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
		hub.RaceID,
		[]bsamdb.Field{{
			Column: "athlete",
			Value2d: utils.StrSliceToAnySlice(
				utils.StrSliceRemove(athletes, userID),
			),
		}},
	)

	return err
}
