package racing

import (
	"bsam-server/utils"
	"errors"
	"time"

	"github.com/shiguredo/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 10 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	markPosPeriod  = 5 * time.Second
	nearSailPeriod = 3 * time.Second
	livePeriod     = 1 * time.Second
	maxMessageSize = 1024
	nearRange      = 5.0
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
	Location    Location
	Send        chan []byte
	Connecting  bool
}

type Position struct {
	Lat float64 `json:"latitude"`
	Lng float64 `json:"longitude"`
}

type Location struct {
	Lat           float64 `json:"latitude"`
	Lng           float64 `json:"longitude"`
	Acc           float64 `json:"accuracy"`
	Heading       float64 `json:"heading"`
	HeadingFixing float64 `json:"heading_fixing"`
	CompassDeg    float64 `json:"compass_degree"`
}

type PositionWithID struct {
	UserID string  `json:"user_id"`
	Lat    float64 `json:"latitude"`
	Lng    float64 `json:"longitude"`
}

type LocationWithDetail struct {
	UserID        string  `json:"user_id"`
	Lat           float64 `json:"latitude"`
	Lng           float64 `json:"longitude"`
	Acc           float64 `json:"accuracy"`
	Heading       float64 `json:"heading"`
	HeadingFixing float64 `json:"heading_fixing"`
	CompassDeg    float64 `json:"compass_degree"`
	NextMarkNo    int     `json:"next_mark_no"`
	CourseLimit   float32 `json:"course_limit"`
}

func NewClient(assocID string, conn *websocket.Conn) *Client {
	return &Client{
		ID:          utils.RandString(8),
		Hub:         rooms[assocID],
		Conn:        conn,
		UserID:      "",
		Role:        "",
		MarkNo:      -1,
		NextMarkNo:  1,
		CourseLimit: 0.0,
		Location:    Location{},
		Send:        make(chan []byte),
		Connecting:  true,
	}
}

// getNearSail returns the sail that is near to the client.
func (c *Client) getNearSail() []PositionWithID {
	var result []PositionWithID

	for _, athlete := range c.Hub.Athletes {
		if c.UserID == athlete.UserID {
			continue
		}

		if utils.CalcDistanceAtoBEarth(c.Location.Lat, c.Location.Lng, athlete.Location.Lat, athlete.Location.Lng) < nearRange {
			result = append(
				result,
				PositionWithID{
					UserID: athlete.UserID,
					Lat:    athlete.Location.Lat,
					Lng:    athlete.Location.Lng,
				},
			)
		}
	}

	return result
}
