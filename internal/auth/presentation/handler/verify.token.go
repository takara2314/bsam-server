package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/takara2314/bsam-server/internal/auth/app"
)

type VerifyTokenPOSTRequest struct {
	Token string `json:"token"`
}

func VerifyTokenPOST(c *gin.Context, req VerifyTokenPOSTRequest) {
	associationID, err := app.ParseToken(req.Token)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "token is invalid",
		})
		return
	}

	newToken, err := app.CreateToken(associationID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "OK",
		"token":   newToken,
	})
}
