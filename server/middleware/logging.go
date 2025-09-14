package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hedgehog125/project-reboot/common"
	"github.com/hedgehog125/project-reboot/server/servercommon"
)

func NewLogger(logger common.Logger) gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		ginCtx.Set(servercommon.LoggerKey, logger.With(
			// TODO: I think gin might already be providing this through the context?
			"url", ginCtx.Request.URL,
			"method", ginCtx.Request.Method,
		))

		ginCtx.Next()
	}
}

// TODO: handle panics and stop the default error handler from being registered
func NewError() gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		logger := servercommon.GetLogger(ginCtx)
		ginCtx.Next()

		statusCode := -1
		mergedDetails := []servercommon.ErrorDetail{}
		for _, ginError := range ginCtx.Errors {
			serverErr := servercommon.NewError(ginError.Err)
			if serverErr.Status != -1 {
				if statusCode == -1 {
					statusCode = serverErr.Status
				} else {
					logger.Warn(
						"server errors have different status codes",
						"previousStatusCode", statusCode,
						"newStatusCode", serverErr.Status,
						"errors", ginCtx.Errors,
					)
					statusCode = http.StatusInternalServerError
				}
			}
			mergedDetails = append(mergedDetails, serverErr.Details...)

			if serverErr.ShouldLog {
				if statusCode >= 500 {
					logger.Error("an internal server error occurred", "error", serverErr)
				} else {
					logger.Info("a HTTP 4xx was returned to a client", "error", serverErr)
				}
			}
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
