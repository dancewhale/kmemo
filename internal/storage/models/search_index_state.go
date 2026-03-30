package models

import "time"

// SearchIndexState 记录全文索引流水线与实体内容之间的同步状态。
//
// 不强依赖外键指向 card/asset/tag 的原因：
// 1) object_type + object_id 是多态引用，若用 FK 需三张可空列或触发器，SQLite 上冗长；
// 2) 索引任务可能先于实体落库重试，或实体已删而索引行需保留失败信息；
// 3) 解耦后搜索服务可独立演进、批量清理，不级联删除业务数据。
//
// 状态机含义（index_status）：
// - pending: 已登记待建索引，或内容变更后待处理；
// - indexed: indexed_hash 与 content_hash 一致，索引与内容同步；
// - stale: 内容 hash 变化，索引落后，需要增量或全量重建；
// - failed: 上次索引失败，可结合 retry_count / last_error 退避重试。
type SearchIndexState struct {
	ObjectType string `gorm:"column:object_type;type:text;primaryKey;size:32"`
	ObjectID   string `gorm:"column:object_id;type:text;primaryKey"`
	// ContentHash 当前实体参与索引的内容指纹；变化时应置 stale 或 pending。
	ContentHash string `gorm:"column:content_hash;type:text;not null;default:''"`
	// IndexedHash 已成功写入搜索侧的指纹；与 ContentHash 比较判断是否 stale。
	IndexedHash string `gorm:"column:indexed_hash;type:text;not null;default:''"`
	IndexStatus string `gorm:"column:index_status;type:text;not null;default:pending;index"`
	LastIndexedAt *time.Time `gorm:"column:last_indexed_at"`
	LastError     *string    `gorm:"column:last_error;type:text"`
	RetryCount    int        `gorm:"column:retry_count;not null;default:0"`

	CreatedAt time.Time `gorm:"column:created_at;not null"`
	UpdatedAt time.Time `gorm:"column:updated_at;not null"`
}

func (SearchIndexState) TableName() string {
	return "search_index_state"
}
