package racing

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/shiguredo/websocket"
)

var (
	ErrInvalidJWT = errors.New("invalid jwt")
)

type AuthInfo struct {
	Token  string `json:"token"`
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
		}
	}
}

// getUserInfoFromJWT returns user_id and role from JWT token.
func getUserInfoFromJWT(t string) (string, string, error) {
	token, err := jwt.Parse(t, func(token *jwt.Token) (any, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if token == nil || err != nil {
		return "", "", ErrInvalidJWT
	}

	if !token.Valid {
		return "", "", ErrInvalidJWT
	}

	userID := token.Claims.(jwt.MapClaims)["user_id"].(string)
	role := token.Claims.(jwt.MapClaims)["role"].(string)

	return userID, role, nil
}
