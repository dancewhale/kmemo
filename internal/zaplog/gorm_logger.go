package zaplog

import (
	"context"
	"errors"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

type GormLogger struct {
	logger        *zap.Logger
	logLevel      gormlogger.LogLevel
	slowThreshold time.Duration
	debug         bool
}

func NewGormLogger(base *zap.Logger, level string, slowThreshold time.Duration, debug bool) gormlogger.Interface {
	if base == nil {
		base = zap.NewNop()
	}
	if slowThreshold <= 0 {
		slowThreshold = 200 * time.Millisecond
	}
	return &GormLogger{
		logger:        base.Named("gorm"),
		logLevel:      parseGormLevel(level),
		slowThreshold: slowThreshold,
		debug:         debug,
	}
}

func parseGormLevel(level string) gormlogger.LogLevel {
	switch level {
	case "silent":
		return gormlogger.Silent
	case "error":
		return gormlogger.Error
	case "warn", "":
		return gormlogger.Warn
	case "info", "debug":
		return gormlogger.Info
	default:
		return gormlogger.Warn
	}
}

func (l *GormLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	clone := *l
	clone.logLevel = level
	return &clone
}

func (l *GormLogger) Info(ctx context.Context, msg string, data ...any) {
	if l.logLevel < gormlogger.Info || !l.debug {
		return
	}
	FromContext(ctx).Named("gorm").Debug(msg, zap.Any("data", data))
}

func (l *GormLogger) Warn(ctx context.Context, msg string, data ...any) {
	if l.logLevel < gormlogger.Warn {
		return
	}
	FromContext(ctx).Named("gorm").Warn(msg, zap.Any("data", data))
}

func (l *GormLogger) Error(ctx context.Context, msg string, data ...any) {
	if l.logLevel < gormlogger.Error {
		return
	}
	FromContext(ctx).Named("gorm").Error(msg, zap.Any("data", data))
}

func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.logLevel == gormlogger.Silent {
		return
	}

	elapsed := time.Since(begin)
	logger := FromContext(ctx).Named("gorm")
	requestID, hasRequestID := RequestIDFromContext(ctx)
	fields := []zap.Field{zap.Duration("duration", elapsed)}
	if hasRequestID {
		fields = append(fields, zap.String("request_id", requestID))
	}

	logSQL := l.debug || err != nil || (l.slowThreshold > 0 && elapsed > l.slowThreshold)
	if logSQL {
		sql, rows := fc()
		fields = append(fields, zap.String("sql", sql), zap.Int64("rows", rows))
	}

	switch {
	case err != nil && l.logLevel >= gormlogger.Error && !errors.Is(err, gorm.ErrRecordNotFound):
		logger.Error("database query failed", append(fields, zap.Error(err))...)
	case l.slowThreshold > 0 && elapsed > l.slowThreshold && l.logLevel >= gormlogger.Warn:
		logger.Warn("slow query", fields...)
	case l.debug && l.logLevel >= gormlogger.Info:
		logger.Debug("database query", fields...)
	}
}
