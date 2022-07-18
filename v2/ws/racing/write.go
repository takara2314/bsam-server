package racing

import (
	"encoding/json"
	"time"

	"github.com/gorilla/websocket"
)

type MarkPosMsg struct {
	Positions []Position `json:"positions"`
}

type NearSailMsg struct {
	Neighbors []PositionWithID `json:"neighbors"`
}

type LiveMsg struct {
	Athletes []PositionWithDetail `json:"athletes"`
	Marks    []PositionWithID     `json:"marks"`
}

func (c *Client) sendMarkPosMsg(msg *MarkPosMsg) {
	c.Send <- insertTypeToJSON(msg, "mark_position")
}

func (c *Client) sendNearSailMsg(msg *NearSailMsg) {
	c.Send <- insertTypeToJSON(msg, "near_sail")
}

func (c *Client) sendLiveMsg(msg *LiveMsg) {
	c.Send <- insertTypeToJSON(msg, "live")
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
	defer func() {
		ticker.Stop()
		c.Hub.Unregister <- c
	}()

	for {
		select {
		case msg, ok := <-c.Send:
			err := c.sendEvent(msg, ok)
			if err != nil {
				return
			}

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
