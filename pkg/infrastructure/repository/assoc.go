package repository

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/samber/oops"
	"github.com/takara2314/bsam-server/pkg/auth"
	"github.com/takara2314/bsam-server/pkg/domain"
)

type User struct {
	ID             string    `firestore:"-"`
	Name           string    `firestore:"name"`
	HashedPassword string    `firestore:"hashedPassword"`
	ContractType   string    `firestore:"contractType"`
	CreatedAt      time.Time `firestore:"createdAt"`
	ExpiredAt      time.Time `firestore:"expiredAt"`
}

func CreateAssoc(
	ctx context.Context,
	client *firestore.Client,
	id string,
	name string,
	password string,
	contractType domain.ContractType,
) error {
	_, err := client.Collection("assocs").Doc(id).Set(ctx, User{
		ID:             id,
		Name:           name,
		HashedPassword: auth.HashPassword(password),
		ContractType:   string(contractType),
		CreatedAt:      time.Now(),
		ExpiredAt:      time.Now().Add(contractType.Duration()),
	})
	if err != nil {
		return oops.
			In("repository.AddAssoc").
			Wrapf(err, "failed to add assoc")
	}

	return nil
}

func FetchAssocByID(ctx context.Context, client *firestore.Client, id string) (*User, error) {
	doc, err := client.Collection("assocs").Doc(id).Get(ctx)
	if err != nil {
		return nil, oops.
			In("repository.FetchAssocByID").
			Wrapf(err, "failed to fetch assoc")
	}

	var user User
	err = doc.DataTo(&user)
	if err != nil {
		return nil, oops.
			In("repository.FetchAssocByID").
			Wrapf(err, "failed to convert data to user")
	}

	user.ID = doc.Ref.ID

	return &user, err
}
