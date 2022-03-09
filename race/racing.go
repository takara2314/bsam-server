package race

import (
	"net/http"
	"sailing-assist-mie-api/abort"
	"sailing-assist-mie-api/bsamdb"
	"sailing-assist-mie-api/message"

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

	// Connect to the database and insert such data.
	db, err := bsamdb.Open()
	if err != nil {
		panic(err)
	}
	defer db.DB.Close()

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

	// If mark device.
	if role == "mark" {
		pointId := c.Param("point")

		switch pointId {
		case "a":
			rooms[raceId].PointA.UserId = userId
		case "b":
			rooms[raceId].PointB.UserId = userId
		case "c":
			rooms[raceId].PointC.UserId = userId
		default:
			abort.BadRequest(c, message.InvalidPointId)
			return
		}
	}

	// Upgrade to WebSocket.
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		abort.BadRequest(c, message.NotSupportWebSocket)
		return
	}

	client := &Client{
		Hub:    rooms[raceId],
		Conn:   conn,
		UserId: userId,
		Role:   role,
		Send:   make(chan *PointNav),
	}

	client.Hub.Register <- client

	go client.readPump()
	go client.writePump()
}
