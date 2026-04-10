package grpcpython

import (
	"context"
	"fmt"

	kmemov1 "kmemo/gen/kmemo/v1"
	"kmemo/internal/contracts/sourceprocess"
)

// Processor implements sourceprocess.Processor via KmemoProcessor gRPC (Python worker).
type Processor struct {
	api kmemov1.KmemoProcessorClient
}

// NewProcessor constructs a Processor. api must be non-nil for live calls.
func NewProcessor(api kmemov1.KmemoProcessorClient) *Processor {
	return &Processor{api: api}
}

var _ sourceprocess.Processor = (*Processor)(nil)

// CleanHTML implements sourceprocess.Processor.
func (p *Processor) CleanHTML(ctx context.Context, in sourceprocess.CleanHTMLInput) (*sourceprocess.CleanHTMLOutput, error) {
	if p == nil || p.api == nil {
		return nil, sourceprocess.ErrUnavailable
	}
	if in.HTML == "" {
		return nil, fmt.Errorf("%w: html is required", sourceprocess.ErrInvalidInput)
	}
	resp, err := p.api.CleanHtml(ctx, cleanHTMLRequest(in))
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, fmt.Errorf("sourceprocess: empty CleanHtml response")
	}
	return cleanHTMLOutput(resp), nil
}

// SubmitImportJob implements sourceprocess.Processor.
func (p *Processor) SubmitImportJob(ctx context.Context, in sourceprocess.SubmitImportJobInput) (*sourceprocess.SubmitImportJobOutput, error) {
	if p == nil || p.api == nil {
		return nil, sourceprocess.ErrUnavailable
	}
	if in.JobID == "" {
		return nil, fmt.Errorf("%w: job_id is required", sourceprocess.ErrInvalidInput)
	}
	if in.SourceType == "" {
		return nil, fmt.Errorf("%w: source_type is required", sourceprocess.ErrInvalidInput)
	}
	if in.WorkspaceDir == "" || in.OutputDir == "" || in.TempDir == "" {
		return nil, fmt.Errorf("%w: workspace_dir, output_dir, and temp_dir are required", sourceprocess.ErrInvalidInput)
	}
	if in.SourcePath == nil && in.SourceURI == nil && in.SourceURL == nil && in.RawHTML == nil {
		return nil, fmt.Errorf("%w: one of source_path, source_uri, source_url, or raw_html is required", sourceprocess.ErrInvalidInput)
	}
	resp, err := p.api.SubmitImportJob(ctx, submitImportJobRequest(in))
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, fmt.Errorf("sourceprocess: empty SubmitImportJob response")
	}
	return submitImportJobOutput(resp), nil
}

// GetJob implements sourceprocess.Processor.
func (p *Processor) GetJob(ctx context.Context, jobID string) (*sourceprocess.Job, error) {
	if p == nil || p.api == nil {
		return nil, sourceprocess.ErrUnavailable
	}
	if jobID == "" {
		return nil, fmt.Errorf("%w: job_id is required", sourceprocess.ErrInvalidInput)
	}
	resp, err := p.api.GetJob(ctx, &kmemov1.GetJobRequest{JobId: jobID})
	if err != nil {
		return nil, err
	}
	if resp == nil || resp.GetJob() == nil {
		return nil, fmt.Errorf("sourceprocess: empty GetJob response")
	}
	return jobFromProto(resp.GetJob()), nil
}

// ListJobEvents implements sourceprocess.Processor.
func (p *Processor) ListJobEvents(ctx context.Context, in sourceprocess.ListJobEventsInput) ([]sourceprocess.JobEvent, error) {
	if p == nil || p.api == nil {
		return nil, sourceprocess.ErrUnavailable
	}
	if in.JobID == "" {
		return nil, fmt.Errorf("%w: job_id is required", sourceprocess.ErrInvalidInput)
	}
	resp, err := p.api.ListJobEvents(ctx, listJobEventsRequest(in))
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, fmt.Errorf("sourceprocess: empty ListJobEvents response")
	}
	return jobEventsFromProto(resp.GetEvents()), nil
}

// CancelJob implements sourceprocess.Processor.
func (p *Processor) CancelJob(ctx context.Context, jobID string) (*sourceprocess.CancelJobOutput, error) {
	if p == nil || p.api == nil {
		return nil, sourceprocess.ErrUnavailable
	}
	if jobID == "" {
		return nil, fmt.Errorf("%w: job_id is required", sourceprocess.ErrInvalidInput)
	}
	resp, err := p.api.CancelJob(ctx, &kmemov1.CancelJobRequest{JobId: jobID})
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, fmt.Errorf("sourceprocess: empty CancelJob response")
	}
	return cancelJobOutput(resp), nil
}

// GetCapabilities implements sourceprocess.Processor.
func (p *Processor) GetCapabilities(ctx context.Context) (*sourceprocess.Capabilities, error) {
	if p == nil || p.api == nil {
		return nil, sourceprocess.ErrUnavailable
	}
	resp, err := p.api.GetCapabilities(ctx, &kmemov1.GetCapabilitiesRequest{})
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, fmt.Errorf("sourceprocess: empty GetCapabilities response")
	}
	return capabilitiesFromProto(resp), nil
}
