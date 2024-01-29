package util

import (
	"github.com/brianvoe/gofakeit"
)

func RandomInt() int64 {
	return int64(gofakeit.Number(0, 1000))
}

func RandomName() string {
	return gofakeit.Name()
}

func RandomCurrency() string{
	return gofakeit.Currency().Short
}
