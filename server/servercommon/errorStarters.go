// Boilerplate to shorten the start of a ContextError chain
package servercommon

func Send404IfNotFound(err error) *ContextError {
	return NewContextError(err).Send404IfNotFound()
}

func SendUnauthorizedIfNotFound(err error) *ContextError {
	return NewContextError(err).SendUnauthorizedIfNotFound()
}

func ExpectError(err error, expectedError error, statusCode int, errorCode string) *ContextError {
	return NewContextError(err).Expect(expectedError, statusCode, errorCode)
}
