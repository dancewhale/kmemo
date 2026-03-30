package models

import "time"

// Knowledge 对应「知识库 / 牌组」层级，贴近 Anki deck 的展示与组织方式。
// 自引用 parent_id 表达树形导航；fsrs_preset_id 绑定 FSRS 参数集供下属卡片继承或覆盖。
type Knowledge struct {
	ID           string `gorm:"column:id;type:text;primaryKey"`
	Name         string `gorm:"column:name;type:text;not null"`
	Description  string `gorm:"column:description;type:text"`
	ParentID     *string `gorm:"column:parent_id;type:text;index"`
	FSRSPresetID *string `gorm:"column:fsrs_preset_id;type:text;index"`
	CreatedAt    time.Time `gorm:"column:created_at;not null"`
	UpdatedAt    time.Time `gorm:"column:updated_at;not null"`
	ArchivedAt   *time.Time `gorm:"column:archived_at"`

	Parent   *Knowledge  `gorm:"foreignKey:ParentID"`
	Children []Knowledge `gorm:"foreignKey:ParentID"`

	FSRSPreset *FSRSParameter `gorm:"foreignKey:FSRSPresetID"`

	SourceDocuments []SourceDocument `gorm:"foreignKey:KnowledgeID"`
	Cards           []Card           `gorm:"foreignKey:KnowledgeID"`
}

func (Knowledge) TableName() string {
	return "knowledge"
}

// SourceDocument 记录内容从何处导入，便于在 UI 展示来源、回溯原文路径或 URL。
type SourceDocument struct {
	ID          string `gorm:"column:id;type:text;primaryKey"`
	KnowledgeID string `gorm:"column:knowledge_id;type:text;not null;index"`
	// SourceType: pdf / epub / text / html 等，供界面显示图标与打开方式。
	SourceType string `gorm:"column:source_type;type:text;not null;index"`
	Title      string `gorm:"column:title;type:text"`
	Author     string `gorm:"column:author;type:text"`
	// OriginalURI 可为 URL、本地路径或外部系统 ID，用于展示「出处」链接。
	OriginalURI string `gorm:"column:original_uri;type:text"`
	// OriginalHash 用于检测来源是否变更、是否需重新转换。
	OriginalHash string `gorm:"column:original_hash;type:text"`
	// FilePath 指向已转换的主 HTML 等资源，供渲染器直接加载。
	FilePath  string `gorm:"column:file_path;type:text"`
	CreatedAt time.Time `gorm:"column:created_at;not null"`
	UpdatedAt time.Time `gorm:"column:updated_at;not null"`

	Knowledge *Knowledge `gorm:"foreignKey:KnowledgeID"`
	Cards     []Card       `gorm:"foreignKey:SourceDocumentID"`
}

func (SourceDocument) TableName() string {
	return "source_document"
}
