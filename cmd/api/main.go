package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/takara2314/bsam-server/internal/api/common"
	"github.com/takara2314/bsam-server/internal/api/presentation"
	"github.com/takara2314/bsam-server/pkg/infrastructure"
	"github.com/takara2314/bsam-server/pkg/logging"
)

func main() {
	var err error
	ctx := context.Background()

	logging.InitSlog()

	common.FirestoreClient, err = infrastructure.NewFirestore(
		ctx,
		"bsam-app",
	)
	if err != nil {
		panic(err)
	}

	// TODO: 後で消す
	_, _, err = common.FirestoreClient.Collection("users").Add(ctx, map[string]interface{}{
		"first": "Ada",
		"last":  "Lovelace",
		"born":  1815,
	})
	if err != nil {
		panic(err)
	}

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
