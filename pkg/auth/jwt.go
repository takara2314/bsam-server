package auth

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/samber/oops"
)

func CreateJWT(assocID string, exp time.Time, secretKey string) string {
	claims := jwt.MapClaims{
		"assoc_id": assocID,
		"iat":      time.Now().Unix(),
		"exp":      exp.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenStr, _ := token.SignedString([]byte(secretKey))

	return tokenStr
}

func ParseJWT(tokenStr string, secretKey string) (string, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return "", oops.
			In("auth.parseJWT").
			Wrapf(err, "failed to parse jwt token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", oops.
			In("auth.parseJWT").
			Wrapf(err, "failed to parse jwt token claims")
	}

	assocID, ok := claims["assoc_id"].(string)
	if !ok {
		return "", oops.
			In("auth.parseJWT").
			Wrapf(err, "failed to parse jwt token claims")
	}

	return assocID, nil
}
