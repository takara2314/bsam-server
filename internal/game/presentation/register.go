package presentation

import (
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/takara2314/bsam-server/internal/game/presentation/handler"
	"github.com/takara2314/bsam-server/internal/game/presentation/handlerValidating"
)

func NewGin() *gin.Engine {
	if os.Getenv("ENVIRONMENT") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	return gin.New()
}

func RegisterRouter(router *gin.Engine) {
	router.Use(cors.Default())

	router.GET("/", handler.IndexGET)
	router.GET("/healthz", handler.HealthzGET)
	router.GET("/:associationID", handlerValidating.AssociationIDWS)
}
