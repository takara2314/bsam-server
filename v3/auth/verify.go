package auth

import (
	"fmt"
	"os"

	"github.com/golang-jwt/jwt"
)

func VerifyJWT(token string) bool {
	info, err := jwt.Parse(token, func(t *jwt.Token) (any, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if info == nil || err != nil {
		fmt.Println(err)
		return false
	}

	if !info.Valid {
		fmt.Println("not valid")
		return false
	}

	return true
}
