package e2e

import (
	"context"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/bytedance/sonic"
	"github.com/gorilla/websocket"
	"github.com/takara2314/bsam-server/pkg/domain"
	"github.com/takara2314/bsam-server/pkg/racehub"
)

func TestWSPostGeolocation(t *testing.T) {
	var (
		url                = "ws://localhost:8081/japan"
		assocID            = "japan"
		password           = "nippon"
		timeoutSec         = 1 * time.Second
		sampleGeolocations = []racehub.PostGeolocationInput{
			{
				Latitude:      35.6895,
				Longitude:     139.6917,
				AccuracyMeter: 10.0,
				Heading:       1.0,
				RecordedAt:    time.Now(),
			},
			{
				Latitude:      36.6895,
				Longitude:     140.6917,
				AccuracyMeter: 11.0,
				Heading:       2.0,
				RecordedAt:    time.Now(),
			},
			{
				Latitude:      37.6895,
				Longitude:     141.6917,
				AccuracyMeter: 12.0,
				Heading:       3.0,
				RecordedAt:    time.Now(),
			},
		}
	)

	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSec)
	defer cancel()

	// トークンを取得
	token, err := FetchTokenFromAPI(assocID, password)
	if err != nil {
		t.Fatalf("トークンの取得に失敗しました: %v", err)
	}

	var wg sync.WaitGroup
	for i, geolocation := range sampleGeolocations {
		wg.Add(1)
		go postMarkGeolocationWS(
			ctx,
			t,
			timeoutSec,
			&wg,
			url,
			token,
			i+1,
			geolocation.Latitude,
			geolocation.Longitude,
			geolocation.AccuracyMeter,
			geolocation.Heading,
			geolocation.RecordedAt,
		)
	}

	done := make(chan bool)
	go func() {
		wg.Wait()

		// 位置情報を受け取る
		go receiveMarkGeolocationsWS(
			ctx,
			t,
			timeoutSec,
			url,
			token,
			len(sampleGeolocations),
			done,
		)
	}()

	select {
	case <-done:
		return
	case <-time.After(timeoutSec):
		t.Fatal("MarkGeolocations待機中にタイムアウトしました")
	}
}

func postMarkGeolocationWS(
	ctx context.Context,
	t *testing.T,
	timeoutSec time.Duration,
	wg *sync.WaitGroup,
	url string,
	token string,
	markNo int,
	latitude float64,
	longitude float64,
	accuracyMeter float64,
	heading float64,
	recordedAt time.Time,
) {
	deviceID := domain.CreateDeviceID(domain.RoleMark, markNo)

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
		MessageType: racehub.HandlerTypeAuth,
		Token:       token,
		DeviceID:    deviceID,
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

	<-authResultReceived

	// MarkGeolocationメッセージを送信
	input := racehub.PostGeolocationInput{
		MessageType:   racehub.HandlerTypePostGeolocation,
		Latitude:      latitude,
		Longitude:     longitude,
		AccuracyMeter: accuracyMeter,
		Heading:       heading,
		RecordedAt:    recordedAt,
	}

	payload, err = sonic.Marshal(&input)
	if err != nil {
		t.Fatalf("メッセージのエンコードに失敗しました: %v", err)
	}

	if err := conn.WriteMessage(
		websocket.TextMessage,
		payload,
	); err != nil {
		t.Fatalf("メッセージの送信に失敗しました: %v", err)
	}

	wg.Done()
}

func receiveMarkGeolocationsWS(
	ctx context.Context,
	t *testing.T,
	timeoutSec time.Duration,
	url string,
	token string,
	wantMarkCounts int,
	done chan bool,
) {
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
		DeviceID:       "athlete9",
		WantMarkCounts: 5,
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

	received := make(chan bool)

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

			var msg *racehub.MarkGeolocationsOutput
			err = sonic.Unmarshal(payload, &msg)
			if err != nil {
				t.Errorf("メッセージのデコードに失敗しました: %v", err)
			}

			if msg.MessageType != racehub.ActionTypeMarkGeolocations {
				continue
			}

			if len(msg.Marks) != 5 {
				t.Errorf(
					"位置情報の数が正しくありません: got %d, want %d",
					len(msg.Marks),
					wantMarkCounts,
				)
			}

			correctStored := []bool{true, true, true, false, false}
			for i, mark := range msg.Marks {
				if mark.Stored != correctStored[i] {
					t.Errorf(
						"マーク%dの位置情報の格納状態が正しくありません: got %v, want %v",
						i+1,
						mark.Stored,
						correctStored[i],
					)
				}
			}

			received <- true
		}
	}()
	<-received
	done <- true
}
