package e2e

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func TestWSHandshake(t *testing.T) {
	var (
		url        = "ws://localhost:8081/nippon"
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

	// レスポンスのステータスコードを確認
	if resp.StatusCode != http.StatusSwitchingProtocols {
		t.Errorf("予期しないステータスコード: got %d, want %d", resp.StatusCode, http.StatusSwitchingProtocols)
	}

	// Upgradeヘッダーを確認
	if resp.Header.Get("Upgrade") != "websocket" {
		t.Errorf(
			"Upgradeヘッダーが正しくありません: got %s, want websocket",
			resp.Header.Get("Upgrade"),
		)
	}

	// Connectionヘッダーを確認
	if resp.Header.Get("Connection") != "Upgrade" {
		t.Errorf(
			"Connectionヘッダーが正しくありません: got %s, want Upgrade",
			resp.Header.Get("Connection"),
		)
	}

	// 正常にクローズ
	if err := conn.WriteMessage(
		websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
	); err != nil {
		t.Errorf("コネクションのクローズ中にエラーが発生しました: %v", err)
	}
}
