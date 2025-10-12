package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hedgehog125/project-reboot/common"
	"github.com/hedgehog125/project-reboot/ratelimiting"
)

func NewRateLimiting(eventName string, limiter common.LimiterService) gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		_, commErr := limiter.RequestSession(eventName, 1, ginCtx.ClientIP())
		if commErr != nil {
			if errors.Is(commErr, ratelimiting.ErrGlobalRateLimitExceeded) ||
				errors.Is(commErr, ratelimiting.ErrUserRateLimitExceeded) {
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
