package middleware

import (
	"bsam-server/v3/abort"
	"bsam-server/v3/auth"

	"github.com/gin-gonic/gin"
)

func AuthJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, valid := getTokenFromAuthHeader(c.GetHeader("Authorization"))
		if !valid {
			abort.BadRequest(c)
			return
		}

		if !auth.VerifyJWT(token) {
			abort.Unauthorized(c)
			return
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
