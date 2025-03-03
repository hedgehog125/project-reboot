package util

func HasErrors(errs []error) bool {
	for _, err := range errs {
		if err != nil {
			return true
		}
	}
	return false
}
