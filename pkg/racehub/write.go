package racehub

import (
	"log/slog"
	"time"

	"github.com/bytedance/sonic"
	"github.com/gorilla/websocket"
	"github.com/takara2314/bsam-server/pkg/domain"
)

const (
	ActionTypeAuthResult = "auth_result"

	AuthResultOK              = "OK"
	AuthResultFailedAuthToken = "failed_auth_token"
	AuthResultOutsideAssoc    = "outside_assoc"
	AuthResultInvalidDeviceID = "invalid_device_id"
)

type Action interface {
	AuthResult(
		c *Client,
		ok bool,
		message string,
	) (*AuthResultOutput, error)

	MarkGeolocations(
		c *Client,
	) (*MarkGeolocationsOutput, error)
}

type UnimplementedAction struct{}

type AuthResultOutput struct {
	MessageType string `json:"type"`
	OK          bool   `json:"ok"`
	DeviceID    string `json:"device_id"`
	Role        string `json:"role"`
	MarkNo      int    `json:"mark_no"`
	Message     string `json:"message"`
}

type MarkGeolocationsOutput struct {
	MessageType string                       `json:"type"`
	Marks       []MarkGeolocationsOutputMark `json:"marks"`
}

type MarkGeolocationsOutputMark struct {
	MarkNo        int       `json:"mark_no"`
	Stored        bool      `json:"stored"`
	Latitude      float64   `json:"latitude"`
	Longitude     float64   `json:"longitude"`
	AccuracyMeter float64   `json:"accuracy_meter"`
	Heading       float64   `json:"heading"`
	RecordedAt    time.Time `json:"recorded_at"`
}

func (c *Client) writePump() {
	sendingMarkGeolocationsTicker := time.NewTicker(
		sendingMarkGeolocationsTickerPeriodSec,
	)
	pingTicker := time.NewTicker(pingPeriodSec)

	defer func() {
		sendingMarkGeolocationsTicker.Stop()
		pingTicker.Stop()
		c.Hub.Unregister(c)
	}()

	for {
		select {
		case message, ok := <-c.Send:
			if err := c.writeMessage(message, ok); err != nil {
				return
			}

		case <-sendingMarkGeolocationsTicker.C:
			if err := c.WriteMarkGeolocations(); err != nil {
				return
			}

		case <-pingTicker.C:
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

func (c *Client) writeMessage(msg any, ok bool) error {
	if !ok {
		slog.Info(
			"write pump stopped",
			"client", c,
		)
		return nil
	}

	slog.Info(
		"writing message",
		"client", c,
		"message", msg,
	)

	payload, err := sonic.Marshal(msg)
	if err != nil {
		slog.Error(
			"failed to marshal message",
			"client", c,
			"error", err,
		)
		return err
	}

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

	if _, err := w.Write(payload); err != nil {
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

func (c *Client) WriteAuthResult(
	ok bool,
	message string,
) error {
	slog.Info(
		"writing auth result",
		"client", c,
	)

	output, err := c.Hub.action.AuthResult(
		c,
		ok,
		message,
	)
	if err != nil {
		slog.Error(
			"failed to create auth result output",
			"client", c,
			"error", err,
		)
		return err
	}

	c.Send <- output
	return nil
}

func (c *Client) WriteMarkGeolocations() error {
	slog.Info(
		"writing mark geolocations",
		"client", c,
	)

	if c.Role != domain.RoleMark {
		return nil
	}

	return nil
}

func (c *Client) writePing() error {
	slog.Info(
		"writing ping",
		"client", c,
	)

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
