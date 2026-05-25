package utils

import (
	"math/rand"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func RandomString(n int) string {
	s := make([]byte, n)
	for i := range s {
		s[i] = alphabet[rand.Intn(len(alphabet))]
	}
	return string(s)
}

func RandomOwner() string {
	return RandomString(6)
}

func RandomMoney() int64 {
	return RandomInt(0, 1000)
}

func RandomCurrency() string {
	currencies := []string{"USD", "EUR", "CAD"}
	return currencies[rand.Intn(len(currencies))]
}
