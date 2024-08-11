package app

import (
	"context"

	"github.com/takara2314/bsam-server/internal/auth/common"
	"github.com/takara2314/bsam-server/pkg/auth"
	repoFirestore "github.com/takara2314/bsam-server/pkg/infrastructure/repository/firestore"
)

func VerifyPassword(associationID string, password string) error {
	ctx := context.Background()

	_, err := repoFirestore.FetchAssociationByIDAndHashedPassword(
		ctx,
		common.FirestoreClient,
		associationID,
		auth.HashPassword(password),
	)
	if err != nil {
		return err
	}

	return nil
}
