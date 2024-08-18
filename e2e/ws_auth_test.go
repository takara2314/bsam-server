package e2e

import (
	"context"
	"encoding/json"
	"errors"
	"net/url"
	"testing"
	"time"

	"github.com/takara2314/bsam-server/e2e/auth"
	"github.com/takara2314/bsam-server/e2e/raceclient"
	"github.com/takara2314/bsam-server/pkg/racehub"
)

func TestWSAuth(t *testing.T) {
	var (
		serverURL = url.URL{
			Scheme: "ws",
			Host:   "localhost:8081",
			Path:   "/japan",
		}
		associationID = "japan"
		password      = "nippon"
		timeout       = 1 * time.Second
	)

	ctx, cancel := context.WithTimeout(
		context.Background(),
		timeout,
	)
	defer cancel()

	// トークンを取得
	token, err := auth.FetchTokenFromAPI(associationID, password)
	if err != nil {
		t.Fatalf("トークンの取得に失敗しました: %v", err)
	}

	client := raceclient.NewClient(serverURL, "athlete1")

	err = client.Connect(ctx, timeout)
	if err != nil {
		t.Fatalf("接続に失敗しました: %v", err)
	}
	defer client.Close()

	// 認証メッセージを送信
	err = client.Send(racehub.AuthInput{
		MessageType:    racehub.HandlerTypeAuth,
		Token:          token,
		DeviceID:       "athlete1",
		WantMarkCounts: 3,
	})
	if err != nil {
		t.Fatalf("メッセージの送信に失敗しました: %v", err)
	}

	it := client.ReceiveStream(ctx)
	for {
		payload, err := it.Read()
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				break
			}
			t.Fatalf("メッセージの受信に失敗しました: %v", err)
		}

		var msg map[string]any
		err = json.Unmarshal(payload, &msg)
		if err != nil {
			t.Fatalf("メッセージのパースに失敗しました: %v", err)
		}

		switch msg["type"] {
		case "auth_result":
			if msg["ok"] != true {
				t.Errorf(
					"認証に失敗しました: got %v want ok",
					msg["message"],
				)
			}
			return
		}
	}
}
