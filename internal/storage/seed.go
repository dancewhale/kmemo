package storage

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"kmemo/internal/storage/models"
)

// newSeedKnowledgeID 生成首次初始化知识库主键（UUID v7）。
func newSeedKnowledgeID() (string, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return "", fmt.Errorf("storage: generate seed knowledge id: %w", err)
	}
	return id.String(), nil
}

// SeedDefaultData 在空库时写入最小可用业务数据；可重复调用，已存在知识库记录时立即返回。
func SeedDefaultData(db *gorm.DB) error {
	if db == nil {
		return fmt.Errorf("storage: seed default data: db is nil")
	}
	var n int64
	if err := db.Model(&models.Knowledge{}).Count(&n).Error; err != nil {
		return fmt.Errorf("storage: seed count knowledge: %w", err)
	}
	if n > 0 {
		return nil
	}

	now := time.Now().UTC()
	id, err := newSeedKnowledgeID()
	if err != nil {
		return err
	}
	row := &models.Knowledge{
		ID:           id,
		Name:         "默认知识库",
		Description:  "",
		ParentID:     nil,
		FSRSPresetID: nil,
		CreatedAt:    now,
		UpdatedAt:    now,
		ArchivedAt:   nil,
	}
	if err := db.Create(row).Error; err != nil {
		return fmt.Errorf("storage: seed default knowledge: %w", err)
	}
	return nil
}
