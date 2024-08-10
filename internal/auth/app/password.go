package app

import (
	"context"

	"github.com/takara2314/bsam-server/internal/auth/common"
	"github.com/takara2314/bsam-server/pkg/auth"
	repoFirestore "github.com/takara2314/bsam-server/pkg/infrastructure/repository/firestore"
)

func VerifyPassword(assocID string, password string) error {
	ctx := context.Background()

	_, err := repoFirestore.FetchAssocByIDAndHashedPassword(
		ctx,
		common.FirestoreClient,
		assocID,
		auth.HashPassword(password),
	)
	if err != nil {
		return err
	}

	return nil
}
