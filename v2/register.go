package v2

import (
	"bsam-server/v2/ws/racing"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Register(e *gin.Engine) *gin.RouterGroup {
	router := e.Group("/v2")

	// Test API
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello 2th version API!")
	})

	// Racing Socket
	router.GET("/racing/:id", racing.Handler)
	go racing.AutoRooming()

	return router
}
