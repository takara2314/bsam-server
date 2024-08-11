package action

import (
	"context"

	"github.com/takara2314/bsam-server/internal/game/common"
	"github.com/takara2314/bsam-server/pkg/domain"
	"github.com/takara2314/bsam-server/pkg/geolocationhub"
	"github.com/takara2314/bsam-server/pkg/racehub"
)

type RaceAction struct {
	racehub.UnimplementedAction
}

// 認証結果メッセージを送信するときの処理
func (r *RaceAction) AuthResult(
	c *racehub.Client,
	ok bool,
	message string,
) (*racehub.AuthResultOutput, error) {
	return &racehub.AuthResultOutput{
		MessageType: racehub.ActionTypeAuthResult,
		OK:          ok,
		DeviceID:    c.DeviceID,
		Role:        c.Role,
		MarkNo:      c.MarkNo,
		Message:     message,
	}, nil
}

// 選手にマークの位置情報を送信するときの処理
func (r *RaceAction) MarkGeolocations(
	c *racehub.Client,
) (*racehub.MarkGeolocationsOutput, error) {
	ctx := context.Background()
	geolocHub := geolocationhub.NewHub(
		c.Hub.AssociationID,
		common.FirestoreClient,
	)

	marks := make(
		[]racehub.MarkGeolocationsOutputMark, c.WantMarkCounts,
	)

	for i := range marks {
		marks[i] = fetchMarkGeolocation(
			ctx,
			geolocHub,
			i+1,
		)
	}

	return &racehub.MarkGeolocationsOutput{
		MessageType: racehub.ActionTypeMarkGeolocations,
		MarkCounts:  c.WantMarkCounts,
		Marks:       marks,
	}, nil
}

// マークの位置情報を取得する
// データが取得できなかった場合は、Storedフィールドをfalseにする
func fetchMarkGeolocation(
	ctx context.Context,
	geolocHub *geolocationhub.Hub,
	markNo int,
) racehub.MarkGeolocationsOutputMark {
	deviceID := domain.CreateDeviceID(
		domain.RoleMark,
		markNo,
	)

	loc, err := geolocHub.FetchLatestGeolocationByDeviceID(
		ctx,
		deviceID,
	)
	if err != nil {
		return racehub.MarkGeolocationsOutputMark{
			Stored: false,
		}
	}

	return racehub.MarkGeolocationsOutputMark{
		MarkNo:        markNo,
		Stored:        true,
		Latitude:      loc.Latitude,
		Longitude:     loc.Longitude,
		AccuracyMeter: loc.AccuracyMeter,
		Heading:       loc.Heading,
		RecordedAt:    loc.RecordedAt,
	}
}
