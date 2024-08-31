package util

import (
	"fmt"
	"math/rand"
	"strings"
)

const alphabets = "abcdefghijklmnopqrstuvwxyz"

func init() {
	fmt.Println("Init called")
}

func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func RandomString(n int) string {

	var sb strings.Builder

	for i := 0; i < n; i++ {
		c := alphabets[rand.Intn(len(alphabets))]
		sb.WriteByte(c)
	}
	return sb.String()
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

func RandomEmail() string {
	return fmt.Sprintf("%s@email.com", RandomString(6))
}
