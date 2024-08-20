package e2e

import (
	"context"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/takara2314/bsam-server/e2e/raceclient"
)

func TestWSHandshake(t *testing.T) {
	var (
		serverURL = url.URL{
			Scheme: "ws",
			Host:   "localhost:8081",
			Path:   "/japan",
		}
		timeout = 1 * time.Second
	)

	ctx, cancel := context.WithTimeout(
		context.Background(),
		timeout,
	)
	defer cancel()

	client := raceclient.NewClient(serverURL, "")

	err := client.Connect(ctx, timeout)
	if err != nil {
		t.Fatalf("接続に失敗しました: %v", err)
	}
	defer client.Close()

	// レスポンスのステータスコードを確認
	if client.Resp.StatusCode != http.StatusSwitchingProtocols {
		t.Errorf(
			"予期しないステータスコード: got %d, want %d",
			client.Resp.StatusCode,
			http.StatusSwitchingProtocols,
		)
	}

	// Upgradeヘッダーを確認
	if client.Resp.Header.Get("Upgrade") != "websocket" {
		t.Errorf(
			"Upgradeヘッダーが正しくありません: got %s, want websocket",
			client.Resp.Header.Get("Upgrade"),
		)
	}

	// Connectionヘッダーを確認
	if client.Resp.Header.Get("Connection") != "Upgrade" {
		t.Errorf(
			"Connectionヘッダーが正しくありません: got %s, want Upgrade",
			client.Resp.Header.Get("Connection"),
		)
	}
}
