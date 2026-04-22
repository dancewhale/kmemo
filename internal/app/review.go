package app

import (
	"fmt"
	"time"

	"go.uber.org/zap"

	"kmemo/internal/actions/review"
	"kmemo/internal/storage/models"
	"kmemo/internal/storage/repository"
	"kmemo/internal/zaplog"
)

type SRSDTO struct {
	CardID       string     `json:"cardId"`
	FSRSState    string     `json:"fsrsState"`
	DueAt        *time.Time `json:"dueAt"`
	LastReviewAt *time.Time `json:"lastReviewAt"`
	Stability    *float64   `json:"stability"`
	Difficulty   *float64   `json:"difficulty"`
	Reps         int        `json:"reps"`
	Lapses       int        `json:"lapses"`
}

type CardWithSRSDTO struct {
	Card CardDTO `json:"card"`
	SRS  SRSDTO  `json:"srs"`
}

type SRSStatisticsDTO struct {
	NewCount        int `json:"newCount"`
	LearningCount   int `json:"learningCount"`
	ReviewCount     int `json:"reviewCount"`
	RelearningCount int `json:"relearningCount"`
	TotalCards      int `json:"totalCards"`
	DueToday        int `json:"dueToday"`
}

type ReviewRequest struct {
	CardID string `json:"cardId"`
	Rating int    `json:"rating"`
}

type ReviewLogDTO struct {
	ID         string    `json:"id"`
	CardID     string    `json:"cardId"`
	ReviewedAt time.Time `json:"reviewedAt"`
	Rating     int       `json:"rating"`
	ReviewKind string    `json:"reviewKind"`
}

func (d *Desktop) GetDueCards(knowledgeID *string, limit int) ([]*CardWithSRSDTO, error) {
	ctx := d.actionContext()
	log := zaplog.L(ctx).Named("review.api")
	items, err := d.actions.Review.GetDueCards(ctx, knowledgeID, limit)
	if err != nil {
		log.Info("GetDueCards failed", zapOptionalString("knowledgeId", knowledgeID), zap.Int("limit", limit), zap.Error(err))
		return nil, err
	}
	result := make([]*CardWithSRSDTO, 0, len(items))
	for _, item := range items {
		result = append(result, &CardWithSRSDTO{Card: *toCardDTO(item, nil), SRS: *toSRSDTO(item.SRS)})
	}
	sample := make([]string, 0, min(5, len(result)))
	for i := range result {
		if i >= 5 {
			break
		}
		sample = append(sample, result[i].Card.ID)
	}
	log.Info("GetDueCards ok",
		zapOptionalString("knowledgeId", knowledgeID),
		zap.Int("limit", limit),
		zap.Int("returned", len(result)),
		zap.Strings("sample_card_ids", sample),
	)
	return result, nil
}

func (d *Desktop) SubmitReview(req ReviewRequest) error {
	ctx := d.actionContext()
	log := zaplog.L(ctx).Named("review.api")
	err := d.actions.Review.Submit(ctx, review.SubmitInput{CardID: req.CardID, Rating: req.Rating})
	if err != nil {
		log.Info("SubmitReview failed", zap.String("cardId", req.CardID), zap.Int("rating", req.Rating), zap.Error(err))
		return err
	}
	log.Info("SubmitReview ok", zap.String("cardId", req.CardID), zap.Int("rating", req.Rating))
	return nil
}

func (d *Desktop) UndoLastReview(cardID string) error {
	ctx := d.actionContext()
	log := zaplog.L(ctx).Named("review.api")
	err := d.actions.Review.UndoLastReview(ctx, cardID)
	if err != nil {
		log.Info("UndoLastReview failed", zap.String("cardId", cardID), zap.Error(err))
		return err
	}
	log.Info("UndoLastReview ok", zap.String("cardId", cardID))
	return nil
}

func (d *Desktop) GetSRSStatistics(knowledgeID *string) (*SRSStatisticsDTO, error) {
	ctx := d.actionContext()
	log := zaplog.L(ctx).Named("review.api")
	stats, err := d.actions.Review.GetStatistics(ctx, knowledgeID)
	if err != nil {
		log.Info("GetSRSStatistics failed", zapOptionalString("knowledgeId", knowledgeID), zap.Error(err))
		return nil, err
	}
	out := &SRSStatisticsDTO{
		NewCount:        stats.NewCount,
		LearningCount:   stats.LearningCount,
		ReviewCount:     stats.ReviewCount,
		RelearningCount: stats.RelearningCount,
		TotalCards:      stats.TotalCards,
		DueToday:        stats.DueToday,
	}
	log.Info("GetSRSStatistics ok",
		zapOptionalString("knowledgeId", knowledgeID),
		zap.Int("new", out.NewCount),
		zap.Int("learning", out.LearningCount),
		zap.Int("review", out.ReviewCount),
		zap.Int("relearning", out.RelearningCount),
		zap.Int("total_cards", out.TotalCards),
		zap.Int("due_today", out.DueToday),
	)
	return out, nil
}

func (d *Desktop) GetReviewHistory(cardID string, limit int) ([]*ReviewLogDTO, error) {
	ctx := d.actionContext()
	log := zaplog.L(ctx).Named("review.api")
	items, err := d.actions.Review.GetHistory(ctx, cardID, limit)
	if err != nil {
		log.Info("GetReviewHistory failed", zap.String("cardId", cardID), zap.Int("limit", limit), zap.Error(err))
		return nil, err
	}
	result := make([]*ReviewLogDTO, 0, len(items))
	for _, item := range items {
		result = append(result, toReviewLogDTO(item))
	}
	head := make([]string, 0, min(3, len(result)))
	for i := range result {
		if i >= 3 {
			break
		}
		head = append(head, fmt.Sprintf("%s@rating=%d", result[i].ID, result[i].Rating))
	}
	log.Info("GetReviewHistory ok", zap.String("cardId", cardID), zap.Int("limit", limit), zap.Int("returned", len(result)), zap.Strings("sample", head))
	return result, nil
}

func (d *Desktop) GetReviewStats(startDate, endDate time.Time) (*repository.ReviewStatistics, error) {
	ctx := d.actionContext()
	log := zaplog.L(ctx).Named("review.api")
	stats, err := d.actions.Review.GetStatsByDate(ctx, startDate, endDate)
	if err != nil {
		log.Info("GetReviewStats failed",
			zap.Time("start", startDate),
			zap.Time("end", endDate),
			zap.Error(err),
		)
		return nil, err
	}
	log.Info("GetReviewStats ok",
		zap.Time("start", startDate),
		zap.Time("end", endDate),
		zap.Int("total_reviews", stats.TotalReviews),
		zap.Float64("avg_rating", stats.AvgRating),
		zap.Int("rating_bucket_keys", len(stats.RatingCounts)),
		zap.Int("reviews_by_day_keys", len(stats.ReviewsByDay)),
	)
	return stats, nil
}

func toSRSDTO(model *models.CardSRS) *SRSDTO {
	if model == nil {
		return nil
	}
	return &SRSDTO{
		CardID:       model.CardID,
		FSRSState:    model.FSRSState,
		DueAt:        model.DueAt,
		LastReviewAt: model.LastReviewAt,
		Stability:    model.Stability,
		Difficulty:   model.Difficulty,
		Reps:         model.Reps,
		Lapses:       model.Lapses,
	}
}

func toReviewLogDTO(model *models.ReviewLog) *ReviewLogDTO {
	if model == nil {
		return nil
	}
	return &ReviewLogDTO{ID: model.ID, CardID: model.CardID, ReviewedAt: model.ReviewedAt, Rating: model.Rating, ReviewKind: model.ReviewKind}
}
