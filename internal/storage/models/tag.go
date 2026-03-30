package models

import "time"

// Tag 面向展示与筛选：颜色、图标、slug 直接服务前端组件。
type Tag struct {
	ID          string `gorm:"column:id;type:text;primaryKey"`
	Name        string `gorm:"column:name;type:text;not null;uniqueIndex:idx_tag_name"`
	Slug        string `gorm:"column:slug;type:text;not null;uniqueIndex:idx_tag_slug"`
	Color       string `gorm:"column:color;type:text"` // 如 #RRGGBB 或设计 token
	Icon        string `gorm:"column:icon;type:text"`  // emoji 或图标 key
	Description string `gorm:"column:description;type:text"`
	SortOrder   int    `gorm:"column:sort_order;not null;default:0"`
	// CardCount 可选缓存计数，列表页展示「标签云」时避免每次 COUNT。
	CardCount int `gorm:"column:card_count;not null;default:0"`

	CreatedAt time.Time `gorm:"column:created_at;not null"`
	UpdatedAt time.Time `gorm:"column:updated_at;not null"`
	SoftDelete

	CardTags []CardTag `gorm:"foreignKey:TagID"`
}

func (Tag) TableName() string {
	return "tag"
}

// CardTag 显式关联表，便于后续增加权重、来源等字段而不破坏 GORM 约定。
type CardTag struct {
	CardID string `gorm:"column:card_id;type:text;primaryKey;index:idx_card_tag_card"`
	TagID  string `gorm:"column:tag_id;type:text;primaryKey;index:idx_card_tag_tag"`
	// 联合唯一：(card_id, tag_id) 由复合主键保证。
	CreatedAt time.Time `gorm:"column:created_at;not null"`

	Card *Card `gorm:"foreignKey:CardID"`
	Tag  *Tag  `gorm:"foreignKey:TagID"`
}

func (CardTag) TableName() string {
	return "card_tag"
}
