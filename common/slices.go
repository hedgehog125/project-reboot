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
	return checkPathPattern(path, simplifyPathPattern(pattern))
}
func simplifyPathPattern(pattern []string) []string {
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
				// TODO: test this branch
				return patternItem == "***" || pathIndex < len(remainingPath)-1
			}
			remainingPattern = remainingPattern[patternIndex+1:]
			patternIndex = 0
			remainingPath = remainingPath[pathIndex:]
			pathIndex = 0

			firstLiteralString := remainingPattern[0]
			if !(firstLiteralString == "*" || firstLiteralString == "**" || firstLiteralString == "***") {
				firstLiteralIndex := slices.Index(remainingPath, firstLiteralString)
				if firstLiteralIndex == -1 {
					return false
				}
				if firstLiteralIndex == 0 && patternItem == "**" { // "**" must match at least one item
					return false
				}
				remainingPath = remainingPath[firstLiteralIndex:]
			}

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
					return checkPathPattern(remainingPath[pathIndex:], remainingPattern[patternIndex:])
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
