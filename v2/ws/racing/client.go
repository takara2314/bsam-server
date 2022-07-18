package racing

import (
	"bsam-server/utils"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 10 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	navPeriod      = 5 * time.Second
	maxMessageSize = 1024
)

var (
	ErrClosedChannel = errors.New("closed channel")
)

type Client struct {
	ID     string
	Hub    *Hub
	Conn   *websocket.Conn
	UserID string
	Send   chan *Message
}

type Message struct {
	Type    bool   `json:"type"`
	Message string `json:"message"`
}

func NewClient(raceID string, conn *websocket.Conn) *Client {
	return &Client{
		ID:     utils.RandString(8),
		Hub:    rooms[raceID],
		Conn:   conn,
		UserID: "",
		Send:   make(chan *Message),
	}
}

func (c *Client) readPump() {
	defer func() {
		c.Hub.Unregister <- c
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, msgRaw, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println(c.ID, "Disconnected:", err)
			if websocket.IsUnexpectedCloseError(
				err,
				websocket.CloseGoingAway,
				websocket.CloseAbnormalClosure,
			) {
				log.Println(err)
			}
			return
		}

		var msg map[string]interface{}
		if err := json.Unmarshal([]byte(msgRaw), &msg); err != nil {
			log.Println(err)
			continue
		}

		switch msg["type"].(string) {
		case "message":
			fmt.Println("Message Received:", msg["message"])
		}
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Hub.Unregister <- c
	}()

	for {
		select {
		case msg, ok := <-c.Send:
			err := c.sendEvent(msg, ok)
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

func (c *Client) sendEvent(msg *Message, ok bool) error {
	if !ok {
		c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
		return ErrClosedChannel
	}

	c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
	return c.Conn.WriteJSON(msg)
}

func (c *Client) pingEvent() error {
	c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
	return c.Conn.WriteMessage(websocket.PingMessage, nil)
}
