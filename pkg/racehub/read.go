package racehub

import (
	"log/slog"
	"time"

	"github.com/bytedance/sonic"
	"github.com/gorilla/websocket"
)

type Handler interface {
	Auth(c *Client)
	PostGeolocation(c *Client)
}

type UnimplementedHandler struct{}

func (UnimplementedHandler) Auth(c *Client) {
	panic("not implemented")
}

func (UnimplementedHandler) PostGeolocation(c *Client) {
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

		c.routeMessage(msg)
	}
}

func (c *Client) routeMessage(msg map[string]any) {
	handlerType, ok := msg["type"].(string)
	if !ok {
		slog.Warn(
			"invalid message type",
			"client", c,
			"message", msg,
		)
		return
	}

	switch handlerType {
	case "auth":
		c.Hub.handler.Auth(c)
	case "post_geolocation":
		c.Hub.handler.PostGeolocation(c)
	default:
		slog.Warn(
			"unknown handler type",
			"client", c,
			"message", msg,
		)
	}
}
