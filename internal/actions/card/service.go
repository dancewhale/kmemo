package card

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"kmemo/internal/contracts"
	"kmemo/internal/file"
	"kmemo/internal/storage/models"
	"kmemo/internal/storage/repository"
)

type Dependencies struct {
	Cards        repository.CardRepository
	Knowledge    repository.KnowledgeRepository
	SRS          repository.SRSRepository
	Tags         repository.TagRepository
	Transactions repository.TransactionManager
	FileStore    contracts.FileStore
}

type Service struct {
	deps Dependencies
}

func NewService(deps Dependencies) *Service {
	return &Service{deps: deps}
}

type CreateInput struct {
	KnowledgeID      string
	SourceDocumentID *string
	ParentID         *string
	Title            string
	CardType         string
	HTMLContent      string
	TagIDs           []string
}

type UpdateInput struct {
	Title       string
	HTMLContent string
	Status      string
}

// GetOutput 为单卡查询结果：包含数据库中的 Card 以及从 FileStore 读取的正文 HTML。
type GetOutput struct {
	Card         *models.Card
	HTMLContent  string
}

func (s *Service) Create(ctx context.Context, input CreateInput) (string, error) {
	if _, err := s.deps.Knowledge.GetByID(ctx, input.KnowledgeID); err != nil {
		return "", err
	}

	uid, genErr := uuid.NewV7()
	if genErr != nil {
		return "", genErr
	}
	id := uid.String()
	now := time.Now().UTC()
	slug := file.NormalizeSlug(input.Title)
	htmlHash := hashContent(input.HTMLContent)
	path := ""
	fileCreated := false

	if s.deps.FileStore != nil {
		loc, err := s.deps.FileStore.CreateFileObject(ctx, contracts.CreateFileObjectInput{
			Kind: contracts.FileObjectKindCard,
			ID:   contracts.FileObjectID(id),
			Slug: slug,
			Ext:  "html",
			Data: []byte(input.HTMLContent),
		})
		if err != nil {
			return "", err
		}
		path = loc.Path
		fileCreated = true
	}

	card := &models.Card{
		ID:               id,
		KnowledgeID:      input.KnowledgeID,
		SourceDocumentID: input.SourceDocumentID,
		ParentID:         input.ParentID,
		Title:            input.Title,
		Path:             path,
		CardType:         input.CardType,
		Slug:             slug,
		HTMLHash:         htmlHash,
		Status:           "active",
		IsRoot:           input.ParentID == nil,
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	if s.deps.Transactions == nil {
		if err := s.deps.Cards.Create(ctx, card); err != nil {
			if fileCreated && s.deps.FileStore != nil {
				_ = s.deps.FileStore.PermanentlyDeleteFileObject(ctx, contracts.FileObjectRef{Kind: contracts.FileObjectKindCard, ID: contracts.FileObjectID(id)})
			}
			return "", err
		}
		if len(input.TagIDs) > 0 {
			if err := s.deps.Cards.AddTags(ctx, id, input.TagIDs); err != nil {
				if fileCreated && s.deps.FileStore != nil {
					_ = s.deps.FileStore.PermanentlyDeleteFileObject(ctx, contracts.FileObjectRef{Kind: contracts.FileObjectKindCard, ID: contracts.FileObjectID(id)})
				}
				return "", err
			}
		}
		if shouldInitializeSRS(input.CardType) {
			srs := &models.CardSRS{
				CardID:    id,
				FSRSState: "new",
				CreatedAt: now,
				UpdatedAt: now,
			}
			if err := s.deps.SRS.CreateOrUpdate(ctx, srs); err != nil {
				if fileCreated && s.deps.FileStore != nil {
					_ = s.deps.FileStore.PermanentlyDeleteFileObject(ctx, contracts.FileObjectRef{Kind: contracts.FileObjectKindCard, ID: contracts.FileObjectID(id)})
				}
				return "", err
			}
		}
		return id, nil
	}

	err := s.deps.Transactions.WithTx(ctx, func(tx *gorm.DB) error {
		cards, ok := s.deps.Cards.WithTx(tx).(repository.CardRepository)
		if !ok {
			return fmt.Errorf("card: repository.WithTx did not return CardRepository")
		}
		srsRepo, ok := s.deps.SRS.WithTx(tx).(repository.SRSRepository)
		if !ok {
			return fmt.Errorf("card: repository.WithTx did not return SRSRepository")
		}

		if err := cards.Create(ctx, card); err != nil {
			return err
		}
		if len(input.TagIDs) > 0 {
			if err := cards.AddTags(ctx, id, input.TagIDs); err != nil {
				return err
			}
		}
		if shouldInitializeSRS(input.CardType) {
			srs := &models.CardSRS{
				CardID:    id,
				FSRSState: "new",
				CreatedAt: now,
				UpdatedAt: now,
			}
			if err := srsRepo.CreateOrUpdate(ctx, srs); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		if fileCreated && s.deps.FileStore != nil {
			_ = s.deps.FileStore.PermanentlyDeleteFileObject(ctx, contracts.FileObjectRef{Kind: contracts.FileObjectKindCard, ID: contracts.FileObjectID(id)})
		}
		return "", err
	}
	return id, nil
}

func (s *Service) Get(ctx context.Context, id string) (*GetOutput, error) {
	c, err := s.deps.Cards.GetByID(ctx, id, "SRS", "Knowledge")
	if err != nil {
		return nil, err
	}
	out := &GetOutput{Card: c}
	if s.deps.FileStore == nil {
		return out, nil
	}
	lookup := contracts.FileObjectLookup{
		Ref: contracts.FileObjectRef{
			Kind: contracts.FileObjectKindCard,
			ID:   contracts.FileObjectID(c.ID),
		},
	}
	if c.Slug != "" {
		lookup.Name = &contracts.FileObjectNameHint{Slug: c.Slug, Ext: "html"}
	}
	data, _, err := s.deps.FileStore.ReadFileObject(ctx, lookup, contracts.FileObjectScopeAny)
	if err != nil {
		if errors.Is(err, file.ErrFileObjectNotFound) {
			return out, nil
		}
		return nil, err
	}
	out.HTMLContent = string(data)
	return out, nil
}

func (s *Service) List(ctx context.Context, opts repository.ListCardOptions) ([]*models.Card, int64, error) {
	return s.deps.Cards.List(ctx, opts)
}

func (s *Service) Update(ctx context.Context, id string, input UpdateInput) error {
	card, err := s.deps.Cards.GetByID(ctx, id)
	if err != nil {
		return err
	}
	card.Title = input.Title
	card.Slug = file.NormalizeSlug(input.Title)
	card.Status = input.Status
	if input.HTMLContent != "" {
		card.HTMLHash = hashContent(input.HTMLContent)
		if s.deps.FileStore != nil {
			loc, err := s.deps.FileStore.OverwriteFileObject(ctx, contracts.OverwriteFileObjectInput{
				Kind: contracts.FileObjectKindCard,
				ID:   contracts.FileObjectID(card.ID),
				Slug: card.Slug,
				Ext:  "html",
				Data: []byte(input.HTMLContent),
			})
			if err != nil {
				return err
			}
			card.Path = loc.Path
		}
	}
	return s.deps.Cards.Update(ctx, card)
}

func (s *Service) Delete(ctx context.Context, id string) error {
	if s.deps.FileStore != nil {
		_, _ = s.deps.FileStore.MoveFileObjectToTrash(ctx, contracts.FileObjectRef{Kind: contracts.FileObjectKindCard, ID: contracts.FileObjectID(id)})
	}
	return s.deps.Cards.Delete(ctx, id)
}

func (s *Service) AddTags(ctx context.Context, cardID string, tagIDs []string) error {
	if err := s.deps.Cards.AddTags(ctx, cardID, tagIDs); err != nil {
		return err
	}
	for _, tagID := range tagIDs {
		if err := s.deps.Tags.UpdateCardCount(ctx, tagID); err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) RemoveTags(ctx context.Context, cardID string, tagIDs []string) error {
	if err := s.deps.Cards.RemoveTags(ctx, cardID, tagIDs); err != nil {
		return err
	}
	for _, tagID := range tagIDs {
		if err := s.deps.Tags.UpdateCardCount(ctx, tagID); err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) GetTags(ctx context.Context, cardID string) ([]*models.Tag, error) {
	return s.deps.Cards.GetTags(ctx, cardID)
}

func (s *Service) CreateTag(ctx context.Context, tag *models.Tag) error {
	return s.deps.Tags.Create(ctx, tag)
}

func (s *Service) GetTag(ctx context.Context, id string) (*models.Tag, error) {
	return s.deps.Tags.GetByID(ctx, id)
}

func (s *Service) ListTags(ctx context.Context) ([]*models.Tag, error) {
	items, _, err := s.deps.Tags.List(ctx, repository.ListTagOptions{})
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (s *Service) UpdateTag(ctx context.Context, tag *models.Tag) error {
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

func (s *Service) DeleteTag(ctx context.Context, id string) error {
	return s.deps.Tags.Delete(ctx, id)
}

func (s *Service) Suspend(ctx context.Context, cardID string) error {
	if err := s.deps.Cards.UpdateStatus(ctx, []string{cardID}, "suspended"); err != nil {
		return err
	}
	return s.deps.SRS.Suspend(ctx, cardID)
}

func (s *Service) Resume(ctx context.Context, cardID string) error {
	if err := s.deps.Cards.UpdateStatus(ctx, []string{cardID}, "active"); err != nil {
		return err
	}
	return s.deps.SRS.Resume(ctx, cardID)
}

func shouldInitializeSRS(cardType string) bool {
	switch cardType {
	case "article", "excerpt", "qa", "cloze":
		return true
	default:
		return false
	}
}

func hashContent(content string) string {
	sum := sha256.Sum256([]byte(content))
	return hex.EncodeToString(sum[:])
}

