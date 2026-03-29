package bootstrap

import (
	"context"
	"fmt"

	"kmemo/internal/app"
	"kmemo/internal/config"
	"kmemo/internal/pyclient"
)

// Headless bundles dependencies for CLI / non-UI entrypoints.
type Headless struct {
	Config config.Config
	Py     *pyclient.Client
}

// NewHeadless wires config and the Python gRPC client. No SQLite / indexing yet.
func NewHeadless(ctx context.Context) (*Headless, error) {
	cfg := config.Load()
	var py *pyclient.Client
	if !cfg.SkipPython {
		var err error
		py, err = pyclient.New(ctx, cfg)
		if err != nil {
			return nil, fmt.Errorf("pyclient: %w", err)
		}
	}
	return &Headless{Config: cfg, Py: py}, nil
}

// NewDesktop builds the object graph for Wails bindings.
func NewDesktop(ctx context.Context) (*app.Desktop, error) {
	h, err := NewHeadless(ctx)
	if err != nil {
		return nil, err
	}
	return app.NewDesktop(h.Config, h.Py), nil
}
