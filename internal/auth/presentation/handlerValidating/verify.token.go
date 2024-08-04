package handlerValidating

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/takara2314/bsam-server/internal/auth/presentation/handler"
)

type VerifyTokenPOSTRequest struct {
	Token string `json:"token"`
}

func VerifyTokenPOST(c *gin.Context) {
	var req VerifyTokenPOSTRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "token is required",
		})
		return
	}

	handler.VerifyTokenPOST(c, req.Token)
}
