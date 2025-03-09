package common

import (
	"errors"

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

const (
	ErrorDatabase = "database"
	ErrorOther    = "other"
)

func CategorizeError(err error) string {
	if _, ok := err.(*sqlite3.Error); ok {
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
