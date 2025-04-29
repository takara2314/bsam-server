package racing

import (
	"errors"
	"time"

	"github.com/takara2314/bsam-server/utils"

	"github.com/gorilla/websocket"
)

//nolint:gomnd
const (
	AutoRoomingInterval = 30 * time.Second
	ReadBufferByte      = 2048
	WriteBufferByte     = 2048
	WriteWait           = 10 * time.Second
	PongWait            = 10 * time.Second
	PingPeriod          = (PongWait * 9) / 10
	MarkNum             = 3
	MarkPosPeriod       = 5 * time.Second
	NearSailPeriod      = 3 * time.Second
	LivePeriod          = 1 * time.Second
	MaxMessageByte      = 1024
	NearRangeMeter      = 5.0
	ClientIDLength      = 8
	GuestUserIDLength   = 8
	AthleteRoleID       = 0
	MarkRoleID          = 1
	ManagerRoleID       = 2
	GuestRoleID         = 3
	UnknownRoleID       = -1
	AthleteRole         = "athlete"
	MarkRole            = "mark"
	ManagerRole         = "manager"
	GuestRole           = "guest"
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

		if distance < NearRangeMeter {
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
	case AthleteRole:
		return AthleteRoleID
	case MarkRole:
		return MarkRoleID
	case ManagerRole:
		return ManagerRoleID
	case GuestRole:
		return GuestRoleID
	default:
		return UnknownRoleID
	}
}
