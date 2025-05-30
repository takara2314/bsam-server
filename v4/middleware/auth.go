package middleware

import (
	"github.com/takara2314/bsam-server/v4/abort"
	"github.com/takara2314/bsam-server/v4/auth"

	"github.com/gin-gonic/gin"
)

const BearerTokenPrefixSize = 7

func AuthJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		value := c.GetHeader("Authorization")

		if value == "" {
			abort.Unauthorized(c)
			return
		}

		token, valid := getTokenFromAuthHeader(value)
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

	if len(value) <= BearerTokenPrefixSize {
		return "", false
	}

	if value[:6] != "Bearer" {
		return "", false
	}

	return value[7:], true
}
