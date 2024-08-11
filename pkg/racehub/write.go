package racehub

import (
	"log/slog"
	"time"

	"github.com/bytedance/sonic"
	"github.com/gorilla/websocket"
)

const (
	ActionTypeAuthResult = "auth_result"

	AuthResultOK              = "OK"
	AuthResultFailedAuthToken = "failed_auth_token"
	AuthResultOutsideAssoc    = "outside_assoc"
	AuthResultInvalidDeviceID = "invalid_device_id"
)

type AuthResultInput struct {
	MessageType string `json:"type"`
	OK          bool   `json:"ok"`
	DeviceID    string `json:"device_id"`
	Role        string `json:"role"`
	MyMarkNo    int    `json:"my_mark_no"`
	Message     string `json:"message"`
}

func (c *Client) writePump() {
	pingTicker := time.NewTicker(pingPeriodSec)
	defer func() {
		pingTicker.Stop()
		c.Hub.Unregister(c)
	}()

	for {
		select {
		case message, ok := <-c.Send:
			if err := c.writeMessage(message, ok); err != nil {
				return
			}

		case <-pingTicker.C:
			slog.Info(
				"sending ping",
				"client", c,
			)
			if err := c.writePing(); err != nil {
				return
			}

		case <-c.StoppingWritePump:
			slog.Info(
				"stopping write pump",
				"client", c,
			)
			return
		}
	}
}

func (c *Client) writeMessage(msg []byte, ok bool) error {
	if err := c.Conn.SetWriteDeadline(time.Now().Add(writeWaitSec)); err != nil {
		slog.Error(
			"failed to set write deadline",
			"client", c,
			"error", err,
		)
		return err
	}
	if !ok {
		return nil
	}

	w, err := c.Conn.NextWriter(websocket.TextMessage)
	if err != nil {
		slog.Error(
			"failed to get next writer",
			"client", c,
			"error", err,
		)
		return err
	}

	if _, err := w.Write(msg); err != nil {
		slog.Error(
			"failed to write message",
			"client", c,
			"error", err,
		)
		return err
	}

	if err := w.Close(); err != nil {
		slog.Error(
			"failed to close writer",
			"client", c,
			"error", err,
		)
		return err
	}

	return nil
}

func (c *Client) writePing() error {
	if err := c.Conn.SetWriteDeadline(time.Now().Add(writeWaitSec)); err != nil {
		slog.Error(
			"failed to set write deadline for ping",
			"client", c,
			"error", err,
		)
		return err
	}
	if err := c.Conn.WriteMessage(websocket.PingMessage, []byte("ping")); err != nil {
		slog.Error(
			"failed to write ping message",
			"error", err,
			"client", c,
		)
		return err
	}

	return nil
}

func (c *Client) SendAuthResult(
	ok bool,
	deviceID string,
	role string,
	myMarkNo int,
	msg string,
) {
	input := AuthResultInput{
		MessageType: ActionTypeAuthResult,
		OK:          ok,
		DeviceID:    deviceID,
		Role:        role,
		MyMarkNo:    myMarkNo,
		Message:     msg,
	}

	payload, err := sonic.Marshal(input)
	if err != nil {
		slog.Error(
			"failed to marshal auth result",
			"client", c,
			"device_id", deviceID,
			"role", role,
			"my_mark_no", myMarkNo,
			"error", err,
		)
		return
	}

	c.Send <- payload

	slog.Info(
		"sent auth result",
		"client", c,
		"input", input,
	)
}
