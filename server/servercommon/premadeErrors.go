package servercommon

import (
	"fmt"
	"net/http"

	"github.com/hedgehog125/project-reboot/common"
)

func NewUnauthorizedError() *ContextError {
	err := common.NewErrorWithCategory("unauthorized", common.ErrTypeClient)
	return &ContextError{
		Err:        err,
		Status:     http.StatusUnauthorized,
		ErrorCodes: []string{},
		Category:   err.Category,
		ShouldLog:  true,
	}
}

func NewNotFoundError() *ContextError {
	err := common.NewErrorWithCategory("not found", common.ErrTypeClient)
	return &ContextError{
		Err:        err,
		Status:     http.StatusNotFound,
		ErrorCodes: []string{},
		Category:   err.Category,
		ShouldLog:  false,
	}
}

func NewBadRequestError(fieldName string, message string) *ContextError {
	err := common.NewErrorWithCategory(fmt.Sprintf("%v: %v", fieldName, message), common.ErrTypeClient)
	return &ContextError{
		Err:        err,
		Status:     http.StatusBadRequest,
		ErrorCodes: []string{}, // TODO: add error code?
		Category:   err.Category,
		ShouldLog:  false,
	}
}
