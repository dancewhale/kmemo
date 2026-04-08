package filestore

import (
	"context"

	"kmemo/internal/contracts"
	"kmemo/internal/file"
)

// adapter wraps the file module to implement the contracts.FileStore interface
type adapter struct {
	store file.FileStore
}

// NewAdapter creates a new FileStore adapter
func NewAdapter(store file.FileStore) contracts.FileStore {
	return &adapter{
		store: store,
	}
}

// CreateFileObject creates a new file object
func (a *adapter) CreateFileObject(ctx context.Context, input contracts.CreateFileObjectInput) (contracts.FileObjectLocation, error) {
	fileInput := file.CreateFileObjectInput{
		Kind: file.FileObjectKind(input.Kind),
		ID:   file.FileObjectID(input.ID),
		Slug: input.Slug,
		Ext:  input.Ext,
		Data: input.Data,
	}

	loc, err := a.store.CreateFileObject(ctx, fileInput)
	if err != nil {
		return contracts.FileObjectLocation{}, err
	}

	return toContractsLocation(loc), nil
}

// OverwriteFileObject overwrites an existing file object
func (a *adapter) OverwriteFileObject(ctx context.Context, input contracts.OverwriteFileObjectInput) (contracts.FileObjectLocation, error) {
	var ifMatch *file.FileObjectMatchCondition
	if input.IfMatch != nil {
		ifMatch = &file.FileObjectMatchCondition{
			ExpectedModTimeUnixNano: input.IfMatch.ExpectedModTimeUnixNano,
			ExpectedContentHash:     input.IfMatch.ExpectedContentHash,
		}
	}

	fileInput := file.OverwriteFileObjectInput{
		Kind:    file.FileObjectKind(input.Kind),
		ID:      file.FileObjectID(input.ID),
		Slug:    input.Slug,
		Ext:     input.Ext,
		Data:    input.Data,
		IfMatch: ifMatch,
	}

	loc, err := a.store.OverwriteFileObject(ctx, fileInput)
	if err != nil {
		return contracts.FileObjectLocation{}, err
	}

	return toContractsLocation(loc), nil
}

// ReadFileObject reads a file object's content and location
func (a *adapter) ReadFileObject(ctx context.Context, lookup contracts.FileObjectLookup, scope contracts.FileObjectScope) ([]byte, contracts.FileObjectLocation, error) {
	fileLookup := toFileLookup(lookup)
	fileScope := file.FileObjectScope(scope)

	data, loc, err := a.store.ReadFileObject(ctx, fileLookup, fileScope)
	if err != nil {
		return nil, contracts.FileObjectLocation{}, err
	}

	return data, toContractsLocation(loc), nil
}

// FindFileObject finds a file object's location without reading content
func (a *adapter) FindFileObject(ctx context.Context, lookup contracts.FileObjectLookup, scope contracts.FileObjectScope) (*contracts.FileObjectLocation, error) {
	fileLookup := toFileLookup(lookup)
	fileScope := file.FileObjectScope(scope)

	loc, err := a.store.FindFileObject(ctx, fileLookup, fileScope)
	if err != nil {
		return nil, err
	}

	if loc == nil {
		return nil, nil
	}

	contractsLoc := toContractsLocation(*loc)
	return &contractsLoc, nil
}

// FileObjectExists checks if a file object exists
func (a *adapter) FileObjectExists(ctx context.Context, lookup contracts.FileObjectLookup, scope contracts.FileObjectScope) (bool, error) {
	fileLookup := toFileLookup(lookup)
	fileScope := file.FileObjectScope(scope)

	return a.store.FileObjectExists(ctx, fileLookup, fileScope)
}

// MoveFileObjectToTrash moves a file object from active to trash
func (a *adapter) MoveFileObjectToTrash(ctx context.Context, ref contracts.FileObjectRef) (contracts.FileObjectLocation, error) {
	fileRef := file.FileObjectRef{
		Kind: file.FileObjectKind(ref.Kind),
		ID:   file.FileObjectID(ref.ID),
	}

	loc, err := a.store.MoveFileObjectToTrash(ctx, fileRef)
	if err != nil {
		return contracts.FileObjectLocation{}, err
	}

	return toContractsLocation(loc), nil
}

// RestoreFileObjectFromTrash restores a file object from trash to active
func (a *adapter) RestoreFileObjectFromTrash(ctx context.Context, ref contracts.FileObjectRef) (contracts.FileObjectLocation, error) {
	fileRef := file.FileObjectRef{
		Kind: file.FileObjectKind(ref.Kind),
		ID:   file.FileObjectID(ref.ID),
	}

	loc, err := a.store.RestoreFileObjectFromTrash(ctx, fileRef)
	if err != nil {
		return contracts.FileObjectLocation{}, err
	}

	return toContractsLocation(loc), nil
}

// PermanentlyDeleteFileObject permanently deletes a file object from trash
func (a *adapter) PermanentlyDeleteFileObject(ctx context.Context, ref contracts.FileObjectRef) error {
	fileRef := file.FileObjectRef{
		Kind: file.FileObjectKind(ref.Kind),
		ID:   file.FileObjectID(ref.ID),
	}

	return a.store.PermanentlyDeleteFileObject(ctx, fileRef)
}

// RenameFileObjectSlug renames the slug portion of a file object's name
func (a *adapter) RenameFileObjectSlug(ctx context.Context, ref contracts.FileObjectRef, newSlug string) (contracts.FileObjectLocation, error) {
	fileRef := file.FileObjectRef{
		Kind: file.FileObjectKind(ref.Kind),
		ID:   file.FileObjectID(ref.ID),
	}

	loc, err := a.store.RenameFileObjectSlug(ctx, fileRef, newSlug)
	if err != nil {
		return contracts.FileObjectLocation{}, err
	}

	return toContractsLocation(loc), nil
}

// ScanFileObjects scans for file objects matching the given options
func (a *adapter) ScanFileObjects(ctx context.Context, options contracts.ScanFileObjectsOptions) (<-chan contracts.FileObjectLocation, <-chan error) {
	var kind *file.FileObjectKind
	if options.Kind != nil {
		k := file.FileObjectKind(*options.Kind)
		kind = &k
	}

	fileOptions := file.ScanFileObjectsOptions{
		Kind:  kind,
		Scope: file.FileObjectScope(options.Scope),
	}

	fileLocs, fileErrs := a.store.ScanFileObjects(ctx, fileOptions)

	// Convert channels
	contractsLocs := make(chan contracts.FileObjectLocation)
	contractsErrs := make(chan error, 1)

	go func() {
		defer close(contractsLocs)
		defer close(contractsErrs)

		for {
			select {
			case loc, ok := <-fileLocs:
				if !ok {
					fileLocs = nil
					if fileErrs == nil {
						return
					}
					continue
				}
				contractsLocs <- toContractsLocation(loc)
			case err, ok := <-fileErrs:
				if !ok {
					fileErrs = nil
					if fileLocs == nil {
						return
					}
					continue
				}
				contractsErrs <- err
				return
			}
		}
	}()

	return contractsLocs, contractsErrs
}

// Helper functions to convert between file and contracts types

func toFileLookup(lookup contracts.FileObjectLookup) file.FileObjectLookup {
	var name *file.FileObjectNameHint
	if lookup.Name != nil {
		name = &file.FileObjectNameHint{
			Slug: lookup.Name.Slug,
			Ext:  lookup.Name.Ext,
		}
	}

	return file.FileObjectLookup{
		Ref: file.FileObjectRef{
			Kind: file.FileObjectKind(lookup.Ref.Kind),
			ID:   file.FileObjectID(lookup.Ref.ID),
		},
		Name: name,
	}
}

func toContractsLocation(loc file.FileObjectLocation) contracts.FileObjectLocation {
	return contracts.FileObjectLocation{
		Ref: contracts.FileObjectRef{
			Kind: contracts.FileObjectKind(loc.Ref.Kind),
			ID:   contracts.FileObjectID(loc.Ref.ID),
		},
		Scope:    contracts.FileObjectScope(loc.Scope),
		Path:     loc.Path,
		FileName: loc.FileName,
		Slug:     loc.Slug,
		Ext:      loc.Ext,
		Size:     loc.Size,
		ModTime:  loc.ModTime,
	}
}
