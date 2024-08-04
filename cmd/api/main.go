package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/takara2314/bsam-server/internal/api/common"
	"github.com/takara2314/bsam-server/internal/api/presentation"
	"github.com/takara2314/bsam-server/pkg/infrastructure"
	"github.com/takara2314/bsam-server/pkg/logging"
)

func main() {
	var err error
	ctx := context.Background()

	logging.InitSlog()

	common.FirestoreClient, err = infrastructure.InitFirestore(
		ctx,
		"bsam-app",
	)
	if err != nil {
		panic(err)
	}

	if os.Getenv("ENVIRONMENT") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
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
