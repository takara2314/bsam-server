package infrastructure

import (
	"context"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/samber/oops"
)

func InitFirestore(ctx context.Context, projectID string) (*firestore.Client, error) {
	conf := &firebase.Config{ProjectID: projectID}
	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		return nil, oops.
			In("infrastructure.initFirestore").
			Wrapf(err, "failed to initialize firebase app")
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		return nil, oops.
			In("infrastructure.initFirestore").
			Wrapf(err, "failed to initialize firestore client")
	}

	return client, err
}
