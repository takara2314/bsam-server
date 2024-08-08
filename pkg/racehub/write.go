package racehub

import (
	"log/slog"
	"time"

	"github.com/gorilla/websocket"
)

func (c *Client) writePump() {
	pingTicker := time.NewTicker(pingPeriodSec)
	defer func() {
		pingTicker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			if err := c.Conn.SetWriteDeadline(time.Now().Add(writeWaitSec)); err != nil {
				slog.Error(
					"failed to set write deadline",
					"client", c,
					"error", err,
				)
				return
			}
			if !ok {
				if err := c.Conn.WriteMessage(websocket.CloseMessage, []byte{}); err != nil {
					slog.Error(
						"failed to write close message",
						"client", c,
						"error", err,
					)
				}
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				slog.Error(
					"failed to get next writer",
					"client", c,
					"error", err,
				)
				return
			}

			if _, err := w.Write(message); err != nil {
				slog.Error(
					"failed to write message",
					"client", c,
					"error", err,
				)
				return
			}

			if err := w.Close(); err != nil {
				slog.Error(
					"failed to close writer",
					"client", c,
					"error", err,
				)
				return
			}

		case <-pingTicker.C:
			if err := c.Conn.SetWriteDeadline(time.Now().Add(writeWaitSec)); err != nil {
				slog.Error(
					"failed to set write deadline for ping",
					"client", c,
					"error", err,
				)
				return
			}
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				slog.Error(
					"failed to write ping message",
					"error", err,
					"client", c,
				)
				return
			}
		}
	}
}
