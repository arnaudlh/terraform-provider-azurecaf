package utils

import (
	"math/rand"
)

var (
	alphagenerator = []rune("abcdefghijklmnopqrstuvwxyz")
)

// RandSeq generates a random value to add to the resource names
func RandSeq(length int, seed int64) string {
	if length <= 0 {
		return ""
	}
	r := rand.New(rand.NewSource(seed))
	b := make([]rune, length)
	for i := range b {
		b[i] = alphagenerator[r.Intn(len(alphagenerator))]
	}
	return string(b)
}
