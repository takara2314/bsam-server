package action

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/takara2314/bsam-server/internal/game/common"
	"github.com/takara2314/bsam-server/pkg/domain"
	"github.com/takara2314/bsam-server/pkg/geolocationlib"
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

	marks := make(
		[]racehub.MarkGeolocationsOutputMark, c.WantMarkCounts,
	)

	for i := range marks {
		marks[i] = fetchMarkGeolocation(
			ctx,
			common.FirestoreClient,
			c.Hub.AssociationID,
			i+1,
		)
	}

	return &racehub.MarkGeolocationsOutput{
		MessageType: racehub.ActionTypeMarkGeolocations,
		MarkCounts:  c.WantMarkCounts,
		Marks:       marks,
	}, nil
}

// レースの状態を管理するときの処理
func (r *RaceAction) ManageRaceStatus(
	c *racehub.Client,
	started bool,
) (*racehub.ManageRaceStatusOutput, error) {
	return &racehub.ManageRaceStatusOutput{
		MessageType: racehub.ActionTypeManageRaceStatus,
		Started:     started,
	}, nil
}

// 次のマークを管理するときの処理
func (r *RaceAction) ManageNextMark(
	c *racehub.Client,
	nextMarkNo int,
) (*racehub.ManageNextMarkOutput, error) {
	return &racehub.ManageNextMarkOutput{
		MessageType: racehub.ActionTypeManageNextMark,
		NextMarkNo:  nextMarkNo,
	}, nil
}

// マークの位置情報を取得する
// データが取得できなかった場合は、Storedフィールドをfalseにする
func fetchMarkGeolocation(
	ctx context.Context,
	firestore *firestore.Client,
	associationID string,
	markNo int,
) racehub.MarkGeolocationsOutputMark {
	deviceID := domain.CreateDeviceID(
		domain.RoleMark,
		markNo,
	)

	loc, err := geolocationlib.FetchLatestGeolocationByDeviceID(
		ctx,
		firestore,
		associationID,
		deviceID,
	)
	if err != nil {
		return racehub.MarkGeolocationsOutputMark{
			MarkNo: markNo,
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
