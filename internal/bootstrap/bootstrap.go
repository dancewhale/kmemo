package bootstrap

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	grpcpython "kmemo/internal/adapters/fsrs/grpc_python"
	"kmemo/internal/adapters/fsrs/noop"
	"kmemo/internal/adapters/grpcworker"
	sourceprocessgrpc "kmemo/internal/adapters/sourceprocess/grpc_python"
	sourceprocessnoop "kmemo/internal/adapters/sourceprocess/noop"
	"kmemo/internal/app"
	"kmemo/internal/config"
	"kmemo/internal/contracts/fsrs"
	"kmemo/internal/contracts/sourceprocess"
	"kmemo/internal/zaplog"
)

// Headless bundles dependencies for CLI / non-UI entrypoints.
type Headless struct {
	Config        config.Config
	Logger        *zap.Logger
	Worker        *grpcworker.Client
	FSRS          fsrs.FSRSScheduler
	SourceProcess sourceprocess.Processor
}

// NewHeadless wires config, optional Python gRPC worker, and FSRS scheduler port.
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

	var worker *grpcworker.Client
	var sched fsrs.FSRSScheduler = &noop.Scheduler{}
	var processor sourceprocess.Processor = &sourceprocessnoop.Processor{}
	if !cfg.SkipPython {
		worker, err = grpcworker.New(ctx, cfg)
		if err != nil {
			logger.Error("grpc worker initialization failed", zap.Error(err))
			_ = logger.Sync()
			return nil, fmt.Errorf("grpc worker: %w", err)
		}
		sched = grpcpython.NewScheduler(worker.Processor())
		processor = sourceprocessgrpc.NewProcessor(worker.Processor())
	} else {
		logger.Warn("python grpc worker skipped")
	}
	return &Headless{Config: cfg, Logger: logger, Worker: worker, FSRS: sched, SourceProcess: processor}, nil
}

// NewDesktop builds the object graph for Wails bindings.
func NewDesktop(ctx context.Context) (*app.Desktop, error) {
	h, err := NewHeadless(ctx)
	if err != nil {
		return nil, err
	}
	return app.NewDesktop(h.Config, h.Logger, h.Worker, h.SourceProcess), nil
}
