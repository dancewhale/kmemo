package file

import "errors"

var (
	ErrInvalidKind                    = errors.New("invalid file object kind")
	ErrInvalidUUIDv7                  = errors.New("invalid uuid v7")
	ErrInvalidFileExtension           = errors.New("invalid file extension")
	ErrInvalidFileName                = errors.New("invalid file name")
	ErrInvalidInput                   = errors.New("invalid input")

	ErrFileObjectAlreadyExists        = errors.New("file object already exists")
	ErrFileObjectAlreadyExistsInTrash = errors.New("file object already exists in trash")
	ErrFileObjectNotFound             = errors.New("file object not found")

	ErrDuplicateFileObjectID          = errors.New("duplicate file object id")
	ErrFileExtensionMismatch          = errors.New("file extension mismatch")
	ErrConcurrentModification         = errors.New("concurrent modification")

	ErrTrashConflict                  = errors.New("trash conflict")
	ErrRestoreConflict                = errors.New("restore conflict")
	ErrUnsupportedOperation           = errors.New("unsupported operation")

	ErrLockTimeout                    = errors.New("lock timeout")
	ErrIO                             = errors.New("io error")
)
