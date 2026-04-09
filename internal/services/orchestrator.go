package services

import (
	"context"

	"kmemo/internal/contracts/fsrs"
)

// Orchestrator is the future façade over storage, HTML, indexing, and Python RPC.
// Skeleton only: no scheduling or import pipelines yet.
type Orchestrator struct {
	scheduler fsrs.FSRSScheduler
}

// NewOrchestrator constructs a placeholder coordinator.
func NewOrchestrator(sched fsrs.FSRSScheduler) *Orchestrator {
	return &Orchestrator{scheduler: sched}
}

// PingPython is a no-op placeholder for health checks.
func (o *Orchestrator) PingPython(ctx context.Context) error {
	// TODO: replace with a cheap RPC or channel health probe when defined.
	_ = ctx
	_ = o.scheduler
	return nil
}
