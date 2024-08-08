package racehub

import (
	"encoding/json"
	"log/slog"
	"time"

	"github.com/gorilla/websocket"
)

type Handler interface {
	Auth(c *Client)
}

type UnimplementedHandler struct{}

func (UnimplementedHandler) Auth(c *Client) {
	panic("not implemented")
}

func (c *Client) readPump() {
	defer c.Hub.Unregister(c)

	c.Conn.SetReadLimit(maxMessageSizeByte)

	if err := c.Conn.SetReadDeadline(time.Now().Add(pongWaitSec)); err != nil {
		slog.Error(
			"failed to set read deadline",
			"client", c,
			"error", err,
		)
		return
	}

	c.Conn.SetPongHandler(func(string) error {
		return c.Conn.SetReadDeadline(time.Now().Add(pongWaitSec))
	})

	for {
		msgType, payload, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				slog.Error(
					"unexpected close error",
					"client", c,
					"error", err,
				)
			}
			break
		}

		slog.Info(
			"message received",
			"client", c,
			"type", msgType,
			"payload", string(payload),
		)

		var msg map[string]any
		if err := json.Unmarshal(payload, &msg); err != nil {
			slog.Error(
				"failed to unmarshal message",
				"client", c,
				"error", err,
			)
			continue
		}

		handlerType, ok := msg["type"].(string)
		if !ok {
			slog.Warn(
				"invalid message type",
				"client", c,
				"message", msg,
			)
			continue
		}

		switch handlerType {
		case "auth":
			c.Hub.handler.Auth(c)
		default:
			slog.Warn(
				"unknown handler type",
				"client", c,
				"message", msg,
			)
		}
	}
}
