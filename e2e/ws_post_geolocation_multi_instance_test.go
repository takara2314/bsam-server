package e2e

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/takara2314/bsam-server/pkg/racehub"
)

func TestWSPostGeolocationMultiInstance(t *testing.T) {
	var (
		url                = "ws://localhost:8081/japan"
		url2               = "ws://localhost:8181/japan"
		associationID      = "japan"
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
	token, err := FetchTokenFromAPI(associationID, password)
	if err != nil {
		t.Fatalf("トークンの取得に失敗しました: %v", err)
	}

	var wg sync.WaitGroup
	for i, geolocation := range sampleGeolocations {
		wg.Add(1)

		chooseURL := url
		if i%2 == 1 {
			chooseURL = url2
		}

		go postMarkGeolocationWS(
			ctx,
			t,
			timeoutSec,
			&wg,
			chooseURL,
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
			url2,
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
