package review

import (
	"context"
	"fmt"
	"time"

	fsrscontract "kmemo/internal/contracts/fsrs"
	"kmemo/internal/storage/models"
	"kmemo/internal/storage/repository"
)

type Dependencies struct {
	Cards      repository.CardRepository
	SRS        repository.SRSRepository
	ReviewLogs repository.ReviewLogRepository
	FSRS       fsrscontract.FSRSScheduler
}

type Service struct {
	deps Dependencies
}

func NewService(deps Dependencies) *Service {
	return &Service{deps: deps}
}

type SubmitInput struct {
	CardID string
	Rating int
}

func (s *Service) GetDueCards(ctx context.Context, knowledgeID *string, limit int) ([]*models.Card, error) {
	srsItems, err := s.deps.SRS.GetDueCards(ctx, repository.DueCardsOptions{KnowledgeID: knowledgeID, Limit: limit})
	if err != nil {
		return nil, err
	}
	cards := make([]*models.Card, 0, len(srsItems))
	for _, item := range srsItems {
		card, err := s.deps.Cards.GetByID(ctx, item.CardID, "SRS", "Knowledge")
		if err != nil {
			return nil, err
		}
		cards = append(cards, card)
	}
	return cards, nil
}

func (s *Service) Submit(ctx context.Context, input SubmitInput) error {
	if input.CardID == "" {
		return repository.ErrInvalidInput
	}
	card, err := s.deps.Cards.GetByID(ctx, input.CardID, "SRS")
	if err != nil {
		return err
	}
	result, err := s.deps.FSRS.Review(ctx, fsrscontract.ReviewInput{
		CardID:     card.ID,
		Prior:      card.SRS,
		Rating:     input.Rating,
		ReviewedAt: time.Now().UTC(),
	})
	if err != nil {
		return err
	}
	if result == nil || result.Next == nil || result.Log == nil {
		return fmt.Errorf("review: empty fsrs review result")
	}
	return s.deps.SRS.UpdateAfterReview(ctx, input.CardID, result.Next, result.Log)
}

func (s *Service) UndoLastReview(ctx context.Context, cardID string) error {
	return s.deps.SRS.UndoLastReview(ctx, cardID)
}

func (s *Service) GetStatistics(ctx context.Context, knowledgeID *string) (*repository.SRSStatistics, error) {
	return s.deps.SRS.GetStatistics(ctx, knowledgeID)
}

func (s *Service) GetHistory(ctx context.Context, cardID string, limit int) ([]*models.ReviewLog, error) {
	items, _, err := s.deps.ReviewLogs.GetByCardID(ctx, cardID, repository.ListReviewLogOptions{Limit: limit, OrderDesc: true})
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (s *Service) GetStatsByDate(ctx context.Context, startDate, endDate time.Time) (*repository.ReviewStatistics, error) {
	return s.deps.ReviewLogs.GetStatsByDateRange(ctx, startDate, endDate, nil)
}
