package passedmarklib

import (
	"context"
	"log/slog"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/samber/oops"
	"github.com/takara2314/bsam-server/pkg/domain"
	"github.com/takara2314/bsam-server/pkg/nextmarklib"
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
	wantMarkCounts int,
	passedAt time.Time,
) error {
	nextMarkNo := domain.CalcNextMarkNo(wantMarkCounts, markNo)

	if err := nextmarklib.StoreNextMark(
		ctx,
		firestoreClient,
		associationID,
		deviceID,
		nextMarkNo,
		passedAt,
	); err != nil {
		return oops.
			In("passedmarklib.StorePassedMark").
			With("association_id", associationID).
			With("device_id", deviceID).
			With("mark_no", markNo).
			With("want_mark_counts", wantMarkCounts).
			With("next_mark_no", nextMarkNo).
			With("passed_at", passedAt).
			Wrapf(err, "failed to set passed_mark as next_mark to firestore")
	}

	slog.Info(
		"passed_mark stored",
		"association_id", associationID,
		"device_id", deviceID,
		"mark_no", markNo,
		"want_mark_counts", wantMarkCounts,
		"next_mark_no", nextMarkNo,
		"passed_at", passedAt,
	)

	return nil
}
