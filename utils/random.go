package utils

import (
	"math/rand"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

func RandString(n int) string {
	strSlice := make([]rune, n)
	for i := range strSlice {
		strSlice[i] = letterRunes[rand.Intn(len(letterRunes))]
	}

	return string(strSlice)
}
