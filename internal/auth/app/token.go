package app

import (
	"context"

	"github.com/samber/oops"
	"github.com/takara2314/bsam-server/internal/auth/common"
	"github.com/takara2314/bsam-server/pkg/auth"
	repoFirestore "github.com/takara2314/bsam-server/pkg/infrastructure/repository/firestore"
)

func ParseToken(token string) (string, error) {
	associationID, err := auth.ParseJWT(token, common.Env.JWTSecretKey)
	if err != nil {
		return "", oops.
			In("app.ParseToken").
			Wrapf(err, "failed to parse token")
	}

	return associationID, nil
}

func CreateToken(associationID string) (string, error) {
	ctx := context.Background()

	association, err := repoFirestore.FetchAssociationByID(
		ctx,
		common.FirestoreClient,
		associationID,
	)
	if err != nil {
		return "", oops.
			In("auth.CreateToken").
			Wrapf(err, "failed to fetch association")
	}

	return auth.CreateJWT(
		associationID,
		association.ExpiresAt,
		common.Env.JWTSecretKey,
	), nil
}
