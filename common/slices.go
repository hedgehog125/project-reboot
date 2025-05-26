package common

import "slices"

func DeleteSliceIndex[T any](slice []T, index int) []T {
	actualIndex := index
	if index < 0 {
		actualIndex = len(slice) + index
	}

	return slices.Delete(slice, actualIndex, actualIndex+1)
}

func CheckPathPattern(path []string, pattern []string) bool {
	remainingPattern := slices.Clone(pattern)
	remainingPath := slices.Clone(path)
	pathIndex := 0
	for patternIndex := 0; patternIndex < len(remainingPattern); {
		patternItem := remainingPattern[patternIndex]
		if pathIndex > len(remainingPath)-1 {
			return patternItem == "***"
		}

		if patternItem == "*" {
			patternIndex++
			pathIndex++
			continue
		}
		if patternItem == "**" || patternItem == "***" {
			if patternIndex == len(remainingPattern)-1 {
				// TODO: test this branch
				return patternItem == "***" || pathIndex < len(remainingPath)-1
			}
			remainingPattern = remainingPattern[patternIndex+1:]
			patternIndex = 0
			combinedPatternItem := patternItem
			foundMatch := false
			for i, value := range remainingPattern {
				if value == "**" {
					combinedPatternItem = "**"
				}

				if value != "**" && value != "***" {
					remainingPattern = remainingPattern[i:]
					patternItem = combinedPatternItem
					foundMatch = true
					break
				}
			}
			if !foundMatch {
				return true
			}
			remainingPath = remainingPath[pathIndex:]
			pathIndex = 0

			firstLiteralIndex := 0
			firstLiteralString := remainingPattern[0]
			if firstLiteralString != "*" {
				firstLiteralIndex = slices.Index(remainingPath, firstLiteralString)
			}
			if firstLiteralIndex == -1 {
				return false
			}
			if firstLiteralIndex == 0 && patternItem == "**" { // "**" must match at least one item
				return false
			}

			remainingPath = remainingPath[firstLiteralIndex:]
			for pathIndex < len(remainingPath) {
				pathItem := remainingPath[pathIndex]
				patternItem = remainingPattern[patternIndex]
				if patternIndex == len(remainingPattern)-1 {
					// TODO: test this branch
					if patternItem == "**" || patternItem == "***" {
						return true
					}
					if pathIndex == len(remainingPath)-1 &&
						(patternItem == "*" || pathItem == patternItem) {
						return true
					}

					// There's still more of the path but no pattern to match it
					return false
				}
				if patternItem == "**" || patternItem == "***" {
					// Recursive so that nested backtracking works correctly
					return CheckPathPattern(remainingPath[pathIndex:], remainingPattern[patternIndex:])
				}

				if patternItem != "*" && pathItem != patternItem {
					patternIndex = 0
					if pathItem != firstLiteralString {
						pathIndex++
					}
					continue
				}
				pathIndex++
				patternIndex++
			}
			continue
		}

		pathItem := remainingPath[pathIndex]
		if pathItem != patternItem {
			return false
		}
		patternIndex++
		pathIndex++
	}
	return true
}
