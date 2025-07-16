package jobs

import (
	"errors"
	"slices"
	"time"

	"github.com/hedgehog125/project-reboot/common"
)

const (
	ErrTypeEncode  = "encode"
	ErrTypeDecode  = "decode" // From Job.Decode() method
	ErrTypeEnqueue = "enqueue"
	// Lower level
	ErrTypeInvalidData = "invalid data"
)

var ErrUnknownJobType = common.NewErrorWithCategories(
	"unknown job type", common.ErrTypeJobs,
)

var ErrWrapperDecode = common.NewErrorWrapper(
	ErrTypeDecode, common.ErrTypeJobs,
)

// TODO: test this
var ErrWrapperDatabase = common.NewErrorWrapper(common.ErrTypeTwoFactorAction).
	SetChild(common.ErrWrapperDatabase)
var ErrWrapperInvalidData = common.NewErrorWrapper(
	ErrTypeInvalidData, common.ErrTypeTwoFactorAction,
)

type CommonError = common.Error
type Error struct {
	CommonError
	RetryBackoffs []time.Duration
}
type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func NewError(err error) *Error {
	serverErr := &Error{}
	if errors.As(err, &serverErr) {
		return serverErr.Clone()
	}

	commErr := &common.Error{}
	if !errors.As(err, &commErr) {
		commErr = common.AutoWrapError(err)
	}
	return &Error{
		CommonError: *commErr,
		Status:      -1,
		Details:     []ErrorDetail{},
		ShouldLog:   true,
	}
}

func (err *Error) Unwrap() error {
	return err.StandardError()
}
func (err *Error) Clone() *Error {
	copiedErr := &Error{
		CommonError: *err.CommonError.Clone(),
		Status:      err.Status,
		Details:     slices.Clone(err.Details),
		ShouldLog:   err.ShouldLog,
	}
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
