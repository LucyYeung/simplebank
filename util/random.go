package util

import (
	"math/rand"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func RandomString(length int) string {
	result := make([]byte, length)

	for i := range result {
		result[i] = alphabet[rand.Intn(len(alphabet))]
	}
	return string(result)
}

func RandomOwner() string {
	return RandomString(6)
}

func RandomMoney() int64 {
	return RandomInt(0, 1000)
}

func RandomCurrency() string {
	currencies := []string{EUR, USD, CAD}
	return currencies[rand.Intn(len(currencies))]
}
