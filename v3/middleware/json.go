package middleware

import (
	"bsam-server/v3/abort"
	"strings"

	"github.com/gin-gonic/gin"
)

func CheckMIME(mime string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !isThisMIME(c.GetHeader("Content-Type"), mime) {
			abort.BadRequest(c)
			return
		}
	}
}

func isThisMIME(contentType string, mime string) bool {
	return strings.Contains(contentType, mime)
}
