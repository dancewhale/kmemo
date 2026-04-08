package file

import (
	"fmt"
	"os"
	"path/filepath"
)

// WriteFileAtomically writes content to a file atomically using temp file + rename
func WriteFileAtomically(targetPath string, content []byte) error {
	// Ensure parent directory exists
	dir := filepath.Dir(targetPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("%w: failed to create directory: %v", ErrIO, err)
	}

	// Create temporary file in the same directory
	tmpFile, err := os.CreateTemp(dir, ".tmp-*")
	if err != nil {
		return fmt.Errorf("%w: failed to create temp file: %v", ErrIO, err)
	}
	tmpPath := tmpFile.Name()

	// Clean up temp file on error
	defer func() {
		if tmpFile != nil {
			tmpFile.Close()
			os.Remove(tmpPath)
		}
	}()

	// Write content
	if _, err := tmpFile.Write(content); err != nil {
		return fmt.Errorf("%w: failed to write content: %v", ErrIO, err)
	}

	// Sync to disk
	if err := tmpFile.Sync(); err != nil {
		return fmt.Errorf("%w: failed to sync: %v", ErrIO, err)
	}

	// Close temp file
	if err := tmpFile.Close(); err != nil {
		return fmt.Errorf("%w: failed to close temp file: %v", ErrIO, err)
	}
	tmpFile = nil // Prevent defer cleanup

	// Atomic rename
	if err := os.Rename(tmpPath, targetPath); err != nil {
		os.Remove(tmpPath) // Clean up on rename failure
		return fmt.Errorf("%w: failed to rename: %v", ErrIO, err)
	}

	// Sync parent directory (best effort, ignore errors on platforms that don't support it)
	if dirFile, err := os.Open(dir); err == nil {
		dirFile.Sync()
		dirFile.Close()
	}

	return nil
}

// ReplaceFileAtomically replaces an existing file atomically
func ReplaceFileAtomically(targetPath string, content []byte) error {
	// Check if target exists
	if _, err := os.Stat(targetPath); os.IsNotExist(err) {
		return fmt.Errorf("%w: target file does not exist", ErrFileObjectNotFound)
	}

	// Use same atomic write logic
	return WriteFileAtomically(targetPath, content)
}
