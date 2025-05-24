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
	// Highest level categories
	ErrTypeDatabase        = "database"
	ErrTypeTwoFactorAction = "two factor action"
	ErrTypeClient          = "client"
	ErrTypeOther           = "other"
)

type Error struct {
	Err                   error
	Categories            []string
	HighestCategory       string
	ErrDuplicatesCategory bool
}

func (err *Error) Error() string {
	message := ""
	if err.HighestCategory != ErrTypeOther {
		message += fmt.Sprintf("%v error: ", err.HighestCategory)
	}

	reversedCategories := slices.Clone(err.Categories)
	slices.Reverse(reversedCategories) // Highest to lowest level
	if err.ErrDuplicatesCategory {
		reversedCategories = DeleteSliceIndex(reversedCategories, -1) // Ignore the lowest (last) category since it duplicates the error
	}

	for _, category := range reversedCategories {
		message += fmt.Sprintf("%v error: ", category)
	}

	return message + err.Err.Error()
}
func (err *Error) Unwrap() error {
	return err.Err
}
func (err *Error) Is(target error) bool {
	// Needed so that errors.Is(err.AddCategory("extra category"), err) == true
	// We don't really care if the properties on this struct are different, only that the underlying error is the same

	if target == nil {
		return false
	}
	targetStruct, ok := target.(*Error)
	if !ok {
		return false
	}
	return err.Err == targetStruct.Err
}

func (err *Error) SetHighestCategory(category string) *Error {
	copiedErr := err.Copy()
	copiedErr.HighestCategory = category
	return copiedErr
}

// Note: not to be confused with `HighestCategory` which is something like "database". This is a level lower, e.g "create user"
func (err *Error) HighestSpecificCategory() string {
	return err.Categories[len(err.Categories)-1]
}
func (err *Error) AllCategories() []string {
	return slices.Concat(err.Categories, []string{err.HighestCategory})
}

// Note: this mutates the error, so ensure it's been wrapped or copied first
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
	copiedErr.Categories = slices.Clone(err.Categories)

	return &copiedErr
}

// e.g err.HasCategories(common.ErrTypeDatabase, "create user")
func (err *Error) HasCategories(requiredCategories ...string) bool {
	allCategories := err.AllCategories()
	if len(requiredCategories) > len(allCategories) {
		return false
	}

	slices.Reverse(allCategories) // Check from the highest level first, so lower level can be implicitly ignored
	for i, requiredCategory := range requiredCategories {
		if requiredCategory != "*" && allCategories[i] != requiredCategory {
			return false
		}
	}
	return true
}

func NewErrorWithCategories(err string, highestCategory string, categories ...string) *Error {
	return &Error{
		Err:             errors.New(err),
		Categories:      categories,
		HighestCategory: highestCategory,
	}
}
func WrapErrorWithCategory(err error, highestCategory string, categories ...string) *Error {
	catErr := &Error{
		Err:             err,
		Categories:      categories,
		HighestCategory: highestCategory,
	}
	if err == nil {
		if len(categories) == 0 {
			panic("you must provide at least one category in addition to the highest category or provide an error")
		}

		catErr.Err = errors.New(categories[0])
		catErr.ErrDuplicatesCategory = true
	}

	return catErr
}

func CategorizeError(err error) string {
	var commErr *Error
	if errors.As(err, &commErr) {
		return commErr.HighestCategory
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
