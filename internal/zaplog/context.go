package zaplog

import (
	"context"

	"go.uber.org/zap"
)

type contextKey string

const (
	loggerKey    contextKey = "zaplog.logger"
	requestIDKey contextKey = "zaplog.request_id"
)

func WithLogger(ctx context.Context, logger *zap.Logger) context.Context {
	if logger == nil {
		logger = zap.NewNop()
	}
	return context.WithValue(ctx, loggerKey, logger)
}

func FromContext(ctx context.Context) *zap.Logger {
	logger, _ := ctx.Value(loggerKey).(*zap.Logger)
	if logger == nil {
		logger = zap.NewNop()
	}
	if requestID, ok := RequestIDFromContext(ctx); ok {
		return logger.With(zap.String("request_id", requestID))
	}
	return logger
}

func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, requestIDKey, requestID)
}

func RequestIDFromContext(ctx context.Context) (string, bool) {
	requestID, ok := ctx.Value(requestIDKey).(string)
	return requestID, ok && requestID != ""
}

func EnsureRequestID(ctx context.Context) (context.Context, string) {
	if requestID, ok := RequestIDFromContext(ctx); ok {
		return ctx, requestID
	}
	requestID := NewRequestID()
	return WithRequestID(ctx, requestID), requestID
}
