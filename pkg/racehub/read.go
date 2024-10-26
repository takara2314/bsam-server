package racehub

import (
	"log/slog"
	"time"

	"github.com/bytedance/sonic"
	"github.com/gorilla/websocket"
	"github.com/takara2314/bsam-server/pkg/domain"
)

const (
	HandlerTypeAuth             = "auth"
	HandlerTypePostGeolocation  = "post_geolocation"
	HandlerTypePassedMark       = "passed_mark"
	HandlerTypeManageRaceStatus = "manage_race_status"
	HandlerTypeManageNextMark   = "manage_next_mark"
)

type Handler interface {
	Auth(*Client, *AuthInput)
	PostGeolocation(*Client, *PostGeolocationInput)
	PassedMark(*Client, *PassedMarkInput)
	ManageRaceStatus(*Client, *ManageRaceStatusInput)
	ManageNextMark(*Client, *ManageNextMarkInput)
}

type UnimplementedHandler struct{}

type AuthInput struct {
	MessageType    string `json:"type"`
	Token          string `json:"token"`
	DeviceID       string `json:"device_id"`
	WantMarkCounts int    `json:"want_mark_counts"`
}

type PostGeolocationInput struct {
	MessageType           string    `json:"type"`
	Latitude              float64   `json:"latitude"`
	Longitude             float64   `json:"longitude"`
	AltitudeMeter         float64   `json:"altitude_meter"`
	AccuracyMeter         float64   `json:"accuracy_meter"`
	AltitudeAccuracyMeter float64   `json:"altitude_accuracy_meter"`
	Heading               float64   `json:"heading"`
	SpeedMeterPerSec      float64   `json:"speed_meter_per_sec"`
	RecordedAt            time.Time `json:"recorded_at"`
}

type PassedMarkInput struct {
	MessageType  string    `json:"type"`
	PassedMarkNo int       `json:"passed_mark_no"`
	PassedAt     time.Time `json:"passed_at"`
}

type ManageRaceStatusInput struct {
	MessageType string    `json:"type"`
	Started     bool      `json:"started"`
	StartedAt   time.Time `json:"started_at"`
	FinishedAt  time.Time `json:"finished_at"`
}

type ManageNextMarkInput struct {
	MessageType    string `json:"type"`
	TargetDeviceID string `json:"target_device_id"`
	NextMarkNo     int    `json:"next_mark_no"`
}

func (c *Client) readPump() {
	defer c.Hub.Unregister(c)

	c.Conn.SetReadLimit(maxIngressMessageBytes)
	if err := c.Conn.SetReadDeadline(time.Now().Add(pongTimeout)); err != nil {
		slog.Error(
			"failed to set read deadline",
			"client", c,
			"error", err,
		)
		return
	}

	for {
		_, payload, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(
				err,
				websocket.CloseNormalClosure,
				websocket.CloseGoingAway,
				websocket.CloseAbnormalClosure,
				websocket.CloseNoStatusReceived,
			) {
				slog.Error(
					"unexpected close error",
					"client", c,
					"error", err,
				)
			} else {
				slog.Info(
					"client disconnected",
					"client", c,
					"reason", err,
				)
			}
			return
		}

		// slog.Info(
		// 	"payload received",
		// 	"client", c,
		// 	"type", msgType,
		// 	"payload", string(payload),
		// )

		var msg map[string]any
		if err := sonic.Unmarshal(payload, &msg); err != nil {
			slog.Error(
				"failed to unmarshal message",
				"client", c,
				"error", err,
			)
			continue
		}

		// slog.Info(
		// 	"payload unmarshaled",
		// 	"client", c,
		// 	"type", msgType,
		// 	"payload", msg,
		// )

		handlerType, ok := msg["type"].(string)
		if !ok {
			slog.Warn(
				"invalid message type",
				"client", c,
				"message", msg,
			)
			continue
		}

		c.routeMessage(handlerType, payload, msg)
	}
}

func (c *Client) routeMessage(
	handlerType string,
	payload []byte,
	msg map[string]any,
) {
	switch handlerType {
	case HandlerTypeAuth:
		var input AuthInput
		if err := sonic.Unmarshal(payload, &input); err != nil {
			slog.Error(
				"failed to unmarshal auth input",
				"client", c,
				"error", err,
			)
			return
		}
		c.Hub.handler.Auth(c, &input)

	case HandlerTypePostGeolocation:
		// 未認証のクライアントからは受け付けない
		if !c.Authed {
			slog.Warn(
				"not authed client tried to post geolocation",
				"client", c,
			)
			return
		}

		var input PostGeolocationInput
		if err := sonic.Unmarshal(payload, &input); err != nil {
			slog.Error(
				"failed to unmarshal post_geolocation input",
				"client", c,
				"error", err,
			)
			return
		}
		c.Hub.handler.PostGeolocation(c, &input)

	case HandlerTypePassedMark:
		var input PassedMarkInput
		if err := sonic.Unmarshal(payload, &input); err != nil {
			slog.Error(
				"failed to unmarshal passed_mark input",
				"client", c,
				"error", err,
			)
			return
		}
		c.Hub.handler.PassedMark(c, &input)

	case HandlerTypeManageRaceStatus:
		// マネージャ以外のクライアントからは受け付けない
		if c.Role != domain.RoleManager {
			slog.Warn(
				"non-manager client tried to manage race status",
				"client", c,
			)
			return
		}

		var input ManageRaceStatusInput
		if err := sonic.Unmarshal(payload, &input); err != nil {
			slog.Error(
				"failed to unmarshal manage_race_status input",
				"client", c,
				"error", err,
			)
			return
		}
		c.Hub.handler.ManageRaceStatus(c, &input)

	case HandlerTypeManageNextMark:
		// マネージャ以外のクライアントからは受け付けない
		if c.Role != domain.RoleManager {
			slog.Warn(
				"non-manager client tried to manage next mark",
				"client", c,
			)
			return
		}

		var input ManageNextMarkInput
		if err := sonic.Unmarshal(payload, &input); err != nil {
			slog.Error(
				"failed to unmarshal manage_next_mark input",
				"client", c,
				"error", err,
			)
			return
		}
		c.Hub.handler.ManageNextMark(c, &input)

	default:
		slog.Warn(
			"unknown handler type",
			"client", c,
			"message", msg,
		)
	}
}
