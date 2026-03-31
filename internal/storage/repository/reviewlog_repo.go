package repository

import (
	"context"
	"kmemo/internal/storage/dao"
	"kmemo/internal/storage/models"
	"time"

	"gorm.io/gorm"
)

// ReviewLogRepository 复习日志仓储
type ReviewLogRepository interface {
	Repository

	Create(ctx context.Context, log *models.ReviewLog) error
	GetByID(ctx context.Context, id string) (*models.ReviewLog, error)
	GetByCardID(ctx context.Context, cardID string, opts ListReviewLogOptions) ([]*models.ReviewLog, int64, error)
	GetLastByCardID(ctx context.Context, cardID string) (*models.ReviewLog, error)
	Delete(ctx context.Context, id string) error
	GetStatsByDateRange(ctx context.Context, startDate, endDate time.Time, knowledgeID *string) (*ReviewStatistics, error)
}

type ListReviewLogOptions struct {
	StartDate *time.Time
	EndDate   *time.Time
	Rating    *int
	Limit     int
	Offset    int
	OrderDesc bool
}

type ReviewStatistics struct {
	TotalReviews int
	AvgRating    float64
	RatingCounts map[int]int
	ReviewsByDay map[string]int
}

type reviewLogRepo struct {
	db *gorm.DB
}

func NewReviewLogRepository(db *gorm.DB) ReviewLogRepository {
	return &reviewLogRepo{db: db}
}

func (r *reviewLogRepo) WithTx(tx *gorm.DB) Repository {
	return &reviewLogRepo{db: tx}
}

func (r *reviewLogRepo) Create(ctx context.Context, log *models.ReviewLog) error {
	return dao.Use(r.db).ReviewLog.WithContext(ctx).Create(log)
}

func (r *reviewLogRepo) GetByID(ctx context.Context, id string) (*models.ReviewLog, error) {
	log, err := dao.Use(r.db).ReviewLog.WithContext(ctx).Where(dao.ReviewLog.ID.Eq(id)).First()
	return log, convertError(err)
}

func (r *reviewLogRepo) GetByCardID(ctx context.Context, cardID string, opts ListReviewLogOptions) ([]*models.ReviewLog, int64, error) {
	q := dao.Use(r.db).ReviewLog.WithContext(ctx).Where(dao.ReviewLog.CardID.Eq(cardID))

	if opts.StartDate != nil {
		q = q.Where(dao.ReviewLog.ReviewedAt.Gte(*opts.StartDate))
	}
	if opts.EndDate != nil {
		q = q.Where(dao.ReviewLog.ReviewedAt.Lte(*opts.EndDate))
	}
	if opts.Rating != nil {
		q = q.Where(dao.ReviewLog.Rating.Eq(*opts.Rating))
	}

	count, err := q.Count()
	if err != nil {
		return nil, 0, convertError(err)
	}

	if opts.OrderDesc {
		q = q.Order(dao.ReviewLog.ReviewedAt.Desc())
	} else {
		q = q.Order(dao.ReviewLog.ReviewedAt)
	}

	if opts.Limit > 0 {
		q = q.Limit(opts.Limit)
	}
	if opts.Offset > 0 {
		q = q.Offset(opts.Offset)
	}

	list, err := q.Find()
	return list, count, convertError(err)
}

func (r *reviewLogRepo) GetLastByCardID(ctx context.Context, cardID string) (*models.ReviewLog, error) {
	log, err := dao.Use(r.db).ReviewLog.WithContext(ctx).Where(
		dao.ReviewLog.CardID.Eq(cardID),
	).Order(dao.ReviewLog.ReviewedAt.Desc()).First()
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return log, convertError(err)
}

func (r *reviewLogRepo) Delete(ctx context.Context, id string) error {
	_, err := dao.Use(r.db).ReviewLog.WithContext(ctx).Where(dao.ReviewLog.ID.Eq(id)).Delete()
	return convertError(err)
}

func (r *reviewLogRepo) GetStatsByDateRange(ctx context.Context, startDate, endDate time.Time, knowledgeID *string) (*ReviewStatistics, error) {
	q := dao.Use(r.db).ReviewLog.WithContext(ctx).Where(
		dao.ReviewLog.ReviewedAt.Between(startDate, endDate),
	)

	logs, err := q.Find()
	if err != nil {
		return nil, convertError(err)
	}

	stats := &ReviewStatistics{
		TotalReviews: len(logs),
		RatingCounts: make(map[int]int),
		ReviewsByDay: make(map[string]int),
	}

	totalRating := 0
	for _, log := range logs {
		totalRating += log.Rating
		stats.RatingCounts[log.Rating]++
		day := log.ReviewedAt.Format("2006-01-02")
		stats.ReviewsByDay[day]++
	}

	if stats.TotalReviews > 0 {
		stats.AvgRating = float64(totalRating) / float64(stats.TotalReviews)
	}

	return stats, nil
}
