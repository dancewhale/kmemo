package app

import (
	"context"

	"go.uber.org/zap"

	"kmemo/internal/adapters/grpcworker"
	"kmemo/internal/config"
	"kmemo/internal/contracts/sourceprocess"
	"kmemo/internal/zaplog"
)

// Desktop is bound to the Wails frontend. Keep methods small and delegate to internal/services later.
type Desktop struct {
	cfg           config.Config
	logger        *zap.Logger
	worker        *grpcworker.Client
	sourceProcess sourceprocess.Processor
}

// NewDesktop constructs the Wails-facing app shell.
func NewDesktop(cfg config.Config, logger *zap.Logger, worker *grpcworker.Client, sourceProcess sourceprocess.Processor) *Desktop {
	if logger == nil {
		logger = zaplog.Nop()
	}
	return &Desktop{cfg: cfg, logger: logger.Named("desktop"), worker: worker, sourceProcess: sourceProcess}
}

// OnStartup is registered with Wails for lifecycle hooks.
func (d *Desktop) OnStartup(ctx context.Context) {
	ctx = zaplog.WithLogger(ctx, d.logger)
	ctx, _ = zaplog.EnsureRequestID(ctx)
	// TODO: warm caches, migrate SQLite, verify Python health, etc.
	zaplog.FromContext(ctx).Info("desktop startup")
}

// GetVersion returns a static label for the skeleton UI.
func (d *Desktop) GetVersion() string {
	return "0.1.0-skeleton"
}

// PythonEndpoint exposes where the UI thinks the worker lives (debug / future settings UI).
func (d *Desktop) PythonEndpoint() string {
	return d.cfg.PythonGRPCAddr
}

// GetSourceProcessCapabilities returns worker-advertised source-process capabilities.
func (d *Desktop) GetSourceProcessCapabilities(ctx context.Context) (*sourceprocess.Capabilities, error) {
	return d.sourceProcess.GetCapabilities(ctx)
}

// OnShutdown releases host resources when the Wails app exits.
func (d *Desktop) OnShutdown(ctx context.Context) {
	ctx = zaplog.WithLogger(ctx, d.logger)
	ctx, _ = zaplog.EnsureRequestID(ctx)
	logger := zaplog.FromContext(ctx)
	if d == nil || d.worker == nil {
		logger.Info("desktop shutdown")
		return
	}
	if err := d.worker.Close(); err != nil {
		logger.Error("desktop shutdown close grpc worker failed", zap.Error(err))
		return
	}
	logger.Info("desktop shutdown")
}
