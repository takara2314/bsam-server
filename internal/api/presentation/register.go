package presentation

import (
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/takara2314/bsam-server/internal/api/common"
	"github.com/takara2314/bsam-server/internal/api/presentation/handler"
	"github.com/takara2314/bsam-server/internal/api/presentation/handlerValidating"
	"github.com/takara2314/bsam-server/pkg/middleware"
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

	router.Use(middleware.AuthToken(common.Env.JWTSecretKey))
	router.POST("/geolocation", handlerValidating.GeolocationPOST)
	router.GET("/races/:associationID", handler.RacesAssociationIDGET)
}
