package racing

import (
	"errors"
	"time"

	"bsam-server/utils"

	"github.com/gorilla/websocket"
)

//nolint:gomnd
const (
	AutoRoomingInterval = 30 * time.Second
	ReadBufferSize      = 2048
	WriteBufferSize     = 2048
	writeWait           = 10 * time.Second
	pongWait            = 10 * time.Second
	pingPeriod          = (pongWait * 9) / 10
	MarkNum             = 3
	markPosPeriod       = 5 * time.Second
	nearSailPeriod      = 3 * time.Second
	livePeriod          = 1 * time.Second
	maxMessageSize      = 1024
	nearRange           = 5.0
	ClientIDLength      = 8
	GuestUserIDLength   = 8
	AthleteRoleID       = 0
	MarkRoleID          = 1
	ManagerRoleID       = 2
	GuestRoleID         = 3
	UnknownRoleID       = -1
)

var ErrClosedChannel = errors.New("closed channel")

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
		ID:           utils.RandString(ClientIDLength),
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

		distance := utils.CalcDistanceAtoBEarth(
			c.Location.Lat,
			c.Location.Lng,
			athlete.Location.Lat,
			athlete.Location.Lng,
		)

		if distance < nearRange {
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
		return AthleteRoleID

	case "mark":
		return MarkRoleID

	case "manager":
		return ManagerRoleID

	case "guest":
		return GuestRoleID

	default:
		return UnknownRoleID
	}
}
