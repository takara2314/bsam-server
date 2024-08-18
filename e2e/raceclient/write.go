package raceclient

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
)

func (c *Client) Send(msg any) error {
	payload, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("メッセージのシリアライズに失敗しました (%s): %w", c.DeviceID, err)
	}

	select {
	case <-c.closeCh:
		return fmt.Errorf("メッセージの送信に失敗しました (%s): %w", c.DeviceID, ErrClientClosed)
	default:
	}

	c.mu.RLock()
	err = c.Conn.WriteMessage(websocket.TextMessage, payload)
	c.mu.RUnlock()
	if err != nil {
		return fmt.Errorf("メッセージの送信に失敗しました (%s): %w", c.DeviceID, err)
	}

	return nil
}
