package e2e

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"sync"
	"testing"
	"time"

	"github.com/takara2314/bsam-server/e2e/auth"
	"github.com/takara2314/bsam-server/e2e/raceclient"
	"github.com/takara2314/bsam-server/pkg/domain"
	"github.com/takara2314/bsam-server/pkg/racehub"
)

func TestWSManageRaceStatusMultiInstance(t *testing.T) {
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

	errCh := make(chan error)

	// 3つの選手デバイスから接続
	// athlete1 -> server1
	// athlete2 -> server2
	// athlete3 -> server1
	allSentCh := make(chan bool)
	var wg sync.WaitGroup
	for i := range 3 {
		serverURL := serverURL1
		if i%2 == 1 {
			serverURL = serverURL2
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := connectAndReceiveStartRace(
				ctx,
				timeout,
				serverURL,
				token,
				domain.CreateDeviceID(domain.RoleAthlete, i+1),
				1,
			); err != nil {
				errCh <- err
			}
		}()
	}

	go func() {
		wg.Wait()
		close(allSentCh)
	}()

	select {
	case err := <-errCh:
		t.Fatalf("選手デバイスでエラーが発生しました: %v", err)
	default:
	}

	// 本部デバイスからレース開始信号を送信
	// manager1 -> server2
	err = connectAndSendStartRace(
		ctx,
		timeout,
		serverURL2,
		token,
		domain.CreateDeviceID(domain.RoleManager, 1),
	)
	if err != nil {
		t.Fatalf("本部デバイスでエラーが発生しました: %v", err)
	}

	<-allSentCh
}

// WebSocketに接続し、レース開始信号を受信する
func connectAndReceiveStartRace(
	ctx context.Context,
	timeout time.Duration,
	serverURL url.URL,
	token string,
	deviceID string,
	wantMarkCounts int,
) error {
	client := raceclient.NewClient(serverURL, deviceID)

	err := client.Connect(ctx, timeout)
	if err != nil {
		return fmt.Errorf("接続に失敗しました: %v", err)
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

		case "manage_race_status":
			// 開始と確認したときのみ終了
			if msg["started"] == true {
				return nil
			}
		}
	}

	return nil
}

// WebSocketに接続し、レース開始信号を送信する
func connectAndSendStartRace(
	ctx context.Context,
	timeout time.Duration,
	serverURL url.URL,
	token string,
	deviceID string,
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

			// レース開始アクションを送信
			err = client.Send(racehub.ManageRaceStatusInput{
				MessageType: racehub.HandlerTypeManageRaceStatus,
				Started:     true,
			})
			if err != nil {
				return fmt.Errorf("メッセージの送信に失敗しました: %v", err)
			}
			return nil
		}
	}

	return nil
}
