package common

func CountBools(bools ...bool) int {
	count := 0
	for _, b := range bools {
		if b {
			count++
		}
	}
	return count
}

func AllOrNone(bools ...bool) bool {
	count := CountBools(bools...)
	return count == 0 || count == len(bools)
}
