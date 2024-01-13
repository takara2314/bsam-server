package racing

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type AuthResultMsg struct {
	Authed     bool   `json:"authed"`
	UserID     string `json:"user_id"`
	Role       string `json:"role"`
	MarkNo     int    `json:"mark_no"`
	NextMarkNo int    `json:"next_mark_no"`
	LinkType   string `json:"link_type"`
}

type MarkPosMsg struct {
	MarkNum int    `json:"mark_num"`
	Marks   []Mark `json:"marks"`
}

type NearSailMsg struct {
	Neighbors []Athlete `json:"neighbors"`
}

type LiveMsg struct {
	Athletes []Athlete `json:"athletes"`
	Marks    []Mark    `json:"marks"`
}

type StartRaceMsg struct {
	IsStarted bool  `json:"started"`
	StartAt   int64 `json:"start_at"`
	EndAt     int64 `json:"end_at"`
}

type SetNextMarkNoMsg struct {
	NextMarkNo int `json:"next_mark_no"`
}

// sendMarkPosMsg sends mark positions to the client.
func (c *Client) sendMarkPosMsg() {
	if c.Role == "" {
		return
	}

	msg := MarkPosMsg{
		MarkNum: len(c.Hub.Marks),
		Marks:   c.Hub.getMarkInfos(),
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

// sendLiveMsg sends live positions to the manager and the guest.
func (c *Client) sendLiveMsg() {
	if !(c.Role == "manager" || c.Role == "guest") {
		return
	}

	msg := c.Hub.generateLiveMsg()
	c.sendLiveMsgEvent(&msg)
}

// sendStartRaceMsg sends what started or not to the client.
func (c *Client) sendStartRaceMsg() {
	if c.Role == "" {
		return
	}

	c.sendStartRaceMsgEvent(&StartRaceMsg{
		IsStarted: c.Hub.IsStarted,
		StartAt:   c.Hub.StartAt.UnixNano(),
		EndAt:     c.Hub.EndAt.UnixNano(),
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

func (c *Client) sendSetNextMarkNoEvent(msg *SetNextMarkNoMsg) {
	c.Send <- insertTypeToJSON(msg, "set_next_mark_no")
}

func (c *Client) sendEvent(msg []byte, ok bool) error {
	if !ok {
		err := c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
		if err != nil {
			return err
		}
		return ErrClosedChannel
	}

	err := c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
	if err != nil {
		return err
	}

	w, err := c.Conn.NextWriter(websocket.TextMessage)
	if err != nil {
		return err
	}
	_, err = w.Write(msg)
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) pingEvent() error {
	err := c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
	if err != nil {
		return err
	}

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
	}()

	for {
		select {
		case msg, ok := <-c.Send:
			err := c.sendEvent(msg, ok)
			if err != nil {
				log.Printf("%s (%s) >> write pump error: %v\n", c.ID, c.UserID, err)
				c.Hub.Disconnect <- c
				return
			}

		case <-tickerMarkPos.C:
			// Send mark positions to the client every 5 seconds
			go c.sendMarkPosMsg()

		case <-tickerNearSail.C:
			// Send near sail positions to the athlete every 3 seconds
			go c.sendNearSailMsg()

		case <-tickerLive.C:
			// Send live positions to the manager and guest every 1 second
			go c.sendLiveMsg()

		case <-ticker.C:
			// Ping every 9 seconds
			err := c.pingEvent()
			if err != nil {
				log.Printf("%s (%s) >> ping error: %v\n", c.ID, c.UserID, err)
				c.Hub.Disconnect <- c
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
	return LiveMsg{
		Athletes: h.getAthleteInfos(),
		Marks:    h.getMarkInfos(),
	}
}
