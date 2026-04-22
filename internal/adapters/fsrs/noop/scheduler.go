package noop

import (
	"context"

	"kmemo/internal/contracts/fsrs"
	"kmemo/internal/storage/models"
)

// Scheduler is a no-op FSRSScheduler for hosts without a Python worker.
type Scheduler struct{}

var _ fsrs.FSRSScheduler = (*Scheduler)(nil)

// Review implements fsrs.FSRSScheduler.
func (s *Scheduler) Review(ctx context.Context, in fsrs.ReviewInput) (*fsrs.ReviewOutput, error) {
	_, _ = ctx, in
	return nil, fsrs.ErrUnavailable
}

// GetRetrievability implements fsrs.FSRSScheduler.
func (s *Scheduler) GetRetrievability(ctx context.Context, in fsrs.RetrievabilityInput) (*fsrs.RetrievabilityOutput, error) {
	_, _ = ctx, in
	return nil, fsrs.ErrUnavailable
}

// Reschedule implements fsrs.FSRSScheduler.
func (s *Scheduler) Reschedule(ctx context.Context, in fsrs.RescheduleInput) (*fsrs.RescheduleOutput, error) {
	_, _ = ctx, in
	return nil, fsrs.ErrUnavailable
}

// SetGlobalSetting implements fsrs.FSRSScheduler.
func (s *Scheduler) SetGlobalSetting(ctx context.Context, param *models.FSRSParameter) error {
	_, _ = ctx, param
	return fsrs.ErrUnavailable
}

// OptimizeParameters implements fsrs.FSRSScheduler.
func (s *Scheduler) OptimizeParameters(ctx context.Context, in fsrs.OptimizeParametersInput) (*models.FSRSParameter, error) {
	_, _ = ctx, in
	return nil, fsrs.ErrUnavailable
}

// GetBuiltinDefaultParameters implements fsrs.FSRSScheduler.
func (s *Scheduler) GetBuiltinDefaultParameters(ctx context.Context) (*fsrs.BuiltinDefaultParameters, error) {
	_ = ctx
	return nil, fsrs.ErrUnavailable
}
