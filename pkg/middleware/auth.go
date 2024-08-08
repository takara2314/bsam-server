package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/takara2314/bsam-server/pkg/auth"
)

const (
	bearerPrefix       = "Bearer "
	bearerPrefixLength = len(bearerPrefix)
)

func AuthToken(jwtSecretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		value := c.GetHeader("Authorization")

		if value == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "authorization header is required",
			})
			return
		}

		token, err := extractTokenFromAuthHeader(value)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		assocID, err := auth.ParseJWT(
			token,
			jwtSecretKey,
		)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "token is invalid",
			})
			return
		}

		c.Set("assoc_id", assocID)
	}
}

func extractTokenFromAuthHeader(value string) (string, error) {
	if value == "" {
		return "", errors.New("authorization header is required")
	}
	if !strings.HasPrefix(value, bearerPrefix) {
		return "", errors.New("authorization header must start with Bearer")
	}
	return strings.TrimPrefix(value, bearerPrefix), nil
}
