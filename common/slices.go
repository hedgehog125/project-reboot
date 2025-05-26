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
	// TODO: ** should match 1 or more
	// *** should match 0 or more

	// **/documents/projects/**/assets/**
	// /home/nico/documents/projects/unity/experiments/assets/hats/coolHat.png
	// Because of **, get index of documents
	// Check the next is projects, it is
	// Because of **, get the index of assets, starting from the current position
	// Next is a ** so it passes

	// fake-path/documents/home/nico/documents/projects/unity/experiments/assets/hats/coolHat.png
	// Because of **, get index of documents
	// Next is home so search for next
	// Found another, which has projects after it
	// Same as previous example...

	remainingPattern := slices.Clone(pattern) // TODO: is this needed?
	remainingPath := slices.Clone(path)
	pathIndex := 0
	for patternIndex := 0; patternIndex < len(remainingPattern); {
		if pathIndex >= len(remainingPath)-1 {
			return false
		}

		patternItem := pattern[patternIndex]
		if patternItem == "*" {
			patternIndex++
			pathIndex++
			continue
		}
		if patternItem == "**" {
			if patternIndex == len(remainingPattern)-1 {
				return true
			}
			remainingPattern = remainingPattern[patternIndex+1:]
			patternIndex = 0
			firstLiteralIndex := slices.IndexFunc(remainingPattern, func(value string) bool {
				return !(value == "*" || value == "**") // Ignore these completely
			})
			if firstLiteralIndex == -1 {
				return true
			}
			remainingPattern = remainingPattern[firstLiteralIndex:]

			remainingPath = remainingPath[pathIndex:]
			pathIndex = 0

			firstLiteralString := remainingPattern[0]
			firstLiteralIndex = slices.Index(remainingPath, firstLiteralString)
			if firstLiteralIndex == -1 {
				return false
			}
			if firstLiteralIndex == len(remainingPath)-1 {
				// TODO: slight optimisation so remainingPath can be safely set to remainingPath[firstLiteralIndex+1:]
			}

			remainingPath = remainingPath[firstLiteralIndex:]
			for pathIndex < len(remainingPath) {
				pathItem := remainingPath[pathIndex]
				if patternIndex >= len(remainingPattern)-1 {
					return false // There's still more of the path
				}
				patternItem = remainingPattern[patternIndex]
				if patternItem == "**" {
					// TODO: isn't this recursive?
					// Maybe this whole branch can be flattened with an isNMatching variable?
					// The == "**" branch would set up the pattern looping like below
					// Rather than isNMatching, probably currentWildcard
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
