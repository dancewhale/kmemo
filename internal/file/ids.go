package file

import (
	"fmt"
	"regexp"
	"strings"
)

var uuidv7Regex = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-7[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)

// NormalizeUUIDv7 normalizes a UUID v7 string to lowercase standard format
func NormalizeUUIDv7(id string) (FileObjectID, error) {
	normalized := strings.ToLower(strings.TrimSpace(id))
	if !ValidateUUIDv7(normalized) {
		return "", fmt.Errorf("%w: %s", ErrInvalidUUIDv7, id)
	}
	return FileObjectID(normalized), nil
}

// ValidateUUIDv7 checks if a string is a valid UUID v7
func ValidateUUIDv7(id string) bool {
	return uuidv7Regex.MatchString(id)
}

// NormalizeSlug normalizes a slug according to kmemo rules
func NormalizeSlug(slug string) string {
	// Trim whitespace
	slug = strings.TrimSpace(slug)

	// Convert to lowercase
	slug = strings.ToLower(slug)

	// Replace whitespace with hyphens
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, "\t", "-")
	slug = strings.ReplaceAll(slug, "\n", "-")

	// Remove dangerous characters (keep only alphanumeric, hyphen, underscore, dot)
	var result strings.Builder
	for _, r := range slug {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' || r == '_' || r == '.' {
			result.WriteRune(r)
		}
	}
	slug = result.String()

	// Merge consecutive hyphens
	for strings.Contains(slug, "--") {
		slug = strings.ReplaceAll(slug, "--", "-")
	}

	// Trim hyphens from edges
	slug = strings.Trim(slug, "-")

	// Limit length to 64 characters
	if len(slug) > 64 {
		slug = slug[:64]
		slug = strings.TrimRight(slug, "-")
	}

	// Use "untitled" if empty
	if slug == "" {
		slug = "untitled"
	}

	return slug
}
