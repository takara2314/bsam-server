package firestore

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/samber/oops"
)

type Assoc struct {
	ID             string    `firestore:"-"`
	Name           string    `firestore:"name"`
	HashedPassword string    `firestore:"hashedPassword"`
	ContractType   string    `firestore:"contractType"`
	ExpiresAt      time.Time `firestore:"expiresAt"`
	UpdatedAt      time.Time `firestore:"updatedAt"`
}

func SetAssoc(
	ctx context.Context,
	client *firestore.Client,
	id string,
	name string,
	hashedPassword string,
	contractType string,
	expiresAt time.Time,
	updatedAt time.Time,
) error {
	_, err := client.Collection("assocs").Doc(id).Set(ctx, Assoc{
		ID:             id,
		Name:           name,
		HashedPassword: hashedPassword,
		ContractType:   contractType,
		ExpiresAt:      expiresAt,
		UpdatedAt:      updatedAt,
	})
	if err != nil {
		return oops.
			In("repository.AddAssoc").
			Wrapf(err, "failed to add assoc")
	}

	return nil
}

func FetchAssocByID(ctx context.Context, client *firestore.Client, id string) (*Assoc, error) {
	doc, err := client.Collection("assocs").Doc(id).Get(ctx)
	if err != nil {
		return nil, oops.
			In("repository.FetchAssocByID").
			Wrapf(err, "failed to fetch assoc")
	}

	var assoc Assoc
	err = doc.DataTo(&assoc)
	if err != nil {
		return nil, oops.
			In("repository.FetchAssocByID").
			Wrapf(err, "failed to convert data to user")
	}

	assoc.ID = doc.Ref.ID

	return &assoc, err
}

func FetchAssocByIDAndHashedPassword(ctx context.Context, client *firestore.Client, id string, hashedPassword string) (*Assoc, error) {
	assoc, err := FetchAssocByID(ctx, client, id)
	if err != nil {
		return nil, oops.
			In("repository.FetchAssocByIDAndHashedPassword").
			Wrapf(err, "not found this assoc id")
	}

	if assoc.HashedPassword != hashedPassword {
		return nil, oops.
			In("repository.FetchAssocByIDAndHashedPassword").
			Wrapf(nil, "hashed password is not matched")
	}

	return assoc, err
}
