package racelib

import (
	"context"
	"log/slog"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/samber/oops"
	"github.com/takara2314/bsam-server/pkg/domain"
	repoFirestore "github.com/takara2314/bsam-server/pkg/infrastructure/repository/firestore"
	"github.com/takara2314/bsam-server/pkg/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Race struct {
	Started    bool
	StartedAt  time.Time
	FinishedAt time.Time
	DeviceIDs  []string
}

func StoreRace(
	ctx context.Context,
	firestoreClient *firestore.Client,
	associationID string,
	started bool,
	startedAt time.Time,
	finishedAt time.Time,
	deviceIDs []string,
) error {
	raceID := associationID

	if err := repoFirestore.SetRace(
		ctx,
		firestoreClient,
		raceID,
		started,
		startedAt,
		finishedAt,
		deviceIDs,
		time.Now(),
	); err != nil {
		return oops.
			In("racelib.StoreRace").
			With("association_id", associationID).
			With("started", started).
			With("started_at", startedAt).
			With("finished_at", finishedAt).
			Wrapf(err, "failed to set race to firestore")
	}

	slog.Info(
		"race stored",
		"association_id", associationID,
		"started", started,
		"started_at", startedAt,
		"finished_at", finishedAt,
	)

	return nil
}

func AddDeviceIDToRace(
	ctx context.Context,
	firestoreClient *firestore.Client,
	associationID string,
	deviceID string,
) error {
	raceID := associationID

	// レースを取得できない場合は、たいていレースが存在しないので、空のインスタンスを入れる
	race, err := repoFirestore.FetchRaceByID(
		ctx,
		firestoreClient,
		raceID,
	)
	if err != nil {
		race = new(repoFirestore.Race)
	}

	if err := repoFirestore.SetRace(
		ctx,
		firestoreClient,
		raceID,
		race.Started,
		race.StartedAt,
		race.FinishedAt,
		util.AddStrToSliceIfNotExists(race.DeviceIDs, deviceID),
		time.Now(),
	); err != nil {
		return oops.
			In("racelib.AddDeviceIDToRace").
			With("race_id", raceID).
			With("device_id", deviceID).
			Wrapf(err, "failed to set race")
	}

	return nil
}

func RemoveDeviceIDFromRace(
	ctx context.Context,
	firestoreClient *firestore.Client,
	associationID string,
	deviceID string,
) error {
	raceID := associationID

	race, err := repoFirestore.FetchRaceByID(
		ctx,
		firestoreClient,
		raceID,
	)
	if err != nil {
		return oops.
			In("racelib.RemoveDeviceIDFromRace").
			With("race_id", raceID).
			Wrapf(err, "failed to fetch race")
	}

	if err := repoFirestore.SetRace(
		ctx,
		firestoreClient,
		raceID,
		race.Started,
		race.StartedAt,
		race.FinishedAt,
		util.RemoveStrFromSliceIfExists(race.DeviceIDs, deviceID),
		time.Now(),
	); err != nil {
		return oops.
			In("racelib.RemoveDeviceIDFromRace").
			With("race_id", raceID).
			Wrapf(err, "failed to set race")
	}

	return nil
}

func FetchLatestRaceByAssociationID(
	ctx context.Context,
	firestoreClient *firestore.Client,
	associationID string,
) (*Race, error) {
	raceID := associationID

	race, err := repoFirestore.FetchRaceByID(
		ctx,
		firestoreClient,
		raceID,
	)
	if err != nil {
		return nil, oops.
			In("racelib.FetchLatestRaceByAssociationID").
			With("association_id", associationID).
			With("race_id", raceID).
			Wrapf(err, "failed to fetch race by association id")
	}

	return &Race{
		Started:    race.Started,
		StartedAt:  race.StartedAt,
		FinishedAt: race.FinishedAt,
		DeviceIDs:  race.DeviceIDs,
	}, nil
}

func FetchLatestRaceDetailByAssociationID(
	ctx context.Context,
	firestoreClient *firestore.Client,
	associationID string,
) (*domain.Race, codes.Code) {
	// 協会情報を取得
	association, err := repoFirestore.FetchAssociationByID(
		ctx,
		firestoreClient,
		associationID,
	)
	if status.Code(err) == codes.NotFound {
		return nil, codes.NotFound
	} else if err != nil {
		return nil, codes.Internal
	}

	// レース情報を取得
	// レースが存在しないなどエラーが起きた場合は、空のインスタンスを入れる
	race, err := repoFirestore.FetchRaceByID(
		ctx,
		firestoreClient,
		associationID,
	)
	if err != nil {
		race = new(repoFirestore.Race)
	}

	athleteIDs, markIDs, _, _ := domain.SeparateDeviceIDByRole(race.DeviceIDs)

	return &domain.Race{
		AssociationID: associationID,
		Name:          association.RaceName,
		Started:       race.Started,
		StartedAt:     race.StartedAt,
		FinishedAt:    race.FinishedAt,
		Association: domain.Association{
			ID:           associationID,
			Name:         association.Name,
			ContractType: association.ContractType,
			ExpiresAt:    association.ExpiresAt,
		},
		AthleteIDs: athleteIDs,
		MarkIDs:    markIDs,
	}, codes.OK
}
