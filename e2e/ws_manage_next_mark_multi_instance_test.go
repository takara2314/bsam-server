package e2e

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"testing"
	"time"

	"github.com/bytedance/sonic"
	"github.com/takara2314/bsam-server/e2e/auth"
	"github.com/takara2314/bsam-server/e2e/raceclient"
	"github.com/takara2314/bsam-server/pkg/domain"
	"github.com/takara2314/bsam-server/pkg/racehub"
)

func TestWSManageNextMarkMultiInstance(t *testing.T) {
	var (
		serverURL1 = url.URL{
			Scheme: "ws",
			Host:   "localhost:8081",
			Path:   "/japan",
		}
		serverURL2 = url.URL{
			Scheme: "ws",
			Host:   "localhost:8181",
			Path:   "/japan",
		}
		associationID      = "japan"
		password           = "nippon"
		expectedNextMarkNo = 2
		timeout            = 1 * time.Second
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

	receivedNextMarkNo := make(chan int)
	errCh := make(chan error)

	// 選手デバイスから接続
	// athlete1 -> server1
	go func() {
		nextMarkNo, err := connectAndReceiveNextMark(
			ctx,
			timeout,
			serverURL1,
			token,
			domain.CreateDeviceID(domain.RoleAthlete, 1),
			3,
		)
		if err != nil {
			errCh <- err
		}

		receivedNextMarkNo <- nextMarkNo
	}()

	time.Sleep(10 * time.Millisecond)

	// 本部デバイスから次のマークを送信
	// manager1 -> server2
	go func() {
		if err := connectAndSendNextMark(
			ctx,
			timeout,
			serverURL2,
			token,
			domain.CreateDeviceID(domain.RoleManager, 1),
			domain.CreateDeviceID(domain.RoleAthlete, 1),
			expectedNextMarkNo,
		); err != nil {
			errCh <- err
		}
	}()

	select {
	case err := <-errCh:
		t.Fatal(err)
	case nextMarkNo := <-receivedNextMarkNo:
		if nextMarkNo != expectedNextMarkNo {
			errCh <- fmt.Errorf(
				"次のマーク番号が正しくありません: got %d, want %d",
				nextMarkNo,
				expectedNextMarkNo,
			)
		}
	}
}

// WebSocketに接続し、次のマークを受信する
func connectAndReceiveNextMark(
	ctx context.Context,
	timeout time.Duration,
	serverURL url.URL,
	token string,
	deviceID string,
	wantMarkCounts int,
) (int, error) {
	client := raceclient.NewClient(serverURL, deviceID)

	err := client.Connect(ctx, timeout)
	if err != nil {
		return -1, fmt.Errorf("接続に失敗しました: %v", err)
	}
	defer client.Close()

	// 認証メッセージを送信
	err = client.Send(racehub.AuthInput{
		MessageType:    racehub.HandlerTypeAuth,
		Token:          token,
		DeviceID:       deviceID,
		WantMarkCounts: wantMarkCounts,
	})
	if err != nil {
		return 0, fmt.Errorf("メッセージの送信に失敗しました: %v", err)
	}

	it := client.ReceiveStream()
	for {
		payload, err := it.Read()
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				break
			}
			return 0, fmt.Errorf("メッセージの受信に失敗しました: %v", err)
		}

		var msg map[string]any
		err = json.Unmarshal(payload, &msg)
		if err != nil {
			return 0, fmt.Errorf("メッセージのパースに失敗しました: %v", err)
		}

		switch msg["type"] {
		case "auth_result":
			if msg["ok"] != true {
				return 0, fmt.Errorf(
					"認証に失敗しました: got %v want ok",
					msg["message"],
				)
			}

		case "manage_next_mark":
			var parsed racehub.ManageNextMarkOutput
			err := sonic.Unmarshal(payload, &parsed)
			if err != nil {
				return 0, fmt.Errorf("メッセージのパースに失敗しました: %v", err)
			}
			return parsed.NextMarkNo, nil
		}
	}

	return 0, nil
}

// WebSocketに接続し、次のマークを送信する
func connectAndSendNextMark(
	ctx context.Context,
	timeout time.Duration,
	serverURL url.URL,
	token string,
	deviceID string,
	targetDeviceID string,
	nextMarkNo int,
) error {
	client := raceclient.NewClient(serverURL, deviceID)

	err := client.Connect(ctx, timeout)
	if err != nil {
		return fmt.Errorf("接続に失敗しました: %v", err)
	}
	defer client.Close()

	// 認証メッセージを送信
	err = client.Send(racehub.AuthInput{
		MessageType: racehub.HandlerTypeAuth,
		Token:       token,
		DeviceID:    deviceID,
	})
	if err != nil {
		return fmt.Errorf("メッセージの送信に失敗しました: %v", err)
	}

	it := client.ReceiveStream()
	for {
		payload, err := it.Read()
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				break
			}
			return fmt.Errorf("メッセージの受信に失敗しました: %v", err)
		}

		var msg map[string]any
		err = json.Unmarshal(payload, &msg)
		if err != nil {
			return fmt.Errorf("メッセージのパースに失敗しました: %v", err)
		}

		switch msg["type"] {
		case "auth_result":
			if msg["ok"] != true {
				return fmt.Errorf(
					"認証に失敗しました: got %v want ok",
					msg["message"],
				)
			}

			// 次のマークを送信
			err = client.Send(racehub.ManageNextMarkInput{
				MessageType:    racehub.HandlerTypeManageNextMark,
				TargetDeviceID: targetDeviceID,
				NextMarkNo:     nextMarkNo,
			})
			if err != nil {
				return fmt.Errorf("メッセージの送信に失敗しました: %v", err)
			}

			return nil
		}
	}

	return fmt.Errorf("WebSocket接続中にタスクを達成できませんでした")
}
