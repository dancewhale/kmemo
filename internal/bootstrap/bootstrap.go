package bootstrap

import (
	"context"
	"fmt"
	"path/filepath"

	"go.uber.org/zap"

	"kmemo/internal/actions"
	filestoreadapter "kmemo/internal/adapters/filestore"
	grpcpython "kmemo/internal/adapters/fsrs/grpc_python"
	"kmemo/internal/adapters/fsrs/noop"
	"kmemo/internal/adapters/grpcworker"
	sourceprocessgrpc "kmemo/internal/adapters/sourceprocess/grpc_python"
	sourceprocessnoop "kmemo/internal/adapters/sourceprocess/noop"
	"kmemo/internal/app"
	"kmemo/internal/config"
	"kmemo/internal/contracts"
	"kmemo/internal/contracts/fsrs"
	"kmemo/internal/contracts/sourceprocess"
	"kmemo/internal/file"
	"kmemo/internal/storage"
	"kmemo/internal/storage/repository"
	"kmemo/internal/zaplog"
)

// Headless bundles dependencies for CLI / non-UI entrypoints.
type Headless struct {
	Config        config.Config
	Logger        *zap.Logger
	Worker        *grpcworker.Client
	FSRS          fsrs.FSRSScheduler
	SourceProcess sourceprocess.Processor
	Storage       *storage.Storage
	Repositories  repository.RepositoryFactory
	Transactions  repository.TransactionManager
	FileStore     contracts.FileStore
	Actions       *actions.Actions
}

// NewHeadless wires config, local storage, file store, optional Python gRPC worker, and action-layer dependencies.
func NewHeadless(ctx context.Context) (*Headless, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}
	logFile := ""
	if cfg.LogFileEnabled {
		logFile = filepath.Join(cfg.LogsDir, cfg.LogFileName)
	}
	logger, err := zaplog.NewForHost(cfg.LogLevel, cfg.LogFileEnabled, logFile)
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
		zap.Bool("log_file_enabled", cfg.LogFileEnabled),
		zap.String("log_file", logFile),
		zap.Bool("repository_debug", cfg.RepositoryDebug),
		zap.Duration("db_slow_threshold", cfg.DBSlowThreshold),
	)

	store, err := storage.New(storage.Options{
		Driver:          cfg.DBDriver,
		DSN:             filepath.Clean(cfg.DBPath),
		LogLevel:        cfg.LogLevel,
		SlowThreshold:   cfg.DBSlowThreshold,
		RepositoryDebug: cfg.RepositoryDebug,
		Logger:          logger,
	}, true)
	if err != nil {
		return nil, fmt.Errorf("init storage: %w", err)
	}

	repos := repository.NewRepositoryFactory(store.DB())
	txManager := repository.NewTransactionManager(store.DB())

	fileBackend, err := file.NewFileStore(file.DefaultConfig(cfg.VaultDir))
	if err != nil {
		return nil, fmt.Errorf("init file store: %w", err)
	}
	fileStore := filestoreadapter.NewAdapter(fileBackend)

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

	actionSet := actions.New(actions.Dependencies{
		Repositories:  repos,
		Transactions:  txManager,
		FileStore:     fileStore,
		FSRS:          sched,
		SourceProcess: processor,
	})

	return &Headless{
		Config:        cfg,
		Logger:        logger,
		Worker:        worker,
		FSRS:          sched,
		SourceProcess: processor,
		Storage:       store,
		Repositories:  repos,
		Transactions:  txManager,
		FileStore:     fileStore,
		Actions:       actionSet,
	}, nil
}

// NewDesktop builds the object graph for Wails bindings.
func NewDesktop(ctx context.Context) (*app.Desktop, error) {
	h, err := NewHeadless(ctx)
	if err != nil {
		return nil, err
	}
	return app.NewDesktop(h.Config, h.Logger, h.Worker, h.SourceProcess, h.Actions), nil
}
