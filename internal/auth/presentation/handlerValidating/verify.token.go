package handlerValidating

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/takara2314/bsam-server/internal/auth/presentation/handler"
)

func VerifyTokenPOST(c *gin.Context) {
	var req handler.VerifyTokenPOSTRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "token is required",
		})
		return
	}

	handler.VerifyTokenPOST(c, handler.VerifyTokenPOSTRequest{
		Token: req.Token,
	})
}
