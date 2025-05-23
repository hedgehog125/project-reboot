package common

import (
	"errors"
	"fmt"
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
			successfulActionIDs = slices.Delete(successfulActionIDs, index, index+1)
		}
	}
	return successfulActionIDs
}

const (
	ErrTypeDatabase = "database"
	ErrTypeClient   = "client"
	ErrTypeOther    = "other"
)

type ErrorWithCategory struct {
	Err                   error
	Category              string
	ErrDuplicatesCategory bool
}

func (err *ErrorWithCategory) Error() string {
	if err.ErrDuplicatesCategory {
		return err.Err.Error()
	} else {
		return fmt.Sprintf("%v error: %v", err.Category, err.Err.Error())
	}
}
func (err *ErrorWithCategory) Unwrap() error {
	return err.Err
}

func NewErrorWithCategory(err string, category string) *ErrorWithCategory {
	return &ErrorWithCategory{
		Err:      errors.New(err),
		Category: category,
	}
}
func WrapErrorWithCategory(err error, category string) *ErrorWithCategory {
	catErr := &ErrorWithCategory{
		Err:      err,
		Category: category,
	}
	if err == nil {
		catErr.Err = errors.New(category)
		catErr.ErrDuplicatesCategory = true
	}

	return catErr
}

func CategorizeError(err error) string {
	var catErr *ErrorWithCategory
	if errors.As(err, &catErr) {
		return catErr.Category
	}
	if errors.As(err, &sqlite3.Error{}) {
		return ErrTypeDatabase
	}
	if ent.IsConstraintError(err) ||
		ent.IsNotFound(err) ||
		ent.IsNotLoaded(err) ||
		ent.IsNotSingular(err) ||
		ent.IsValidationError(err) ||
		errors.Is(err, ent.ErrTxStarted) {
		return ErrTypeDatabase
	}

	return ErrTypeOther
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
