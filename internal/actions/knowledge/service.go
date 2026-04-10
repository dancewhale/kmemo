package knowledge

import (
	"context"
	"time"

	"github.com/google/uuid"

	"kmemo/internal/storage/models"
	"kmemo/internal/storage/repository"
)

// Service handles knowledge-oriented user actions.
type Service struct {
	repo repository.KnowledgeRepository
}

func NewService(repo repository.KnowledgeRepository) *Service {
	return &Service{repo: repo}
}

type CreateInput struct {
	Name        string
	Description string
	ParentID    *string
}

type UpdateInput struct {
	Name        string
	Description string
}

func (s *Service) Create(ctx context.Context, input CreateInput) (string, error) {
	now := time.Now().UTC()
	knowledge := &models.Knowledge{
		ID:          uuid.NewString(),
		Name:        input.Name,
		Description: input.Description,
		ParentID:    input.ParentID,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := s.repo.Create(ctx, knowledge); err != nil {
		return "", err
	}
	return knowledge.ID, nil
}

func (s *Service) Get(ctx context.Context, id string) (*models.Knowledge, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) List(ctx context.Context, parentID *string) ([]*models.Knowledge, error) {
	list, _, err := s.repo.List(ctx, repository.ListKnowledgeOptions{ParentID: parentID})
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (s *Service) ListAll(ctx context.Context) ([]*models.Knowledge, error) {
	list, _, err := s.repo.List(ctx, repository.ListKnowledgeOptions{})
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (s *Service) GetTree(ctx context.Context, rootID string) (*models.Knowledge, error) {
	return s.repo.GetTree(ctx, rootID)
}

func (s *Service) Update(ctx context.Context, id string, input UpdateInput) error {
	knowledge, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	knowledge.Name = input.Name
	knowledge.Description = input.Description
	return s.repo.Update(ctx, knowledge)
}

func (s *Service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *Service) Move(ctx context.Context, id string, newParentID *string) error {
	return s.repo.Move(ctx, id, newParentID)
}

func (s *Service) Archive(ctx context.Context, id string) error {
	return s.repo.Archive(ctx, id)
}

func (s *Service) Unarchive(ctx context.Context, id string) error {
	return s.repo.Unarchive(ctx, id)
}
