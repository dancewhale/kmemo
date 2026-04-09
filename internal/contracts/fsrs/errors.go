package fsrs

import "errors"

var (
	// ErrUnavailable is returned when the FSRS backend is not configured (e.g. Python skipped).
	ErrUnavailable = errors.New("fsrs: processor unavailable")

	// ErrInvalidInput indicates missing or invalid port-level input.
	ErrInvalidInput = errors.New("fsrs: invalid input")
)
