package common

import "math/rand"

func RandPositiveNegativeRange(max float64) float64 {
	return (rand.Float64() * max * 2) - max
}
