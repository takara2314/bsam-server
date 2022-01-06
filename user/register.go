package user

import (
	"github.com/gin-gonic/gin"
)

// Register registers handler to assigned router.
func Register(router *gin.RouterGroup) {
	router.GET("/:username", InfoGET)
}
