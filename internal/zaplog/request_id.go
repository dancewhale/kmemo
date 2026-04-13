package zaplog

import "github.com/google/uuid"

const MetadataRequestIDKey = "x-request-id"

func NewRequestID() string {
    uid, genErr := uuid.NewV7()
    if genErr != nil {
            return ""
    }
    id := uid.String()
	return id
}
