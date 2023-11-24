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
	ID           string
	Hub          *Hub
	Conn         *websocket.Conn
	UserID       string
	Role         string
	MarkNo       int
	NextMarkNo   int
	CompassDeg   float64
	CourseLimit  float32
	Location     Location
	BatteryLevel int
	Send         chan []byte
}

type Position struct {
	Lat float64 `json:"latitude"`
	Lng float64 `json:"longitude"`
	Acc float64 `json:"accuracy"`
}

type Location struct {
	Lat           float64 `json:"latitude"`
	Lng           float64 `json:"longitude"`
	Acc           float64 `json:"accuracy"`
	Heading       float64 `json:"heading"`
	HeadingFixing float64 `json:"heading_fixing"`
}

type Athlete struct {
	UserID       string   `json:"user_id"`
	NextMarkNo   int      `json:"next_mark_no"`
	CourseLimit  float32  `json:"course_limit"`
	BatteryLevel int      `json:"battery_level"`
	CompassDeg   float64  `json:"compass_degree"`
	Location     Location `json:"location"`
}

type Mark struct {
	UserID       string   `json:"user_id"`
	MarkNo       int      `json:"mark_no"`
	BatteryLevel int      `json:"battery_level"`
	Position     Position `json:"position"`
}

func NewClient(assocID string, conn *websocket.Conn) *Client {
	return &Client{
		ID:           utils.RandString(8),
		Hub:          rooms[assocID],
		Conn:         conn,
		UserID:       "",
		Role:         "",
		MarkNo:       -1,
		NextMarkNo:   1,
		CourseLimit:  0.0,
		Location:     Location{},
		BatteryLevel: -1,
		Send:         make(chan []byte),
	}
}

// getNearSail returns the sail that is near to the client.
func (c *Client) getNearSail() []Athlete {
	var result []Athlete

	for _, athlete := range c.Hub.Athletes {
		if c.UserID == athlete.UserID {
			continue
		}

		if utils.CalcDistanceAtoBEarth(c.Location.Lat, c.Location.Lng, athlete.Location.Lat, athlete.Location.Lng) < nearRange {
			result = append(
				result,
				Athlete{
					UserID:       athlete.UserID,
					NextMarkNo:   athlete.NextMarkNo,
					CourseLimit:  athlete.CourseLimit,
					BatteryLevel: athlete.BatteryLevel,
					CompassDeg:   athlete.CompassDeg,
					Location:     athlete.Location,
				},
			)
		}
	}

	return result
}

func (c *Client) getRoleID() int {
	switch c.Role {
	case "athlete":
		return 0
	case "mark":
		return 1
	case "manager":
		return 2
	case "guest":
		return 3
	}
	return -1
}
