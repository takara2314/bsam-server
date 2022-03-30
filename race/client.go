package race

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 3 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 1024
)

type Client struct {
	Hub         *Hub
	Conn        *websocket.Conn
	UserId      string
	Role        string
	NextPoint   int
	LatestPoint int
	Position    Position
	CourseLimit float32
	Send        chan *PointNav
	SendManage  chan *ManageInfo
}

type Position struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type PointNav struct {
	Begin  bool  `json:"is_begin"`
	Next   Point `json:"next"`
	Latest int   `json:"latest"`
}

type Point struct {
	Point     int     `json:"point"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type PointDevice struct {
	UserId    string
	Latitude  float64
	Longitude float64
}

type ManageInfo struct {
	UserId    string   `json:"user_id"`
	Latitude  float64  `json:"latitude"`
	Longitude float64  `json:"longitude"`
	Next      PointNav `json:"next"`
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
			fmt.Println(c.UserId, "障害発生 >>", err)
			if websocket.IsUnexpectedCloseError(
				err,
				websocket.CloseGoingAway,
				websocket.CloseAbnormalClosure,
			) {
				fmt.Println(err)
			}
			break
		}

		// Obtain a position info into a client instance.
		fmt.Println(c.UserId, "message >>", string(message))
		err = json.Unmarshal(message, &c.Position)
		if err != nil {
			panic(err)
		}

		// Check that the user passed the mark.
		c.passCheck()

		// message = bytes.TrimSpace(bytes.Replace(message, []byte{'\n'}, []byte{' '}, -1))
		// c.Hub.Boardcast <- message
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

	for {
		select {
		case message, isOpen := <-c.Send:
			c.sendEvent(message, isOpen)

		case message, isOpen := <-c.SendManage:
			c.sendManageEvent(message, isOpen)

		case <-ticker.C:
			c.pingEvent()

			// Do not send next nav info to a manage user and a point user.
			if !(c.Role == "manage" || c.Role == "admin") {
				err := c.sendNextNav()
				if err != nil {
					return
				}
			}
		}
	}
}

// sendEvent sends client a navigation infomation.
// SEARCH: always looping without a send signal?
func (c *Client) sendEvent(message *PointNav, isOpen bool) {
	c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
	if !isOpen {
		c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
		return
	}

	w, err := c.Conn.NextWriter(websocket.TextMessage)
	if err != nil {
		return
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
		return
	}
}

// sendManageEvent sends manage clients a manage infomation.
func (c *Client) sendManageEvent(message *ManageInfo, isOpen bool) {
	c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
	if !isOpen {
		c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
		return
	}

	w, err := c.Conn.NextWriter(websocket.TextMessage)
	if err != nil {
		return
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
		return
	}
}

func (c *Client) pingEvent() {
	c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
	if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
		return
	}
}

func (c *Client) sendNextNav() error {
	fmt.Println("ping to", c.UserId)
	fmt.Println(c.Hub.Clients)
	c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
	err := c.Conn.WriteMessage(websocket.PingMessage, nil)
	if err != nil {
		fmt.Println("nilじゃなかった1", err)
		return err
	}

	w, err := c.Conn.NextWriter(websocket.TextMessage)
	if err != nil {
		return err
	}

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
		Latest: c.LatestPoint,
	}

	encoded, _ := json.Marshal(nav)
	w.Write(encoded)

	// Broadcast for manage users and admin users.
	c.Hub.Managecast <- &ManageInfo{
		UserId:    c.UserId,
		Latitude:  c.Position.Latitude,
		Longitude: c.Position.Longitude,
		Next:      nav,
	}

	return nil
}
