package common

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/hedgehog125/project-reboot/ent" // Note: will have to reorganise if I end up needing to use the common module in schemas
	"github.com/mattn/go-sqlite3"
)

// TODO: rename to ErrType1 and ErrType2?
const (
	// General categories
	ErrTypeDatabase = "database [general]"
	ErrTypeTimeout  = "timeout [general]"
	ErrTypeNetwork  = "network [general]"
	ErrTypeDisk     = "disk [general]"
	ErrTypeMemory   = "memory [general]"
	ErrTypeAPI      = "api [general]"
	ErrTypeClient   = "client [general]"
	// If there's no applicable general category, there should be no [general] category on the error. Functions that return a category should return an empty string
	ErrTypeOther = "other" // Used to maintain the hierarchy, but doesn't have a [general] tag because the lower level error is more useful

	// Package categories
	ErrTypeCore            = "core [package]"
	ErrTypeJobs            = "jobs [package]"
	ErrTypeTwoFactorAction = "two factor action [package]"
	ErrTypeMessengers      = "messengers [package]"
	ErrTypeDbCommon        = "db common [package]"
	ErrTypeServerCommon    = "server common [package]"
	// Similar idea here if it's unknown
)

var ErrWrapperDatabase = NewDynamicErrorWrapper(func(err error) *Error {
	wrappedErr := WrapErrorWithCategories(err)
	sqliteErr := sqlite3.Error{}
	if errors.Is(err, context.DeadlineExceeded) {
		return wrappedErr.AddCategories(ErrTypeTimeout, ErrTypeDatabase)
	} else if errors.As(err, &sqliteErr) {
		if slices.Index([]sqlite3.ErrNo{
			sqlite3.ErrFull,
			sqlite3.ErrAuth,
			sqlite3.ErrReadonly,
			sqlite3.ErrBusy,
			sqlite3.ErrNoLFS,
			sqlite3.ErrCantOpen,
			sqlite3.ErrIoErr,
			sqlite3.ErrLocked,
		}, sqliteErr.Code) != -1 {
			return wrappedErr.AddCategories(ErrTypeDisk, ErrTypeDatabase)
		} else if sqliteErr.Code == sqlite3.ErrNomem {
			return wrappedErr.AddCategories(ErrTypeMemory, ErrTypeDatabase)
		}
	}

	return wrappedErr.AddCategories(ErrTypeOther, ErrTypeDatabase)
})

var ErrNoTxInContext = NewErrorWithCategories("no db transaction found in context")

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

type Error struct {
	Err                   error
	Categories            []string
	ErrDuplicatesCategory bool
}

func (err *Error) Error() string {
	message := ""

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
func (err *Error) StandardError() error {
	if err == nil {
		return nil
	}
	return err
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

func (err *Error) GeneralCategory() string {
	category, _ := GetLastCategoryWithTag(err.Categories, CategoryTagGeneral)
	return category
}

// requiredCategories is highest to lowest level e.g "auth [package]", "create user", common.ErrTypeDatabase
func (err *Error) HasCategories(requiredCategories ...string) bool {
	reversedCategories := slices.Clone(err.Categories)
	slices.Reverse(reversedCategories) // Highest to lowest level
	return CheckPathPattern(reversedCategories, slices.Concat(requiredCategories, []string{"***"}))
}
func (err *Error) HighestCategory() string {
	return err.Categories[len(err.Categories)-1]
}
func (err *Error) LowestCategory() string {
	return err.Categories[0]
}
func (err *Error) AddCategories(categories ...string) *Error {
	copiedErr := err.Clone()
	for _, category := range categories {
		hasCategoryTag := slices.Contains(ParseCategoryTags(category), CategoryTagPackage)
		insertIndex := -1
		if !hasCategoryTag {
			_, insertIndex = GetLastCategoryWithTag(err.Categories, CategoryTagPackage)
		}

		if insertIndex == -1 {
			copiedErr.Categories = append(copiedErr.Categories, category)
		} else {
			copiedErr.Categories = slices.Insert(copiedErr.Categories, insertIndex, category)
		}
	}
	return copiedErr
}
func (err *Error) AddCategory(category string) *Error {
	return err.AddCategories(category)
}
func (err *Error) RemoveHighestCategory() *Error {
	// TODO: should highest category include packages for this and HighestCategory()?
	copiedErr := err.Clone()
	catCount := len(copiedErr.Categories)
	if catCount != 0 {
		copiedErr.Categories = slices.Delete(copiedErr.Categories, catCount-1, catCount)
	}

	return copiedErr
}
func (err *Error) RemoveLowestCategory() *Error {
	copiedErr := err.Clone()
	if len(copiedErr.Categories) != 0 {
		copiedErr.Categories = slices.Delete(copiedErr.Categories, 0, 1)
	}

	return copiedErr
}

func (err Error) Clone() *Error {
	copiedErr := err
	copiedErr.Categories = slices.Clone(err.Categories)

	return &copiedErr
}

// categories is lowest to highest level, e.g. "constraint", common.ErrTypeDatabase, "create profile", "create user", "auth [package]"
func NewErrorWithCategories(message string, categories ...string) *Error {
	return &Error{
		Err:                   errors.New(message),
		Categories:            slices.Concat([]string{message}, categories),
		ErrDuplicatesCategory: true,
	}
}

// categories is lowest to highest level, e.g. "constraint", common.ErrTypeDatabase, "create profile", "create user", "auth [package]"
func WrapErrorWithCategories(err error, categories ...string) *Error {
	commErr, ok := err.(*Error)
	if ok {
		return commErr.AddCategories(categories...)
	}

	return &Error{
		Err:        err,
		Categories: categories,
	}
}

const (
	CategoryTagGeneral = "general"
	CategoryTagPackage = "package"
)

func ParseCategoryTags(category string) []string {
	rawTags := strings.Split(GetStringBetween(category, "[", "]"), ",")

	knownTags := []string{CategoryTagGeneral, CategoryTagPackage}
	tags := []string{}
	for _, rawTag := range rawTags {
		tag := strings.Trim(rawTag, " ")
		if tag == "" {
			continue
		}

		if !slices.Contains(knownTags, tag) {
			panic(fmt.Sprintf("ParseCategoryTags: %v is not a valid tag. category string:\n%v", tag, category))
		}
		tags = append(tags, tag)
	}
	return tags
}
func GetCategoryType(categoryTags []string) string {
	knownCategories := []string{CategoryTagGeneral, CategoryTagPackage}

	for _, tag := range categoryTags {
		if slices.Contains(knownCategories, tag) {
			return tag
		}
	}
	return ""
}
func GetLastCategoryWithTag(categories []string, requiredTag string) (string, int) {
	for i, category := range slices.Backward(categories) {
		tags := ParseCategoryTags(category)
		if slices.Contains(tags, requiredTag) {
			return category, i
		}
	}
	return "", -1
}

func AutoWrapError(err error) *Error {
	var commErr *Error
	if errors.As(err, &commErr) {
		return commErr
	}

	commErr = WrapErrorWithCategories(err, "auto wrapped")
	if errors.As(err, &sqlite3.Error{}) {
		return ErrWrapperDatabase.Wrap(commErr)
	}
	if ent.IsConstraintError(err) ||
		ent.IsNotFound(err) ||
		ent.IsNotLoaded(err) ||
		ent.IsNotSingular(err) ||
		ent.IsValidationError(err) ||
		errors.Is(err, ent.ErrTxStarted) {
		return ErrWrapperDatabase.Wrap(commErr)
	}

	return commErr
}

type ErrorWrapper interface {
	Wrap(err error) *Error
}
type ConstantErrorWrapper struct {
	Categories []string
	Child      ErrorWrapper
}

func NewErrorWrapper(categories ...string) *ConstantErrorWrapper {
	return &ConstantErrorWrapper{
		Categories: categories,
	}
}

func (errWrapper *ConstantErrorWrapper) Wrap(err error) *Error {
	if errWrapper.Child != nil {
		return errWrapper.Child.Wrap(err).AddCategories(errWrapper.Categories...)
	} else {
		return WrapErrorWithCategories(err, errWrapper.Categories...)
	}
}
func (errWrapper *ConstantErrorWrapper) HasWrapped(err error) bool {
	var commErr *Error
	if !errors.As(err, &commErr) {
		return false
	}

	return CheckPathPattern(commErr.Categories, slices.Concat([]string{"***"}, errWrapper.Categories, []string{"***"}))
}

func (errWrapper *ConstantErrorWrapper) SetChild(child ErrorWrapper) *ConstantErrorWrapper {
	newErrWrapper := errWrapper.Clone()
	newErrWrapper.Child = child
	return newErrWrapper
}

// TODO: add some kind of AddCategory method?
func (errWrapper ConstantErrorWrapper) Clone() *ConstantErrorWrapper {
	copiedErrWrapper := errWrapper
	copiedErrWrapper.Categories = slices.Clone(copiedErrWrapper.Categories)

	return &copiedErrWrapper
}

type DynamicErrorWrapper struct {
	callback func(err error) *Error // Mostly just private so IDEs autocomplete to Wrap instead
}

func NewDynamicErrorWrapper(callback func(err error) *Error) *DynamicErrorWrapper {
	return &DynamicErrorWrapper{
		callback: callback,
	}
}
func (errWrapper *DynamicErrorWrapper) Wrap(err error) *Error {
	return errWrapper.callback(err)
}

// TODO: rename and probably rework
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
