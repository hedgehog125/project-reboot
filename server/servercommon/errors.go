package servercommon

import (
	"errors"

	"github.com/hedgehog125/project-reboot/common"
	"github.com/hedgehog125/project-reboot/ent"
)

type ContextError struct {
	Err        error
	Status     int // Set to -1 to keep the current code
	ErrorCodes []string
	Category   string
	ShouldLog  bool
}

func (err *ContextError) Error() string {
	return err.Err.Error()
}

func (err *ContextError) Unwrap() error {
	return err.Err
}

func NewContextError(err error) *ContextError {
	return &ContextError{
		Err:        err,
		Status:     500,
		ErrorCodes: []string{},
		Category:   common.CategorizeError(err),
		ShouldLog:  true,
	}
}

// Adds the final defaults
func (err *ContextError) Finish() *ContextError {
	if err.Status == 500 {
		if err.ErrorCodes != nil {
			err.ErrorCodes = append(err.ErrorCodes, "INTERNAL")
		}
	}

	return err
}

func (err *ContextError) Send404IfNotFound() *ContextError {
	return sendStatusIfNotFound(err, 404, "", true)
}

func (err *ContextError) SendUnauthorizedIfNotFound() *ContextError {
	return sendStatusIfNotFound(err, 401, "", false)
}

func (err *ContextError) Expect(
	expectedError error,
	statusCode int, errorCode string,
) *ContextError {
	return sendStatusAndCodeIfCondition(err, errors.Is(err, expectedError), statusCode, errorCode, true)
}
func (err *ContextError) ExpectAnyOf(
	expectedErrors []error,
	statusCode int, errorCode string,
) *ContextError {
	isExpected := false
	for _, expectedErr := range expectedErrors {
		if errors.Is(err, expectedErr) {
			isExpected = true
			break
		}
	}
	return sendStatusAndCodeIfCondition(err, isExpected, statusCode, errorCode, true)
}

func sendStatusIfNotFound(
	err *ContextError, statusCode int,
	errorCode string, preventLog bool,
) *ContextError {
	return sendStatusAndCodeIfCondition(err, ent.IsNotFound(err.Err), statusCode, errorCode, preventLog)
}

func sendStatusAndCodeIfCondition(
	err *ContextError, condition bool, statusCode int,
	errorCode string, preventLog bool,
) *ContextError {
	if condition {
		if statusCode != -1 {
			err.Status = statusCode
		}
		if errorCode != "" {
			err.ErrorCodes = append(err.ErrorCodes, errorCode)
		}
		if preventLog {
			err.ShouldLog = false
		}
	}
	return err
}
