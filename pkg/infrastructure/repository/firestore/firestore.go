package firestore

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/samber/oops"
)

func NewClient(ctx context.Context, projectID string) (*firestore.Client, error) {
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return nil, oops.
			In("firestore.NewClient").
			Wrapf(err, "failed to initialize firestore client")
	}

	return client, err
}
