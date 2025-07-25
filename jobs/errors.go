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
	ErrTypeRunJob  = "run job"
	// Lower level
	ErrTypeInvalidData = "invalid data"
)

var ErrUnknownJobType = common.NewErrorWithCategories(
	"unknown job type", common.ErrTypeJobs,
)
var ErrNoTxInContext = common.ErrNoTxInContext.AddCategory(common.ErrTypeJobs)

var ErrWrapperDecode = common.NewErrorWrapper(
	ErrTypeDecode, common.ErrTypeJobs,
)

// TODO: test this
var ErrWrapperDatabase = common.NewErrorWrapper(common.ErrTypeJobs).
	SetChild(common.ErrWrapperDatabase)
var ErrWrapperInvalidData = common.NewErrorWrapper(
	ErrTypeInvalidData, common.ErrTypeJobs,
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
	jobErr := &Error{}
	if errors.As(err, &jobErr) {
		return jobErr.Clone()
	}

	commErr := &common.Error{}
	if !errors.As(err, &commErr) {
		commErr = common.AutoWrapError(err)
	}
	return &Error{
		CommonError:   *commErr,
		RetryBackoffs: []time.Duration{},
	}
}

func (err *Error) Unwrap() error {
	return err.StandardError()
}
func (err *Error) Clone() *Error {
	copiedErr := &Error{
		CommonError:   *err.CommonError.Clone(),
		RetryBackoffs: slices.Clone(err.RetryBackoffs),
	}
	return copiedErr
}

func (err *Error) AddCategory(category string) *Error {
	copiedErr := err.Clone()
	copiedErr.Err = err.CommonError.AddCategory(category)
	return copiedErr
}
