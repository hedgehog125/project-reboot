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

var ErrWrapperEncode = common.NewErrorWrapper(
	common.ErrTypeJobs, ErrTypeEncode,
)
var ErrWrapperDecode = common.NewErrorWrapper(
	common.ErrTypeJobs, ErrTypeDecode,
)

// TODO: test this
var ErrWrapperDatabase = common.NewErrorWrapper(common.ErrTypeJobs).
	SetChild(common.ErrWrapperDatabase)
var ErrWrapperInvalidData = common.NewErrorWrapper(
	common.ErrTypeJobs, ErrTypeInvalidData,
)

type CommonError = common.Error
type Error struct { // TODO: could the common error just be used instead?
	CommonError      // TODO: use pointer?
	JobRetryBackoffs []time.Duration
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
		CommonError:      *commErr,
		JobRetryBackoffs: []time.Duration{},
	}
}

func (err *Error) StandardError() error {
	if err == nil {
		return nil
	}
	return err
}
func (err *Error) Unwrap() error {
	return &err.CommonError
}
func (err *Error) Clone() *Error {
	copiedErr := &Error{
		CommonError:      *err.CommonError.Clone(),
		JobRetryBackoffs: slices.Clone(err.JobRetryBackoffs),
	}
	return copiedErr
}
func (err *Error) SetRetries(backoffs []time.Duration) *Error {
	copiedErr := err.Clone()
	copiedErr.JobRetryBackoffs = backoffs
	return copiedErr
}

func (err *Error) AddCategory(category string) *Error {
	copiedErr := err.Clone()
	copiedErr.Err = err.CommonError.AddCategory(category)
	return copiedErr
}
