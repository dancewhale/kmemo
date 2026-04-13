package zaplog

import (
	"fmt"
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// New creates the project default zap logger (console only, stderr).
func New(level string) (*zap.Logger, error) {
	return NewForHost(level, false, "")
}

// NewForHost builds a logger with the same console encoding as [New]. When fileEnabled
// is true and logFilePath is non-empty, logs are also appended to that file (tee).
func NewForHost(level string, fileEnabled bool, logFilePath string) (*zap.Logger, error) {
	zapLevel, err := parseLevel(level)
	if err != nil {
		return nil, err
	}

	atom := zap.NewAtomicLevelAt(zapLevel)
	encCfg := productionConsoleEncoderConfig()

	cores := []zapcore.Core{
		zapcore.NewCore(zapcore.NewConsoleEncoder(encCfg), zapcore.AddSync(os.Stderr), atom),
	}
	if fileEnabled && logFilePath != "" {
		if err := os.MkdirAll(filepath.Dir(logFilePath), 0o755); err != nil {
			return nil, fmt.Errorf("create log directory: %w", err)
		}
		f, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
		if err != nil {
			return nil, fmt.Errorf("open log file: %w", err)
		}
		cores = append(cores, zapcore.NewCore(zapcore.NewConsoleEncoder(encCfg), zapcore.AddSync(f), atom))
	}

	core := zapcore.NewTee(cores...)
	return zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel)), nil
}

func productionConsoleEncoderConfig() zapcore.EncoderConfig {
	cfg := zap.NewProductionEncoderConfig()
	cfg.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.TimeKey = "ts"
	cfg.MessageKey = "msg"
	cfg.LevelKey = "level"
	return cfg
}

func parseLevel(level string) (zapcore.Level, error) {
	switch level {
	case "", "info":
		return zapcore.InfoLevel, nil
	case "debug":
		return zapcore.DebugLevel, nil
	case "warn":
		return zapcore.WarnLevel, nil
	case "error":
		return zapcore.ErrorLevel, nil
	default:
		return zapcore.InfoLevel, fmt.Errorf("zaplog: unsupported log level %q", level)
	}
}
