package file

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// fileStore is the default implementation of FileStore
type fileStore struct {
	config    FileStoreConfig
	kindSpecs map[FileObjectKind]FileObjectKindSpec
}

// NewFileStore creates a new FileStore instance
func NewFileStore(config FileStoreConfig) (FileStore, error) {
	// Build kind specs map
	kindSpecs := make(map[FileObjectKind]FileObjectKindSpec)
	for _, spec := range config.Kinds {
		kindSpecs[spec.Kind] = spec
	}

	return &fileStore{
		config:    config,
		kindSpecs: kindSpecs,
	}, nil
}

// getKindSpec retrieves the kind spec for a given kind
func (s *fileStore) getKindSpec(kind FileObjectKind) (FileObjectKindSpec, error) {
	spec, ok := s.kindSpecs[kind]
	if !ok {
		return FileObjectKindSpec{}, fmt.Errorf("%w: %s", ErrInvalidKind, kind)
	}
	return spec, nil
}

// CreateFileObject creates a new file object
func (s *fileStore) CreateFileObject(ctx context.Context, input CreateFileObjectInput) (FileObjectLocation, error) {
	// Validate and normalize input
	kindSpec, err := s.getKindSpec(input.Kind)
	if err != nil {
		return FileObjectLocation{}, err
	}

	id, err := NormalizeUUIDv7(string(input.ID))
	if err != nil {
		return FileObjectLocation{}, err
	}

	if input.Ext == "" {
		return FileObjectLocation{}, fmt.Errorf("%w: extension is required", ErrInvalidFileExtension)
	}

	slug := NormalizeSlug(input.Slug)

	// Build lock path
	lockPath := BuildLockDirectoryPath(s.config.RootDir, kindSpec, id)

	// Execute with lock
	var location FileObjectLocation
	err = WithObjectLock(ctx, lockPath, s.config.LockWait, s.config.LockRetry, func() error {
		// Check if object already exists in active scope
		activeLoc, err := s.findFileObjectInScope(ctx, FileObjectRef{Kind: input.Kind, ID: id}, kindSpec, FileObjectScopeActive)
		if err != nil {
			return err
		}
		if activeLoc != nil {
			return fmt.Errorf("%w: object already exists in active scope", ErrFileObjectAlreadyExists)
		}

		// Check if object already exists in trash scope
		trashLoc, err := s.findFileObjectInScope(ctx, FileObjectRef{Kind: input.Kind, ID: id}, kindSpec, FileObjectScopeTrash)
		if err != nil {
			return err
		}
		if trashLoc != nil {
			return fmt.Errorf("%w: object already exists in trash scope", ErrFileObjectAlreadyExistsInTrash)
		}

		// Build target path
		targetPath := BuildExpectedObjectFilePath(s.config.RootDir, kindSpec, id, slug, input.Ext, FileObjectScopeActive)

		// Write file atomically
		if err := WriteFileAtomically(targetPath, input.Data); err != nil {
			return err
		}

		// Get file info
		info, err := os.Stat(targetPath)
		if err != nil {
			return fmt.Errorf("%w: failed to stat created file: %v", ErrIO, err)
		}

		// Build location result
		location = FileObjectLocation{
			Ref:      FileObjectRef{Kind: input.Kind, ID: id},
			Scope:    FileObjectScopeActive,
			Path:     targetPath,
			FileName: filepath.Base(targetPath),
			Slug:     slug,
			Ext:      input.Ext,
			Size:     info.Size(),
			ModTime:  info.ModTime(),
		}

		return nil
	})

	if err != nil {
		return FileObjectLocation{}, err
	}

	return location, nil
}

// OverwriteFileObject overwrites an existing file object
func (s *fileStore) OverwriteFileObject(ctx context.Context, input OverwriteFileObjectInput) (FileObjectLocation, error) {
	// Validate and normalize input
	kindSpec, err := s.getKindSpec(input.Kind)
	if err != nil {
		return FileObjectLocation{}, err
	}

	id, err := NormalizeUUIDv7(string(input.ID))
	if err != nil {
		return FileObjectLocation{}, err
	}

	if input.Ext == "" {
		return FileObjectLocation{}, fmt.Errorf("%w: extension is required", ErrInvalidFileExtension)
	}

	newSlug := NormalizeSlug(input.Slug)

	// Build lock path
	lockPath := BuildLockDirectoryPath(s.config.RootDir, kindSpec, id)

	// Execute with lock
	var location FileObjectLocation
	err = WithObjectLock(ctx, lockPath, s.config.LockWait, s.config.LockRetry, func() error {
		// Find existing object in active scope
		existingLoc, err := s.findFileObjectInScope(ctx, FileObjectRef{Kind: input.Kind, ID: id}, kindSpec, FileObjectScopeActive)
		if err != nil {
			return err
		}
		if existingLoc == nil {
			return fmt.Errorf("%w: object not found in active scope", ErrFileObjectNotFound)
		}

		// Check extension match
		if existingLoc.Ext != input.Ext {
			return fmt.Errorf("%w: cannot change extension from %s to %s", ErrFileExtensionMismatch, existingLoc.Ext, input.Ext)
		}

		// Check IfMatch condition if provided
		if input.IfMatch != nil {
			if input.IfMatch.ExpectedModTimeUnixNano != 0 && existingLoc.ModTime.UnixNano() != input.IfMatch.ExpectedModTimeUnixNano {
				return fmt.Errorf("%w: modification time mismatch", ErrConcurrentModification)
			}
			// Note: ContentHash check would require reading the file, skipping for now
		}

		// Handle slug change
		var targetPath string
		if existingLoc.Slug != newSlug {
			// Slug changed, need to rename file
			targetPath = BuildExpectedObjectFilePath(s.config.RootDir, kindSpec, id, newSlug, input.Ext, FileObjectScopeActive)

			// Rename existing file to new slug
			if err := os.Rename(existingLoc.Path, targetPath); err != nil {
				return fmt.Errorf("%w: failed to rename file: %v", ErrIO, err)
			}
		} else {
			// Slug unchanged, use existing path
			targetPath = existingLoc.Path
		}

		// Replace file content atomically
		if err := ReplaceFileAtomically(targetPath, input.Data); err != nil {
			return err
		}

		// Get updated file info
		info, err := os.Stat(targetPath)
		if err != nil {
			return fmt.Errorf("%w: failed to stat updated file: %v", ErrIO, err)
		}

		// Build location result
		location = FileObjectLocation{
			Ref:      FileObjectRef{Kind: input.Kind, ID: id},
			Scope:    FileObjectScopeActive,
			Path:     targetPath,
			FileName: filepath.Base(targetPath),
			Slug:     newSlug,
			Ext:      input.Ext,
			Size:     info.Size(),
			ModTime:  info.ModTime(),
		}

		return nil
	})

	if err != nil {
		return FileObjectLocation{}, err
	}

	return location, nil
}

// findFileObjectInScope finds a file object in a specific scope
func (s *fileStore) findFileObjectInScope(ctx context.Context, ref FileObjectRef, kindSpec FileObjectKindSpec, scope FileObjectScope) (*FileObjectLocation, error) {
	// Build directory path
	dir := BuildObjectDirectory(s.config.RootDir, kindSpec, ref.ID, scope)

	// Check if directory exists
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return nil, nil
	}

	// Read directory entries
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to read directory: %v", ErrIO, err)
	}

	// Find matching files
	var matches []FileObjectLocation
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		// Parse file name
		parsed, ok := ParseObjectFileName(kindSpec, entry.Name())
		if !ok {
			continue
		}

		// Check if ID matches
		if parsed.ID != ref.ID {
			continue
		}

		// Get file info
		fullPath := filepath.Join(dir, entry.Name())
		info, err := entry.Info()
		if err != nil {
			continue
		}

		// Add to matches
		matches = append(matches, FileObjectLocation{
			Ref:      ref,
			Scope:    scope,
			Path:     fullPath,
			FileName: entry.Name(),
			Slug:     parsed.Slug,
			Ext:      parsed.Ext,
			Size:     info.Size(),
			ModTime:  info.ModTime(),
		})
	}

	// Check for duplicates
	if len(matches) > 1 {
		return nil, fmt.Errorf("%w: found %d files with same ID", ErrDuplicateFileObjectID, len(matches))
	}

	if len(matches) == 0 {
		return nil, nil
	}

	return &matches[0], nil
}

// FindFileObject finds a file object's location
func (s *fileStore) FindFileObject(ctx context.Context, lookup FileObjectLookup, scope FileObjectScope) (*FileObjectLocation, error) {
	// Validate and normalize input
	kindSpec, err := s.getKindSpec(lookup.Ref.Kind)
	if err != nil {
		return nil, err
	}

	id, err := NormalizeUUIDv7(string(lookup.Ref.ID))
	if err != nil {
		return nil, err
	}

	// If name hint is provided, try exact file name first
	if lookup.Name != nil && lookup.Name.Slug != "" && lookup.Name.Ext != "" {
		slug := NormalizeSlug(lookup.Name.Slug)
		ext := lookup.Name.Ext

		// Try each scope based on scope parameter
		scopes := []FileObjectScope{}
		switch scope {
		case FileObjectScopeActive:
			scopes = []FileObjectScope{FileObjectScopeActive}
		case FileObjectScopeTrash:
			scopes = []FileObjectScope{FileObjectScopeTrash}
		case FileObjectScopeAny:
			scopes = []FileObjectScope{FileObjectScopeActive, FileObjectScopeTrash}
		}

		for _, scope := range scopes {
			expectedPath := BuildExpectedObjectFilePath(s.config.RootDir, kindSpec, id, slug, ext, scope)
			if info, err := os.Stat(expectedPath); err == nil {
				return &FileObjectLocation{
					Ref:      FileObjectRef{Kind: lookup.Ref.Kind, ID: id},
					Scope:    scope,
					Path:     expectedPath,
					FileName: filepath.Base(expectedPath),
					Slug:     slug,
					Ext:      ext,
					Size:     info.Size(),
					ModTime:  info.ModTime(),
				}, nil
			}
		}
	}

	// Fall back to directory scan
	switch scope {
	case FileObjectScopeActive:
		return s.findFileObjectInScope(ctx, FileObjectRef{Kind: lookup.Ref.Kind, ID: id}, kindSpec, FileObjectScopeActive)
	case FileObjectScopeTrash:
		return s.findFileObjectInScope(ctx, FileObjectRef{Kind: lookup.Ref.Kind, ID: id}, kindSpec, FileObjectScopeTrash)
	case FileObjectScopeAny:
		// Try active first
		loc, err := s.findFileObjectInScope(ctx, FileObjectRef{Kind: lookup.Ref.Kind, ID: id}, kindSpec, FileObjectScopeActive)
		if err != nil {
			return nil, err
		}
		if loc != nil {
			return loc, nil
		}
		// Try trash
		return s.findFileObjectInScope(ctx, FileObjectRef{Kind: lookup.Ref.Kind, ID: id}, kindSpec, FileObjectScopeTrash)
	default:
		return nil, fmt.Errorf("%w: invalid scope", ErrInvalidInput)
	}
}

// ReadFileObject reads a file object's content and location
func (s *fileStore) ReadFileObject(ctx context.Context, lookup FileObjectLookup, scope FileObjectScope) ([]byte, FileObjectLocation, error) {
	// Find the object
	loc, err := s.FindFileObject(ctx, lookup, scope)
	if err != nil {
		return nil, FileObjectLocation{}, err
	}
	if loc == nil {
		return nil, FileObjectLocation{}, fmt.Errorf("%w: object not found", ErrFileObjectNotFound)
	}

	// Read file content
	content, err := os.ReadFile(loc.Path)
	if err != nil {
		return nil, FileObjectLocation{}, fmt.Errorf("%w: failed to read file: %v", ErrIO, err)
	}

	return content, *loc, nil
}

// FileObjectExists checks if a file object exists
func (s *fileStore) FileObjectExists(ctx context.Context, lookup FileObjectLookup, scope FileObjectScope) (bool, error) {
	loc, err := s.FindFileObject(ctx, lookup, scope)
	if err != nil {
		return false, err
	}
	return loc != nil, nil
}

// MoveFileObjectToTrash moves a file object from active to trash
func (s *fileStore) MoveFileObjectToTrash(ctx context.Context, ref FileObjectRef) (FileObjectLocation, error) {
	// Validate and normalize input
	kindSpec, err := s.getKindSpec(ref.Kind)
	if err != nil {
		return FileObjectLocation{}, err
	}

	id, err := NormalizeUUIDv7(string(ref.ID))
	if err != nil {
		return FileObjectLocation{}, err
	}

	// Build lock path
	lockPath := BuildLockDirectoryPath(s.config.RootDir, kindSpec, id)

	// Execute with lock
	var location FileObjectLocation
	err = WithObjectLock(ctx, lockPath, s.config.LockWait, s.config.LockRetry, func() error {
		// Find object in active scope
		activeLoc, err := s.findFileObjectInScope(ctx, FileObjectRef{Kind: ref.Kind, ID: id}, kindSpec, FileObjectScopeActive)
		if err != nil {
			return err
		}
		if activeLoc == nil {
			return fmt.Errorf("%w: object not found in active scope", ErrFileObjectNotFound)
		}

		// Check if object already exists in trash
		trashLoc, err := s.findFileObjectInScope(ctx, FileObjectRef{Kind: ref.Kind, ID: id}, kindSpec, FileObjectScopeTrash)
		if err != nil {
			return err
		}
		if trashLoc != nil {
			return fmt.Errorf("%w: object already exists in trash", ErrTrashConflict)
		}

		// Build trash target path
		trashDir := BuildTrashObjectDirectory(s.config.RootDir, kindSpec, id)
		if err := os.MkdirAll(trashDir, 0755); err != nil {
			return fmt.Errorf("%w: failed to create trash directory: %v", ErrIO, err)
		}

		trashPath := filepath.Join(trashDir, activeLoc.FileName)

		// Move file to trash
		if err := os.Rename(activeLoc.Path, trashPath); err != nil {
			return fmt.Errorf("%w: failed to move file to trash: %v", ErrIO, err)
		}

		// Get file info
		info, err := os.Stat(trashPath)
		if err != nil {
			return fmt.Errorf("%w: failed to stat moved file: %v", ErrIO, err)
		}

		// Build location result
		location = FileObjectLocation{
			Ref:      FileObjectRef{Kind: ref.Kind, ID: id},
			Scope:    FileObjectScopeTrash,
			Path:     trashPath,
			FileName: activeLoc.FileName,
			Slug:     activeLoc.Slug,
			Ext:      activeLoc.Ext,
			Size:     info.Size(),
			ModTime:  info.ModTime(),
		}

		return nil
	})

	if err != nil {
		return FileObjectLocation{}, err
	}

	return location, nil
}

// RestoreFileObjectFromTrash restores a file object from trash to active
func (s *fileStore) RestoreFileObjectFromTrash(ctx context.Context, ref FileObjectRef) (FileObjectLocation, error) {
	// Validate and normalize input
	kindSpec, err := s.getKindSpec(ref.Kind)
	if err != nil {
		return FileObjectLocation{}, err
	}

	id, err := NormalizeUUIDv7(string(ref.ID))
	if err != nil {
		return FileObjectLocation{}, err
	}

	// Build lock path
	lockPath := BuildLockDirectoryPath(s.config.RootDir, kindSpec, id)

	// Execute with lock
	var location FileObjectLocation
	err = WithObjectLock(ctx, lockPath, s.config.LockWait, s.config.LockRetry, func() error {
		// Find object in trash scope
		trashLoc, err := s.findFileObjectInScope(ctx, FileObjectRef{Kind: ref.Kind, ID: id}, kindSpec, FileObjectScopeTrash)
		if err != nil {
			return err
		}
		if trashLoc == nil {
			return fmt.Errorf("%w: object not found in trash", ErrFileObjectNotFound)
		}

		// Check if object already exists in active scope
		activeLoc, err := s.findFileObjectInScope(ctx, FileObjectRef{Kind: ref.Kind, ID: id}, kindSpec, FileObjectScopeActive)
		if err != nil {
			return err
		}
		if activeLoc != nil {
			return fmt.Errorf("%w: object already exists in active scope", ErrRestoreConflict)
		}

		// Build active target path
		activeDir := BuildActiveObjectDirectory(s.config.RootDir, kindSpec, id)
		if err := os.MkdirAll(activeDir, 0755); err != nil {
			return fmt.Errorf("%w: failed to create active directory: %v", ErrIO, err)
		}

		activePath := filepath.Join(activeDir, trashLoc.FileName)

		// Move file to active
		if err := os.Rename(trashLoc.Path, activePath); err != nil {
			return fmt.Errorf("%w: failed to restore file: %v", ErrIO, err)
		}

		// Get file info
		info, err := os.Stat(activePath)
		if err != nil {
			return fmt.Errorf("%w: failed to stat restored file: %v", ErrIO, err)
		}

		// Build location result
		location = FileObjectLocation{
			Ref:      FileObjectRef{Kind: ref.Kind, ID: id},
			Scope:    FileObjectScopeActive,
			Path:     activePath,
			FileName: trashLoc.FileName,
			Slug:     trashLoc.Slug,
			Ext:      trashLoc.Ext,
			Size:     info.Size(),
			ModTime:  info.ModTime(),
		}

		return nil
	})

	if err != nil {
		return FileObjectLocation{}, err
	}

	return location, nil
}

// PermanentlyDeleteFileObject permanently deletes a file object from trash
func (s *fileStore) PermanentlyDeleteFileObject(ctx context.Context, ref FileObjectRef) error {
	// Validate and normalize input
	kindSpec, err := s.getKindSpec(ref.Kind)
	if err != nil {
		return err
	}

	id, err := NormalizeUUIDv7(string(ref.ID))
	if err != nil {
		return err
	}

	// Build lock path
	lockPath := BuildLockDirectoryPath(s.config.RootDir, kindSpec, id)

	// Execute with lock
	return WithObjectLock(ctx, lockPath, s.config.LockWait, s.config.LockRetry, func() error {
		// Find object in trash scope
		trashLoc, err := s.findFileObjectInScope(ctx, FileObjectRef{Kind: ref.Kind, ID: id}, kindSpec, FileObjectScopeTrash)
		if err != nil {
			return err
		}
		if trashLoc == nil {
			return fmt.Errorf("%w: object not found in trash", ErrFileObjectNotFound)
		}

		// Delete file
		if err := os.Remove(trashLoc.Path); err != nil {
			return fmt.Errorf("%w: failed to delete file: %v", ErrIO, err)
		}

		return nil
	})
}

// RenameFileObjectSlug renames the slug portion of a file object's name
func (s *fileStore) RenameFileObjectSlug(ctx context.Context, ref FileObjectRef, newSlug string) (FileObjectLocation, error) {
	// Validate and normalize input
	kindSpec, err := s.getKindSpec(ref.Kind)
	if err != nil {
		return FileObjectLocation{}, err
	}

	// Check if kind supports slug
	if kindSpec.FileNameStyle != FileNameStyleWithSlug {
		return FileObjectLocation{}, fmt.Errorf("%w: kind does not support slug", ErrUnsupportedOperation)
	}

	id, err := NormalizeUUIDv7(string(ref.ID))
	if err != nil {
		return FileObjectLocation{}, err
	}

	normalizedSlug := NormalizeSlug(newSlug)

	// Build lock path
	lockPath := BuildLockDirectoryPath(s.config.RootDir, kindSpec, id)

	// Execute with lock
	var location FileObjectLocation
	err = WithObjectLock(ctx, lockPath, s.config.LockWait, s.config.LockRetry, func() error {
		// Find object in active scope
		activeLoc, err := s.findFileObjectInScope(ctx, FileObjectRef{Kind: ref.Kind, ID: id}, kindSpec, FileObjectScopeActive)
		if err != nil {
			return err
		}
		if activeLoc == nil {
			return fmt.Errorf("%w: object not found in active scope", ErrFileObjectNotFound)
		}

		// Check if slug is unchanged
		if activeLoc.Slug == normalizedSlug {
			location = *activeLoc
			return nil
		}

		// Build new file name
		newFileName := BuildObjectFileName(kindSpec, id, normalizedSlug, activeLoc.Ext)
		newPath := filepath.Join(filepath.Dir(activeLoc.Path), newFileName)

		// Rename file
		if err := os.Rename(activeLoc.Path, newPath); err != nil {
			return fmt.Errorf("%w: failed to rename file: %v", ErrIO, err)
		}

		// Get file info
		info, err := os.Stat(newPath)
		if err != nil {
			return fmt.Errorf("%w: failed to stat renamed file: %v", ErrIO, err)
		}

		// Build location result
		location = FileObjectLocation{
			Ref:      FileObjectRef{Kind: ref.Kind, ID: id},
			Scope:    FileObjectScopeActive,
			Path:     newPath,
			FileName: newFileName,
			Slug:     normalizedSlug,
			Ext:      activeLoc.Ext,
			Size:     info.Size(),
			ModTime:  info.ModTime(),
		}

		return nil
	})

	if err != nil {
		return FileObjectLocation{}, err
	}

	return location, nil
}

// ScanFileObjects scans for file objects matching the given options
func (s *fileStore) ScanFileObjects(ctx context.Context, options ScanFileObjectsOptions) (<-chan FileObjectLocation, <-chan error) {
	locationCh := make(chan FileObjectLocation)
	errorCh := make(chan error, 1)

	go func() {
		defer close(locationCh)
		defer close(errorCh)

		// Determine which kinds to scan
		kinds := []FileObjectKindSpec{}
		if options.Kind != nil {
			spec, err := s.getKindSpec(*options.Kind)
			if err != nil {
				errorCh <- err
				return
			}
			kinds = append(kinds, spec)
		} else {
			kinds = s.config.Kinds
		}

		// Determine which scopes to scan
		scopes := []FileObjectScope{}
		switch options.Scope {
		case FileObjectScopeActive:
			scopes = []FileObjectScope{FileObjectScopeActive}
		case FileObjectScopeTrash:
			scopes = []FileObjectScope{FileObjectScopeTrash}
		case FileObjectScopeAny:
			scopes = []FileObjectScope{FileObjectScopeActive, FileObjectScopeTrash}
		default:
			scopes = []FileObjectScope{FileObjectScopeActive}
		}

		// Scan each kind and scope
		for _, kindSpec := range kinds {
			for _, scope := range scopes {
				// Build root directory for this kind and scope
				var rootDir string
				switch scope {
				case FileObjectScopeActive:
					rootDir = filepath.Join(s.config.RootDir, kindSpec.DirectoryName)
				case FileObjectScopeTrash:
					rootDir = filepath.Join(s.config.RootDir, "trash", kindSpec.DirectoryName)
				}

				// Walk directory tree
				err := filepath.WalkDir(rootDir, func(path string, d fs.DirEntry, err error) error {
					if err != nil {
						// Skip directories that don't exist
						if os.IsNotExist(err) {
							return nil
						}
						return err
					}

					// Skip directories
					if d.IsDir() {
						return nil
					}

					// Check context cancellation
					select {
					case <-ctx.Done():
						return ctx.Err()
					default:
					}

					// Parse file name
					parsed, ok := ParseObjectFileName(kindSpec, d.Name())
					if !ok {
						return nil
					}

					// Get file info
					info, err := d.Info()
					if err != nil {
						return nil
					}

					// Send location
					locationCh <- FileObjectLocation{
						Ref:      FileObjectRef{Kind: kindSpec.Kind, ID: parsed.ID},
						Scope:    scope,
						Path:     path,
						FileName: d.Name(),
						Slug:     parsed.Slug,
						Ext:      parsed.Ext,
						Size:     info.Size(),
						ModTime:  info.ModTime(),
					}

					return nil
				})

				if err != nil && err != context.Canceled {
					errorCh <- fmt.Errorf("%w: failed to scan directory: %v", ErrIO, err)
					return
				}
			}
		}
	}()

	return locationCh, errorCh
}
