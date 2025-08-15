package servercommon

import (
	"errors"
	"slices"
	"time"

	"github.com/hedgehog125/project-reboot/common"
	"github.com/hedgehog125/project-reboot/common/dbcommon"
	"github.com/hedgehog125/project-reboot/ent"
)

const (
	ErrTypeParseBodyJson = "parse body json"
)

var ErrCancelTransaction = NewError(dbcommon.ErrCancelTransaction).DisableLogging()

var ErrWrapperParseBodyJson = common.NewErrorWrapper(
	common.ErrTypeServerCommon,
	ErrTypeParseBodyJson, common.ErrTypeClient,
)

type CommonError = common.Error
type Error struct {
	*CommonError
	Status    int // Set to -1 to keep the current code
	Details   []ErrorDetail
	ShouldLog bool
}
type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func NewError(err error) *Error {
	if err == nil {
		return nil
	}
	serverErr := &Error{}
	if errors.As(err, &serverErr) {
		return serverErr.Clone()
	}

	commErr := common.AutoWrapError(err)
	if commErr == nil {
		return nil
	}
	return &Error{
		CommonError: commErr,
		Status:      -1,
		Details:     []ErrorDetail{},
		ShouldLog:   true,
	}
}

func (err *Error) StandardError() error {
	if err == nil {
		return nil
	}
	return err
}
func (err *Error) Unwrap() error {
	return err.CommonError
}

func (err *Error) Clone() *Error {
	if err == nil {
		return nil
	}
	return &Error{
		CommonError: err.CommonError.Clone(),
		Status:      err.Status,
		Details:     slices.Clone(err.Details),
		ShouldLog:   err.ShouldLog,
	}
}
func (err *Error) ConfigureRetries(maxRetries int, baseBackoff time.Duration, backoffMultiplier float64) *Error {
	copiedErr := err.Clone()
	copiedErr.Err = err.CommonError.ConfigureRetries(
		maxRetries, baseBackoff, backoffMultiplier,
	)
	return copiedErr
}

func (err *Error) AddDetails(details ...ErrorDetail) *Error {
	copiedErr := err.Clone()
	copiedErr.Details = append(copiedErr.Details, details...)
	return copiedErr
}
func (err *Error) AddDetail(detail ErrorDetail) *Error {
	return err.AddDetails(detail)
}
func (err *Error) SetStatus(code int) *Error {
	copiedErr := err.Clone()
	copiedErr.Status = code
	return copiedErr
}
func (err *Error) SetShouldLog(shouldLog bool) *Error {
	copiedErr := err.Clone()
	copiedErr.ShouldLog = shouldLog
	return copiedErr
}
func (err *Error) DisableLogging() *Error {
	return err.SetShouldLog(false)
}
func (err *Error) EnableLogging() *Error {
	return err.SetShouldLog(true)
}

func (err *Error) Send404IfNotFound() *Error {
	return err.sendStatusIfNotFound(404, nil, true)
}
func (err *Error) SendUnauthorizedIfNotFound() *Error {
	return err.sendStatusIfNotFound(401, nil, false)
}
func (err *Error) Expect(
	expectedError error,
	statusCode int, detail *ErrorDetail,
) *Error {
	return err.sendStatusAndDetailIfCondition(errors.Is(err, expectedError), statusCode, detail, true)
}
func (err *Error) ExpectAnyOf(
	expectedErrors []error,
	statusCode int, detail *ErrorDetail,
) *Error {
	isExpected := false
	for _, expectedErr := range expectedErrors {
		if errors.Is(err, expectedErr) {
			isExpected = true
			break
		}
	}
	return err.sendStatusAndDetailIfCondition(isExpected, statusCode, detail, true)
}
func (err *Error) sendStatusIfNotFound(
	statusCode int, detail *ErrorDetail,
	preventLog bool,
) *Error {
	return err.sendStatusAndDetailIfCondition(ent.IsNotFound(err.Err), statusCode, detail, preventLog)
}

func (err *Error) sendStatusAndDetailIfCondition(
	condition bool, statusCode int,
	detail *ErrorDetail, preventLog bool,
) *Error {
	copiedErr := err.Clone()
	if condition {
		copiedErr = copiedErr.ConfigureRetries(0, 0, 0)
		if statusCode != -1 {
			copiedErr.Status = statusCode
		}
		if detail != nil {
			copiedErr.Details = append(copiedErr.Details, *detail)
		}
		if preventLog {
			copiedErr.ShouldLog = false
		}
	}
	return copiedErr
}
