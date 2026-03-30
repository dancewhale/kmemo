package models

import "time"

// Asset 表示渲染所需的媒体与附件（图/音/视频/封面等）。
//
// storage_path 存相对路径的原因：
// 1) 用户数据目录可能整体迁移、云同步或换机，绝对路径会立即失效；
// 2) 多平台（macOS/Windows）路径格式不同，相对应用配置的「数据根」可移植；
// 3) 备份与打包时只需保留根目录下的相对树结构。
// 应用层解析时：filepath.Join(dataRoot, asset.StoragePath)。
type Asset struct {
	ID     string `gorm:"column:id;type:text;primaryKey"`
	CardID string `gorm:"column:card_id;type:text;not null;index:idx_asset_card;index:idx_asset_card_kind"`
	// Kind: image / audio / video / file / cover / thumbnail，供 UI 选图标与播放器。
	Kind string `gorm:"column:kind;type:text;not null;index:idx_asset_kind;index:idx_asset_card_kind"`
	// StoragePath 相对数据根目录，见包注释。
	StoragePath string `gorm:"column:storage_path;type:text;not null"`
	// Checksum 用于去重与完整性校验；索引用以快速查找重复资源。
	Checksum string `gorm:"column:checksum;type:text;index:idx_asset_checksum"`
	// Status: active / missing / deleted — missing 表示文件丢失但仍保留 DB 记录以便修复。
	Status string `gorm:"column:status;type:text;not null;default:active;index:idx_asset_status"`

	CreatedAt time.Time `gorm:"column:created_at;not null"`
	UpdatedAt time.Time `gorm:"column:updated_at;not null"`
	SoftDelete

	Card *Card `gorm:"foreignKey:CardID"`
}

func (Asset) TableName() string {
	return "asset"
}
