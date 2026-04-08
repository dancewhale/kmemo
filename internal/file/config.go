package file

import "time"

// FileStoreConfig contains configuration for the file store
type FileStoreConfig struct {
	RootDir   string
	Kinds     []FileObjectKindSpec
	LockWait  time.Duration
	LockRetry time.Duration
}

// KmemoKinds defines the default file object kinds for kmemo
var KmemoKinds = []FileObjectKindSpec{
	{Kind: FileObjectKindCard, DirectoryName: "cards", FileNameStyle: FileNameStyleWithSlug},
	{Kind: FileObjectKindAsset, DirectoryName: "assets", FileNameStyle: FileNameStyleWithoutSlug},
	{Kind: FileObjectKindSource, DirectoryName: "sources", FileNameStyle: FileNameStyleWithSlug},
}

// DefaultConfig returns the default configuration for kmemo
func DefaultConfig(vaultDir string) FileStoreConfig {
	return FileStoreConfig{
		RootDir:   vaultDir,
		Kinds:     KmemoKinds,
		LockWait:  30 * time.Second,
		LockRetry: 100 * time.Millisecond,
	}
}
