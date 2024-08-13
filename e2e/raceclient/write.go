package raceclient

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
)

func (c *Client) writePump() {
	defer func() {
		c.Close()
	}()

	for {
		select {
		case message := <-c.sendCh:
			c.mu.RLock()
			err := c.Conn.WriteMessage(websocket.TextMessage, message)
			c.mu.RUnlock()
			if err != nil {
				return
			}
		case <-c.closeCh:
			return
		}
	}
}

func (c *Client) Send(msg any) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("メッセージのシリアライズに失敗しました: %w", err)
	}

	select {
	case c.sendCh <- data:
		return nil
	case <-c.closeCh:
		return ErrClientClosed
	}
}
