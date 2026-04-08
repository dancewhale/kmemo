package models

import "time"

// Asset 表示渲染所需的媒体与附件（图/音/视频/封面等）。
//
// 文件路径由 file 模块根据 ID + Slug + Ext 计算得出，不存储在数据库中。
// 这样设计的原因：
// 1) 数据目录可以整体迁移、云同步或换机，路径由运行时计算；
// 2) 多平台（macOS/Windows）路径格式不同，派生路径策略保证可移植；
// 3) 备份与打包时只需保留文件对象本身，路径可重建。
type Asset struct {
	ID     string `gorm:"column:id;type:text;primaryKey"`
	CardID string `gorm:"column:card_id;type:text;not null;index:idx_asset_card;index:idx_asset_card_kind"`
	// Kind: image / audio / video / file / cover / thumbnail，供 UI 选图标与播放器。
	Kind string `gorm:"column:kind;type:text;not null;index:idx_asset_kind;index:idx_asset_card_kind"`
	// Slug 用于可读文件名（可选）
	Slug string `gorm:"column:slug;type:text"`
	// Ext 文件扩展名
	Ext string `gorm:"column:ext;type:text;not null"`
	// MimeType MIME 类型
	MimeType string `gorm:"column:mime_type;type:text"`
	// Size 文件大小（字节）
	Size int64 `gorm:"column:size;type:integer"`
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
