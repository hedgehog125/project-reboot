package common

import (
	"errors"
	"slices"

	"github.com/hedgehog125/project-reboot/ent" // Note: will have to reorganise if I end up needing to use the common module in schemas
	"github.com/mattn/go-sqlite3"
)

func HasErrors(errs []error) bool {
	for _, err := range errs {
		if err != nil {
			return true
		}
	}
	return false
}

func GetSuccessfulActionIDs(actionIDs []string, errs []*ErrWithStrId) []string {
	successfulActionIDs := make([]string, len(actionIDs))
	copy(successfulActionIDs, actionIDs)

	for _, err := range errs {
		index := slices.Index(successfulActionIDs, err.Id)
		if index != -1 {
			successfulActionIDs = slices.Delete(successfulActionIDs, index, index)
		}
	}
	return successfulActionIDs
}

const (
	ErrorDatabase     = "database"
	ErrorNotFound     = "not-found"
	ErrorUnauthorized = "unauthorized"
	ErrorOther        = "other"
)

func CategorizeError(err error) string {
	if errors.As(err, &sqlite3.Error{}) {
		return ErrorDatabase
	} else if ent.IsConstraintError(err) ||
		ent.IsNotFound(err) ||
		ent.IsNotLoaded(err) ||
		ent.IsNotSingular(err) ||
		ent.IsValidationError(err) ||
		errors.Is(err, ent.ErrTxStarted) {
		return ErrorDatabase
	}

	return ErrorOther
}

type ContextPanic struct {
	Message       string
	ShouldRecover bool
}

// Crashes the whole server rather than just sending a 500
func UnrecoverablePanic(message string) {
	panic(&ContextPanic{
		Message:       message,
		ShouldRecover: false,
	})
}
