package common

import "math/rand"

func RandPositiveNegativeRange(maxValue float64) float64 {
	return (rand.Float64() * maxValue * 2) - maxValue
}
