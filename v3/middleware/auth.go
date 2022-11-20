package middleware

import (
	"bsam-server/v3/auth"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, valid := getTokenFromAuthHeader(c.GetHeader("Authorization"))
		if !valid {
			c.Status(http.StatusBadRequest)
			c.Abort()
		}

		if !auth.VerifyJWT(token) {
			c.Status(http.StatusUnauthorized)
			c.Abort()
		}
	}
}

func getTokenFromAuthHeader(value string) (string, bool) {
	if value == "" {
		return "", false
	}

	if len(value) < 8 {
		return "", false
	}

	if value[:6] != "Bearer" {
		return "", false
	}

	return value[7:], true
}
