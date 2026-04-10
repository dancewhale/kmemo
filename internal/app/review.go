package app

import (
	"time"

	"kmemo/internal/actions/review"
	"kmemo/internal/storage/models"
	"kmemo/internal/storage/repository"
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
	items, err := d.actions.Review.GetDueCards(d.actionContext(), knowledgeID, limit)
	if err != nil {
		return nil, err
	}
	result := make([]*CardWithSRSDTO, 0, len(items))
	for _, item := range items {
		result = append(result, &CardWithSRSDTO{Card: *toCardDTO(item, nil), SRS: *toSRSDTO(item.SRS)})
	}
	return result, nil
}

func (d *Desktop) SubmitReview(req ReviewRequest) error {
	return d.actions.Review.Submit(d.actionContext(), review.SubmitInput{CardID: req.CardID, Rating: req.Rating})
}

func (d *Desktop) UndoLastReview(cardID string) error {
	return d.actions.Review.UndoLastReview(d.actionContext(), cardID)
}

func (d *Desktop) GetSRSStatistics(knowledgeID *string) (*SRSStatisticsDTO, error) {
	stats, err := d.actions.Review.GetStatistics(d.actionContext(), knowledgeID)
	if err != nil {
		return nil, err
	}
	return &SRSStatisticsDTO{
		NewCount:        stats.NewCount,
		LearningCount:   stats.LearningCount,
		ReviewCount:     stats.ReviewCount,
		RelearningCount: stats.RelearningCount,
		TotalCards:      stats.TotalCards,
		DueToday:        stats.DueToday,
	}, nil
}

func (d *Desktop) GetReviewHistory(cardID string, limit int) ([]*ReviewLogDTO, error) {
	items, err := d.actions.Review.GetHistory(d.actionContext(), cardID, limit)
	if err != nil {
		return nil, err
	}
	result := make([]*ReviewLogDTO, 0, len(items))
	for _, item := range items {
		result = append(result, toReviewLogDTO(item))
	}
	return result, nil
}

func (d *Desktop) GetReviewStats(startDate, endDate time.Time) (*repository.ReviewStatistics, error) {
	return d.actions.Review.GetStatsByDate(d.actionContext(), startDate, endDate)
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
