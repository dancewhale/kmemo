package repository

import (
	"context"
	"kmemo/internal/storage/dao"
	"kmemo/internal/storage/models"
	"time"

	"gorm.io/gorm"
)

// SRSRepository SRS 调度仓储
type SRSRepository interface {
	Repository

	CreateOrUpdate(ctx context.Context, srs *models.CardSRS) error
	GetByCardID(ctx context.Context, cardID string) (*models.CardSRS, error)
	GetDueCards(ctx context.Context, opts DueCardsOptions) ([]*models.CardSRS, error)
	UpdateAfterReview(ctx context.Context, cardID string, srs *models.CardSRS, log *models.ReviewLog) error
	GetStatistics(ctx context.Context, knowledgeID *string) (*SRSStatistics, error)

	// CountDueCardsByKnowledgeIDs 按知识库统计「当前到期」的 SRS 条目数（与 GetDueCards 过滤条件一致；未出现在结果中的 knowledge_id 视为 0）
	CountDueCardsByKnowledgeIDs(ctx context.Context, knowledgeIDs []string) (map[string]int64, error)
	Suspend(ctx context.Context, cardID string) error
	Resume(ctx context.Context, cardID string) error
	UndoLastReview(ctx context.Context, cardID string) error
	GetLastReviewLog(ctx context.Context, cardID string) (*models.ReviewLog, error)
}

type DueCardsOptions struct {
	KnowledgeID *string
	Limit       int
	States      []string
}

type SRSStatistics struct {
	NewCount        int
	LearningCount   int
	ReviewCount     int
	RelearningCount int
	TotalCards      int
	DueToday        int
}

type srsRepo struct {
	db *gorm.DB
}

func NewSRSRepository(db *gorm.DB) SRSRepository {
	return &srsRepo{db: db}
}

func (r *srsRepo) WithTx(tx *gorm.DB) Repository {
	return &srsRepo{db: tx}
}

func (r *srsRepo) CreateOrUpdate(ctx context.Context, srs *models.CardSRS) error {
	existing, err := dao.Use(r.db).CardSRS.WithContext(ctx).Where(dao.CardSRS.CardID.Eq(srs.CardID)).First()
	if err != nil && err != gorm.ErrRecordNotFound {
		return convertError(err)
	}

	if existing != nil {
		srs.UpdatedAt = time.Now()
		_, err = dao.Use(r.db).CardSRS.WithContext(ctx).Where(dao.CardSRS.CardID.Eq(srs.CardID)).Updates(srs)
		return convertError(err)
	}

	return dao.Use(r.db).CardSRS.WithContext(ctx).Create(srs)
}

func (r *srsRepo) GetByCardID(ctx context.Context, cardID string) (*models.CardSRS, error) {
	srs, err := dao.Use(r.db).CardSRS.WithContext(ctx).Where(dao.CardSRS.CardID.Eq(cardID)).First()
	return srs, convertError(err)
}

func (r *srsRepo) GetDueCards(ctx context.Context, opts DueCardsOptions) ([]*models.CardSRS, error) {
	q := dao.Use(r.db).CardSRS.WithContext(ctx).Where(dao.CardSRS.Suspended.Is(false))

	if opts.KnowledgeID != nil {
		q = q.Joins(dao.CardSRS.Card).Where(dao.Card.KnowledgeID.Eq(*opts.KnowledgeID))
	}

	if len(opts.States) > 0 {
		q = q.Where(dao.CardSRS.FSRSState.In(opts.States...))
	}

	now := time.Now()
	q = q.Where(dao.CardSRS.DueAt.Lte(now))

	if opts.Limit > 0 {
		q = q.Limit(opts.Limit)
	}

	list, err := q.Find()
	return list, convertError(err)
}

func (r *srsRepo) UpdateAfterReview(ctx context.Context, cardID string, srs *models.CardSRS, log *models.ReviewLog) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		srs.UpdatedAt = time.Now()
		if _, err := dao.Use(tx).CardSRS.WithContext(ctx).Where(dao.CardSRS.CardID.Eq(cardID)).Updates(srs); err != nil {
			return convertError(err)
		}
		return dao.Use(tx).ReviewLog.WithContext(ctx).Create(log)
	})
}

func (r *srsRepo) CountDueCardsByKnowledgeIDs(ctx context.Context, knowledgeIDs []string) (map[string]int64, error) {
	if len(knowledgeIDs) == 0 {
		return map[string]int64{}, nil
	}
	now := time.Now()
	var rows []struct {
		KnowledgeID string `gorm:"column:knowledge_id"`
		N           int64  `gorm:"column:n"`
	}
	err := r.db.WithContext(ctx).
		Table("card_srs").
		Select("card.knowledge_id as knowledge_id, count(*) as n").
		Joins("JOIN card ON card.id = card_srs.card_id AND card.deleted_at IS NULL").
		Where("card.knowledge_id IN ?", knowledgeIDs).
		Where("card_srs.suspended = ?", false).
		Where("card_srs.due_at <= ?", now).
		Group("card.knowledge_id").
		Scan(&rows).Error
	if err != nil {
		return nil, convertError(err)
	}
	out := make(map[string]int64, len(rows))
	for _, row := range rows {
		out[row.KnowledgeID] = row.N
	}
	return out, nil
}

func (r *srsRepo) GetStatistics(ctx context.Context, knowledgeID *string) (*SRSStatistics, error) {
	q := dao.Use(r.db).CardSRS.WithContext(ctx)

	if knowledgeID != nil {
		q = q.Joins(dao.CardSRS.Card).Where(dao.Card.KnowledgeID.Eq(*knowledgeID))
	}

	total, err := q.Count()
	if err != nil {
		return nil, convertError(err)
	}
	newCount, err := q.Where(dao.CardSRS.FSRSState.Eq("new")).Count()
	if err != nil {
		return nil, convertError(err)
	}
	learningCount, err := q.Where(dao.CardSRS.FSRSState.Eq("learning")).Count()
	if err != nil {
		return nil, convertError(err)
	}
	reviewCount, err := q.Where(dao.CardSRS.FSRSState.Eq("review")).Count()
	if err != nil {
		return nil, convertError(err)
	}
	relearningCount, err := q.Where(dao.CardSRS.FSRSState.Eq("relearning")).Count()
	if err != nil {
		return nil, convertError(err)
	}

	now := time.Now()
	dueToday, err := q.Where(dao.CardSRS.DueAt.Lte(now)).Count()
	if err != nil {
		return nil, convertError(err)
	}

	return &SRSStatistics{
		NewCount:        int(newCount),
		LearningCount:   int(learningCount),
		ReviewCount:     int(reviewCount),
		RelearningCount: int(relearningCount),
		TotalCards:      int(total),
		DueToday:        int(dueToday),
	}, nil
}

func (r *srsRepo) Suspend(ctx context.Context, cardID string) error {
	_, err := dao.Use(r.db).CardSRS.WithContext(ctx).Where(dao.CardSRS.CardID.Eq(cardID)).Update(dao.CardSRS.Suspended, true)
	return convertError(err)
}

func (r *srsRepo) Resume(ctx context.Context, cardID string) error {
	_, err := dao.Use(r.db).CardSRS.WithContext(ctx).Where(dao.CardSRS.CardID.Eq(cardID)).Update(dao.CardSRS.Suspended, false)
	return convertError(err)
}

func (r *srsRepo) GetLastReviewLog(ctx context.Context, cardID string) (*models.ReviewLog, error) {
	log, err := dao.Use(r.db).ReviewLog.WithContext(ctx).Where(
		dao.ReviewLog.CardID.Eq(cardID),
	).Order(dao.ReviewLog.ReviewedAt.Desc()).First()
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return log, convertError(err)
}

func (r *srsRepo) UndoLastReview(ctx context.Context, cardID string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		log, err := dao.Use(tx).ReviewLog.WithContext(ctx).Where(
			dao.ReviewLog.CardID.Eq(cardID),
		).Order(dao.ReviewLog.ReviewedAt.Desc()).First()
		if err != nil {
			return convertError(err)
		}

		updates := map[string]interface{}{
			"fsrs_state":      log.ReviewKind,
			"stability":       log.PreviousStability,
			"difficulty":      log.PreviousDifficulty,
			"last_review_at":  nil,
			"updated_at":      time.Now(),
		}
		if _, err := dao.Use(tx).CardSRS.WithContext(ctx).Where(dao.CardSRS.CardID.Eq(cardID)).Updates(updates); err != nil {
			return convertError(err)
		}

		_, err = dao.Use(tx).ReviewLog.WithContext(ctx).Where(dao.ReviewLog.ID.Eq(log.ID)).Delete()
		return convertError(err)
	})
}
