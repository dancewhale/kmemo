package main

import (
	"context"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"kmemo/internal/bootstrap"
	"kmemo/internal/zaplog"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	h, err := bootstrap.NewHeadless(ctx)
	if err != nil {
		fallback := zaplog.Nop()
		if logger, loggerErr := zaplog.New("error"); loggerErr == nil {
			fallback = logger.Named("headless")
		}
		fallback.Error("bootstrap failed", zap.Error(err))
		_ = fallback.Sync()
		return
	}
	defer func() {
		if h.Worker != nil {
			_ = h.Worker.Close()
		}
		_ = h.Logger.Sync()
	}()

	logger := h.Logger.Named("headless")
	if h.Config.SkipPython {
		logger.Info("kmemo headless host ready", zap.Bool("skip_python", true))
	} else {
		logger.Info("kmemo headless host ready", zap.String("python_grpc", h.Config.PythonGRPCAddr))
	}
	logger.Info("headless orchestration placeholder ready")

	<-ctx.Done()
	logger.Info("headless shutting down")
}
