package logger

import (
	"context"

	"go.uber.org/zap"
)

type CtxKey string

const CtxLoggerKey CtxKey = "logger"

func ContextWithLogger(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, CtxLoggerKey, logger)
}

func LoggerFromContext(ctx context.Context) *zap.Logger {
	logger, ok := ctx.Value(CtxLoggerKey).(*zap.Logger)
	if !ok {
		logger, _ = zap.NewProduction()
	}
	return logger
}
