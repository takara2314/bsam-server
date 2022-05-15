package race

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var (
	ErrClosedChannel = errors.New("closed channel")
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 10 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	navPeriod      = 5 * time.Second
	maxMessageSize = 1024
)

type Client struct {
	Hub         *Hub
	Conn        *websocket.Conn
	UserID      string
	Role        string
	PointNo     int
	NextPoint   int
	LatestPoint int
	Position    Position
	CourseLimit float32
	Send        chan *PointNav
	SendManage  chan *ManageInfo
	SendLive    chan *LiveInfo
	Mux         sync.RWMutex
}

type Position struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type PointNav struct {
	Begin       bool    `json:"is_begin"`
	Next        Point   `json:"next"`
	Latest      int     `json:"latest"`
	DebugNowLat float64 `json:"debug_now_lat"`
	DebugNowLng float64 `json:"debug_now_lng"`
}

type Point struct {
	Point     int     `json:"point"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type PointDevice struct {
	DeviceID  string  `json:"device_id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type ManageInfo struct {
	UserID    string   `json:"user_id"`
	Latitude  float64  `json:"latitude"`
	Longitude float64  `json:"longitude"`
	Next      PointNav `json:"next"`
}

type LiveInfo struct {
	Begin  bool        `json:"is_begin"`
	PointA PointDevice `json:"point_a"`
	PointB PointDevice `json:"point_b"`
	PointC PointDevice `json:"point_c"`
}

// readPump waits message
// and when receive message, send it as boardcast.
func (c *Client) readPump() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println(c.UserID, "Disconnected >>", err)
			if websocket.IsUnexpectedCloseError(
				err,
				websocket.CloseGoingAway,
				websocket.CloseAbnormalClosure,
			) {
				log.Println(err)
			}
			break
		}

		// Obtain a position info into a client instance.
		var tmp Position
		err = json.Unmarshal(message, &tmp)
		if !(tmp.Latitude == 0.0 || tmp.Longitude == 0.0) {
			fmt.Println(c.UserID, c.Position)
			// Update client position
			c.Position = tmp
		}
		if err != nil {
			panic(err)
		}

		// Check that the user passed the mark.
		c.passCheck()
	}
}

// writePump pings once every 2.7s
// and when Send channel sends data, send it to the client.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	go c.sendNextNavEvent()

	for {
		select {
		case message, ok := <-c.Send:
			err := c.sendEvent(message, ok)
			if err != nil {
				return
			}

		case message, ok := <-c.SendManage:
			err := c.sendManageEvent(message, ok)
			if err != nil {
				return
			}

		case message, ok := <-c.SendLive:
			err := c.sendLiveEvent(message, ok)
			if err != nil {
				return
			}

		case <-ticker.C:
			err := c.pingEvent()
			if err != nil {
				return
			}
		}
	}
}

// sendNextNavEvent sends next nav info every 5.4s
func (c *Client) sendNextNavEvent() {
	ticker := time.NewTicker(navPeriod)
	defer ticker.Stop()

	for {
		<-ticker.C
		// Do not send next nav info to a manage user and a point user.
		if !(c.Role == "manage" || c.Role == "admin") {
			err := c.sendNextNav()
			if err != nil {
				return
			}
		}
	}
}

// sendEvent sends client a navigation infomation.
// SEARCH: always looping without a send signal?
func (c *Client) sendEvent(message *PointNav, isOpen bool) error {
	c.Mux.Lock()
	defer c.Mux.Unlock()

	c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
	if !isOpen {
		c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
		return ErrClosedChannel
	}

	w, err := c.Conn.NextWriter(websocket.TextMessage)
	if err != nil {
		return err
	}
	encoded, _ := json.Marshal(message)
	w.Write(encoded)

	for i := 0; i < len(c.Send); i++ {
		w.Write([]byte{'\n'})
		encoded, _ = json.Marshal(<-c.Send)
		w.Write(encoded)
	}

	err = w.Close()
	if err != nil {
		return err
	}

	return nil
}

// sendManageEvent sends manage clients a manage infomation.
func (c *Client) sendManageEvent(message *ManageInfo, isOpen bool) error {
	c.Mux.Lock()
	defer c.Mux.Unlock()

	c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
	if !isOpen {
		c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
		return ErrClosedChannel
	}

	w, err := c.Conn.NextWriter(websocket.TextMessage)
	if err != nil {
		return err
	}
	encoded, _ := json.Marshal(message)
	w.Write(encoded)

	for i := 0; i < len(c.SendManage); i++ {
		w.Write([]byte{'\n'})
		encoded, _ = json.Marshal(<-c.SendManage)
		w.Write(encoded)
	}

	err = w.Close()
	if err != nil {
		return err
	}

	return nil
}

// sendLiveEvent sends live message.
func (c *Client) sendLiveEvent(message *LiveInfo, isOpen bool) error {
	c.Mux.Lock()
	defer c.Mux.Unlock()

	c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
	if !isOpen {
		c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
		return ErrClosedChannel
	}

	w, err := c.Conn.NextWriter(websocket.TextMessage)
	if err != nil {
		return err
	}
	encoded, _ := json.Marshal(message)
	w.Write(encoded)

	for i := 0; i < len(c.SendLive); i++ {
		w.Write([]byte{'\n'})
		encoded, _ = json.Marshal(<-c.SendLive)
		w.Write(encoded)
	}

	err = w.Close()
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) pingEvent() error {
	c.Mux.Lock()
	defer c.Mux.Unlock()

	c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
	if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
		return err
	}

	return nil
}

func (c *Client) sendNextNav() error {
	// Announce next point info.
	var nextLat, nextLng float64
	switch c.NextPoint {
	case 1:
		nextLat = c.Hub.PointA.Latitude
		nextLng = c.Hub.PointA.Longitude
	case 2:
		nextLat = c.Hub.PointB.Latitude
		nextLng = c.Hub.PointB.Longitude
	case 3:
		nextLat = c.Hub.PointC.Latitude
		nextLng = c.Hub.PointC.Longitude
	}

	nav := PointNav{
		Begin: c.Hub.Begin,
		Next: Point{
			Point:     c.NextPoint,
			Latitude:  nextLat,
			Longitude: nextLng,
		},
		Latest:      c.LatestPoint,
		DebugNowLat: c.Position.Latitude,
		DebugNowLng: c.Position.Longitude,
	}

	if IsClosedSendChan(c.Send) {
		c.Send <- &nav
	} else {
		return ErrClosedChannel
	}

	// Broadcast for manage users and admin users.
	c.Hub.Managecast <- &ManageInfo{
		UserID:    c.UserID,
		Latitude:  c.Position.Latitude,
		Longitude: c.Position.Longitude,
		Next:      nav,
	}

	return nil
}

// IsClosedSendChan returns true when that send channel is opened.
func IsClosedSendChan(c chan *PointNav) bool {
	var ok bool

	select {
	case _, ok = <-c:
	default:
		ok = true
	}

	return ok
}

// IsClosedSendManageChan returns true when that send manage channel is opened.
func IsClosedSendManageChan(c chan *ManageInfo) bool {
	var ok bool

	select {
	case _, ok = <-c:
	default:
		ok = true
	}

	return ok
}

// IsClosedSendLiveChan returns true when that send live channel is opened.
func IsClosedSendLiveChan(c chan *LiveInfo) bool {
	var ok bool

	select {
	case _, ok = <-c:
	default:
		ok = true
	}

	return ok
}
