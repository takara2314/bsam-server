package presentation

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/takara2314/bsam-server/internal/game/presentation/handler"
)

func RegisterRouter(router *gin.Engine) {
	router.Use(cors.Default())

	router.GET("/healthz", handler.HealthzGET)
}
