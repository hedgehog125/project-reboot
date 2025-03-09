package servercommon

import (
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
	return err.SendStatusIfNotFound(404, "NOT_FOUND")
}

// 401 is HTTP unauthorized
func (err *ContextError) Send401IfNotFound() *ContextError {
	return err.SendStatusIfNotFound(401, "UNAUTHORIZED")
}

func (err *ContextError) SendStatusIfNotFound(statusCode int, errorCode string) *ContextError {
	if ent.IsNotFound(err.Err) {
		if statusCode != -1 {
			err.Status = statusCode
		}
		if errorCode != "" {
			err.ErrorCodes = append(err.ErrorCodes, errorCode)
		}
		err.ShouldLog = false
	}
	return err
}
