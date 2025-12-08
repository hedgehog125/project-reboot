package ratelimiting

import "github.com/NicoClack/cryptic-stash/common"

const (
	ErrTypeRequestSession = "request session"
	ErrTypeAdjustTo       = "adjust to"
	// Lower level
	ErrTypeRateLimitExceeded = "rate limit exceeded"
)

var ErrWrapperRequestSession = common.NewErrorWrapper(common.ErrTypeRateLimiting, ErrTypeRequestSession)
var ErrWrapperAdjustTo = common.NewErrorWrapper(common.ErrTypeRateLimiting, ErrTypeAdjustTo)

var ErrWrapperRateLimitExceeded = common.NewErrorWrapper(common.ErrTypeRateLimiting, ErrTypeRateLimitExceeded)

var ErrInvalidEventName = common.NewErrorWithCategories("invalid event name", common.ErrTypeRateLimiting)
var ErrGlobalRateLimitExceeded = ErrWrapperRateLimitExceeded.Wrap(
	common.NewErrorWithCategories("global rate limit exceeded"),
)
var ErrUserRateLimitExceeded = ErrWrapperRateLimitExceeded.Wrap(
	common.NewErrorWithCategories("user rate limit exceeded"),
)
