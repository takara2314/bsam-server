package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/takara2314/bsam-server/internal/api/common"
	"github.com/takara2314/bsam-server/internal/api/presentation"
	"github.com/takara2314/bsam-server/pkg/environment"
	"github.com/takara2314/bsam-server/pkg/infrastructure/repository"
	"github.com/takara2314/bsam-server/pkg/logging"
)

func main() {
	var err error
	ctx := context.Background()

	logging.InitSlog()

	common.Env, err = environment.LoadVariables(false)
	if err != nil {
		slog.Error(
			"failed to load env",
			"error", err,
		)
		panic(err)
	}

	common.FirestoreClient, err = repository.NewFirestore(
		ctx,
		common.Env.GoogleCloudProjectID,
	)
	if err != nil {
		panic(err)
	}
	defer common.FirestoreClient.Close()

	router := presentation.NewGin()
	presentation.RegisterRouter(router)

	slog.Info(
		"api server started",
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
