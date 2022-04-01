package race

import (
	"fmt"
	"net/http"
	"sailing-assist-mie-api/abort"
	"sailing-assist-mie-api/bsamdb"
	"sailing-assist-mie-api/message"
	"sailing-assist-mie-api/utils"
	"strconv"

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

func RacingWS(c *gin.Context) {
	raceId := c.Param("id")
	userId := c.Query("user")

	if _, exist := rooms[raceId]; !exist {
		abort.NotFound(c, message.RaceNotFound)
		return
	}

	fmt.Println("A")

	// Connect to the database and insert such data.
	db, err := bsamdb.Open()
	if err != nil {
		panic(err)
	}

	// User ID must be contain.
	if userId == "" {
		abort.BadRequest(c, message.NoUserIdContain)
		return
	}

	// The ID correct check.
	exist, err := db.IsExist(
		"users",
		"id",
		userId,
	)
	if err != nil {
		panic(err)
	}
	if !exist {
		abort.BadRequest(c, message.UserNotFound)
		return
	}

	// Obtain a role.
	rows, err := db.SelectSpecified(
		"users",
		[]bsamdb.Field{
			{Column: "id", Value: userId},
		},
		[]string{"role"},
	)
	if err != nil {
		panic(err)
	}

	rows.Next()
	var role string
	rows.Scan(&role)

	// If mark device, register as it.
	if role == "mark" {
		pointNoStr := c.Param("point")
		pointNo, err := strconv.Atoi(pointNoStr)
		if err != nil {
			abort.BadRequest(c, message.InvalidPointId)
			return
		}

		switch pointNo {
		case 1:
			rooms[raceId].PointA.UserId = userId
		case 2:
			rooms[raceId].PointB.UserId = userId
		case 3:
			rooms[raceId].PointC.UserId = userId
		default:
			abort.BadRequest(c, message.InvalidPointId)
			return
		}
	}

	// Close
	db.DB.Close()

	fmt.Println("B")

	// Upgrade to WebSocket.
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		abort.BadRequest(c, message.NotSupportWebSocket)
		return
	}

	client := &Client{
		Hub:         rooms[raceId],
		Conn:        conn,
		UserId:      userId,
		Role:        role,
		NextPoint:   1,
		LatestPoint: 0,
		CourseLimit: 20.0,
		Send:        make(chan *PointNav),
	}

	fmt.Println("C")

	var isOpened bool
	select {
	case _, isOpened = <-client.Hub.Register:
	default:
		isOpened = true
	}
	if !isOpened {
		fmt.Println("しまってます…！")
	}

	client.Hub.Register <- client

	go client.readPump()
	go client.writePump()
}

// passCheck checks that the user passed the mark point.
func (c *Client) passCheck() {
	switch c.NextPoint {
	case 1:
		distance := utils.CalcDistanceAtoBEarth(
			c.Position.Latitude,
			c.Position.Longitude,
			c.Hub.PointA.Latitude,
			c.Hub.PointA.Longitude,
		)

		if distance < float64(c.CourseLimit) {
			c.NextPoint = 2
			c.LatestPoint = 1
		}

	case 2:
		distance := utils.CalcDistanceAtoBEarth(
			c.Position.Latitude,
			c.Position.Longitude,
			c.Hub.PointB.Latitude,
			c.Hub.PointB.Longitude,
		)

		if distance < float64(c.CourseLimit) {
			c.NextPoint = 3
			c.LatestPoint = 2
		}

	case 3:
		distance := utils.CalcDistanceAtoBEarth(
			c.Position.Latitude,
			c.Position.Longitude,
			c.Hub.PointC.Latitude,
			c.Hub.PointC.Longitude,
		)

		if distance < float64(c.CourseLimit) {
			c.NextPoint = 1
			c.LatestPoint = 2
		}
	}
}
