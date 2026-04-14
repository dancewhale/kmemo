package card

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"kmemo/internal/contracts"
	"kmemo/internal/file"
	"kmemo/internal/storage/models"
	"kmemo/internal/storage/repository"
	"kmemo/internal/zaplog"
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

type MoveInput struct {
	CardID         string
	TargetParentID *string
	TargetIndex    int
}

type ReorderChildrenInput struct {
	KnowledgeID     string
	ParentID        *string
	OrderedChildIDs []string
}

// GetOutput 为单卡查询结果：包含数据库中的 Card 以及从 FileStore 读取的正文 HTML。
type GetOutput struct {
	Card        *models.Card
	HTMLContent string
}

func (s *Service) Create(ctx context.Context, input CreateInput) (string, error) {
	started := time.Now()
	log := zaplog.L(ctx).Named("card")
	log.Debug("card.create.start",
		zap.String("knowledge_id", input.KnowledgeID),
		zap.String("card_type", input.CardType),
		zap.Int("tag_count", len(input.TagIDs)),
	)

	if _, err := s.deps.Knowledge.GetByID(ctx, input.KnowledgeID); err != nil {
		log.Debug("card.create.fail", zap.String("phase", "knowledge"), zap.Error(err))
		return "", err
	}
	if input.ParentID != nil {
		parent, err := s.deps.Cards.GetByID(ctx, *input.ParentID)
		if err != nil {
			log.Debug("card.create.fail", zap.String("phase", "parent"), zap.Error(err))
			return "", err
		}
		if parent.KnowledgeID != input.KnowledgeID {
			log.Debug("card.create.fail", zap.String("phase", "parent_knowledge_mismatch"))
			return "", repository.ErrInvalidInput
		}
	}

	uid, genErr := uuid.NewV7()
	if genErr != nil {
		log.Debug("card.create.fail", zap.String("phase", "id"), zap.Error(genErr))
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
			log.Debug("card.create.fail", zap.String("phase", "file_create"), zap.Error(err))
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
		SortOrder:        0,
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	if s.deps.Transactions == nil {
		maxSort, err := s.deps.Cards.GetMaxSortOrder(ctx, input.KnowledgeID, input.ParentID)
		if err != nil {
			return "", err
		}
		card.SortOrder = maxSort + 1
		if err := s.deps.Cards.Create(ctx, card); err != nil {
			if fileCreated && s.deps.FileStore != nil {
				_ = s.deps.FileStore.PermanentlyDeleteFileObject(ctx, contracts.FileObjectRef{Kind: contracts.FileObjectKindCard, ID: contracts.FileObjectID(id)})
			}
			log.Debug("card.create.fail", zap.String("phase", "card_create"), zap.Error(err))
			return "", err
		}
		if len(input.TagIDs) > 0 {
			if err := s.deps.Cards.AddTags(ctx, id, input.TagIDs); err != nil {
				if fileCreated && s.deps.FileStore != nil {
					_ = s.deps.FileStore.PermanentlyDeleteFileObject(ctx, contracts.FileObjectRef{Kind: contracts.FileObjectKindCard, ID: contracts.FileObjectID(id)})
				}
				log.Debug("card.create.fail", zap.String("phase", "add_tags"), zap.Error(err))
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
				log.Debug("card.create.fail", zap.String("phase", "srs_init"), zap.Error(err))
				return "", err
			}
		}
		log.Debug("card.create.success", zap.String("card_id", id), zap.Duration("duration", time.Since(started)))
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
		maxSort, err := cards.GetMaxSortOrder(ctx, input.KnowledgeID, input.ParentID)
		if err != nil {
			return err
		}
		card.SortOrder = maxSort + 1

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
		log.Debug("card.create.fail", zap.String("phase", "tx"), zap.Error(err))
		return "", err
	}
	log.Debug("card.create.success", zap.String("card_id", id), zap.Duration("duration", time.Since(started)))
	return id, nil
}

func (s *Service) Move(ctx context.Context, input MoveInput) error {
	if input.CardID == "" || input.TargetIndex < 0 {
		return repository.ErrInvalidInput
	}
	cardModel, err := s.deps.Cards.GetByID(ctx, input.CardID)
	if err != nil {
		return err
	}
	if input.TargetParentID != nil {
		if *input.TargetParentID == input.CardID {
			return repository.ErrInvalidInput
		}
		targetParent, err := s.deps.Cards.GetByID(ctx, *input.TargetParentID)
		if err != nil {
			return err
		}
		if targetParent.KnowledgeID != cardModel.KnowledgeID {
			return repository.ErrInvalidInput
		}
		if err := s.ensureNotDescendant(ctx, input.CardID, *input.TargetParentID); err != nil {
			return err
		}
	}

	run := func(cards repository.CardRepository) error {
		sourceParentID := cardModel.ParentID
		sourceSiblings, err := cards.ListSiblings(ctx, cardModel.KnowledgeID, sourceParentID)
		if err != nil {
			return err
		}
		targetSiblings, err := cards.ListSiblings(ctx, cardModel.KnowledgeID, input.TargetParentID)
		if err != nil {
			return err
		}

		if sameParent(sourceParentID, input.TargetParentID) {
			targetSiblings = sourceSiblings
		}

		sourceWithout := removeCardByID(sourceSiblings, input.CardID)
		targetWithout := removeCardByID(targetSiblings, input.CardID)
		targetIndex := clampIndex(input.TargetIndex, len(targetWithout))
		targetAfter := insertCardAt(targetWithout, cardModel, targetIndex)

		updates := make([]repository.CardSortOrderUpdate, 0, len(targetAfter)+len(sourceWithout))
		for i, item := range targetAfter {
			updates = append(updates, repository.CardSortOrderUpdate{
				CardID:    item.ID,
				SortOrder: i,
				ParentID:  input.TargetParentID,
				IsRoot:    input.TargetParentID == nil,
			})
		}
		if !sameParent(sourceParentID, input.TargetParentID) {
			for i, item := range sourceWithout {
				updates = append(updates, repository.CardSortOrderUpdate{
					CardID:    item.ID,
					SortOrder: i,
					ParentID:  sourceParentID,
					IsRoot:    sourceParentID == nil,
				})
			}
		}
		return cards.BatchUpdateSortOrders(ctx, updates)
	}

	if s.deps.Transactions == nil {
		return run(s.deps.Cards)
	}
	return s.deps.Transactions.WithTx(ctx, func(tx *gorm.DB) error {
		cards, ok := s.deps.Cards.WithTx(tx).(repository.CardRepository)
		if !ok {
			return fmt.Errorf("card: repository.WithTx did not return CardRepository")
		}
		return run(cards)
	})
}

func (s *Service) ReorderChildren(ctx context.Context, input ReorderChildrenInput) error {
	if input.KnowledgeID == "" || len(input.OrderedChildIDs) == 0 {
		return repository.ErrInvalidInput
	}
	run := func(cards repository.CardRepository) error {
		siblings, err := cards.ListSiblings(ctx, input.KnowledgeID, input.ParentID)
		if err != nil {
			return err
		}
		if len(siblings) != len(input.OrderedChildIDs) {
			return repository.ErrInvalidInput
		}
		exists := make(map[string]struct{}, len(siblings))
		for _, s := range siblings {
			exists[s.ID] = struct{}{}
		}
		seen := make(map[string]struct{}, len(input.OrderedChildIDs))
		updates := make([]repository.CardSortOrderUpdate, 0, len(input.OrderedChildIDs))
		for i, id := range input.OrderedChildIDs {
			if _, ok := exists[id]; !ok {
				return repository.ErrInvalidInput
			}
			if _, ok := seen[id]; ok {
				return repository.ErrInvalidInput
			}
			seen[id] = struct{}{}
			updates = append(updates, repository.CardSortOrderUpdate{
				CardID:    id,
				SortOrder: i,
				ParentID:  input.ParentID,
				IsRoot:    input.ParentID == nil,
			})
		}
		return cards.BatchUpdateSortOrders(ctx, updates)
	}

	if s.deps.Transactions == nil {
		return run(s.deps.Cards)
	}
	return s.deps.Transactions.WithTx(ctx, func(tx *gorm.DB) error {
		cards, ok := s.deps.Cards.WithTx(tx).(repository.CardRepository)
		if !ok {
			return fmt.Errorf("card: repository.WithTx did not return CardRepository")
		}
		return run(cards)
	})
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
	started := time.Now()
	log := zaplog.L(ctx).Named("card")
	log.Debug("card.update.start", zap.String("card_id", id), zap.String("status", input.Status))
	card, err := s.deps.Cards.GetByID(ctx, id)
	if err != nil {
		log.Debug("card.update.fail", zap.String("phase", "load"), zap.String("card_id", id), zap.Error(err))
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
				log.Debug("card.update.fail", zap.String("phase", "file_overwrite"), zap.String("card_id", id), zap.Error(err))
				return err
			}
			card.Path = loc.Path
		}
	}
	if err := s.deps.Cards.Update(ctx, card); err != nil {
		log.Debug("card.update.fail", zap.String("phase", "db_update"), zap.String("card_id", id), zap.Error(err))
		return err
	}
	log.Debug("card.update.success", zap.String("card_id", id), zap.Duration("duration", time.Since(started)))
	return nil
}

func (s *Service) Delete(ctx context.Context, id string) error {
	started := time.Now()
	log := zaplog.L(ctx).Named("card")
	log.Debug("card.delete.start", zap.String("card_id", id))
	if s.deps.FileStore != nil {
		_, _ = s.deps.FileStore.MoveFileObjectToTrash(ctx, contracts.FileObjectRef{Kind: contracts.FileObjectKindCard, ID: contracts.FileObjectID(id)})
	}
	if err := s.deps.Cards.Delete(ctx, id); err != nil {
		log.Debug("card.delete.fail", zap.String("card_id", id), zap.Error(err))
		return err
	}
	log.Debug("card.delete.success", zap.String("card_id", id), zap.Duration("duration", time.Since(started)))
	return nil
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

func sameParent(a, b *string) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}

func clampIndex(idx, size int) int {
	if idx < 0 {
		return 0
	}
	if idx > size {
		return size
	}
	return idx
}

func removeCardByID(items []*models.Card, cardID string) []*models.Card {
	result := make([]*models.Card, 0, len(items))
	for _, item := range items {
		if item.ID == cardID {
			continue
		}
		result = append(result, item)
	}
	return result
}

func insertCardAt(items []*models.Card, target *models.Card, index int) []*models.Card {
	result := make([]*models.Card, 0, len(items)+1)
	result = append(result, items[:index]...)
	result = append(result, target)
	result = append(result, items[index:]...)
	return result
}

func (s *Service) ensureNotDescendant(ctx context.Context, cardID, targetParentID string) error {
	parentID := &targetParentID
	for parentID != nil {
		if *parentID == cardID {
			return repository.ErrInvalidInput
		}
		node, err := s.deps.Cards.GetByID(ctx, *parentID)
		if err != nil {
			return err
		}
		parentID = node.ParentID
	}
	return nil
}
