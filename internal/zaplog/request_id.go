package zaplog

import "github.com/google/uuid"

const MetadataRequestIDKey = "x-request-id"

func NewRequestID() string {
	// Use a dedicated request-id namespace to avoid confusion with domain IDs
	// (knowledge/card IDs are plain UUIDs in this project).
	return "req_" + uuid.NewString()
}
