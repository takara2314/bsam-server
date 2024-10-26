package action

import (
	"context"
	"log/slog"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/samber/oops"
	"github.com/takara2314/bsam-server/internal/game/common"
	"github.com/takara2314/bsam-server/pkg/domain"
	"github.com/takara2314/bsam-server/pkg/geolocationlib"
	"github.com/takara2314/bsam-server/pkg/nextmarklib"
	"github.com/takara2314/bsam-server/pkg/racehub"
	"github.com/takara2314/bsam-server/pkg/racelib"
	"google.golang.org/grpc/codes"
)

type RaceAction struct {
	racehub.UnimplementedAction
}

// 接続結果メッセージを送信するときの処理
func (r *RaceAction) ConnectResult(
	c *racehub.Client,
	ok bool,
	hubID string,
) (*racehub.ConnectResultOutput, error) {
	return &racehub.ConnectResultOutput{
		MessageType: racehub.ActionTypeConnectResult,
		OK:          ok,
		HubID:       hubID,
	}, nil
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
		Authed:      c.Authed,
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

// マネージャーにマークの位置情報と選手のレース情報を送信するときの処理
func (r *RaceAction) ParticipantsInfo(
	c *racehub.Client,
) (*racehub.ParticipantsInfoOutput, error) {
	ctx := context.Background()

	marks := make(
		[]racehub.MarkGeolocationsOutputMark, c.WantMarkCounts,
	)
	athletes := []racehub.ParticipantsInfoOutputAthlete{}

	for i := range marks {
		marks[i] = fetchMarkGeolocation(
			ctx,
			common.FirestoreClient,
			c.Hub.AssociationID,
			i+1,
		)
	}

	raceDetail, code := racelib.FetchLatestRaceDetailByAssociationID(
		ctx,
		common.FirestoreClient,
		c.Hub.AssociationID,
	)
	if code != codes.OK {
		return nil, oops.
			In("action.ParticipantsInfo").
			Errorf("failed to fetch latest race detail by association id")
	}

	for _, athleteID := range raceDetail.AthleteIDs {
		athlete, err := fetchAthleteInfo(
			ctx,
			common.FirestoreClient,
			c.Hub.AssociationID,
			athleteID,
			raceDetail.StartedAt,
		)
		if err != nil {
			slog.Error(
				"failed to fetch athlete info",
				"error", err,
				"athlete_id", athleteID,
			)
			continue
		}
		athletes = append(athletes, athlete)
	}

	return &racehub.ParticipantsInfoOutput{
		MessageType: racehub.ActionTypeParticipantsInfo,
		MarkCounts:  c.WantMarkCounts,
		Marks:       marks,
		Athletes:    athletes,
	}, nil
}

// レースの状態を管理するときの処理
func (r *RaceAction) ManageRaceStatus(
	c *racehub.Client,
	started bool,
	startedAt time.Time,
	finishedAt time.Time,
) (*racehub.ManageRaceStatusOutput, error) {
	return &racehub.ManageRaceStatusOutput{
		MessageType: racehub.ActionTypeManageRaceStatus,
		Started:     started,
		StartedAt:   startedAt,
		FinishedAt:  finishedAt,
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

// 選手の情報を取得する
func fetchAthleteInfo(
	ctx context.Context,
	firestore *firestore.Client,
	associationID string,
	deviceID string,
	raceStartedAt time.Time,
) (racehub.ParticipantsInfoOutputAthlete, error) {
	loc, err := geolocationlib.FetchLatestGeolocationByDeviceID(
		ctx,
		firestore,
		associationID,
		deviceID,
	)
	if err != nil {
		return racehub.ParticipantsInfoOutputAthlete{}, oops.
			In("action.fetchAthleteInfo").
			Wrapf(err, "failed to fetch latest geolocation by device id")
	}

	nextMark, err := nextmarklib.FetchNextMarkOnlyAfterThisDT(
		ctx,
		firestore,
		associationID,
		deviceID,
		raceStartedAt,
	)
	if err != nil {
		return racehub.ParticipantsInfoOutputAthlete{}, oops.
			In("action.fetchAthleteInfo").
			Wrapf(err, "failed to fetch next mark only after this dt")
	}

	return racehub.ParticipantsInfoOutputAthlete{
		DeviceID:      deviceID,
		NextMarkNo:    nextMark.NextMarkNo,
		Latitude:      loc.Latitude,
		Longitude:     loc.Longitude,
		AccuracyMeter: loc.AccuracyMeter,
		Heading:       loc.Heading,
		RecordedAt:    loc.RecordedAt,
	}, nil
}
