package auth

import (
	"os"

	"github.com/golang-jwt/jwt/v5"
)

func GetPayloadFromJWT(token string) (map[string]any, bool) {
	info, err := jwt.Parse(token, func(_ *jwt.Token) (any, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}))

	if info == nil || err != nil {
		return nil, false
	}

	payload, ok := info.Claims.(jwt.MapClaims)
	if !ok {
		return nil, false
	}

	return payload, true
}
