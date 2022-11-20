package auth

import (
	"os"

	"github.com/golang-jwt/jwt"
)

func GetPayloadFromJWT(token string) (map[string]interface{}, bool) {
	info, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if info == nil || err != nil {
		return nil, false
	}

	if !info.Valid {
		return nil, false
	}

	return info.Claims.(jwt.MapClaims), true
}
