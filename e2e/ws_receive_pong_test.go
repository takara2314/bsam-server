package e2e

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func TestWSReceivePong(t *testing.T) {
	var (
		url        = "ws://localhost:8081/japan"
		timeoutSec = 1 * time.Second
	)

	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSec)
	defer cancel()

	// WebSocket Dialerの設定
	dialer := websocket.Dialer{
		HandshakeTimeout: timeoutSec,
	}

	// WebSocket接続を試みる
	conn, resp, err := dialer.DialContext(ctx, url, nil)
	if err != nil {
		t.Fatalf("WebSocket接続に失敗しました: %v", err)
	}
	defer conn.Close()

	if resp.StatusCode != http.StatusSwitchingProtocols {
		t.Errorf(
			"予期しないステータスコード: got %d, want %d",
			resp.StatusCode, http.StatusSwitchingProtocols,
		)
	}

	pongReceived := make(chan bool)
	conn.SetPongHandler(func(appData string) error {
		close(pongReceived)
		return nil
	})

	err = conn.WriteControl(
		websocket.PingMessage,
		[]byte("test ping"),
		time.Now().Add(time.Second),
	)
	if err != nil {
		t.Fatalf("Ping送信エラー: %v", err)
	}

	go func() {
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(
					err,
					websocket.CloseNormalClosure,
					websocket.CloseGoingAway,
					websocket.CloseAbnormalClosure,
					websocket.CloseNoStatusReceived,
				) {
					t.Errorf("予期しない接続クローズエラー: %v", err)
				}
				return
			}
		}
	}()

	select {
	case <-pongReceived:
		return
	case <-time.After(timeoutSec):
		t.Fatal("Pong待機中にタイムアウトしました")
	}
}
