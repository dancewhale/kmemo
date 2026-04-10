package sourceprocess

import "context"

// Processor is the application port for source processing via the Python worker or future backends.
type Processor interface {
	CleanHTML(ctx context.Context, in CleanHTMLInput) (*CleanHTMLOutput, error)
	SubmitImportJob(ctx context.Context, in SubmitImportJobInput) (*SubmitImportJobOutput, error)
	GetJob(ctx context.Context, jobID string) (*Job, error)
	ListJobEvents(ctx context.Context, in ListJobEventsInput) ([]JobEvent, error)
	CancelJob(ctx context.Context, jobID string) (*CancelJobOutput, error)
	GetCapabilities(ctx context.Context) (*Capabilities, error)
}
