package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"kmemo/internal/bootstrap"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	h, err := bootstrap.NewHeadless(ctx)
	if err != nil {
		log.Fatalf("bootstrap: %v", err)
	}
	defer func() {
		if h.Py != nil {
			_ = h.Py.Close()
		}
	}()

	if h.Config.SkipPython {
		fmt.Fprintln(os.Stdout, "kmemo headless host ready (Python gRPC skipped: KMEMO_SKIP_PYTHON=1)")
	} else {
		fmt.Fprintf(os.Stdout, "kmemo headless host ready (python gRPC: %s)\n", h.Config.PythonGRPCAddr)
	}
	fmt.Fprintln(os.Stdout, "TODO: CLI commands, storage, and orchestration will attach here.")

	<-ctx.Done()
	fmt.Fprintln(os.Stdout, "shutting down")
}
