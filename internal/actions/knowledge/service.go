package knowledge

import (
	"context"
	"time"

	"github.com/google/uuid"

	"kmemo/internal/storage/models"
	"kmemo/internal/storage/repository"
)

// Dependencies 聚合 knowledge 用例所需的仓储。
type Dependencies struct {
	Repo  repository.KnowledgeRepository
	Cards repository.CardRepository
	SRS   repository.SRSRepository
}

// Service handles knowledge-oriented user actions.
type Service struct {
	deps Dependencies
}

func NewService(deps Dependencies) *Service {
	return &Service{deps: deps}
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
	if err := s.deps.Repo.Create(ctx, knowledge); err != nil {
		return "", err
	}
	return knowledge.ID, nil
}

func (s *Service) Get(ctx context.Context, id string) (*models.Knowledge, error) {
	return s.deps.Repo.GetByID(ctx, id)
}

func (s *Service) List(ctx context.Context, parentID *string) ([]*models.Knowledge, error) {
	list, _, err := s.deps.Repo.List(ctx, repository.ListKnowledgeOptions{ParentID: parentID})
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (s *Service) ListAll(ctx context.Context) ([]*models.Knowledge, error) {
	list, _, err := s.deps.Repo.List(ctx, repository.ListKnowledgeOptions{})
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (s *Service) GetTree(ctx context.Context, rootID string) (*models.Knowledge, error) {
	return s.deps.Repo.GetTree(ctx, rootID)
}

func (s *Service) Update(ctx context.Context, id string, input UpdateInput) error {
	knowledge, err := s.deps.Repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	knowledge.Name = input.Name
	knowledge.Description = input.Description
	return s.deps.Repo.Update(ctx, knowledge)
}

func (s *Service) Delete(ctx context.Context, id string) error {
	return s.deps.Repo.Delete(ctx, id)
}

func (s *Service) Move(ctx context.Context, id string, newParentID *string) error {
	return s.deps.Repo.Move(ctx, id, newParentID)
}

func (s *Service) Archive(ctx context.Context, id string) error {
	return s.deps.Repo.Archive(ctx, id)
}

func (s *Service) Unarchive(ctx context.Context, id string) error {
	return s.deps.Repo.Unarchive(ctx, id)
}

// CountMapsByKnowledgeIDs 批量返回各知识库下的卡片总数与「当前到期」卡片数（与 Card 列表 / GetDueCards 语义对齐）。
func (s *Service) CountMapsByKnowledgeIDs(ctx context.Context, knowledgeIDs []string) (cardCounts map[string]int64, dueCounts map[string]int64, err error) {
	cardCounts, err = s.deps.Cards.CountByKnowledgeIDs(ctx, knowledgeIDs)
	if err != nil {
		return nil, nil, err
	}
	dueCounts, err = s.deps.SRS.CountDueCardsByKnowledgeIDs(ctx, knowledgeIDs)
	if err != nil {
		return nil, nil, err
	}
	return cardCounts, dueCounts, nil
}
