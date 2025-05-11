package middleware

// TODO: return ContextErrors

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/timeout"
	"github.com/gin-gonic/gin"
	"github.com/hedgehog125/project-reboot/common"
	"github.com/hedgehog125/project-reboot/core"
	"github.com/hedgehog125/project-reboot/server/servercommon"
)

// TODO: handle panics and stop the default error handler from being registered
func NewErrorMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()

		for _, ginError := range ctx.Errors {
			err := &servercommon.ContextError{}
			if !errors.As(ginError.Err, &err) {
				err = servercommon.NewContextError(ginError.Err)
			}
			err.Finish()

			if err.ErrorCodes == nil {
				if err.Status != -1 {
					ctx.Status(err.Status)
				}
			} else {
				ctx.JSON(err.Status, gin.H{
					"errors": err.ErrorCodes,
				})
			}

			common.DumpJSON(err)
			fmt.Printf("request error:\n%v\n\n", err.Error())

			// TODO:
		}
	}
}

func NewTimeoutMiddleware() gin.HandlerFunc {
	return timeout.New(
		timeout.WithTimeout(30*time.Second),
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
			ctx.AbortWithStatusJSON(http.StatusRequestTimeout, gin.H{
				"errors": []string{"REQUEST_TIMED_OUT"},
			})
		}),
	)
}

func NewAdminProtectedMiddleware(state *common.State) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		headerValue := ctx.GetHeader("authorization")
		if headerValue == "" {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"errors": []string{"MISSING_AUTHORIZATION_HEADER"},
			})
			return
		}
		headerParts := strings.SplitN(headerValue, " ", 2)

		if len(headerParts) != 2 {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"errors": []string{"MALFORMED_AUTHORIZATION_HEADER"},
			})
			return
		}
		if headerParts[0] != "AdminCode" {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"errors": []string{"UNSUPPORTED_AUTHORIZATION_SCHEME"},
			})
			return
		}

		if core.CheckAdminCode(headerParts[1], state) {
			ctx.Next()
		} else {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"errors": []string{"INVALID_ADMIN_CODE"},
			})
			return
		}
	}
}
