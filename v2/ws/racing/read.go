package racing

import (
	"bsam-server/utils"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type AuthInfo struct {
	Token string `json:"token"`
	Role  string `json:"role"`
}

func (c *Client) auth(msg *AuthInfo) {
	userID, err := utils.GetUserIDFromJWT(msg.Token)
	if err != nil {
		c.Hub.Unregister <- c
		return
	}

	fmt.Println(userID, "認証されました")

	c.UserID = userID
	c.Role = msg.Role

	switch msg.Role {
	case "athlete":
		c.Hub.Athletes[c.ID] = c
	case "mark":
		c.Hub.Marks[c.ID] = c
	}
}

func (c *Client) receivePos(msg *Position) {
	c.Position = *msg
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
			if websocket.IsUnexpectedCloseError(
				err,
				websocket.CloseGoingAway,
				websocket.CloseAbnormalClosure,
			) {
				log.Println("UnexpectedCloseError:", err)
			}
			return
		}

		var msg map[string]interface{}
		if err := json.Unmarshal([]byte(msgRaw), &msg); err != nil {
			log.Println(err)
			continue
		}

		switch msg["type"].(string) {
		case "auth":
			var msg AuthInfo
			json.Unmarshal([]byte(msgRaw), &msg)
			c.auth(&msg)

		case "position":
			var msg Position
			json.Unmarshal([]byte(msgRaw), &msg)
			c.receivePos(&msg)
		}
	}
}
