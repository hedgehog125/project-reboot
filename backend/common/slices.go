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
	return checkPathPattern(path, SimplifyPathPattern(pattern))
}
func SimplifyPathPattern(pattern []string) []string {
	simplifiedPattern := []string{}

	inDoubleWildcardChain := false

	inTripleWildcardChain := false
	tripleWildcardReplacement := ""
	finishTripleWildcardChain := func() {
		simplifiedPattern = append(simplifiedPattern, tripleWildcardReplacement)
		inTripleWildcardChain = false
	}

	for _, item := range pattern {
		if inDoubleWildcardChain {
			if item == "***" {
				continue
			}

			inDoubleWildcardChain = false
		}
		if inTripleWildcardChain {
			if item == "***" {
				continue
			}
			if item == "**" {
				if tripleWildcardReplacement == "**" {
					finishTripleWildcardChain()
					inTripleWildcardChain = true
				} else {
					tripleWildcardReplacement = "**"
				}
				continue
			}
			if item == "*" && tripleWildcardReplacement == "***" {
				tripleWildcardReplacement = "**"
				continue
			}

			finishTripleWildcardChain()
		}

		if item == "***" {
			tripleWildcardReplacement = "***"
			inTripleWildcardChain = true
			continue
		}
		if item == "**" {
			inDoubleWildcardChain = true
		}

		simplifiedPattern = append(simplifiedPattern, item)
	}
	if inTripleWildcardChain {
		finishTripleWildcardChain()
	}
	return simplifiedPattern
}
func checkPathPattern(path []string, pattern []string) bool {
	remainingPattern := slices.Clone(pattern)
	remainingPath := slices.Clone(path)
	pathIndex := 0
	for patternIndex := 0; patternIndex < len(remainingPattern); {
		patternItem := remainingPattern[patternIndex]
		if pathIndex > len(remainingPath)-1 {
			return patternItem == "***" && patternIndex == len(remainingPattern)-1
		}

		if patternItem == "*" {
			patternIndex++
			pathIndex++
			continue
		}
		if patternItem == "**" || patternItem == "***" {
			if patternIndex == len(remainingPattern)-1 {
				return true // ** matches this item and any after it
			}

			remainingPattern = remainingPattern[patternIndex:]
			// patternIndex should be set to 0, but it's not read
			firstPatternLiteralIndex := slices.IndexFunc(remainingPattern, func(value string) bool {
				return value != "**" && value != "***" // Find first literal or *
			})
			if firstPatternLiteralIndex == -1 { // Pattern ends with 2 or more **s
				// Check there are enough items left for the number of **s
				return len(remainingPath)-pathIndex >= len(remainingPattern)
			}

			minItemsBefore := 0
			outerPattern := remainingPattern[:firstPatternLiteralIndex]
			if outerPattern[0] == "**" {
				minItemsBefore = len(outerPattern)
			}
			remainingPattern = remainingPattern[firstPatternLiteralIndex:] // It might loop back around to the start of this

			remainingPath = remainingPath[pathIndex+minItemsBefore:]
			pathIndex = 0

			for pathIndex < len(remainingPath) {
				firstLiteralString := remainingPattern[0]
				if firstLiteralString != "*" {
					firstPathMatchIndex := slices.Index(remainingPath, firstLiteralString)
					if firstPathMatchIndex == -1 {
						return false
					}
					remainingPath = remainingPath[firstPathMatchIndex:]
				}
				if checkPathPattern(remainingPath, remainingPattern) {
					return true
				}
				if len(remainingPath) == 1 {
					return false
				}
				remainingPath = remainingPath[1:]
			}
			return false
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
