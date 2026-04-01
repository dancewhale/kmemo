package bootstrap

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"kmemo/internal/app"
	"kmemo/internal/config"
	"kmemo/internal/pyclient"
	"kmemo/internal/zaplog"
)

// Headless bundles dependencies for CLI / non-UI entrypoints.
type Headless struct {
	Config config.Config
	Logger *zap.Logger
	Py     *pyclient.Client
}

// NewHeadless wires config and the Python gRPC client. No SQLite / indexing yet.
func NewHeadless(ctx context.Context) (*Headless, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}
	logger, err := zaplog.New(cfg.LogLevel)
	if err != nil {
		return nil, fmt.Errorf("create logger: %w", err)
	}

	ctx = zaplog.WithLogger(ctx, logger)
	ctx, requestID := zaplog.EnsureRequestID(ctx)
	logger = zaplog.FromContext(ctx).Named("bootstrap")
	logger.Info("bootstrap config loaded",
		zap.String("request_id", requestID),
		zap.String("python_grpc", cfg.PythonGRPCAddr),
		zap.Bool("skip_python", cfg.SkipPython),
		zap.Duration("dial_timeout", cfg.DialTimeout),
		zap.String("log_level", cfg.LogLevel),
		zap.Bool("repository_debug", cfg.RepositoryDebug),
		zap.Duration("db_slow_threshold", cfg.DBSlowThreshold),
	)

	var py *pyclient.Client
	if !cfg.SkipPython {
		py, err = pyclient.New(ctx, cfg)
		if err != nil {
			logger.Error("python client initialization failed", zap.Error(err))
			_ = logger.Sync()
			return nil, fmt.Errorf("pyclient: %w", err)
		}
	} else {
		logger.Warn("python client skipped")
	}
	return &Headless{Config: cfg, Logger: logger, Py: py}, nil
}

// NewDesktop builds the object graph for Wails bindings.
func NewDesktop(ctx context.Context) (*app.Desktop, error) {
	h, err := NewHeadless(ctx)
	if err != nil {
		return nil, err
	}
	return app.NewDesktop(h.Config, h.Logger, h.Py), nil
}
