package services

import (
	"context"

	"kmemo/internal/pyclient"
)

// Orchestrator is the future façade over storage, HTML, indexing, and Python RPC.
// Skeleton only: no scheduling or import pipelines yet.
type Orchestrator struct {
	py *pyclient.Client
}

// NewOrchestrator constructs a placeholder coordinator.
func NewOrchestrator(py *pyclient.Client) *Orchestrator {
	return &Orchestrator{py: py}
}

// PingPython is a no-op placeholder for health checks.
func (o *Orchestrator) PingPython(ctx context.Context) error {
	// TODO: replace with a cheap RPC or channel health probe when defined.
	_ = ctx
	_ = o.py
	return nil
}
