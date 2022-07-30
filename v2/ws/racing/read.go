package racing

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/websocket"
)

var (
	ErrInvalidJWT = errors.New("invalid jwt")
)

type AuthInfo struct {
	Token string `json:"token"`
}

type PassedInfo struct {
	MarkNo     int `json:"mark_no"`
	NextMarkNo int `json:"next_mark_no"`
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

		case "passed":
			var msg PassedInfo
			json.Unmarshal([]byte(msgRaw), &msg)
			c.handlerPassed(&msg)
		}
	}
}

func getUserInfoFromJWT(t string) (string, string, int, error) {
	token, err := jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if token == nil || err != nil {
		return "", "", -1, ErrInvalidJWT
	}

	if !token.Valid {
		return "", "", -1, ErrInvalidJWT
	}

	fmt.Println(token.Claims.(jwt.MapClaims))

	userID := token.Claims.(jwt.MapClaims)["user_id"].(string)
	role := token.Claims.(jwt.MapClaims)["role"].(string)
	markNo := int(token.Claims.(jwt.MapClaims)["mark_no"].(float64))

	return userID, role, markNo, nil
}
