package file

import (
	"context"
	"fmt"
	"os"
	"time"
)

// AcquireObjectLock acquires a lock for the given object
func AcquireObjectLock(ctx context.Context, lockPath string, lockWait time.Duration, lockRetry time.Duration) error {
	deadline := time.Now().Add(lockWait)

	for {
		// Try to create lock directory
		err := os.Mkdir(lockPath, 0755)
		if err == nil {
			// Lock acquired successfully
			return nil
		}

		// Check if we've timed out
		if time.Now().After(deadline) {
			return fmt.Errorf("%w: failed to acquire lock after %v", ErrLockTimeout, lockWait)
		}

		// Check context cancellation
		select {
		case <-ctx.Done():
			return fmt.Errorf("%w: context cancelled", ErrLockTimeout)
		case <-time.After(lockRetry):
			// Retry after delay
		}
	}
}

// ReleaseObjectLock releases a lock for the given object
func ReleaseObjectLock(lockPath string) error {
	if err := os.Remove(lockPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("%w: failed to release lock: %v", ErrIO, err)
	}
	return nil
}

// WithObjectLock executes a function while holding an object lock
func WithObjectLock(ctx context.Context, lockPath string, lockWait time.Duration, lockRetry time.Duration, fn func() error) error {
	// Ensure parent directory exists
	if err := os.MkdirAll(lockPath[:len(lockPath)-len("/"+lockPath[len(lockPath)-1:])], 0755); err != nil {
		return fmt.Errorf("%w: failed to create lock directory parent: %v", ErrIO, err)
	}

	// Acquire lock
	if err := AcquireObjectLock(ctx, lockPath, lockWait, lockRetry); err != nil {
		return err
	}

	// Ensure lock is released
	defer ReleaseObjectLock(lockPath)

	// Execute function
	return fn()
}
