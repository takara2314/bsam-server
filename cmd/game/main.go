package main

import (
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/takara2314/bsam-server/internal/game/presentation"
	"github.com/takara2314/bsam-server/pkg/logging"
)

func main() {
	logging.InitSlog()

	if os.Getenv("ENVIRONMENT") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	presentation.RegisterRouter(router)

	slog.Info(
		"game server started",
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
