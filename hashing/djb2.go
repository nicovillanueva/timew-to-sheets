package hashing

import (
// "math"
)

// Djb2 serializes a string into int64
// If forcePositive is true, whenever a negative int64 is to be returned, it just flips it. Works for me, deal with it.
func Djb2(s string, forcePositive bool) int64 {
	var hash int64 = 5381
	for _, ch := range s {
		hash = ((hash << 5) + hash) + int64(ch)
	}
	if hash < 0 && forcePositive {
		hash *= -1
	}
	return hash
}
