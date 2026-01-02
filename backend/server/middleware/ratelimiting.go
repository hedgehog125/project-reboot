package middleware

import (
	"errors"
	"net/http"

	"github.com/NicoClack/cryptic-stash/backend/common"
	"github.com/NicoClack/cryptic-stash/backend/ratelimiting"
	"github.com/gin-gonic/gin"
)

func NewRateLimiting(eventName string, limiter common.LimiterService) gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		_, wrappedErr := limiter.RequestSession(eventName, 1, ginCtx.ClientIP())
		if wrappedErr != nil {
			if errors.Is(wrappedErr, ratelimiting.ErrGlobalRateLimitExceeded) ||
				errors.Is(wrappedErr, ratelimiting.ErrUserRateLimitExceeded) {
				// TODO: add retry-after header
				ginCtx.Negotiate(http.StatusTooManyRequests, gin.Negotiate{
					Offered: []string{gin.MIMEJSON, gin.MIMEHTML},
					Data: gin.H{
						"errors": []struct{}{},
					},
					HTMLName: "429.html",
				})
				ginCtx.Abort()
				return
			}
		}
		ginCtx.Next()
	}
}
