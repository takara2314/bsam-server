package race

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 5 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 1024
)

type Client struct {
	Hub      *Hub
	Conn     *websocket.Conn
	UserId   string
	Role     string
	Position Position
	Send     chan *PointNav
}

type Position struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type PointNav struct {
	Begin  bool  `json:"is_begin"`
	Now    Point `json:"now"`
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
			if websocket.IsUnexpectedCloseError(
				err,
				websocket.CloseGoingAway,
				websocket.CloseAbnormalClosure,
			) {
				fmt.Println(err)
			}
			break
		}

		fmt.Println(string(message))

		err = json.Unmarshal(message, &c.Position)
		if err != nil {
			panic(err)
		}

		c.Send <- &PointNav{
			Begin: false,
			Now: Point{
				Point:     2,
				Latitude:  34.29387,
				Longitude: 136.7622367,
			},
			Latest: 1,
		}

		// message = bytes.TrimSpace(bytes.Replace(message, []byte{'\n'}, []byte{' '}, -1))
		// c.Hub.Boardcast <- message
	}
}

// writePump pings once every 49.5s
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

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			err := c.Conn.WriteMessage(websocket.PingMessage, nil)
			if err != nil {
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			encoded, _ := json.Marshal(PointNav{
				Begin: true,
				Now: Point{
					Point:     1,
					Latitude:  20.2,
					Longitude: 20.1,
				},
				Latest: -1,
			})
			w.Write(encoded)
		}
	}
}
