package racing

import (
	"fmt"
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
	userID, role, markNo, err := getUserInfoFromJWT(msg.Token)
	if err != nil {
		fmt.Println("認証に失敗しました。")
		c.Hub.Unregister <- c
		return
	}

	c.UserID = userID
	c.Role = role
	c.MarkNo = markNo

	fmt.Println(c.UserID, "さん:", c.Role, c.MarkNo)

	switch role {
	case "athlete":
		c.Hub.Athletes[c.ID] = c
	case "mark":
		c.Hub.Marks[c.ID] = c
	}

	c.sendMarkPosMsg()
}

func (c *Client) receivePos(msg *Position) {
	c.Position = *msg
}

func (c *Client) handlerPassed(msg *PassedInfo) {
	c.MarkNo = msg.MarkNo
	c.NextMarkNo = msg.NextMarkNo
}
