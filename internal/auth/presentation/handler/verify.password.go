package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func VerifyPasswordPOST(c *gin.Context, assocID string, password string) {
	c.JSON(http.StatusOK, gin.H{
		"message": "OK",
	})
}
