package common

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"time"
)

const BackoffJitter = float64(0.05)

// TODO: enforce a timeout or log a warning if it's exceeded? Some contexts don't have a deadline but instead can just be cancelled after a while
func WithRetries(
	ctx context.Context, fn func() error,
) error {
	maxObservedRunTime := time.Duration(0)
	retriedFraction := float64(0) // When >= 1, max retries is reached
	retriesByCategory := map[string]int{}
	errs := []error{}
	getDebugValue := func() DebugValue {
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

	for {
		startTime := time.Now()
		stdErr := fn()
		if stdErr == nil {
			return nil
		}
		commErr := AutoWrapError(stdErr)
		if commErr.MaxRetries > 0 {
			retriedFraction += 1 / float64(commErr.MaxRetries)
		}
		if retriedFraction >= 1 || (commErr.MaxRetries < 1 && commErr.MaxRetries != -1) {
			return commErr.AddDebugValue(getDebugValue())
		}
		errs = append(errs, stdErr)

		retries := retriesByCategory[commErr.GeneralCategory()]
		jitterMultiplier := ((rand.Float64() * BackoffJitter * 2) - BackoffJitter) + 1
		backoff := time.Duration(math.Round(
			float64(commErr.RetryBackoffBase) *
				math.Pow(commErr.RetryBackoffMultiplier, float64(retries+1)) *
				jitterMultiplier,
		))
		if commErr.RetryBackoffMultiplier > 1 { // Errors with a multiplier of 1 shouldn't increase the backoff for other errors with the same category
			retriesByCategory[commErr.GeneralCategory()] = retries + 1
		}

		runTime := time.Since(startTime)
		maxObservedRunTime = max(maxObservedRunTime, runTime)
		deadline, hasDeadline := ctx.Deadline()
		if hasDeadline && time.Until(deadline) < maxObservedRunTime+backoff {
			return commErr.AddDebugValue(getDebugValue())
		}

		fmt.Printf("waiting %vms\n", backoff.Milliseconds())

		select {
		case <-time.After(backoff):
		case <-ctx.Done():
			return WrapErrorWithCategories(
				context.Canceled,
				ErrTypeTimeout, "with retries", ErrTypeCommon,
			).AddDebugValue(getDebugValue())
		}
	}
}
