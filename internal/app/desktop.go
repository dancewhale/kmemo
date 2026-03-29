package app

import (
	"context"

	"kmemo/internal/config"
	"kmemo/internal/pyclient"
)

// Desktop is bound to the Wails frontend. Keep methods small and delegate to internal/services later.
type Desktop struct {
	cfg config.Config
	py  *pyclient.Client
}

// NewDesktop constructs the Wails-facing app shell.
func NewDesktop(cfg config.Config, py *pyclient.Client) *Desktop {
	return &Desktop{cfg: cfg, py: py}
}

// OnStartup is registered with Wails for lifecycle hooks.
func (d *Desktop) OnStartup(ctx context.Context) {
	// TODO: warm caches, migrate SQLite, verify Python health, etc.
	_ = d.py
	_ = ctx
}

// GetVersion returns a static label for the skeleton UI.
func (d *Desktop) GetVersion() string {
	return "0.1.0-skeleton"
}

// PythonEndpoint exposes where the UI thinks the worker lives (debug / future settings UI).
func (d *Desktop) PythonEndpoint() string {
	return d.cfg.PythonGRPCAddr
}

// OnShutdown releases host resources when the Wails app exits.
func (d *Desktop) OnShutdown(ctx context.Context) {
	_ = ctx
	if d == nil || d.py == nil {
		return
	}
	_ = d.py.Close()
}
