package models

import (
	"time"

	"gorm.io/gorm"
)

// Timestamps 为需要 created_at / updated_at 的表提供统一嵌入字段。
// GORM 会在写入时维护 UpdatedAt。
type Timestamps struct {
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
}

// SoftDelete 嵌入到 Card / Asset / Tag 等模型末尾，以启用 GORM 软删除。
// Knowledge、SourceDocument 等用业务字段归档，不使用本嵌入。
type SoftDelete struct {
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
