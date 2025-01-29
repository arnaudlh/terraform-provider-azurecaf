package utils

import (
	"math/rand"
)

var (
	alphagenerator = []rune("abcdefghijklmnopqrstuvwxyz")
)

// RandSeq generates a random value to add to the resource names
func RandSeq(length int, seed int64) string {
	if length == 0 {
		return ""
	}
	// Create a new random source with the given seed
	// Using New(NewSource()) as recommended since Go 1.20
	r := rand.New(rand.NewSource(seed))
	// generate at least one random character
	b := make([]rune, length)
	for i := range b {
		// We need the random generated string to start with a letter
		b[i] = alphagenerator[r.Intn(len(alphagenerator)-1)]
	}
	return string(b)
}
