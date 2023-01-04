package racing

import (
	"encoding/json"
	"sort"
	"time"

	"github.com/shiguredo/websocket"
)

type AuthResultMsg struct {
	Authed   bool   `json:"authed"`
	UserID   string `json:"user_id"`
	Role     string `json:"role"`
	MarkNo   int    `json:"mark_no"`
	LinkType string `json:"link_type"`
}

type MarkPosMsg struct {
	MarkNum   int        `json:"mark_num"`
	Positions []Position `json:"positions"`
}

type NearSailMsg struct {
	Neighbors []PositionWithID `json:"neighbors"`
}

type LiveMsg struct {
	Athletes []LocationWithDetail `json:"athletes"`
	Marks    []PositionWithID     `json:"marks"`
}

type StartRaceMsg struct {
	IsStarted bool `json:"started"`
}

type SetMarkNoMsg struct {
	MarkNo     int `json:"mark_no"`
	NextMarkNo int `json:"next_mark_no"`
}

// sendMarkPosMsg sends mark positions to the client.
func (c *Client) sendMarkPosMsg() {
	msg := MarkPosMsg{
		MarkNum:   len(c.Hub.Marks),
		Positions: c.Hub.getMarkPositions(),
	}
	c.sendMarkPosMsgEvent(&msg)
}

// sendNearSailMsg sends near sail positions to the athlete.
func (c *Client) sendNearSailMsg() {
	if c.Role != "athlete" {
		return
	}

	msg := NearSailMsg{
		Neighbors: c.getNearSail(),
	}
	c.sendNearSailMsgEvent(&msg)
}

// sendLiveMsg sends live positions to the manager.
func (c *Client) sendLiveMsg() {
	if c.Role != "manager" {
		return
	}

	msg := c.Hub.generateLiveMsg()
	c.sendLiveMsgEvent(&msg)
}

func (c *Client) sendStartRaceMsg() {
	c.sendStartRaceMsgEvent(&StartRaceMsg{
		IsStarted: c.Hub.IsStarted,
	})
}

func (c *Client) sendAuthResultMsgEvent(msg *AuthResultMsg) {
	c.Send <- insertTypeToJSON(msg, "auth_result")
}

func (c *Client) sendMarkPosMsgEvent(msg *MarkPosMsg) {
	c.Send <- insertTypeToJSON(msg, "mark_position")
}

func (c *Client) sendNearSailMsgEvent(msg *NearSailMsg) {
	c.Send <- insertTypeToJSON(msg, "near_sail")
}

func (c *Client) sendLiveMsgEvent(msg *LiveMsg) {
	c.Send <- insertTypeToJSON(msg, "live")
}

func (c *Client) sendStartRaceMsgEvent(msg *StartRaceMsg) {
	c.Send <- insertTypeToJSON(msg, "start_race")
}

func (c *Client) sendSetMarkNoEvent(msg *SetMarkNoMsg) {
	c.Send <- insertTypeToJSON(msg, "set_mark_no")
}

func (c *Client) sendEvent(msg []byte, ok bool) error {
	if !ok {
		c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
		return ErrClosedChannel
	}

	c.Conn.SetWriteDeadline(time.Now().Add(writeWait))

	w, err := c.Conn.NextWriter(websocket.TextMessage)
	if err != nil {
		return err
	}
	w.Write(msg)

	err = w.Close()
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) pingEvent() error {
	c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
	return c.Conn.WriteMessage(websocket.PingMessage, nil)
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	tickerMarkPos := time.NewTicker(markPosPeriod)
	tickerNearSail := time.NewTicker(nearSailPeriod)
	tickerLive := time.NewTicker(livePeriod)

	defer func() {
		ticker.Stop()
		tickerMarkPos.Stop()
		tickerLive.Stop()

		if c.Connecting {
			c.Hub.Unregister <- c
		}
	}()

	for {
		// If the client is not connecting, stop the loop
		if !c.Connecting {
			return
		}

		select {
		case msg, ok := <-c.Send:
			err := c.sendEvent(msg, ok)
			if err != nil {
				return
			}

		case <-tickerMarkPos.C:
			// Send mark positions to the client every 5 seconds
			go c.sendMarkPosMsg()

		case <-tickerNearSail.C:
			// Send near sail positions to the athlete every 3 seconds
			go c.sendNearSailMsg()

		case <-tickerLive.C:
			// Send live positions to the manager every 1 second
			go c.sendLiveMsg()

		case <-ticker.C:
			// Ping every 9 seconds
			err := c.pingEvent()
			if err != nil {
				return
			}
		}
	}
}

// insertTypeToJSON inserts message type to rhe JSON which is returned.
func insertTypeToJSON(msg any, typeStr string) []byte {
	encoded, _ := json.Marshal(msg)

	text := []byte("\"type\":\"" + typeStr + "\",")

	return append(encoded[:1], append(text, encoded[1:]...)...)
}

// generateLiveMsg generates live messages.
func (h *Hub) generateLiveMsg() LiveMsg {
	athletes := make([]LocationWithDetail, len(h.Athletes))
	marks := make([]PositionWithID, h.MarkNum)

	cnt := 0
	for _, c := range h.Athletes {
		athletes[cnt] = LocationWithDetail{
			UserID:        c.UserID,
			Lat:           c.Location.Lat,
			Lng:           c.Location.Lng,
			Acc:           c.Location.Acc,
			Heading:       c.Location.Heading,
			HeadingFixing: c.Location.HeadingFixing,
			CompassDeg:    c.Location.CompassDeg,
			NextMarkNo:    c.NextMarkNo,
			CourseLimit:   c.CourseLimit,
		}
		cnt++
	}

	// Sort by user id asc
	sort.Slice(athletes, func(i int, j int) bool {
		return athletes[i].UserID > athletes[j].UserID
	})

	for _, c := range h.Marks {
		marks[c.MarkNo-1] = PositionWithID{
			UserID: c.UserID,
			Lat:    c.Location.Lat,
			Lng:    c.Location.Lng,
		}
	}

	return LiveMsg{
		Athletes: athletes,
		Marks:    marks,
	}
}
