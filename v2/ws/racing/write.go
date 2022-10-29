package racing

import (
	"encoding/json"
	"time"

	"github.com/gorilla/websocket"
)

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

func (c *Client) sendMarkPosMsg() {
	msg := MarkPosMsg{
		MarkNum:   len(c.Hub.Marks),
		Positions: c.Hub.getMarkPositions(),
	}
	c.sendMarkPosMsgEvent(&msg)
}

func (c *Client) sendNearSailMsg() {
	if c.Role != "athlete" {
		return
	}

	msg := NearSailMsg{
		Neighbors: c.getNearSail(),
	}
	c.sendNearSailMsgEvent(&msg)
}

func (c *Client) sendLiveMsg() {
	if c.Role != "manage" {
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
	tickerNearSail := time.NewTicker(markPosPeriod)
	tickerLive := time.NewTicker(livePeriod)

	defer func() {
		ticker.Stop()
		tickerMarkPos.Stop()
		tickerLive.Stop()
		c.Hub.Unregister <- c
	}()

	for {
		select {
		case msg, ok := <-c.Send:
			err := c.sendEvent(msg, ok)
			if err != nil {
				return
			}

		case <-tickerMarkPos.C:
			go c.sendMarkPosMsg()

		case <-tickerNearSail.C:
			go c.sendNearSailMsg()

		case <-tickerLive.C:
			go c.sendLiveMsg()

		case <-ticker.C:
			err := c.pingEvent()
			if err != nil {
				return
			}
		}
	}
}

func insertTypeToJSON(msg interface{}, typeStr string) []byte {
	encoded, _ := json.Marshal(msg)

	text := []byte("\"type\":\"" + typeStr + "\",")

	return append(encoded[:1], append(text, encoded[1:]...)...)
}

func (h *Hub) generateLiveMsg() LiveMsg {
	athletes := make([]LocationWithDetail, len(h.Athletes))
	marks := make([]PositionWithID, len(h.Marks))

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
			MarkNo:        c.MarkNo,
			NextMarkNo:    c.NextMarkNo,
			CourseLimit:   c.CourseLimit,
		}
		cnt++
	}

	for _, c := range h.Marks {
		marks[c.MarkNo-1] = PositionWithID{
			UserID: c.UserID,
			Lat:    c.Position.Lat,
			Lng:    c.Position.Lng,
		}
	}

	return LiveMsg{
		Athletes: athletes,
		Marks:    marks,
	}
}
