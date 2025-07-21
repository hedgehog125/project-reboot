// Boilerplate to shorten the start of a servercommon.Error chain
package servercommon

import (
	"errors"
)

func NewRollbackError() *Error {
	return NewError(errors.New("rollback")).DisableLogging()
}

func Send404IfNotFound(err error) *Error {
	return NewError(err).Send404IfNotFound()
}

func SendUnauthorizedIfNotFound(err error) *Error {
	return NewError(err).SendUnauthorizedIfNotFound()
}

func ExpectError(
	err error, expectedError error,
	statusCode int, detail *ErrorDetail,
) *Error {
	return NewError(err).Expect(expectedError, statusCode, detail)
}
func ExpectAnyOfErrors(
	err error, expectedErrors []error,
	statusCode int, detail *ErrorDetail,
) *Error {
	return NewError(err).ExpectAnyOf(expectedErrors, statusCode, detail)
}
