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
	ErrTypeDatabase        = "database"
	ErrTypeTwoFactorAction = "two factor action"
	ErrTypeClient          = "client"
	ErrTypeOther           = "other"
)

// TODO: root categories are really more of a top level category. Maybe should be kept separately? You have to set it at creation but it can be overridden later in the chain?

type Error struct {
	Err                   error
	Categories            []string
	ErrDuplicatesCategory bool
}

func (err *Error) Error() string {
	if err.ErrDuplicatesCategory {
		return err.Err.Error()
	} else {
		return fmt.Sprintf("%v error: %v", err.Categories, err.Err.Error()) // TODO: update
	}
}
func (err *Error) Unwrap() error {
	return err.Err
}

func (err *Error) HighestCategory() string {
	return err.Categories[len(err.Categories)-1]
}
func (err *Error) LowestCategory() string {
	return err.Categories[0]
}
func (err *Error) SetLowestCategory(category string) *Error {
	err.Categories[0] = category
	return err
}
func (err *Error) PopCategory() string {
	if len(err.Categories) == 0 {
		return ""
	}

	highestCategory := err.Categories[len(err.Categories)-1]
	err.Categories = slices.Delete(err.Categories, len(err.Categories)-1, len(err.Categories))
	return highestCategory
}
func (err *Error) AddCategory(category string) *Error {
	copiedErr := err.Copy()
	copiedErr.Categories = append(copiedErr.Categories, category)
	return copiedErr
}
func (err Error) Copy() *Error {
	copiedErr := err
	copiedErr.Categories = make([]string, len(err.Categories))
	copy(copiedErr.Categories, err.Categories)

	return &copiedErr
}

// e.g err.HasCategories(common.ErrTypeDatabase, "create user")
func (err *Error) HasCategories(requiredCategories ...string) bool {
	for i, requiredCategory := range requiredCategories {
		if i >= len(err.Categories) {
			return false
		}
		if requiredCategory != "*" && err.Categories[i] != requiredCategory {
			return false
		}
	}
	return true
}

func NewErrorWithCategories(err string, rootCategory string, categories ...string) *Error {
	return &Error{
		Err:        errors.New(err),
		Categories: append([]string{rootCategory}, categories...),
	}
}
func WrapErrorWithCategory(err error, rootCategory string, categories ...string) *Error {
	catErr := &Error{
		Err:        err,
		Categories: append([]string{rootCategory}, categories...),
	}
	if err == nil {
		if len(categories) == 0 {
			panic("you must provide at least one category in addition to the root category or provide an error")
		}

		catErr.Err = errors.New(rootCategory)
		catErr.ErrDuplicatesCategory = true
	}

	return catErr
}

func CategorizeError(err error) string {
	var commErr *Error
	if errors.As(err, &commErr) {
		return commErr.LowestCategory()
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
