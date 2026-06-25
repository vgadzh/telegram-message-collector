package logctx

import (
	"context"

	"go.uber.org/zap"
)

type loggerKey struct{}

var ctxLoggerKey = loggerKey{}

func WithLogger(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, ctxLoggerKey, logger)
}

func FromContext(ctx context.Context) *zap.Logger {
	logger, ok := ctx.Value(ctxLoggerKey).(*zap.Logger)
	if !ok {
		return zap.NewNop()
	}
	return logger
}
