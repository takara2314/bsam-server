package racing

import (
	"bsam-server/utils"
	"errors"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 10 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	navPeriod      = 5 * time.Second
	maxMessageSize = 1024
)

var (
	ErrClosedChannel = errors.New("closed channel")
)

type Client struct {
	ID          string
	Hub         *Hub
	Conn        *websocket.Conn
	UserID      string
	Role        string
	MarkNo      int
	NextMarkNo  int
	CourseLimit float32
	Position    Position
	Send        chan []byte
}

type Position struct {
	Lat float64 `json:"latitude"`
	Lng float64 `json:"longitude"`
}

type PositionWithID struct {
	UserID string  `json:"user_id"`
	Lat    float64 `json:"latitude"`
	Lng    float64 `json:"longitude"`
}

type PositionWithDetail struct {
	UserID      string  `json:"user_id"`
	Lat         float64 `json:"latitude"`
	Lng         float64 `json:"longitude"`
	LastMarkNo  int     `json:"last_mark_no"`
	CourseLimit float32 `json:"course_limit"`
}

func NewClient(raceID string, conn *websocket.Conn) *Client {
	return &Client{
		ID:          utils.RandString(8),
		Hub:         rooms[raceID],
		Conn:        conn,
		UserID:      "",
		Role:        "",
		MarkNo:      -1,
		NextMarkNo:  -1,
		CourseLimit: 0.0,
		Position:    Position{Lat: 0.0, Lng: 0.0},
		Send:        make(chan []byte),
	}
}
