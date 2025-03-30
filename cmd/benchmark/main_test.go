package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddIntArray(t *testing.T) {
	t.Parallel()
	intArray := make([]int32, 3)
	maxNum := int32(10)

	addIntArray(intArray, 1, maxNum)
	assert.Equal(t, []int32{0, 0, 1}, intArray)

	addIntArray(intArray, 10, maxNum)
	assert.Equal(t, []int32{0, 1, 1}, intArray)

	addIntArray(intArray, 9, maxNum)
	assert.Equal(t, []int32{0, 2, 0}, intArray)

	intArray = make([]int32, 3)
	maxNum = int32(3)

	addIntArray(intArray, 2, maxNum)
	assert.Equal(t, []int32{0, 0, 2}, intArray)

	addIntArray(intArray, 3, maxNum)
	assert.Equal(t, []int32{0, 1, 2}, intArray)

	addIntArray(intArray, 9, maxNum)
	assert.Equal(t, []int32{1, 1, 2}, intArray)
}

func TestAddIntArray_Overflows(t *testing.T) {
	t.Parallel()
	intArray := make([]int32, 3)

	assert.True(t, addIntArray(intArray, 1001, 10))
	assert.Equal(t, []int32{0, 0, 1}, intArray)
}
