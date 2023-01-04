package racing

import (
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
	c.Position = *msg
	c.Location = Location{Lat: msg.Lat, Lng: msg.Lng}
}

// receiveLoc receives the location from the client.
func (c *Client) receiveLoc(msg *Location) {
	c.Position = Position{Lat: msg.Lat, Lng: msg.Lng}
	c.Location = *msg
}

// handlerPassed handles the passed message from the client.
func (c *Client) handlerPassed(msg *PassedInfo) {
	log.Printf("Passed: %s -> [%d]\n", c.UserID, msg.NextMarkNo)

	c.NextMarkNo = msg.NextMarkNo
}
