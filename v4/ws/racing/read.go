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
//nolint:funlen,cyclop
func (c *Client) readPump() {
	c.Conn.SetReadLimit(MaxMessageByte)

	err := c.Conn.SetReadDeadline(time.Now().Add(PongWait))
	if err != nil {
		return
	}

	c.Conn.SetPongHandler(func(string) error {
		return c.Conn.SetReadDeadline(time.Now().Add(PongWait))
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

		err = json.Unmarshal(msgRaw, &msg)
		if err != nil {
			log.Println(err)
			continue
		}

		msgType, ok := msg["type"].(string)
		if !ok {
			log.Printf("%s (%s) >> message type assertion failed\n", c.ID, c.UserID)
			continue
		}

		// If unauthenticated or the guest, only accept auth message
		if (c.Role == "" || c.Role == GuestRole) && msgType != "auth" {
			continue
		}

		// Call handler by message type
		switch msgType {
		case "auth":
			c.handleAuthMessage(msgRaw)

		case "position":
			c.handlePositionMessage(msgRaw)

		case "location":
			c.handleLocationMessage(msgRaw)

		case "passed":
			c.handlePassedMessage(msgRaw)

		case "start":
			c.handleStartMessage(msgRaw)

		case "set_next_mark_no":
			c.handleSetNextMarkNoMessage(msgRaw)

		case "battery":
			c.handleBatteryMessage(msgRaw)

		case "debug":
			c.handleDebugMessage(msgRaw)
		}
	}
}

func (c *Client) handleAuthMessage(msgRaw []byte) {
	var msg AuthInfo

	if !c.decodeMessage(msgRaw, &msg) {
		return
	}

	c.auth(&msg)
}

func (c *Client) handlePositionMessage(msgRaw []byte) {
	var msg Position

	if !c.decodeMessage(msgRaw, &msg) {
		return
	}

	c.receivePos(&msg)
}

func (c *Client) handleLocationMessage(msgRaw []byte) {
	var msg Location

	if !c.decodeMessage(msgRaw, &msg) {
		return
	}

	c.receiveLoc(&msg)
}

func (c *Client) handlePassedMessage(msgRaw []byte) {
	var msg PassedInfo

	if !c.decodeMessage(msgRaw, &msg) {
		return
	}

	c.handlerPassed(&msg)
}

func (c *Client) handleStartMessage(msgRaw []byte) {
	var msg StartInfo

	if !c.decodeMessage(msgRaw, &msg) {
		return
	}

	c.Hub.startRace(msg.IsStarted)
}

func (c *Client) handleSetNextMarkNoMessage(msgRaw []byte) {
	var msg SetMarkNoInfo

	if !c.decodeMessage(msgRaw, &msg) {
		return
	}

	c.Hub.setNextMarkNoForce(&msg)
}

func (c *Client) handleBatteryMessage(msgRaw []byte) {
	var msg BatteryInfo

	if !c.decodeMessage(msgRaw, &msg) {
		return
	}

	c.receiveBattery(&msg)
}

func (c *Client) handleDebugMessage(msgRaw []byte) {
	var msg DebugInfo

	if !c.decodeMessage(msgRaw, &msg) {
		return
	}

	log.Printf("Debug <%s>: %s\n", c.UserID, msg.Message)
}

func (c *Client) decodeMessage(msgRaw []byte, target any) bool {
	err := json.Unmarshal(msgRaw, target)
	if err != nil {
		log.Printf("Error <%s>: %s\n", c.UserID, err)
		return false
	}

	return true
}
