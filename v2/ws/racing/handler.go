package racing

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  2048,
	WriteBufferSize: 2048,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func Handler(c *gin.Context) {
	raceID := c.Param("id")

	// Upgrade to WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	client := NewClient(raceID, conn)

	client.Hub.Register <- client

	go client.readPump()
	go client.writePump()
}

func (c *Client) auth(msg *AuthInfo) {
	userID, role, err := getUserInfoFromJWT(msg.Token)
	if err != nil {
		log.Println("Unauthorized:", c.ID)
		c.Hub.Unregister <- c
		return
	}

	log.Printf("Linked: %s <=> %s (%s)\n", c.ID, userID, role)

	c.UserID = userID
	c.Role = role

	switch role {
	case "athlete":
		c.Hub.Athletes[c.ID] = c
	case "mark":
		c.MarkNo = msg.MarkNo
		c.Hub.Marks[c.ID] = c
	}

	c.sendMarkPosMsg()
}

func (c *Client) receivePos(msg *Position) {
	c.Position = *msg
}

func (c *Client) handlerPassed(msg *PassedInfo) {
	log.Printf("Passed: [%d] -> %s -> [%d]\n", msg.MarkNo, c.UserID, msg.NextMarkNo)

	c.MarkNo = msg.MarkNo
	c.NextMarkNo = msg.NextMarkNo
}
