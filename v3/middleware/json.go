package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func CheckMIME(mime string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !isThisMIME(c.GetHeader("Content-Type"), mime) {
			c.JSON(http.StatusBadRequest, nil)
		}
	}
}

func isThisMIME(contentType string, mime string) bool {
	return strings.Contains(contentType, mime)
}
