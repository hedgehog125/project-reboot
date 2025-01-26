package endpoints

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/timeout"
	"github.com/gin-gonic/gin"
	"github.com/hedgehog125/project-reboot/core"
	"github.com/hedgehog125/project-reboot/intertypes"
)

func NewTimeoutMiddleware() gin.HandlerFunc {
	return timeout.New(
		timeout.WithTimeout(5*time.Second),
		timeout.WithHandler(func(ctx *gin.Context) {
			ctx.Next()
		}),
		timeout.WithResponse(func(ctx *gin.Context) {
			if ctx.Writer.Written() {
				conn, _, err := ctx.Writer.Hijack()
				if err != nil {
					_ = conn.Close()
				}
				return
			}
			ctx.JSON(http.StatusRequestTimeout, gin.H{
				"errors": []string{"REQUEST_TIMED_OUT"},
			})
		}),
	)
}

func NewAdminProtectedMiddleware(state *intertypes.State) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		headerValue := ctx.GetHeader("authorization")
		if headerValue == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"errors": []string{"MISSING_AUTHORIZATION_HEADER"},
			})
			return
		}
		headerParts := strings.SplitN(headerValue, " ", 2)

		if len(headerParts) != 2 {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"errors": []string{"MALFORMED_AUTHORIZATION_HEADER"},
			})
			return
		}
		if headerParts[0] != "AdminCode" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"errors": []string{"UNSUPPORTED_AUTHORIZATION_SCHEME"},
			})
			return
		}

		if core.CheckAdminCode(headerParts[1], state) {
			ctx.Next()
		} else {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"errors": []string{"INVALID_ADMIN_CODE"},
			})
			return
		}
	}
}
