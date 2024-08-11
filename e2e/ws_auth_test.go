package e2e

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/bytedance/sonic"
	"github.com/gorilla/websocket"
	"github.com/takara2314/bsam-server/pkg/racehub"
)

func TestWSAuth(t *testing.T) {
	var (
		url        = "ws://localhost:8081/japan"
		assocID    = "japan"
		password   = "nippon"
		timeoutSec = 1 * time.Second
	)

	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSec)
	defer cancel()

	// トークンを取得
	token, err := FetchTokenFromAPI(assocID, password)
	if err != nil {
		t.Fatalf("トークンの取得に失敗しました: %v", err)
	}

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

	// Authメッセージを送信
	authInput := racehub.AuthInput{
		MessageType:    racehub.HandlerTypeAuth,
		Token:          token,
		DeviceID:       "athlete1",
		WantMarkCounts: 3,
	}

	payload, err := sonic.Marshal(&authInput)
	if err != nil {
		t.Fatalf("メッセージのエンコードに失敗しました: %v", err)
	}

	if err := conn.WriteMessage(
		websocket.TextMessage,
		payload,
	); err != nil {
		t.Fatalf("メッセージの送信に失敗しました: %v", err)
	}

	authResultReceived := make(chan bool)

	go func() {
		for {
			_, payload, err := conn.ReadMessage()
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

			var msg *racehub.AuthResultOutput
			err = sonic.Unmarshal(payload, &msg)
			if err != nil {
				t.Errorf("メッセージのデコードに失敗しました: %v", err)
			}

			if msg.MessageType != racehub.ActionTypeAuthResult {
				t.Errorf(
					"メッセージタイプが正しくありません: got %s, want %s",
					msg.MessageType,
					racehub.ActionTypeAuthResult,
				)
			}

			if msg.Message != racehub.AuthResultOK {
				t.Errorf(
					"認証メッセージが正しくありません: got %v, want %s",
					msg.Message,
					racehub.AuthResultOK,
				)
			}

			authResultReceived <- true
		}
	}()

	select {
	case <-authResultReceived:
		return
	case <-time.After(timeoutSec):
		t.Fatal("AuthResult待機中にタイムアウトしました")
	}
}
