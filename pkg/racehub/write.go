package racehub

import (
	"log/slog"
	"time"

	"github.com/bytedance/sonic"
	"github.com/gorilla/websocket"
	"github.com/takara2314/bsam-server/pkg/domain"
)

const (
	ActionTypeAuthResult       = "auth_result"
	ActionTypeMarkGeolocations = "mark_geolocations"
	ActionTypeManageRaceStatus = "manage_race_status"
	ActionTypeManageNextMark   = "manage_next_mark"

	AuthResultOK                    = "OK"
	AuthResultFailedAuthToken       = "failed_auth_token"
	AuthResultOutsideAssociation    = "outside_association"
	AuthResultInvalidDeviceID       = "invalid_device_id"
	AuthResultInvalidWantMarkCounts = "invalid_want_mark_counts"
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

	ManageRaceStatus(
		c *Client,
		started bool,
	) (*ManageRaceStatusOutput, error)

	ManageNextMark(
		c *Client,
		nextMarkNo int,
	) (*ManageNextMarkOutput, error)
}

type UnimplementedAction struct{}

type AuthResultOutput struct {
	MessageType string `json:"type"`
	OK          bool   `json:"ok"`
	DeviceID    string `json:"device_id"`
	Role        string `json:"role"`
	MarkNo      int    `json:"mark_no"`
	Authed      bool   `json:"authed"`
	Message     string `json:"message"`
}

type MarkGeolocationsOutput struct {
	MessageType string                       `json:"type"`
	MarkCounts  int                          `json:"mark_counts"`
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

type ManageRaceStatusOutput struct {
	MessageType string `json:"type"`
	Started     bool   `json:"started"`
}

type ManageNextMarkOutput struct {
	MessageType string `json:"type"`
	NextMarkNo  int    `json:"next_mark_no"`
}

func (c *Client) writePump() {
	sendingMarkGeolocationsTicker := time.NewTicker(
		sendingMarkGeolocationsTickerInterval,
	)
	pingTicker := time.NewTicker(pingInterval)

	defer func() {
		sendingMarkGeolocationsTicker.Stop()
		pingTicker.Stop()
		c.Hub.Unregister(c)
	}()

	for {
		select {
		case message, ok := <-c.SendCh:
			if err := c.writeMessage(message, ok); err != nil {
				return
			}

		case <-sendingMarkGeolocationsTicker.C:
			// 選手ロールのみ送信する
			if c.Role != domain.RoleAthlete {
				continue
			}
			if err := c.WriteMarkGeolocations(); err != nil {
				return
			}

		case <-pingTicker.C:
			if err := c.writePing(); err != nil {
				return
			}

		case <-c.StoppingWritePumpCh:
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

	if err := c.Conn.SetWriteDeadline(time.Now().Add(writeTimeout)); err != nil {
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

	slog.Info(
		"sent payload",
		"client", c,
		"payload", string(payload),
	)

	return nil
}

func (c *Client) WriteAuthResult(
	ok bool,
	message string,
) error {
	slog.Info(
		"writing auth_result",
		"client", c,
	)

	output, err := c.Hub.action.AuthResult(
		c,
		ok,
		message,
	)
	if err != nil {
		slog.Error(
			"failed to create auth_result output",
			"client", c,
			"error", err,
		)
		return err
	}

	c.SendCh <- output
	return nil
}

func (c *Client) WriteMarkGeolocations() error {
	slog.Info(
		"writing mark_geolocations",
		"client", c,
	)

	output, err := c.Hub.action.MarkGeolocations(c)
	if err != nil {
		slog.Error(
			"failed to create mark_geolocations output",
			"client", c,
			"error", err,
		)
		return err
	}

	c.SendCh <- output
	return nil
}

func (c *Client) WriteManageRaceStatus(started bool) error {
	slog.Info(
		"writing start race",
		"client", c,
	)

	output, err := c.Hub.action.ManageRaceStatus(c, started)
	if err != nil {
		slog.Error(
			"failed to create manage_race_status output",
			"client", c,
			"error", err,
		)
		return err
	}

	c.SendCh <- output
	return nil
}

func (c *Client) WriteManageNextMark(nextMarkNo int) error {
	slog.Info(
		"writing manage_next_mark",
		"client", c,
	)

	output, err := c.Hub.action.ManageNextMark(c, nextMarkNo)
	if err != nil {
		slog.Error(
			"failed to create manage_next_mark output",
			"client", c,
			"error", err,
		)
		return err
	}

	c.SendCh <- output
	return nil
}

func (c *Client) writePing() error {
	slog.Info(
		"writing ping",
		"client", c,
	)

	if err := c.Conn.SetWriteDeadline(time.Now().Add(writeTimeout)); err != nil {
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
