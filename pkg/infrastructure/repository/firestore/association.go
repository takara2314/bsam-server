package firestore

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/samber/oops"
)

type Association struct {
	ID             string    `firestore:"-"`
	Name           string    `firestore:"name"`
	HashedPassword string    `firestore:"hashedPassword"`
	ContractType   string    `firestore:"contractType"`
	ExpiresAt      time.Time `firestore:"expiresAt"`
	UpdatedAt      time.Time `firestore:"updatedAt"`
}

func SetAssociation(
	ctx context.Context,
	client *firestore.Client,
	id string,
	name string,
	hashedPassword string,
	contractType string,
	expiresAt time.Time,
	updatedAt time.Time,
) error {
	_, err := client.Collection("associations").Doc(id).Set(ctx, Association{
		ID:             id,
		Name:           name,
		HashedPassword: hashedPassword,
		ContractType:   contractType,
		ExpiresAt:      expiresAt,
		UpdatedAt:      updatedAt,
	})
	if err != nil {
		return oops.
			In("firestore.SetAssociation").
			Wrapf(err, "failed to set association")
	}

	return nil
}

func FetchAssociationByID(
	ctx context.Context,
	client *firestore.Client,
	id string,
) (*Association, error) {
	doc, err := client.Collection("associations").Doc(id).Get(ctx)
	if err != nil {
		return nil, oops.
			In("firestore.FetchAssociationByID").
			Wrapf(err, "failed to fetch association")
	}

	var association Association
	err = doc.DataTo(&association)
	if err != nil {
		return nil, oops.
			In("firestore.FetchAssociationByID").
			Wrapf(err, "failed to convert data to user")
	}

	association.ID = doc.Ref.ID

	return &association, err
}

func FetchAssociationByIDAndHashedPassword(
	ctx context.Context,
	client *firestore.Client,
	id string,
	hashedPassword string,
) (*Association, error) {
	association, err := FetchAssociationByID(ctx, client, id)
	if err != nil {
		return nil, oops.
			In("firestore.FetchAssociationByIDAndHashedPassword").
			Wrapf(err, "not found this association id")
	}

	if association.HashedPassword != hashedPassword {
		return nil, oops.
			In("firestore.FetchAssociationByIDAndHashedPassword").
			Wrapf(nil, "hashed password is not matched")
	}

	return association, err
}
