package helpers

import (
	"crypto/rand"
	"fmt"
	"time"

	rand_exp "golang.org/x/exp/rand"
)

var rng *rand_exp.Rand

func init() {
	rng = rand_exp.New(rand_exp.NewSource(uint64(time.Now().UnixNano())))
}

func RandomString(length int) string {
	b := make([]byte, length)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func RandomChar(length int) string {
	chars := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	code := make([]byte, length)
	for i := range code {
		code[i] = chars[rng.Intn(len(chars))]
	}
	return string(code)
}
