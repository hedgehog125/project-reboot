package servercommon

import (
	"context"
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

type Error struct {
	child     common.WrappedError
	status    int // Set to -1 to keep the current code. TODO: change to 0?
	details   []ErrorDetail
	shouldLog bool
}
type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func NewError(stdErr error) *Error {
	if stdErr == nil {
		return nil
	}
	serverErr := &Error{}
	if errors.As(stdErr, &serverErr) {
		return serverErr.Clone()
	}

	wrappedErr := common.AutoWrapError(stdErr)
	if wrappedErr == nil {
		return nil
	}

	serverErr = &Error{
		child:     wrappedErr,
		status:    -1,
		details:   []ErrorDetail{},
		shouldLog: true,
	}
	if errors.Is(stdErr, context.DeadlineExceeded) {
		serverErr.status = 408
	}
	return serverErr
}

func (err *Error) Error() string {
	return err.child.Error()
}
func (err *Error) StandardError() error {
	if err == nil {
		return nil
	}
	return err
}
func (err *Error) Unwrap() error {
	return err.child
}

func (err *Error) Clone() *Error {
	if err == nil {
		return nil
	}
	return &Error{
		child:     err.child.CloneAsWrappedError(),
		status:    err.status,
		details:   slices.Clone(err.details),
		shouldLog: err.shouldLog,
	}
}
func (err *Error) SetChild(wrappedErr common.WrappedError) *Error {
	if wrappedErr == nil {
		return nil
	}
	copiedErr := err.Clone()
	copiedErr.child = wrappedErr
	return copiedErr
}

func (err *Error) ConfigureRetries(maxRetries int, baseBackoff time.Duration, backoffMultiplier float64) *Error {
	newChild := err.child.CloneAsWrappedError()
	newChild.ConfigureRetriesMut(maxRetries, baseBackoff, backoffMultiplier)
	return err.SetChild(newChild)
}

func (err *Error) AddDetails(details ...ErrorDetail) *Error {
	copiedErr := err.Clone()
	copiedErr.details = append(copiedErr.details, details...)
	return copiedErr
}
func (err *Error) AddDetail(detail ErrorDetail) *Error {
	return err.AddDetails(detail)
}
func (err *Error) SetStatus(code int) *Error {
	copiedErr := err.Clone()
	copiedErr.status = code
	return copiedErr
}
func (err *Error) SetShouldLog(shouldLog bool) *Error {
	copiedErr := err.Clone()
	copiedErr.shouldLog = shouldLog
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
	return err.sendStatusAndDetailIfCondition(ent.IsNotFound(err.child.Unwrap()), statusCode, detail, preventLog)
}

func (err *Error) sendStatusAndDetailIfCondition(
	condition bool, statusCode int,
	detail *ErrorDetail, preventLog bool,
) *Error {
	copiedErr := err.Clone()
	if condition {
		copiedErr = copiedErr.ConfigureRetries(0, 0, 0)
		if statusCode != -1 {
			copiedErr.status = statusCode
		}
		if detail != nil {
			copiedErr.details = append(copiedErr.details, *detail)
		}
		if preventLog {
			copiedErr.shouldLog = false
		}
	}
	return copiedErr
}

func (err *Error) Status() int {
	return err.status
}
func (err *Error) Details() []ErrorDetail {
	return err.details
}
func (err *Error) ShouldLog() bool {
	return err.shouldLog
}
