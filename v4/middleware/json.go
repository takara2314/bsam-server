package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func CheckMIME(mime string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !isThisMIME(c.GetHeader("Content-Type"), mime) {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "Content-Type '" + mime + "' required.",
			})
			return
		}
	}
}

func isThisMIME(contentType string, mime string) bool {
	return strings.Contains(contentType, mime)
}
