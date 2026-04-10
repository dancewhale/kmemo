package sourceprocess

// CleanHTMLInput is the current transport-aligned input for synchronous HTML cleanup.
type CleanHTMLInput struct {
	SourceID string
	HTML     string
}

// CleanHTMLOutput is the current transport-aligned output for synchronous HTML cleanup.
type CleanHTMLOutput struct {
	OK          bool
	Message     string
	CleanedHTML string
}

// ImportOptions controls worker-side import processing.
type ImportOptions struct {
	ConversionMode       string
	FallbackModes        []string
	ExtractMainContent   bool
	SanitizeHTML         bool
	PreserveSemanticTags bool
	DownloadRemoteAssets bool
	InlineSmallImages    bool
	GenerateTOC          bool
	AnalyzeStructure     bool
	KeepSourceCopy       bool
	EnabledCleaners      []string
	ConverterParamsJSON  *string
}

// SubmitImportJobInput submits an async import job to the worker.
type SubmitImportJobInput struct {
	JobID          string
	SourceType     string
	SourcePath     *string
	SourceURI      *string
	SourceURL      *string
	RawHTML        *string
	WorkspaceDir   string
	OutputDir      string
	TempDir        string
	Options        ImportOptions
	Metadata       map[string]string
	IdempotencyKey *string
}

// SubmitImportJobOutput is the worker acknowledgement for job submission.
type SubmitImportJobOutput struct {
	JobID    string
	Status   string
	Accepted bool
}

// Job is the current source-process job view returned by the worker.
type Job struct {
	JobID        string
	Status       string
	Stage        string
	Progress     float32
	ResultPath   *string
	ErrorCode    *string
	ErrorMessage *string
	Result       *ImportResult
}

// ListJobEventsInput requests ordered job events after an optional sequence.
type ListJobEventsInput struct {
	JobID         string
	AfterSequence *int64
}

// JobEvent is a worker-emitted source-process event.
type JobEvent struct {
	JobID         string
	Sequence      int64
	Stage         string
	Message       string
	CreatedAtUnix int64
}

// CancelJobOutput is the worker acknowledgement for cancellation.
type CancelJobOutput struct {
	JobID  string
	Status string
}

// ImportResult holds typed output paths and metadata from source processing.
type ImportResult struct {
	EntryHTMLPath           string
	CleanedHTMLPath         string
	RawTextPath             string
	Assets                  []string
	ExtractedMetadata       map[string]string
	ManifestPath            string
	ContentHash             string
	EffectiveConversionMode string
	ConverterName           string
	ConverterVersion        string
	CleanerVersion          string
}

// ConverterCapability describes a registered converter.
type ConverterCapability struct {
	SourceType       string
	ConversionMode   string
	ConverterName    string
	ConverterVersion string
	Description      string
}

// CleanerCapability describes a registered cleaner.
type CleanerCapability struct {
	CleanerName    string
	CleanerVersion string
	Description    string
}

// Capabilities exposes worker-supported source types, modes, converters, and cleaners.
type Capabilities struct {
	SourceTypes     []string
	ConversionModes []string
	Converters      []ConverterCapability
	Cleaners        []CleanerCapability
}
