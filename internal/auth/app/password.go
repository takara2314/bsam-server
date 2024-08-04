package app

import (
	"context"

	"github.com/takara2314/bsam-server/internal/auth/common"
	"github.com/takara2314/bsam-server/pkg/auth"
	"github.com/takara2314/bsam-server/pkg/infrastructure/repository"
)

func VerifyPassword(assocID string, password string) error {
	ctx := context.Background()

	_, err := repository.FetchAssocByIDAndHashedPassword(
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
