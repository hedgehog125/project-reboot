package servercommon

import (
	"fmt"
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/hedgehog125/project-reboot/common"
)

const (
	LoggerKey = "logger"
)

func GetLogger(ginCtx *gin.Context) common.Logger {
	logger, ok := ginCtx.Get(LoggerKey)
	if ok {
		logger, ok := logger.(common.Logger)
		if ok {
			return logger
		}
	}

	message := "used default logger as no logger was found in Gin context. this shouldn't happen!"
	defaultLogger := slog.Default()
	defaultLogger.Warn(message)
	return defaultLogger.With(
		"loggerError", fmt.Sprintf("warning: %v", message),
	)
}
