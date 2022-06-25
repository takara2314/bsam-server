package race

import (
	"fmt"
	"net/http"
	"sailing-assist-mie-api/abort"
	"sailing-assist-mie-api/bsamdb"
	"sailing-assist-mie-api/message"
	"sailing-assist-mie-api/utils"
	"strconv"
	"strings"

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
	raceID := c.Param("id")
	userID := c.Query("user")

	if _, exist := rooms[raceID]; !exist {
		abort.NotFound(c, message.RaceNotFound)
		return
	}

	// Connect to the database and insert such data.
	db, err := bsamdb.Open()
	if err != nil {
		panic(err)
	}

	var role string

	if !strings.HasPrefix(userID, "NPC") {
		// User ID must be contain.
		if userID == "" {
			abort.BadRequest(c, message.NoUserIDContain)
			return
		}

		// The ID correct check.
		exist, err := db.IsExist(
			"users",
			"id",
			userID,
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
				{Column: "id", Value: userID},
			},
			[]string{"role"},
		)
		if err != nil {
			panic(err)
		}

		rows.Next()
		rows.Scan(&role)
	} else {
		role = "mark"
	}

	// If mark device, register as it.
	pointNo := -1
	if role == "mark" {
		pointNoStr := c.Query("point")

		if pointNoStr != "" {
			pointNo, err = strconv.Atoi(pointNoStr)
			if err != nil {
				abort.BadRequest(c, message.InvalidPointID)
				return
			}

			switch pointNo {
			case 1:
				rooms[raceID].PointA.DeviceID = userID
			case 2:
				rooms[raceID].PointB.DeviceID = userID
			case 3:
				rooms[raceID].PointC.DeviceID = userID
			default:
				abort.BadRequest(c, message.InvalidPointID)
				return
			}
		}
	}

	// Close
	db.DB.Close()

	// Upgrade to WebSocket.
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		abort.BadRequest(c, message.NotSupportWebSocket)
		return
	}

	client := &Client{
		Hub:         rooms[raceID],
		Conn:        conn,
		UserID:      userID,
		Role:        role,
		PointNo:     pointNo,
		NextPoint:   1,
		LatestPoint: 0,
		CourseLimit: 20.0,
		Send:        make(chan *PointNav),
		SendManage:  make(chan *ManageInfo),
		SendLive:    make(chan *LiveInfo),
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
		fmt.Println("distance:", distance)

		if distance < float64(c.CourseLimit) {
			fmt.Println(c.UserID, ">> passed 1")
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
		fmt.Println("distance:", distance)

		if distance < float64(c.CourseLimit) {
			fmt.Println(c.UserID, ">> passed 2")
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
		fmt.Println("distance:", distance)

		if distance < float64(c.CourseLimit) {
			fmt.Println(c.UserID, ">> passed 3")
			c.NextPoint = 1
			c.LatestPoint = 2
		}
	}
}
