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
	return func(ginCtx *gin.Context) {
		ginCtx.Next()

		statusCode := -1
		mergedDetails := []servercommon.ErrorDetail{}
		for _, ginError := range ginCtx.Errors {
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

		if len(ginCtx.Errors) != 0 {
			if statusCode == -1 {
				statusCode = http.StatusInternalServerError
			}
			if !ginCtx.Writer.Written() {
				ginCtx.JSON(statusCode, gin.H{
					"errors": mergedDetails,
				})
			}
		}
	}
}

func NewTimeoutMiddleware() gin.HandlerFunc {
	return timeout.New(
		timeout.WithTimeout(30*time.Second),
		timeout.WithHandler(func(ginCtx *gin.Context) {
			ginCtx.Next()
		}),
		timeout.WithResponse(func(ginCtx *gin.Context) {
			if ginCtx.Writer.Written() {
				conn, _, err := ginCtx.Writer.Hijack()
				if err != nil {
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

func NewAdminProtectedMiddleware(state *common.State) gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		headerValue := ginCtx.GetHeader("authorization")
		if headerValue == "" {
			ginCtx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
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
			ginCtx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
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
			ginCtx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
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
			ginCtx.Next()
		} else {
			ginCtx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
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
