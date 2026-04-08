package contracts

import (
	"context"
	"time"
)

// FileObjectKind represents the type of file object
type FileObjectKind string

const (
	FileObjectKindCard   FileObjectKind = "card"
	FileObjectKindAsset  FileObjectKind = "asset"
	FileObjectKindSource FileObjectKind = "source"
)

// FileObjectID is a UUID v7 string
type FileObjectID string

// FileObjectRef uniquely identifies a file object
type FileObjectRef struct {
	Kind FileObjectKind
	ID   FileObjectID
}

// FileObjectScope represents where the file object is located
type FileObjectScope string

const (
	FileObjectScopeActive FileObjectScope = "active"
	FileObjectScopeTrash  FileObjectScope = "trash"
	FileObjectScopeAny    FileObjectScope = "any"
)

// FileObjectNameHint provides optional name information for faster lookup
type FileObjectNameHint struct {
	Slug string
	Ext  string
}

// FileObjectLookup specifies how to find a file object
type FileObjectLookup struct {
	Ref  FileObjectRef
	Name *FileObjectNameHint
}

// FileObjectLocation contains the full information about a located file object
type FileObjectLocation struct {
	Ref      FileObjectRef
	Scope    FileObjectScope
	Path     string
	FileName string
	Slug     string
	Ext      string
	Size     int64
	ModTime  time.Time
}

// FileObjectMatchCondition is used for optimistic concurrency control
type FileObjectMatchCondition struct {
	ExpectedModTimeUnixNano int64
	ExpectedContentHash     string
}

// CreateFileObjectInput contains parameters for creating a new file object
type CreateFileObjectInput struct {
	Kind FileObjectKind
	ID   FileObjectID
	Slug string
	Ext  string
	Data []byte
}

// OverwriteFileObjectInput contains parameters for overwriting an existing file object
type OverwriteFileObjectInput struct {
	Kind    FileObjectKind
	ID      FileObjectID
	Slug    string
	Ext     string
	Data    []byte
	IfMatch *FileObjectMatchCondition
}

// ScanFileObjectsOptions specifies options for scanning file objects
type ScanFileObjectsOptions struct {
	Kind  *FileObjectKind
	Scope FileObjectScope
}

// FileStore is the contract interface for file object management
type FileStore interface {
	// CreateFileObject creates a new file object
	CreateFileObject(ctx context.Context, input CreateFileObjectInput) (FileObjectLocation, error)

	// OverwriteFileObject overwrites an existing file object
	OverwriteFileObject(ctx context.Context, input OverwriteFileObjectInput) (FileObjectLocation, error)

	// ReadFileObject reads a file object's content and location
	ReadFileObject(ctx context.Context, lookup FileObjectLookup, scope FileObjectScope) ([]byte, FileObjectLocation, error)

	// FindFileObject finds a file object's location without reading content
	FindFileObject(ctx context.Context, lookup FileObjectLookup, scope FileObjectScope) (*FileObjectLocation, error)

	// FileObjectExists checks if a file object exists
	FileObjectExists(ctx context.Context, lookup FileObjectLookup, scope FileObjectScope) (bool, error)

	// MoveFileObjectToTrash moves a file object from active to trash
	MoveFileObjectToTrash(ctx context.Context, ref FileObjectRef) (FileObjectLocation, error)

	// RestoreFileObjectFromTrash restores a file object from trash to active
	RestoreFileObjectFromTrash(ctx context.Context, ref FileObjectRef) (FileObjectLocation, error)

	// PermanentlyDeleteFileObject permanently deletes a file object from trash
	PermanentlyDeleteFileObject(ctx context.Context, ref FileObjectRef) error

	// RenameFileObjectSlug renames the slug portion of a file object's name
	RenameFileObjectSlug(ctx context.Context, ref FileObjectRef, newSlug string) (FileObjectLocation, error)

	// ScanFileObjects scans for file objects matching the given options
	ScanFileObjects(ctx context.Context, options ScanFileObjectsOptions) (<-chan FileObjectLocation, <-chan error)
}
