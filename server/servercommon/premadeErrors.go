package servercommon

import (
	"net/http"

	"github.com/hedgehog125/project-reboot/common"
)

func NewUnauthorizedError() *ContextError {
	return &ContextError{
		Err:        nil,
		Status:     http.StatusUnauthorized,
		ErrorCodes: []string{},
		Category:   common.ErrorUnauthorized,
		ShouldLog:  true,
	}
}

func NewNotFoundError() *ContextError {
	return &ContextError{
		Err:        nil,
		Status:     http.StatusNotFound,
		ErrorCodes: []string{},
		Category:   common.ErrorNotFound,
		ShouldLog:  false,
	}
}
