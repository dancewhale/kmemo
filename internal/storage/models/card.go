package models

import "time"

// Card 存储面向渲染的主内容：路径、类型、哈希等均以「如何显示」为第一考量。
type Card struct {
	ID               string `gorm:"column:id;type:text;primaryKey"`
	KnowledgeID      string `gorm:"column:knowledge_id;type:text;not null;index:idx_card_knowledge"`
	SourceDocumentID *string `gorm:"column:source_document_id;type:text;index"`
	ParentID         *string `gorm:"column:parent_id;type:text;index"`
	Title            string `gorm:"column:title;type:text;not null;default:''"`
	SortOrder        int    `gorm:"column:sort_order;not null;default:0"`
	// Path 物化路径，便于列表/面包屑展示而无需递归查询。
	Path string `gorm:"column:path;type:text;not null;default:'';index:idx_card_path"`
	// CardType: article / excerpt / qa / cloze / note，驱动前端模板与交互。
	CardType string `gorm:"column:card_type;type:text;not null;index"`
	// Slug 从 Title 派生，用于可读文件名
	Slug string `gorm:"column:slug;type:text;not null;default:''"`
	// HTMLHash / AnswerHTMLHash 用于内容完整性校验
	HTMLHash       string  `gorm:"column:html_hash;type:text;not null;default:''"`
	AnswerHTMLHash *string `gorm:"column:answer_html_hash;type:text"`
	SourceRef      *string `gorm:"column:source_ref;type:text"`
	// Status: active / suspended，控制是否参与学习与展示。
	Status string `gorm:"column:status;type:text;not null;default:active;index"`
	IsRoot    bool `gorm:"column:is_root;not null;default:false"`
	IsExtract bool `gorm:"column:is_extract;not null;default:false"`
	// CreateFromCardID 标识从哪张卡拆分/摘录而来，用于 UI 溯源。
	CreateFromCardID *string `gorm:"column:create_from_card_id;type:text;index"`
	// CoverAssetID 可选封面，一对一引用 Asset（与正文媒体区分）。
	CoverAssetID *string `gorm:"column:cover_asset_id;type:text;index"`

	CreatedAt time.Time `gorm:"column:created_at;not null"`
	UpdatedAt time.Time `gorm:"column:updated_at;not null"`
	SoftDelete

	Knowledge      *Knowledge      `gorm:"foreignKey:KnowledgeID"`
	SourceDocument *SourceDocument `gorm:"foreignKey:SourceDocumentID"`
	Parent         *Card           `gorm:"foreignKey:ParentID"`
	Children       []Card          `gorm:"foreignKey:ParentID"`

	CoverAsset *Asset  `gorm:"foreignKey:CoverAssetID"`
	Assets     []Asset `gorm:"foreignKey:CardID"`

	// CardTag 显式关联；查标签请通过 CardTag 再 Preload Tag。
	CardTags []CardTag `gorm:"foreignKey:CardID"`

	SRS *CardSRS `gorm:"foreignKey:CardID"`
}

func (Card) TableName() string {
	return "card"
}
