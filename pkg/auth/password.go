package auth

import "github.com/matthewhartstonge/argon2"

func HashPassword(password string) string {
	if password == "" {
		return ""
	}

	argon := argon2.DefaultConfig()

	hashed, _ := argon.HashEncoded([]byte(password))

	return string(hashed)
}

func VerifyPassword(password string, hashedPassword string) bool {
	passwordMatched, err := argon2.VerifyEncoded(
		[]byte(password),
		[]byte(hashedPassword),
	)

	if err != nil {
		return false
	}

	return passwordMatched
}
