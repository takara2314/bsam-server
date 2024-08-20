package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/takara2314/bsam-server/internal/auth/app"
)

type VerifyPasswordPOSTRequest struct {
	AssociationID string `json:"association_id"`
	Password      string `json:"password"`
}

func VerifyPasswordPOST(c *gin.Context, req VerifyPasswordPOSTRequest) {
	if err := app.VerifyPassword(req.AssociationID, req.Password); err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "association_id or password is incorrect",
		})
		return
	}

	token, err := app.CreateToken(req.AssociationID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "OK",
		"token":   token,
	})
}
