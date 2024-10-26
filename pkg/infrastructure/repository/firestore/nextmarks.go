package firestore

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/samber/oops"
)

type NextMark struct {
	ID         string    `firestore:"-"`
	NextMarkNo int       `firestore:"nextMarkNo"`
	UpdatedAt  time.Time `firestore:"updatedAt"`
}

func SetNextMark(
	ctx context.Context,
	client *firestore.Client,
	id string,
	nextMarkNo int,
	updatedAt time.Time,
) error {
	_, err := client.Collection("next_marks").Doc(id).Set(ctx, NextMark{
		ID:         id,
		NextMarkNo: nextMarkNo,
		UpdatedAt:  updatedAt,
	})
	if err != nil {
		return oops.
			In("firestore.SetNextMark").
			Wrapf(err, "failed to set next_mark")
	}

	return nil
}

func FetchNextMarkByID(
	ctx context.Context,
	client *firestore.Client,
	id string,
) (*NextMark, error) {
	doc, err := client.Collection("next_marks").Doc(id).Get(ctx)
	if err != nil {
		return nil, oops.
			In("firestore.FetchNextMarkByID").
			Wrapf(err, "failed to fetch next_mark")
	}

	var nextMark NextMark
	err = doc.DataTo(&nextMark)
	if err != nil {
		return nil, oops.
			In("firestore.FetchNextMarkByID").
			Wrapf(err, "failed to convert data to next_mark")
	}

	return &nextMark, nil
}
