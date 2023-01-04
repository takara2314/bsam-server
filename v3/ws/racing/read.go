package racing

import (
	"encoding/json"
	"log"
	"time"

	"github.com/shiguredo/websocket"
)

type AuthInfo struct {
	Token  string `json:"token"`
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	MarkNo int    `json:"mark_no"`
}

type PassedInfo struct {
	MarkNo     int `json:"mark_no"`
	NextMarkNo int `json:"next_mark_no"`
}

type StartInfo struct {
	IsStarted bool `json:"started"`
}

type SetMarkNoInfo struct {
	UserID     string `json:"user_id"`
	MarkNo     int    `json:"mark_no"`
	NextMarkNo int    `json:"next_mark_no"`
}

type DebugInfo struct {
	Message string `json:"message"`
}

func (c *Client) readPump() {
	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, msgRaw, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(
				err,
				websocket.CloseGoingAway,
				websocket.CloseAbnormalClosure,
			) {
				log.Println("UnexpectedCloseError:", err)
			}
			return
		}

		var msg map[string]any
		if err := json.Unmarshal([]byte(msgRaw), &msg); err != nil {
			log.Println(err)
			continue
		}

		// Call handler by message type
		switch msg["type"].(string) {
		case "auth":
			var msg AuthInfo
			json.Unmarshal([]byte(msgRaw), &msg)
			c.auth(&msg)

		case "position":
			var msg Position
			json.Unmarshal([]byte(msgRaw), &msg)
			c.receivePos(&msg)

		case "location":
			var msg Location
			json.Unmarshal([]byte(msgRaw), &msg)
			c.receiveLoc(&msg)

		case "passed":
			var msg PassedInfo
			json.Unmarshal([]byte(msgRaw), &msg)
			c.handlerPassed(&msg)

		case "start":
			var msg StartInfo
			json.Unmarshal([]byte(msgRaw), &msg)
			c.Hub.startRace(msg.IsStarted)

		case "set_mark_no":
			var msg SetMarkNoInfo
			json.Unmarshal([]byte(msgRaw), &msg)
			c.Hub.setMarkNoForce(&msg)

		case "debug":
			var msg DebugInfo
			json.Unmarshal([]byte(msgRaw), &msg)
			log.Printf("Debug <%s>: %s\n", c.UserID, msg.Message)
		}
	}
}
