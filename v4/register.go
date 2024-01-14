package v4

import (
	"net/http"

	"bsam-server/v4/api/associations"
	"bsam-server/v4/api/reboot"
	"bsam-server/v4/api/status"
	"bsam-server/v4/middleware"
	"bsam-server/v4/ws/racing"

	"github.com/gin-gonic/gin"
)

func Register(e *gin.Engine) *gin.RouterGroup {
	router := e.Group("/v4")

	// Test API
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello 4th version API!")
	})

	// Associations API
	router.GET("/associations/:id", associations.AssociationGET)

	// Authorized only
	authorized := router.Group("/",
		middleware.AuthJWT(),
	)
	{
		// Associations API
		authorized.GET("/associations", associations.AssociationGETAll)

		// Server Status API
		authorized.GET("/status", status.StatusGET)

		// Server Reboot API
		authorized.POST("/reboot", reboot.RebootPOST)
	}

	// // Authorized and JSON only
	// authorizedJSON := router.Group("/",
	// 	middleware.AuthJWT(),
	// 	middleware.CheckMIME("application/json"),
	// )
	// { }

	// Racing Socket
	router.GET("/racing/:id", racing.Handler)
	go racing.AutoRooming()

	return router
}
