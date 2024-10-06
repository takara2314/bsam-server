package firestore

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/samber/oops"
)

type PassedMark struct {
	ID       string    `firestore:"-"`
	MarkNo   int       `firestore:"markNo"`
	PassedAt time.Time `firestore:"passedAt"`
}

func SetPassedMark(
	ctx context.Context,
	client *firestore.Client,
	id string,
	markNo int,
	passedAt time.Time,
) error {
	_, err := client.Collection("passe_marks").Doc(id).Set(ctx, PassedMark{
		ID:       id,
		MarkNo:   markNo,
		PassedAt: passedAt,
	})
	if err != nil {
		return oops.
			In("firestore.SetPassedMark").
			Wrapf(err, "failed to set passed_mark")
	}

	return nil
}

func FetchPassedMarkByID(
	ctx context.Context,
	client *firestore.Client,
	id string,
) (*PassedMark, error) {
	doc, err := client.Collection("passe_marks").Doc(id).Get(ctx)
	if err != nil {
		return nil, oops.
			In("firestore.FetchPassedMarkByID").
			Wrapf(err, "failed to fetch passed_mark")
	}

	var passedMark PassedMark
	err = doc.DataTo(&passedMark)
	if err != nil {
		return nil, oops.
			In("firestore.FetchPassedMarkByID").
			Wrapf(err, "failed to convert data to passed_mark")
	}

	return &passedMark, nil
}
