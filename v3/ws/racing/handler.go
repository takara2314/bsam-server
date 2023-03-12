package racing

import (
	"bsam-server/v3/abort"
	"fmt"
	"log"
	"net/http"

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

// Handler is a Gin handler for HTTP.
func Handler(c *gin.Context) {
	assocID := c.Param("id")

	// if the room does not exist, return 404
	if _, ok := rooms[assocID]; !ok {
		abort.NotFound(c)
		return
	}

	// Upgrade to WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	client := NewClient(assocID, conn)

	client.Hub.Register <- client

	go client.readPump()
	go client.writePump()
}

// receivePos receives the position from the client.
func (c *Client) receivePos(msg *Position) {
	c.Location = Location{Lat: msg.Lat, Lng: msg.Lng}
}

// receiveLoc receives the location from the client.
func (c *Client) receiveLoc(msg *Location) {
	c.Location = *msg
}

// handlerPassed handles the passed message from the client.
func (c *Client) handlerPassed(msg *PassedInfo) {
	log.Printf("Passed: %s -> [%d]\n", c.UserID, msg.PassedMarkNo)

	c.NextMarkNo = msg.NextMarkNo
}

// receiveBattery receives the battery level from the client.
func (c *Client) receiveBattery(msg *BatteryInfo) {
	c.BatteryLevel = msg.Level
}
