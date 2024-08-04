package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/takara2314/bsam-server/internal/api/common"
	"github.com/takara2314/bsam-server/internal/api/presentation"
	"github.com/takara2314/bsam-server/pkg/infrastructure/repository"
	"github.com/takara2314/bsam-server/pkg/logging"
)

func main() {
	var err error
	ctx := context.Background()

	logging.InitSlog()

	common.FirestoreClient, err = repository.NewFirestore(
		ctx,
		"bsam-app",
	)
	if err != nil {
		panic(err)
	}
	defer common.FirestoreClient.Close()

	router := presentation.NewGin()
	presentation.RegisterRouter(router)

	slog.Info(
		"auth server started",
		"is_production", os.Getenv("ENVIRONMENT") == "production",
	)

	if err := router.Run(":8080"); err != nil {
		slog.Error(
			"failed to run Gin router",
			"error", err,
		)
		panic(err)
	}
}
