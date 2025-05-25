package servercommon

import (
	"fmt"
	"net/http"

	"github.com/hedgehog125/project-reboot/common"
)

const (
	ErrTypeBadRequest = "bad request"
)

var ErrUnauthorized = common.NewErrorWithCategories("unauthorized", common.ErrTypeClient)
var ErrNotFound = common.NewErrorWithCategories("not found", common.ErrTypeClient)

func NewUnauthorizedError() *ContextError {
	return &ContextError{
		Err:        ErrUnauthorized,
		Status:     http.StatusUnauthorized,
		ErrorCodes: []string{},
		Category:   ErrUnauthorized.GeneralCategory(),
		ShouldLog:  true,
	}
}

func NewNotFoundError() *ContextError {
	return &ContextError{
		Err:        ErrNotFound,
		Status:     http.StatusNotFound,
		ErrorCodes: []string{},
		Category:   ErrNotFound.GeneralCategory(),
		ShouldLog:  false,
	}
}

func NewBadRequestError(fieldName string, message string) *ContextError {
	err := common.NewErrorWithCategories(fmt.Sprintf("%v: %v", fieldName, message), common.ErrTypeClient, ErrTypeBadRequest)
	return &ContextError{
		Err:        err,
		Status:     http.StatusBadRequest,
		ErrorCodes: []string{}, // TODO: add error code?
		Category:   err.GeneralCategory(),
		ShouldLog:  false,
	}
}
