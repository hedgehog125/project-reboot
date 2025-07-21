package middleware

// TODO: return ContextErrors

import (
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

		statusCode := -1
		mergedDetails := []servercommon.ErrorDetail{}
		for _, ginError := range ctx.Errors {
			serverErr := servercommon.NewError(ginError.Err)
			if serverErr.Status != -1 {
				if statusCode == -1 {
					statusCode = serverErr.Status
				} else {
					statusCode = http.StatusInternalServerError
					fmt.Printf(
						"warning: API errors have different status codes: %d and %d\n",
						statusCode, serverErr.Status,
					)
				}
			}
			mergedDetails = append(mergedDetails, serverErr.Details...)

			// TODO: use slog
			common.DumpJSON(serverErr)
			fmt.Printf("request error:\n%v\n\n", serverErr.Error())
		}

		if len(ctx.Errors) != 0 {
			if statusCode == -1 {
				statusCode = http.StatusInternalServerError
			}
			if !ctx.Writer.Written() {
				ctx.JSON(statusCode, gin.H{
					"errors": mergedDetails,
				})
			}
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

func NewAdminProtectedMiddleware(state *common.State) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		headerValue := ctx.GetHeader("authorization")
		if headerValue == "" {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"errors": []servercommon.ErrorDetail{
					{
						Message: "authorization header is required",
						Code:    "MISSING_AUTHORIZATION_HEADER",
					},
				},
			})
			return
		}
		headerParts := strings.SplitN(headerValue, " ", 2)
		if len(headerParts) != 2 {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"errors": []servercommon.ErrorDetail{
					{
						Message: "malformed authorization header",
						Code:    "MALFORMED_AUTHORIZATION_HEADER",
					},
				},
			})
			return
		}
		if headerParts[0] != "AdminCode" {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"errors": []servercommon.ErrorDetail{
					{
						Message: "unsupported authorization scheme",
						Code:    "UNSUPPORTED_AUTHORIZATION_SCHEME",
					},
				},
			})
			return
		}

		if core.CheckAdminCode(headerParts[1], state) {
			ctx.Next()
		} else {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"errors": []servercommon.ErrorDetail{
					{
						Message: "invalid admin code",
						Code:    "INVALID_ADMIN_CODE",
					},
				},
			})
			return
		}
	}
}
