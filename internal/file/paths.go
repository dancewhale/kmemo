package file

import (
	"fmt"
	"path/filepath"
	"strings"
)

// BuildBucketPath calculates the bucket path from a UUID v7
// Example: 0195f3f0-a8b7-7c8d-b1b9-d45c8f23a1ee -> 019/5/f/3
func BuildBucketPath(id FileObjectID) string {
	// Remove hyphens to get hex string
	hex := strings.ReplaceAll(string(id), "-", "")

	// Build 4-layer bucket: h1h2h3/h4/h5/h6
	return fmt.Sprintf("%s/%s/%s/%s",
		hex[0:3], // h1h2h3
		hex[3:4], // h4
		hex[4:5], // h5
		hex[5:6], // h6
	)
}

// BuildActiveObjectDirectory builds the full directory path for an active object
func BuildActiveObjectDirectory(rootDir string, kindSpec FileObjectKindSpec, id FileObjectID) string {
	bucket := BuildBucketPath(id)
	return filepath.Join(rootDir, kindSpec.DirectoryName, bucket)
}

// BuildTrashObjectDirectory builds the full directory path for a trash object
func BuildTrashObjectDirectory(rootDir string, kindSpec FileObjectKindSpec, id FileObjectID) string {
	bucket := BuildBucketPath(id)
	return filepath.Join(rootDir, "trash", kindSpec.DirectoryName, bucket)
}

// BuildLockDirectoryPath builds the lock directory path for an object
func BuildLockDirectoryPath(rootDir string, kindSpec FileObjectKindSpec, id FileObjectID) string {
	bucket := BuildBucketPath(id)
	return filepath.Join(rootDir, "locks", kindSpec.DirectoryName, bucket, string(id)+".lock")
}

// BuildObjectDirectory builds the directory path for an object based on scope
func BuildObjectDirectory(rootDir string, kindSpec FileObjectKindSpec, id FileObjectID, scope FileObjectScope) string {
	switch scope {
	case FileObjectScopeActive:
		return BuildActiveObjectDirectory(rootDir, kindSpec, id)
	case FileObjectScopeTrash:
		return BuildTrashObjectDirectory(rootDir, kindSpec, id)
	default:
		return ""
	}
}
