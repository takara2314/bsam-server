package handlerValidating

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/takara2314/bsam-server/internal/auth/presentation/handler"
)

func VerifyPasswordPOST(c *gin.Context) {
	var req handler.VerifyPasswordPOSTRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "assoc_id and password are required",
		})
		return
	}

	handler.VerifyPasswordPOST(c, req)
}
