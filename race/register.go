package race

import (
	"github.com/gin-gonic/gin"
)

// Register registers handler to assigned router.
func Register(router *gin.RouterGroup) {
	router.GET("/:id", infoGET)
	router.PUT("/:id", infoPUT)
}
