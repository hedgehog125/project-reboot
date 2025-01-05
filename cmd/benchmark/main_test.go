package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddIntArray(t *testing.T) {
	intArray := make([]int32, 3)
	max := int32(10)

	addIntArray(&intArray, 1, max)
	assert.Equal(t, []int32{0, 0, 1}, intArray)

	addIntArray(&intArray, 10, max)
	assert.Equal(t, []int32{0, 1, 1}, intArray)

	addIntArray(&intArray, 9, max)
	assert.Equal(t, []int32{0, 2, 0}, intArray)

	intArray = make([]int32, 3)
	max = int32(3)

	addIntArray(&intArray, 2, max)
	assert.Equal(t, []int32{0, 0, 2}, intArray)

	addIntArray(&intArray, 3, max)
	assert.Equal(t, []int32{0, 1, 2}, intArray)

	addIntArray(&intArray, 9, max)
	assert.Equal(t, []int32{1, 1, 2}, intArray)
}
