package noop

import (
	"context"

	"kmemo/internal/contracts/sourceprocess"
)

// Processor is a no-op sourceprocess.Processor for hosts without a Python worker.
type Processor struct{}

var _ sourceprocess.Processor = (*Processor)(nil)

func (p *Processor) CleanHTML(ctx context.Context, in sourceprocess.CleanHTMLInput) (*sourceprocess.CleanHTMLOutput, error) {
	_, _ = ctx, in
	return nil, sourceprocess.ErrUnavailable
}

func (p *Processor) SubmitImportJob(ctx context.Context, in sourceprocess.SubmitImportJobInput) (*sourceprocess.SubmitImportJobOutput, error) {
	_, _ = ctx, in
	return nil, sourceprocess.ErrUnavailable
}

func (p *Processor) GetJob(ctx context.Context, jobID string) (*sourceprocess.Job, error) {
	_, _ = ctx, jobID
	return nil, sourceprocess.ErrUnavailable
}

func (p *Processor) ListJobEvents(ctx context.Context, in sourceprocess.ListJobEventsInput) ([]sourceprocess.JobEvent, error) {
	_, _ = ctx, in
	return nil, sourceprocess.ErrUnavailable
}

func (p *Processor) CancelJob(ctx context.Context, jobID string) (*sourceprocess.CancelJobOutput, error) {
	_, _ = ctx, jobID
	return nil, sourceprocess.ErrUnavailable
}

func (p *Processor) GetCapabilities(ctx context.Context) (*sourceprocess.Capabilities, error) {
	_ = ctx
	return nil, sourceprocess.ErrUnavailable
}
