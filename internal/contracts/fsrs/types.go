package fsrs

import (
	"time"

	"kmemo/internal/storage/models"
)

// ReviewInput is input for a single graded review (rating at a point in time).
type ReviewInput struct {
	CardID     string
	Prior      *models.CardSRS // nil means treat as "new" card state
	Rating     int
	ReviewedAt time.Time
}

// ReviewOutput holds the next persisted SRS snapshot and review log line after Review.
type ReviewOutput struct {
	Next           *models.CardSRS
	Log            *models.ReviewLog
	Retrievability float64
	Warnings       []string
	Diagnostics    string
}

// RetrievabilityInput requests retrievability for a card state at Now.
type RetrievabilityInput struct {
	CardID string
	Prior  *models.CardSRS
	Now    time.Time
}

// RetrievabilityOutput is the worker result without protobuf types.
type RetrievabilityOutput struct {
	Retrievability float64
	EvaluatedAt    *time.Time
	Warnings       []string
	Diagnostics    string
}

// RescheduleInput requests a new schedule given history and a reschedule time.
type RescheduleInput struct {
	CardID       string
	Prior        *models.CardSRS
	ReviewLogs   []*models.ReviewLog
	RescheduleAt time.Time
	Mode         string
}

// RescheduleOutput is the outcome of Reschedule.
type RescheduleOutput struct {
	Next           *models.CardSRS
	Log            *models.ReviewLog
	Retrievability float64
	Warnings       []string
	Diagnostics    string
}

// OptimizeParametersInput drives remote FSRS parameter training.
type OptimizeParametersInput struct {
	RequestID   string
	DatasetID   string
	KnowledgeID string
	ReviewLogs  []*models.ReviewLog
	Baseline    *models.FSRSParameter
	// ResultID and ResultName identify the returned FSRSParameter row when persisting.
	ResultID   string
	ResultName string
}
