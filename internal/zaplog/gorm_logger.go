package zaplog

import (
	"context"
	"errors"
	"strings"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

type GormLogger struct {
	logger        *zap.Logger
	logLevel      gormlogger.LogLevel
	slowThreshold time.Duration
	traceEnabled  bool
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
		traceEnabled:  ShouldEnableRepositoryTrace(level, debug),
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
	if l.logLevel < gormlogger.Info || !l.traceEnabled {
		return
	}
	L(ctx).Named("gorm").Debug(msg, zap.Any("data", data))
}

func (l *GormLogger) Warn(ctx context.Context, msg string, data ...any) {
	if l.logLevel < gormlogger.Warn {
		return
	}
	L(ctx).Named("gorm").Warn(msg, zap.Any("data", data))
}

func (l *GormLogger) Error(ctx context.Context, msg string, data ...any) {
	if l.logLevel < gormlogger.Error {
		return
	}
	L(ctx).Named("gorm").Error(msg, zap.Any("data", data))
}

func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.logLevel == gormlogger.Silent || !l.traceEnabled {
		return
	}

	elapsed := time.Since(begin)
	fields := []zap.Field{zap.Duration("duration", elapsed)}

	logSQL := true
	if logSQL {
		sql, rows := fc()
		op, table := summarizeSQL(sql)
		fields = append(fields,
			zap.String("sql", sql),
			zap.Int64("rows", rows),
			zap.String("op", op),
			zap.String("table", table),
		)
	}

	switch {
	case err != nil && l.logLevel >= gormlogger.Error && !errors.Is(err, gorm.ErrRecordNotFound):
		L(ctx).Named("gorm").Error("database sql failed", append(fields, zap.Error(err))...)
	case l.slowThreshold > 0 && elapsed > l.slowThreshold && l.logLevel >= gormlogger.Warn:
		L(ctx).Named("gorm").Warn("database sql slow", fields...)
	case l.traceEnabled && l.logLevel >= gormlogger.Info:
		L(ctx).Named("gorm").Debug("database sql", fields...)
	}
}

func summarizeSQL(sql string) (op string, table string) {
	trimmed := strings.TrimSpace(sql)
	if trimmed == "" {
		return "UNKNOWN", ""
	}
	upper := strings.ToUpper(trimmed)
	switch {
	case strings.HasPrefix(upper, "SELECT"):
		return "SELECT", extractTableAfterKeyword(trimmed, "FROM")
	case strings.HasPrefix(upper, "INSERT"):
		return "INSERT", extractTableAfterKeyword(trimmed, "INTO")
	case strings.HasPrefix(upper, "UPDATE"):
		return "UPDATE", extractTableAfterKeyword(trimmed, "UPDATE")
	case strings.HasPrefix(upper, "DELETE"):
		return "DELETE", extractTableAfterKeyword(trimmed, "FROM")
	default:
		return "OTHER", ""
	}
}

func extractTableAfterKeyword(sql string, keyword string) string {
	upper := strings.ToUpper(sql)
	idx := strings.Index(upper, keyword)
	if idx < 0 {
		return ""
	}
	rest := strings.TrimSpace(sql[idx+len(keyword):])
	if rest == "" {
		return ""
	}
	token := strings.Fields(rest)[0]
	return strings.Trim(token, "`\"")
}
