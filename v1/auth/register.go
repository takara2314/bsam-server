package auth

import (
	"github.com/gin-gonic/gin"
)

// Register registers handler to assigned router.
func Register(router *gin.RouterGroup) {
	router.POST("/password", passwordPOST)
	router.POST("/token", tokenPOST)
}
