package device

import (
	"github.com/gin-gonic/gin"
)

// Register registers handler to assigned router.
func Register(router *gin.RouterGroup) {
	router.POST("/:id", infoPOST)
	router.PUT("/:id", infoPUT)
}
