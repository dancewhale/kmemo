package zaplog

import "github.com/google/uuid"

const MetadataRequestIDKey = "x-request-id"

func NewRequestID() string {
	return uuid.NewString()
}
