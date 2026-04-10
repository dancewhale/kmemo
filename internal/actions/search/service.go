package search

import (
	"context"

	"kmemo/internal/storage/models"
	"kmemo/internal/storage/repository"
)

type Dependencies struct {
	Cards repository.CardRepository
	Tags  repository.TagRepository
}

type Service struct {
	deps Dependencies
}

func NewService(deps Dependencies) *Service {
	return &Service{deps: deps}
}

func (s *Service) SearchCardsByTags(ctx context.Context, tagIDs []string) ([]*models.Card, error) {
	items, _, err := s.deps.Cards.List(ctx, repository.ListCardOptions{TagIDs: tagIDs, Preload: []string{"SRS", "Knowledge"}})
	if err != nil {
		return nil, err
	}
	return items, nil
}
