package util

import (
	"math/rand"
	"time"

	"github.com/brianvoe/gofakeit"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandomInt() int64 {
	return int64(gofakeit.Number(0, 1000))
}

func RandomUserName() string {
	return gofakeit.Username()
}

func RandomName() string {
	return gofakeit.Name()
}

func RandomCurrency() string {
	return "USD"
}

func RandomEmail() string {
	return gofakeit.Email()
}

func RandomPassword() string {
	return gofakeit.Password(true, true, true, true, false, 12)
}

func RandomString(amount int) string {
	var r = rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]rune, amount)
	for i := range b {
		b[i] = letters[r.Intn(len(letters))]
	}
	return string(b)
}
