package nextmarklib

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/samber/oops"
	repoFirestore "github.com/takara2314/bsam-server/pkg/infrastructure/repository/firestore"
)

type NextMark struct {
	NextMarkNo int
	SetAt      time.Time
}

func StoreNextMark(
	ctx context.Context,
	firestoreClient *firestore.Client,
	associationID string,
	deviceID string,
	nextMarkNo int,
	setAt time.Time,
) error {
	firestoreDeviceID := associationID + "_" + deviceID

	if err := repoFirestore.SetNextMark(
		ctx,
		firestoreClient,
		firestoreDeviceID,
		nextMarkNo,
		setAt,
	); err != nil {
		return oops.
			In("nextmarklib.StoreNextMark").
			With("association_id", associationID).
			With("device_id", deviceID).
			With("next_mark_no", nextMarkNo).
			With("set_at", setAt).
			Wrapf(err, "failed to set next_mark to firestore")
	}

	slog.Info(
		"next_mark stored",
		"association_id", associationID,
		"device_id", deviceID,
		"next_mark_no", nextMarkNo,
		"set_at", setAt,
	)

	return nil
}

func FetchNextMarkOnlyAfterThisDT(
	ctx context.Context,
	firestoreClient *firestore.Client,
	associationID string,
	deviceID string,
	dt time.Time,
) (*NextMark, error) {
	firestoreDeviceID := associationID + "_" + deviceID

	nextMark, err := repoFirestore.FetchNextMarkByID(
		ctx,
		firestoreClient,
		firestoreDeviceID,
	)
	if err != nil {
		return nil, oops.
			In("passedmarklib.FetchPassedMarkOnlyThisDT").
			With("association_id", associationID).
			With("device_id", deviceID).
			With("dt", dt).
			Wrapf(err, "failed to fetch passed_mark")
	}

	// dtより前ならnilを返す
	if nextMark.UpdatedAt.Before(dt) {
		return nil, errors.New("next_mark is before this_dt")
	}

	return &NextMark{
		NextMarkNo: nextMark.NextMarkNo,
		SetAt:      nextMark.UpdatedAt,
	}, nil
}
