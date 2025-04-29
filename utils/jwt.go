package utils

import (
	"errors"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

var ErrInvalidJWT = errors.New("invalid jwt")

func GetUserIDFromJWT(t string) (string, error) {
	token, err := jwt.Parse(t, func(_ *jwt.Token) (any, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}))

	if token == nil || err != nil {
		return "", ErrInvalidJWT
	}

	payload, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", ErrInvalidJWT
	}

	userID, ok := payload["user_id"].(string)
	if !ok {
		return "", ErrInvalidJWT
	}

	return userID, nil
}
