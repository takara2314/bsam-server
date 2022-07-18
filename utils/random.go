package utils

import (
	"math/rand"
	"time"
)

var (
	letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
)

func RandString(n int) string {
	rand.Seed(time.Now().UnixNano())

	strSlice := make([]rune, n)
	for i := range strSlice {
		strSlice[i] = letterRunes[rand.Intn(len(letterRunes))]
	}

	return string(strSlice)
}
