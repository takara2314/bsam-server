package passedmarklib

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/samber/oops"
	repoFirestore "github.com/takara2314/bsam-server/pkg/infrastructure/repository/firestore"
)

type PassedMark struct {
	MarkNo   int
	PassedAt time.Time
}

func StorePassedMark(
	ctx context.Context,
	firestoreClient *firestore.Client,
	associationID string,
	deviceID string,
	markNo int,
	passedAt time.Time,
) error {
	firestoreDeviceID := associationID + "_" + deviceID

	if err := repoFirestore.SetPassedMark(
		ctx,
		firestoreClient,
		firestoreDeviceID,
		markNo,
		passedAt,
	); err != nil {
		return oops.
			In("passedmarklib.StorePassedMark").
			With("association_id", associationID).
			With("device_id", deviceID).
			With("mark_no", markNo).
			With("passed_at", passedAt).
			Wrapf(err, "failed to set passed_mark to firestore")
	}

	slog.Info(
		"passed_mark stored",
		"association_id", associationID,
		"device_id", deviceID,
		"mark_no", markNo,
		"passed_at", passedAt,
	)

	return nil
}

func FetchPassedMarkOnlyAfterThisDT(
	ctx context.Context,
	firestoreClient *firestore.Client,
	associationID string,
	deviceID string,
	dt time.Time,
) (*PassedMark, error) {
	firestoreDeviceID := associationID + "_" + deviceID

	passedMarks, err := repoFirestore.FetchPassedMarkByID(
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
	if passedMarks.PassedAt.Before(dt) {
		return nil, errors.New("passed_mark is before this_dt")
	}

	return &PassedMark{
		MarkNo:   passedMarks.MarkNo,
		PassedAt: passedMarks.PassedAt,
	}, nil
}
