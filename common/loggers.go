package common

import (
	"context"
	"log/slog"
)

func GetLogger(ctx context.Context, service HasDefaultLogger) Logger {
	if ctx != nil {
		logger, ok := ctx.Value(LoggerKey{}).(LoggerService)
		if ok {
			return logger
		}
	}
	if service != nil {
		return service.DefaultLogger()
	}
	return slog.Default()
}
