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

func TestWSPostGeolocation(t *testing.T) {
	var (
		serverURL = url.URL{
			Scheme: "ws",
			Host:   "localhost:8081",
			Path:   "/japan",
		}
		associationID = "japan"
		password      = "nippon"
		geolocations  = []racehub.PostGeolocationInput{
			{
				MessageType:   racehub.HandlerTypePostGeolocation,
				Latitude:      35.6895,
				Longitude:     139.6917,
				AccuracyMeter: 10.0,
				Heading:       1.0,
				RecordedAt:    time.Now(),
			},
			{
				MessageType:   racehub.HandlerTypePostGeolocation,
				Latitude:      36.6895,
				Longitude:     140.6917,
				AccuracyMeter: 11.0,
				Heading:       2.0,
				RecordedAt:    time.Now(),
			},
			{
				MessageType:   racehub.HandlerTypePostGeolocation,
				Latitude:      37.6895,
				Longitude:     141.6917,
				AccuracyMeter: 12.0,
				Heading:       3.0,
				RecordedAt:    time.Now(),
			},
		}
		wantMarkCounts = 5
		timeout        = 1 * time.Second
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

	// 3つのマークデバイスから位置情報を送信
	var wg sync.WaitGroup
	for i, geolocation := range geolocations {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := connectAndPostGeolocation(
				ctx,
				timeout,
				serverURL,
				token,
				domain.CreateDeviceID(domain.RoleMark, i+1),
				geolocation,
			); err != nil {
				errCh <- err
			}
		}()
	}

	wg.Wait()

	select {
	case err := <-errCh:
		t.Fatalf("マークから位置情報を送信しているときにエラーが発生しました: %v", err)
	default:
	}

	// マークの位置情報を取得
	marks, err := connectAndReceiveMarkGeolocations(
		ctx,
		timeout,
		serverURL,
		token,
		domain.CreateDeviceID(domain.RoleAthlete, 1),
		wantMarkCounts,
	)
	if err != nil {
		t.Fatalf("マークの位置情報の取得に失敗しました: %v", err)
	}

	if len(marks) != wantMarkCounts {
		t.Errorf(
			"位置情報の数が正しくありません: got %d, want %d",
			len(marks),
			wantMarkCounts,
		)
	}

	// 位置情報の格納状態を確認
	// mark1, mark2, mark3 のみ送信したため、それ以降は格納されていない
	correctStored := []bool{true, true, true, false, false}
	for i, mark := range marks {
		if mark.Stored != correctStored[i] {
			t.Errorf(
				"マーク%dの位置情報の格納状態が正しくありません: got %v, want %v",
				i+1,
				mark.Stored,
				correctStored[i],
			)
		}
	}
}

// WebSocketに接続し、位置情報を送信する
func connectAndPostGeolocation(
	ctx context.Context,
	timeout time.Duration,
	serverURL url.URL,
	token string,
	deviceID string,
	geolocation racehub.PostGeolocationInput,
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

	it := client.ReceiveStream(ctx)
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

			// 位置情報を送信
			err = client.Send(geolocation)
			if err != nil {
				return fmt.Errorf("メッセージの送信に失敗しました: %v", err)
			}
			return nil
		}
	}

	return errors.New("不正に閉じられました")
}

// WebSocketに接続し、位置情報を受信する
func connectAndReceiveMarkGeolocations(
	ctx context.Context,
	timeout time.Duration,
	serverURL url.URL,
	token string,
	deviceID string,
	wantMarkCounts int,
) ([]racehub.MarkGeolocationsOutputMark, error) {
	client := raceclient.NewClient(serverURL, deviceID)

	err := client.Connect(ctx, timeout)
	if err != nil {
		return nil, fmt.Errorf("接続に失敗しました: %v", err)
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
		return nil, fmt.Errorf("メッセージの送信に失敗しました: %v", err)
	}

	it := client.ReceiveStream(ctx)
	for {
		payload, err := it.Read()
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				break
			}
			return nil, fmt.Errorf("メッセージの受信に失敗しました: %v", err)
		}

		var msg map[string]any
		err = json.Unmarshal(payload, &msg)
		if err != nil {
			return nil, fmt.Errorf("メッセージのパースに失敗しました: %v", err)
		}

		switch msg["type"] {
		case "auth_result":
			if msg["ok"] != true {
				return nil, fmt.Errorf(
					"認証に失敗しました: got %v want ok",
					msg["message"],
				)
			}

		case "mark_geolocations":
			var output racehub.MarkGeolocationsOutput
			err = json.Unmarshal(payload, &output)
			if err != nil {
				return nil, fmt.Errorf("メッセージのパースに失敗しました: %v", err)
			}

			return output.Marks, nil
		}
	}

	return nil, nil
}
