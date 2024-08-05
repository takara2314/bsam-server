package app

import (
	"context"

	"github.com/samber/oops"
	"github.com/takara2314/bsam-server/internal/auth/common"
	"github.com/takara2314/bsam-server/pkg/auth"
	"github.com/takara2314/bsam-server/pkg/infrastructure/repository/firestore"
)

func ParseToken(token string) (string, error) {
	assocID, err := auth.ParseJWT(token, common.Env.JWTSecretKey)
	if err != nil {
		return "", oops.
			In("app.ParseToken").
			Wrapf(err, "failed to parse token")
	}

	return assocID, nil
}

func CreateToken(assocID string) (string, error) {
	ctx := context.Background()

	assoc, err := firestore.FetchAssocByID(ctx, common.FirestoreClient, assocID)
	if err != nil {
		return "", oops.
			In("auth.CreateToken").
			Wrapf(err, "failed to fetch assoc")
	}

	return auth.CreateJWT(
		assocID,
		assoc.ExpiredAt,
		common.Env.JWTSecretKey,
	), nil
}
