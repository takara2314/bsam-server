package utils

import (
	"errors"
	"os"

	"github.com/golang-jwt/jwt"
)

var (
	ErrInvalidJWT = errors.New("invalid jwt")
)

func GetUserIDFromJWT(t string) (string, error) {
	token, err := jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if token == nil || err != nil {
		return "", ErrInvalidJWT
	}

	if !token.Valid {
		return "", ErrInvalidJWT
	}

	return token.Claims.(jwt.MapClaims)["user_id"].(string), nil
}
