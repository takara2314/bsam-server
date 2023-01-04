package racing

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/shiguredo/websocket"
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

	nextMarkNo, err := strconv.Atoi(c.Param("next_mark_no"))
	if err != nil {
		nextMarkNo = 1
	}

	// Upgrade to WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	client := NewClient(raceID, conn, nextMarkNo)

	client.Hub.Register <- client

	go client.readPump()
	go client.writePump()
}

func (c *Client) receivePos(msg *Position) {
	c.Position = *msg
	c.Location = Location{Lat: msg.Lat, Lng: msg.Lng}
}

func (c *Client) receiveLoc(msg *Location) {
	c.Position = Position{Lat: msg.Lat, Lng: msg.Lng}
	c.Location = *msg
}

func (c *Client) handlerPassed(msg *PassedInfo) {
	log.Printf("Passed: [%d] -> %s -> [%d]\n", msg.MarkNo, c.UserID, msg.NextMarkNo)

	c.MarkNo = msg.MarkNo
	c.NextMarkNo = msg.NextMarkNo
}
