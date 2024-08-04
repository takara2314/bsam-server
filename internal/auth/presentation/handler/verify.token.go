package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/takara2314/bsam-server/internal/auth/app"
)

func VerifyTokenPOST(c *gin.Context, token string) {
	assocID, err := app.ParseToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "token is invalid",
		})
		return
	}

	newToken, err := app.CreateToken(assocID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "OK",
		"token":   newToken,
	})
}
