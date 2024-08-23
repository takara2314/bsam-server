package racelib

import (
	"context"
	"log/slog"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/samber/oops"
	repoFirestore "github.com/takara2314/bsam-server/pkg/infrastructure/repository/firestore"
)

type Race struct {
	Started    bool
	StartedAt  time.Time
	FinishedAt time.Time
}

func StoreRace(
	ctx context.Context,
	firestoreClient *firestore.Client,
	associationID string,
	started bool,
	startedAt time.Time,
	finishedAt time.Time,
) error {
	raceID := associationID

	if err := repoFirestore.SetRace(
		ctx,
		firestoreClient,
		raceID,
		started,
		startedAt,
		finishedAt,
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
	}, nil
}
