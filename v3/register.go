package v3

import (
	"bsam-server/v3/api/status"
	"bsam-server/v3/middleware"
	"bsam-server/v3/ws/racing"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Register(e *gin.Engine) *gin.RouterGroup {
	router := e.Group("/v3")

	// Test API
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello 3rd version API!")
	})

	// authorized and JSON only
	router.Use(
		middleware.AuthJWT(),
		middleware.CheckMIME("application/json"),
	)
	{
		// Server Status API
		router.GET("/status", status.StatusGET)
	}

	// Server Status API (non-authorized)
	router.GET("/status-non-auth", status.StatusGET)

	// Racing Socket
	router.GET("/racing/:id", racing.Handler)
	go racing.AutoRooming()

	return router
}
