package race

import (
	"fmt"
	"net/http"
	"sailing-assist-mie-api/abort"
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
	fmt.Println("接続を検知しました！", raceId)

	if _, exist := rooms[raceId]; !exist {
		abort.NotFound(c, message.RaceNotFound)
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		abort.BadRequest(c, message.NotSupportWebSocket)
		return
	}

	client := &Client{
		Hub:  rooms[raceId],
		Conn: conn,
		Id:   c.Query("device"),
		Send: make(chan *Point),
	}

	client.Hub.Register <- client

	go client.readPump()
	go client.writePump()
}
