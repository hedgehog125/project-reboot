package common

import "slices"

func DeleteSliceIndex[T any](slice []T, index int) []T {
	actualIndex := index
	if index < 0 {
		actualIndex = len(slice) + index
	}

	return slices.Delete(slice, actualIndex, actualIndex+1)
}
