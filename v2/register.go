package v2

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Register(e *gin.Engine) *gin.RouterGroup {
	router := e.Group("/v2")

	// Test API
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello 2th version API!")
	})

	return router
}
