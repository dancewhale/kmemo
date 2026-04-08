package file

import (
	"fmt"
	"path/filepath"
	"strings"
)

// ParsedFileName contains the parsed components of a file name
type ParsedFileName struct {
	ID   FileObjectID
	Slug string
	Ext  string
}

// BuildObjectFileName builds a file name according to the kind's naming style
func BuildObjectFileName(kindSpec FileObjectKindSpec, id FileObjectID, slug string, ext string) string {
	// Normalize slug
	normalizedSlug := NormalizeSlug(slug)

	// Clean extension (remove leading dot if present)
	ext = strings.TrimPrefix(ext, ".")

	switch kindSpec.FileNameStyle {
	case FileNameStyleWithSlug:
		// Format: <uuid>--<slug>.<ext>
		return fmt.Sprintf("%s--%s.%s", id, normalizedSlug, ext)
	case FileNameStyleWithoutSlug:
		// Format: <uuid>.<ext>
		return fmt.Sprintf("%s.%s", id, ext)
	default:
		return ""
	}
}

// ParseObjectFileName parses a file name according to the kind's naming style
func ParseObjectFileName(kindSpec FileObjectKindSpec, fileName string) (ParsedFileName, bool) {
	var parsed ParsedFileName

	// Extract extension
	ext := filepath.Ext(fileName)
	if ext == "" {
		return parsed, false
	}
	parsed.Ext = strings.TrimPrefix(ext, ".")

	// Remove extension from name
	nameWithoutExt := strings.TrimSuffix(fileName, ext)

	switch kindSpec.FileNameStyle {
	case FileNameStyleWithSlug:
		// Format: <uuid>--<slug>
		parts := strings.SplitN(nameWithoutExt, "--", 2)
		if len(parts) != 2 {
			return parsed, false
		}

		// Validate UUID
		if !ValidateUUIDv7(parts[0]) {
			return parsed, false
		}

		parsed.ID = FileObjectID(parts[0])
		parsed.Slug = parts[1]
		return parsed, true

	case FileNameStyleWithoutSlug:
		// Format: <uuid>
		if !ValidateUUIDv7(nameWithoutExt) {
			return parsed, false
		}

		parsed.ID = FileObjectID(nameWithoutExt)
		parsed.Slug = ""
		return parsed, true

	default:
		return parsed, false
	}
}

// BuildExpectedObjectFilePath builds the expected full file path for an object
func BuildExpectedObjectFilePath(rootDir string, kindSpec FileObjectKindSpec, id FileObjectID, slug string, ext string, scope FileObjectScope) string {
	dir := BuildObjectDirectory(rootDir, kindSpec, id, scope)
	fileName := BuildObjectFileName(kindSpec, id, slug, ext)
	return filepath.Join(dir, fileName)
}
