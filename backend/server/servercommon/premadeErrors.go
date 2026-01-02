package servercommon

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/NicoClack/cryptic-stash/backend/common"
)

const (
	ErrTypeBadRequest = "bad request"
)

var ErrUnauthorized = NewError(common.NewErrorWithCategories(
	"unauthorized", common.ErrTypeServerCommon, common.ErrTypeClient,
)).SetStatus(http.StatusUnauthorized)
var ErrNotFound = NewError(common.NewErrorWithCategories(
	"not found", common.ErrTypeServerCommon, common.ErrTypeClient,
)).SetStatus(http.StatusNotFound).DisableLogging()

var ErrWrapperBadRequest = common.NewErrorWrapper(common.ErrTypeServerCommon, ErrTypeBadRequest, common.ErrTypeClient)

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
