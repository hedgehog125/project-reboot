package common

import "math/rand"

// Note: not cryptographically random
func RandPositiveNegativeRange(maxValue float64) float64 {
	//nolint: gosec
	return (rand.Float64() * maxValue * 2) - maxValue
}
