package ctx

import (
	"companies/pkg/logger"
	"context"

	"go.uber.org/zap"
	"log"
)

type key string

const (
	loggerKey    key = "logger"
	requestIdKey key = "xRequestId"
)

func withValue(ctx context.Context, key key, val interface{}) context.Context {
	return context.WithValue(ctx, key, val)
}

func value(ctx context.Context, key key) interface{} {
	return ctx.Value(key)
}

func WithLogger(ctx context.Context, logger *logger.Logger) context.Context {
	return withValue(ctx, loggerKey, logger)
}

func Logger(ctx context.Context) *logger.Logger {
	if logger, ok := value(ctx, loggerKey).(*logger.Logger); ok {
		return logger
	}
	log.Printf("no logger in ctx, no-op provided")
	return &logger.Logger{SugaredLogger: zap.NewNop().Sugar()}
}

func WithRequestId(ctx context.Context, requestId string) context.Context {
	return withValue(ctx, requestIdKey, requestId)
}

func RequestId(ctx context.Context) string {
	if requestId, ok := value(ctx, requestIdKey).(string); ok {
		return requestId
	}

	return ""
}
