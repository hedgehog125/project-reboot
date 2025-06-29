package servercommon

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/hedgehog125/project-reboot/common"
)

const (
	ErrTypeBadRequest = "bad request"
)

var ErrUnauthorized = NewError(common.NewErrorWithCategories(
	"unauthorized", common.ErrTypeClient, common.ErrTypeServerCommon,
)).SetStatus(http.StatusUnauthorized)
var ErrNotFound = NewError(common.NewErrorWithCategories(
	"not found", common.ErrTypeClient, common.ErrTypeServerCommon,
)).SetStatus(http.StatusNotFound).DisableLogging()
var ErrWrapperBadRequest = common.NewErrorWrapper(ErrTypeBadRequest, common.ErrTypeClient, common.ErrTypeServerCommon)

func NewUnauthorizedError() *Error {
	return ErrUnauthorized.Clone()
}
func NewNotFoundError() *Error {
	return ErrNotFound.Clone()
}

func NewBadRequestError(fieldName string, message string, errorCode string) *Error {
	fullMessage := fmt.Sprintf("%v: %v", fieldName, message)
	return NewError(ErrWrapperBadRequest.Wrap(errors.New(fullMessage))).
		SetStatus(http.StatusBadRequest).
		AddDetail(ErrorDetail{
			Message: fullMessage,
			Code:    errorCode,
		}).DisableLogging()
}
