// Boilerplate to shorten the start of a ContextError chain
package servercommon

func Send404IfNotFound(err error) *ContextError {
	return NewContextError(err).Send404IfNotFound()
}

// 401 is HTTP unauthorized
func Send401IfNotFound(err error) *ContextError {
	return NewContextError(err).Send401IfNotFound()
}

func SendStatusIfNotFound(err error, statusCode int, errorCode string) *ContextError {
	return NewContextError(err).SendStatusIfNotFound(statusCode, errorCode)
}
