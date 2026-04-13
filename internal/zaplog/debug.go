package zaplog

import (
	"context"
	"strings"

	"go.uber.org/zap/zapcore"
)

// IsDebugLevel returns true when runtime log level should include debug logs.
func IsDebugLevel(level string) bool {
	return strings.EqualFold(strings.TrimSpace(level), "debug")
}

// ShouldEnableRepositoryTrace gates SQL trace output.
// Requirement: only enable when log-level is debug AND repository_debug is true.
func ShouldEnableRepositoryTrace(level string, repositoryDebug bool) bool {
	return repositoryDebug && IsDebugLevel(level)
}

// DebugEnabled reports whether the logger in context emits debug.
func DebugEnabled(ctx context.Context) bool {
	return FromContext(ctx).Core().Enabled(zapcore.DebugLevel)
}

