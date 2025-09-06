package common

import (
	"context"
	"fmt"
	"math"
	"time"
)

const BackoffJitter = float64(0.05)
const MaxBackoffJitter = 500 * time.Millisecond
const BackoffMaxRetriesEpsilon = 1e-9

// TODO: enforce a timeout or log a warning if it's exceeded? Some contexts don't have a deadline but instead can just be cancelled after a while
func WithRetries(
	ctx context.Context, logger Logger, fn func() error,
) error {
	maxObservedRunTime := time.Duration(0)
	retriedFraction := float64(0) // When >= 1, max retries is reached
	retriesByCategory := map[string]int{}
	errs := []error{}
	getPreviousErrorsDebugValue := func() DebugValue {
		message := "no previous errors"
		if len(errs) > 0 {
			message = "from oldest to newest:"
			for _, prevErr := range errs {
				message += "\n" + prevErr.Error()
			}
		}
		return DebugValue{
			Name:    "previous retry errors (WithRetries)",
			Message: message,
			Value:   errs,
		}
	}
	wrapError := func(commErr *Error) *Error {
		return commErr.AddDebugValues(
			getPreviousErrorsDebugValue(),
			DebugValue{
				Name: "retries reset by WithRetries from...",
				Message: fmt.Sprintf(
					"max retries: %v, base backoff: %v, backoff multiplier: %v",
					commErr.MaxRetries, commErr.RetryBackoffBase, commErr.RetryBackoffMultiplier,
				),
			},
		).DisableRetries()
	}

	for {
		startTime := time.Now()
		stdErr := fn()
		if stdErr == nil {
			return nil
		}
		commErr := AutoWrapError(stdErr)
		if commErr.MaxRetries > 0 {
			retriedFraction += 1 / float64(commErr.MaxRetries+1)
		}
		if retriedFraction >= 1-BackoffMaxRetriesEpsilon || (commErr.MaxRetries < 1 && commErr.MaxRetries != -1) {
			return wrapError(commErr)
		}
		errs = append(errs, stdErr)

		retries := retriesByCategory[commErr.GeneralCategory()]
		backoff := CalculateBackoff(retries, commErr.RetryBackoffBase, commErr.RetryBackoffMultiplier)
		if commErr.RetryBackoffMultiplier > 1 { // Errors with a multiplier of 1 shouldn't increase the backoff for other errors with the same category
			retriesByCategory[commErr.GeneralCategory()] = retries + 1
		}

		runTime := time.Since(startTime)
		maxObservedRunTime = max(maxObservedRunTime, runTime)
		deadline, hasDeadline := ctx.Deadline()
		if hasDeadline && time.Until(deadline) < maxObservedRunTime+backoff {
			return wrapError(commErr)
		}

		logger.Debug("[WithRetries] waiting %vms", backoff.Milliseconds())

		select {
		case <-time.After(backoff):
		case <-ctx.Done():
			return WrapErrorWithCategories(
				context.Canceled,
				ErrTypeTimeout, "with retries", ErrTypeCommon,
			).AddDebugValue(getPreviousErrorsDebugValue())
		}
	}
}

func CalculateBackoff(retries int, base time.Duration, multiplier float64) time.Duration {
	withoutJitter := float64(base) * math.Pow(multiplier, float64(retries))
	jitter := RandPositiveNegativeRange(min(withoutJitter*BackoffJitter, float64(MaxBackoffJitter)))
	return time.Duration(math.Round(
		withoutJitter + jitter,
	))
}
