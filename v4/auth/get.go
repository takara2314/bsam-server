package auth

import (
	"os"

	"github.com/golang-jwt/jwt"
)

func GetPayloadFromJWT(token string) (map[string]any, bool) {
	info, err := jwt.Parse(token, func(t *jwt.Token) (any, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if info == nil || err != nil {
		return nil, false
	}

	if !info.Valid {
		return nil, false
	}

	payload, ok := info.Claims.(jwt.MapClaims)
	if !ok {
		return nil, false
	}

	return payload, true
}
