package middleware

import (
	"net/http"
	"time"

	"github.com/gin-contrib/timeout"
	"github.com/gin-gonic/gin"
	"github.com/hedgehog125/project-reboot/server/servercommon"
)

func NewTimeout() gin.HandlerFunc {
	return timeout.New(
		timeout.WithTimeout(20*time.Second),
		timeout.WithResponse(func(ginCtx *gin.Context) {
			if ginCtx.Writer.Written() {
				conn, _, stdErr := ginCtx.Writer.Hijack()
				if stdErr != nil {
					_ = conn.Close()
				}
				return
			}
			ginCtx.AbortWithStatusJSON(http.StatusRequestTimeout, gin.H{
				"errors": []servercommon.ErrorDetail{
					{
						Message: "request timed out",
						Code:    "REQUEST_TIMEOUT",
					},
				},
			})
		}),
	)
}
