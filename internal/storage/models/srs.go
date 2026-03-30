package models

import "time"

// CardSRS 存储单张卡片的 FSRS 调度状态；与业务 Card 一对一。
type CardSRS struct {
	CardID string `gorm:"column:card_id;type:text;primaryKey"`
	// FSRSState: new / learning / review / relearning
	FSRSState string `gorm:"column:fsrs_state;type:text;not null;index"`
	DueAt     *time.Time `gorm:"column:due_at;index"`
	LastReviewAt *time.Time `gorm:"column:last_review_at"`

	Stability     *float64 `gorm:"column:stability"`
	Difficulty    *float64 `gorm:"column:difficulty"`
	ElapsedDays   *float64 `gorm:"column:elapsed_days"`
	ScheduledDays *float64 `gorm:"column:scheduled_days"`
	Reps          int      `gorm:"column:reps;not null;default:0"`
	Lapses        int      `gorm:"column:lapses;not null;default:0"`

	LearningStepsJSON *string `gorm:"column:learning_steps_json;type:text"`
	Suspended         bool    `gorm:"column:suspended;not null;default:false"`

	CreatedAt time.Time `gorm:"column:created_at;not null"`
	UpdatedAt time.Time `gorm:"column:updated_at;not null"`

	Card *Card `gorm:"foreignKey:CardID"`
}

func (CardSRS) TableName() string {
	return "card_srs"
}

// FSRSParameter 全局或命名预设，parameters_json 存 FSRS 权重等。
type FSRSParameter struct {
	ID              string `gorm:"column:id;type:text;primaryKey"`
	Name            string `gorm:"column:name;type:text;not null"`
	ParametersJSON  string `gorm:"column:parameters_json;type:text;not null"`
	DesiredRetention *float64 `gorm:"column:desired_retention"`
	MaximumInterval *int     `gorm:"column:maximum_interval"`

	CreatedAt time.Time `gorm:"column:created_at;not null"`
	UpdatedAt time.Time `gorm:"column:updated_at;not null"`
}

func (FSRSParameter) TableName() string {
	return "fsrs_parameter"
}

// ReviewLog 追加式复习历史，供统计与调试 FSRS；snapshot_json 存当时完整快照。
//
// 为便于搜索/分析服务按时间增量拉取，复合索引 idx_review_card_time 覆盖 (card_id, reviewed_at)。
// 主键 id 使用应用层生成的 ULID/UUID，保证全局唯一。
type ReviewLog struct {
	ID       string `gorm:"column:id;type:text;primaryKey"`
	CardID   string `gorm:"column:card_id;type:text;not null;index:idx_review_card_time"`
	ReviewedAt time.Time `gorm:"column:reviewed_at;not null;index:idx_review_card_time"`
	Rating       int    `gorm:"column:rating;not null"`
	ReviewKind   string `gorm:"column:review_kind;type:text;not null"` // learn / review / relearn 等
	ScheduledDays float64 `gorm:"column:scheduled_days"`
	ElapsedDays   float64 `gorm:"column:elapsed_days"`
	PreviousStability  *float64 `gorm:"column:previous_stability"`
	PreviousDifficulty *float64 `gorm:"column:previous_difficulty"`
	NewStability       *float64 `gorm:"column:new_stability"`
	NewDifficulty      *float64 `gorm:"column:new_difficulty"`
	SnapshotJSON       string   `gorm:"column:snapshot_json;type:text;not null;default:'{}'"`

	Card *Card `gorm:"foreignKey:CardID"`
}

func (ReviewLog) TableName() string {
	return "review_log"
}
