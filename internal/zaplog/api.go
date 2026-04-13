package zaplog

import (
	"context"

	"go.uber.org/zap"
)

type Entry struct {
	logger       *zap.Logger
	debugEnabled bool
}

func L(ctx context.Context) Entry {
	return Entry{
		logger:       FromContext(ctx),
		debugEnabled: DebugEnabled(ctx),
	}
}

func (e Entry) Named(name string) Entry {
	return Entry{
		logger:       e.logger.Named(name),
		debugEnabled: e.debugEnabled,
	}
}

func (e Entry) Debug(msg string, fields ...zap.Field) {
	if !e.debugEnabled {
		return
	}
	e.logger.Debug(msg, fields...)
}

func (e Entry) Info(msg string, fields ...zap.Field) {
	e.logger.Info(msg, fields...)
}

func (e Entry) Warn(msg string, fields ...zap.Field) {
	e.logger.Warn(msg, fields...)
}

func (e Entry) Error(msg string, fields ...zap.Field) {
	e.logger.Error(msg, fields...)
}

func Logger(ctx context.Context) *zap.Logger {
	return FromContext(ctx)
}

func LoggerNamed(ctx context.Context, name string) *zap.Logger {
	return FromContext(ctx).Named(name)
}

