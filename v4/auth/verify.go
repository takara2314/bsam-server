package auth

import (
	"log"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

func VerifyJWT(token string) bool {
	info, err := jwt.Parse(token, func(_ *jwt.Token) (any, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}))

	if info == nil || err != nil {
		return false
	}

	if !info.Valid {
		log.Println("this token is not valid")
		return false
	}

	return true
}
