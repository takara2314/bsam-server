package e2e

import (
	"context"
	"net/url"
	"sync"
	"testing"
	"time"

	"github.com/takara2314/bsam-server/e2e/auth"
	"github.com/takara2314/bsam-server/pkg/domain"
	"github.com/takara2314/bsam-server/pkg/racehub"
)

func TestWSPostGeolocationMultiInstance(t *testing.T) {
	t.Parallel()

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
		geolocations  = []racehub.PostGeolocationInput{
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
	// mark1 -> server1
	// mark2 -> server2
	// mark3 -> server1
	var wg sync.WaitGroup
	for i, geolocation := range geolocations {
		serverURL := serverURL1
		if i%2 == 1 {
			serverURL = serverURL2
		}

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
	// athlete1 -> server2
	marks, err := connectAndReceiveMarkGeolocations(
		ctx,
		timeout,
		serverURL2,
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
