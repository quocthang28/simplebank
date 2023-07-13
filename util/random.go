package util

import (
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {

}

func RandomInt(min, max int64) int64 {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return min + r.Int63n(max-min+1)
}

func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

func RandomMoney() int64 {
	return RandomInt(500, 1000)
}

func RandomOwner() string {
	return RandomString(6)
}

func RandomCurrency() string {
	return []string{"EUR", "USD", "VND"}[rand.Intn(3)]
}
