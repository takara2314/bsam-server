package racing

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type AuthInfo struct {
	Token  string `json:"token"`
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	MarkNo int    `json:"mark_no"`
}

type PassedInfo struct {
	PassedMarkNo int `json:"passed_mark_no"`
	NextMarkNo   int `json:"next_mark_no"`
}

type StartInfo struct {
	IsStarted bool `json:"started"`
}

type SetMarkNoInfo struct {
	UserID     string `json:"user_id"`
	MarkNo     int    `json:"mark_no"`
	NextMarkNo int    `json:"next_mark_no"`
}

type BatteryInfo struct {
	Level int `json:"level"`
}

type DebugInfo struct {
	Message string `json:"message"`
}

// TODO: 関数を細かく分ける
//
//nolint:funlen,gocognit,cyclop
func (c *Client) readPump() {
	c.Conn.SetReadLimit(maxMessageSize)

	err := c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	if err != nil {
		return
	}

	c.Conn.SetPongHandler(func(string) error {
		return c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	})

	for {
		_, msgRaw, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(
				err,
				websocket.CloseGoingAway,
				websocket.CloseAbnormalClosure,
			) {
				log.Printf("%s (%s) >> unexpected close error: %v", c.ID, c.UserID, err)
			}

			return
		}

		var msg map[string]any
		if err := json.Unmarshal(msgRaw, &msg); err != nil {
			log.Println(err)
			continue
		}

		// If unauthenticated or the guest, only accept auth message
		if (c.Role == "" || c.Role == "guest") && msg["type"].(string) != "auth" {
			continue
		}

		// Call handler by message type
		switch msg["type"].(string) {
		case "auth":
			var msg AuthInfo

			err := json.Unmarshal(msgRaw, &msg)
			if err != nil {
				log.Printf("Error <%s>: %s\n", c.UserID, err)
				continue
			}

			c.auth(&msg)

		case "position":
			var msg Position

			err := json.Unmarshal(msgRaw, &msg)
			if err != nil {
				log.Printf("Error <%s>: %s\n", c.UserID, err)
				continue
			}

			c.receivePos(&msg)

		case "location":
			var msg Location

			err := json.Unmarshal(msgRaw, &msg)
			if err != nil {
				log.Printf("Error <%s>: %s\n", c.UserID, err)
				continue
			}

			c.receiveLoc(&msg)

		case "passed":
			var msg PassedInfo

			err := json.Unmarshal(msgRaw, &msg)
			if err != nil {
				log.Printf("Error <%s>: %s\n", c.UserID, err)
				continue
			}

			c.handlerPassed(&msg)

		case "start":
			var msg StartInfo

			err := json.Unmarshal(msgRaw, &msg)
			if err != nil {
				log.Printf("Error <%s>: %s\n", c.UserID, err)
				continue
			}

			c.Hub.startRace(msg.IsStarted)

		case "set_next_mark_no":
			var msg SetMarkNoInfo

			err := json.Unmarshal(msgRaw, &msg)
			if err != nil {
				log.Printf("Error <%s>: %s\n", c.UserID, err)
				continue
			}

			c.Hub.setNextMarkNoForce(&msg)

		case "battery":
			var msg BatteryInfo

			err := json.Unmarshal(msgRaw, &msg)
			if err != nil {
				log.Printf("Error <%s>: %s\n", c.UserID, err)
				continue
			}

			c.receiveBattery(&msg)

		case "debug":
			var msg DebugInfo

			err := json.Unmarshal(msgRaw, &msg)
			if err != nil {
				log.Printf("Error <%s>: %s\n", c.UserID, err)
				continue
			}

			log.Printf("Debug <%s>: %s\n", c.UserID, msg.Message)
		}
	}
}
