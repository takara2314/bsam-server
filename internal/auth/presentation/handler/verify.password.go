package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/takara2314/bsam-server/internal/auth/app"
)

func VerifyPasswordPOST(c *gin.Context, assocID string, password string) {
	if err := app.VerifyPassword(assocID, password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "assoc_id or password is incorrect",
		})
		return
	}

	token, err := app.CreateToken(assocID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "OK",
		"token":   token,
	})
}
