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
	// Skip Entry.{Debug,Info,Warn,Error} so caller/file keys point at real call sites, not api.go.
	base := FromContext(ctx)
	return Entry{
		logger:       base.WithOptions(zap.AddCallerSkip(1)),
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

