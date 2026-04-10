package sourceprocess

import "errors"

var (
	// ErrUnavailable is returned when the source-process backend is not configured (e.g. Python skipped).
	ErrUnavailable = errors.New("sourceprocess: processor unavailable")

	// ErrInvalidInput indicates missing or invalid port-level input.
	ErrInvalidInput = errors.New("sourceprocess: invalid input")

	// ErrNotFound indicates a requested job or resource does not exist.
	ErrNotFound = errors.New("sourceprocess: not found")

	// ErrConflict indicates a state conflict in remote source-process execution.
	ErrConflict = errors.New("sourceprocess: conflict")

	// ErrCanceled indicates a job was canceled.
	ErrCanceled = errors.New("sourceprocess: canceled")
)
