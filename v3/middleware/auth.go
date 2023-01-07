package middleware

import (
	"bsam-server/v3/abort"
	"bsam-server/v3/auth"
	"fmt"

	"github.com/gin-gonic/gin"
)

func AuthJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		value := c.GetHeader("Authorization")
		fmt.Println("Authorization:", value)

		if value == "" {
			abort.Unauthorized(c)
			fmt.Println("A")
			return
		}

		token, valid := getTokenFromAuthHeader(value)
		if !valid {
			abort.BadRequest(c)
			fmt.Println("B")
			return
		}

		if !auth.VerifyJWT(token) {
			abort.Unauthorized(c)
			fmt.Println("C")
			return
		}

		fmt.Println("Pass")
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
