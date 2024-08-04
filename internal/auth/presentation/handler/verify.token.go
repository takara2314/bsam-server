package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func VerifyTokenPOST(c *gin.Context, token string) {
	c.JSON(http.StatusOK, gin.H{
		"message": "OK",
	})
}
