package e2e

import (
	"context"
	"errors"
	"net/url"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/takara2314/bsam-server/e2e/raceclient"
)

func TestWSReceivePong(t *testing.T) {
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

	// Pongを受信したら成功
	pongReceivedCh := make(chan bool)
	client.Conn.SetPongHandler(func(appData string) error {
		close(pongReceivedCh)
		return nil
	})

	// Pingを送信
	err = client.Conn.WriteControl(
		websocket.PingMessage,
		[]byte("ping"),
		time.Now().Add(time.Second),
	)
	if err != nil {
		t.Fatalf("Ping送信エラー: %v", err)
	}

	errCh := make(chan error)

	// Pongを受信するまで待機
	go func(ctx context.Context, errCh chan error) {
		it := client.ReceiveStream()
		for {
			_, err := it.Read()
			if err != nil {
				if errors.Is(err, context.DeadlineExceeded) {
					break
				}
				errCh <- err
			}
		}
	}(ctx, errCh)

	select {
	case <-pongReceivedCh:
		return
	case err := <-errCh:
		t.Fatalf("メッセージの受信に失敗しました: %v", err)
	case <-time.After(timeout):
		t.Fatal("タイムアウトしました")
	}
}
