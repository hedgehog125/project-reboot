package jobs

import (
	"github.com/hedgehog125/project-reboot/common"
)

const (
	ErrTypeEncode  = "encode"
	ErrTypeDecode  = "decode" // From Job.Decode() method
	ErrTypeEnqueue = "enqueue"
	ErrTypeRunJob  = "run job"
	ErrTypeListen  = "listen"
	// Lower level
	ErrTypeInvalidBody = "invalid data"
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
var ErrWrapperEnqueue = common.NewErrorWrapper(common.ErrTypeJobs, ErrTypeEnqueue)
var ErrWrapperRunJob = common.NewErrorWrapper(
	common.ErrTypeJobs, ErrTypeRunJob,
)
var ErrWrapperListen = common.NewErrorWrapper(
	common.ErrTypeJobs, ErrTypeListen,
)

// TODO: test this
var ErrWrapperDatabase = common.NewErrorWrapper(common.ErrTypeJobs).
	SetChild(common.ErrWrapperDatabase)
var ErrWrapperInvalidBody = common.NewErrorWrapper(
	common.ErrTypeJobs, ErrTypeInvalidBody,
)
