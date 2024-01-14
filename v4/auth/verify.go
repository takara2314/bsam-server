package auth

import (
	"log"
	"os"

	"github.com/golang-jwt/jwt"
)

func VerifyJWT(token string) bool {
	info, err := jwt.Parse(token, func(t *jwt.Token) (any, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if info == nil || err != nil {
		return false
	}

	if !info.Valid {
		log.Println("this token is not valid")
		return false
	}

	return true
}
