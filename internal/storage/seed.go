package storage

import (
	"fmt"
	"time"

	"gorm.io/gorm"

	"kmemo/internal/storage/models"
)

// DefaultSeedKnowledgeID 为首次初始化写入的默认知识库主键，便于排查与支持引用。
// 与随机 UUID 并存于表中无冲突（用户自建知识库通常不会使用该 nil-UUID）。
const DefaultSeedKnowledgeID = "00000000-0000-0000-0000-000000000001"

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
	row := &models.Knowledge{
		ID:           DefaultSeedKnowledgeID,
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
