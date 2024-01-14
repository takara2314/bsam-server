package racing

import (
	"log"
	"net/http"

	"bsam-server/utils"
	"bsam-server/v4/abort"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

//nolint:gochecknoglobals
var upgrader = websocket.Upgrader{
	ReadBufferSize:  ReadBufferByte,
	WriteBufferSize: WriteBufferByte,
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
		log.Println("Upgrader error:", err)
		return
	}

	client := NewClient(assocID, conn)

	client.Hub.Register <- client

	go client.readPump()
	go client.writePump()
}

// receivePos receives the position from the client.
func (c *Client) receivePos(msg *Position) {
	c.Location = Location{
		Lat: msg.Lat,
		Lng: msg.Lng,
		Acc: msg.Acc,
	}
	//nolint:errcheck
	go c.Hub.Logger.logLocation(c)
}

// receiveLoc receives the location from the client.
func (c *Client) receiveLoc(msg *Location) {
	c.Location = *msg
	c.CompassDeg = c.calcCompassDeg()
	//nolint:errcheck
	go c.Hub.Logger.logLocation(c)
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

//nolint:gomnd
func (c *Client) calcCompassDeg() float64 {
	if c.NextMarkNo == 0 || c.Location.Acc == 0.0 {
		return 0.0
	}

	marks := c.Hub.getMarkInfos()

	lat1 := c.Location.Lat
	lng1 := c.Location.Lng
	lat2 := marks[c.NextMarkNo-1].Position.Lat
	lng2 := marks[c.NextMarkNo-1].Position.Lng

	bearingDeg := utils.CalcBearingBetweenEarth(lat1, lng1, lat2, lng2)
	diff := bearingDeg - c.Location.Heading

	if diff > 180.0 {
		diff -= 360.0
	} else if diff < -180.0 {
		diff += 360
	}

	return diff
}
