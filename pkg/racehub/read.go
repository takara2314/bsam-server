package racehub

import (
	"log/slog"
	"time"

	"github.com/bytedance/sonic"
	"github.com/gorilla/websocket"
)

type Handler interface {
	Auth(*Client, *AuthInput)
	PostGeolocation(*Client, *PostGeolocationInput)
}

type UnimplementedHandler struct{}

type AuthInput struct {
	MessageType string `json:"type"`
	Token       string `json:"token"`
	DeviceID    string `json:"device_id"`
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

func (UnimplementedHandler) Auth(*Client, *AuthInput) {
	panic("not implemented")
}

func (UnimplementedHandler) PostGeolocation(*Client, *PostGeolocationInput) {
	panic("not implemented")
}

func (c *Client) readPump() {
	defer c.Hub.Unregister(c)

	c.Conn.SetReadLimit(maxIngressMessageBytes)
	if err := c.Conn.SetReadDeadline(time.Now().Add(pongWaitSec)); err != nil {
		slog.Error(
			"failed to set read deadline",
			"client", c,
			"error", err,
		)
		return
	}

	for {
		msgType, payload, err := c.Conn.ReadMessage()
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

		slog.Info(
			"payload received",
			"client", c,
			"type", msgType,
			"payload", string(payload),
		)

		var msg map[string]any
		if err := sonic.Unmarshal(payload, &msg); err != nil {
			slog.Error(
				"failed to unmarshal message",
				"client", c,
				"error", err,
			)
			continue
		}

		slog.Info(
			"payload unmarshaled",
			"client", c,
			"type", msgType,
			"payload", msg,
		)

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
	case "auth":
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

	case "post_geolocation":
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

	default:
		slog.Warn(
			"unknown handler type",
			"client", c,
			"message", msg,
		)
	}
}
