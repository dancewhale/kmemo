package tag

import (
	"context"

	"kmemo/internal/storage/models"
	"kmemo/internal/storage/repository"
)

type Dependencies struct {
	Tags repository.TagRepository
}

type Service struct {
	deps Dependencies
}

func NewService(deps Dependencies) *Service {
	return &Service{deps: deps}
}

func (s *Service) Create(ctx context.Context, tag *models.Tag) error {
	return s.deps.Tags.Create(ctx, tag)
}

func (s *Service) Get(ctx context.Context, id string) (*models.Tag, error) {
	return s.deps.Tags.GetByID(ctx, id)
}

func (s *Service) List(ctx context.Context) ([]*models.Tag, error) {
	items, _, err := s.deps.Tags.List(ctx, repository.ListTagOptions{})
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (s *Service) Update(ctx context.Context, tag *models.Tag) error {
	existing, err := s.deps.Tags.GetByID(ctx, tag.ID)
	if err != nil {
		return err
	}
	existing.Name = tag.Name
	existing.Color = tag.Color
	existing.Icon = tag.Icon
	existing.Description = tag.Description
	return s.deps.Tags.Update(ctx, existing)
}

func (s *Service) Delete(ctx context.Context, id string) error {
	return s.deps.Tags.Delete(ctx, id)
}
