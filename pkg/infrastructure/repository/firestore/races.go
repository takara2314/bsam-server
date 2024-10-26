package firestore

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/samber/oops"
)

type Race struct {
	ID         string    `firestore:"-"`
	Started    bool      `firestore:"started"`
	StartedAt  time.Time `firestore:"startedAt"`
	FinishedAt time.Time `firestore:"finishedAt"`
	DeviceIDs  []string  `firestore:"deviceIDs"`
	UpdatedAt  time.Time `firestore:"updatedAt"`
}

func SetRace(
	ctx context.Context,
	client *firestore.Client,
	id string,
	started bool,
	startedAt time.Time,
	finishedAt time.Time,
	deviceIDs []string,
	updatedAt time.Time,
) error {
	_, err := client.Collection("races").Doc(id).Set(ctx, Race{
		ID:         id,
		Started:    started,
		StartedAt:  startedAt,
		FinishedAt: finishedAt,
		DeviceIDs:  deviceIDs,
		UpdatedAt:  updatedAt,
	})
	if err != nil {
		return oops.
			In("firestore.SetRace").
			Wrapf(err, "failed to set race")
	}

	return nil
}

func FetchRaceByID(
	ctx context.Context,
	client *firestore.Client,
	id string,
) (*Race, error) {
	doc, err := client.Collection("races").Doc(id).Get(ctx)
	if err != nil {
		return nil, oops.
			In("firestore.FetchRaceByID").
			Wrapf(err, "failed to fetch race")
	}

	var race Race
	err = doc.DataTo(&race)
	if err != nil {
		return nil, oops.
			In("firestore.FetchRaceByID").
			Wrapf(err, "failed to convert data to race")
	}

	race.ID = doc.Ref.ID

	return &race, err
}
