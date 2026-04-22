package fsrs

import (
	"context"

	"kmemo/internal/storage/models"
)

// FSRSScheduler is the application port for spaced-repetition scheduling via the Python worker or future backends.
type FSRSScheduler interface {
	// Review applies a rating and returns the next SRS state and a review log snapshot to persist.
	Review(ctx context.Context, in ReviewInput) (*ReviewOutput, error)

	GetRetrievability(ctx context.Context, in RetrievabilityInput) (*RetrievabilityOutput, error)

	Reschedule(ctx context.Context, in RescheduleInput) (*RescheduleOutput, error)

	// SetGlobalSetting pushes FSRS parameters to the remote worker (no-op for backends without global state).
	SetGlobalSetting(ctx context.Context, param *models.FSRSParameter) error

	OptimizeParameters(ctx context.Context, in OptimizeParametersInput) (*models.FSRSParameter, error)

	// GetBuiltinDefaultParameters returns the fsrs library defaults from the Python worker (unavailable without gRPC).
	GetBuiltinDefaultParameters(ctx context.Context) (*BuiltinDefaultParameters, error)
}
