package app

import (
	"strings"
	"time"

	"go.uber.org/zap"

	"kmemo/internal/actions/card"
	"kmemo/internal/storage/models"
	"kmemo/internal/storage/repository"
	"kmemo/internal/zaplog"
)

type CardDTO struct {
	ID               string    `json:"id"`
	KnowledgeID      string    `json:"knowledgeId"`
	KnowledgeName    string    `json:"knowledgeName"`
	SourceDocumentID *string   `json:"sourceDocumentId"`
	ParentID         *string   `json:"parentId"`
	Title            string    `json:"title"`
	CardType         string    `json:"cardType"`
	HTMLPath         string    `json:"htmlPath"`
	HTMLContent      string    `json:"htmlContent"`
	Status           string    `json:"status"`
	Tags             []*TagDTO `json:"tags"`
	SRS              *SRSDTO   `json:"srs"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
}

type CardSummaryDTO struct {
	ID        string    `json:"id"`
	ParentID  *string   `json:"parentId"`
	Title     string    `json:"title"`
	CardType  string    `json:"cardType"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type CardDetailDTO struct {
	CardDTO
	Parent   *CardSummaryDTO   `json:"parent"`
	Children []*CardSummaryDTO `json:"children"`
}

// ListCardsResult 为 Wails 绑定使用：v2 仅支持单返回值 + error，不能返回 ([]*CardDTO, int64, error)。
type ListCardsResult struct {
	Items []*CardDTO `json:"items"`
	Total int64      `json:"total"`
}

type CardFilters struct {
	KnowledgeID *string  `json:"knowledgeId"`
	CardType    string   `json:"cardType"`
	Status      string   `json:"status"`
	TagIDs      []string `json:"tagIds"`
	Keyword     string   `json:"keyword"`
	ParentID    *string  `json:"parentId"`
	IsRoot      *bool    `json:"isRoot"`
	OrderBy     string   `json:"orderBy"` // title, created_at, updated_at, sort_order；空则按仓储默认（无显式排序）
	OrderDesc   bool     `json:"orderDesc"`
	Limit       int      `json:"limit"`
	Offset      int      `json:"offset"`
}

type CreateCardRequest struct {
	KnowledgeID      string   `json:"knowledgeId"`
	SourceDocumentID *string  `json:"sourceDocumentId"`
	ParentID         *string  `json:"parentId"`
	Title            string   `json:"title"`
	CardType         string   `json:"cardType"`
	HTMLContent      string   `json:"htmlContent"`
	TagIDs           []string `json:"tagIds"`
}

type UpdateCardRequest struct {
	Title       string `json:"title"`
	HTMLContent string `json:"htmlContent"`
	Status      string `json:"status"`
}

type MoveCardRequest struct {
	CardID         string  `json:"cardId"`
	TargetParentID *string `json:"targetParentId"`
	TargetIndex    int     `json:"targetIndex"`
}

type ReorderCardChildrenRequest struct {
	KnowledgeID     string   `json:"knowledgeId"`
	ParentID        *string  `json:"parentId"`
	OrderedChildIDs []string `json:"orderedChildIds"`
}

func (d *Desktop) CreateCard(req CreateCardRequest) (string, error) {
	ctx := d.actionContext()
	log := zaplog.L(ctx).Named("card.api")
	if strings.TrimSpace(req.KnowledgeID) == "" || strings.TrimSpace(req.Title) == "" || strings.TrimSpace(req.CardType) == "" {
		log.Info("CreateCard rejected", zap.String("reason", "invalid_input"))
		return "", repository.ErrInvalidInput
	}
	id, err := d.actions.Card.Create(ctx, card.CreateInput{
		KnowledgeID:      strings.TrimSpace(req.KnowledgeID),
		SourceDocumentID: req.SourceDocumentID,
		ParentID:         req.ParentID,
		Title:            strings.TrimSpace(req.Title),
		CardType:         strings.TrimSpace(req.CardType),
		HTMLContent:      req.HTMLContent,
		TagIDs:           req.TagIDs,
	})
	if err != nil {
		log.Info("CreateCard failed", zap.String("title", req.Title), zap.Error(err))
		return "", err
	}
	log.Info("CreateCard ok",
		zap.String("id", id),
		zap.String("knowledgeId", strings.TrimSpace(req.KnowledgeID)),
		zap.String("title", strings.TrimSpace(req.Title)),
		zap.String("cardType", strings.TrimSpace(req.CardType)),
		zap.Int("html_len", len(req.HTMLContent)),
		zap.Int("tag_id_count", len(req.TagIDs)),
	)
	return id, nil
}

func (d *Desktop) GetCard(id string) (*CardDTO, error) {
	ctx := d.actionContext()
	log := zaplog.L(ctx).Named("card.api")
	out, err := d.actions.Card.Get(ctx, id)
	if err != nil {
		log.Info("GetCard failed", zap.String("id", id), zap.Error(err))
		return nil, err
	}
	dto := toCardDTO(out.Card, nil)
	dto.HTMLContent = out.HTMLContent
	log.Info("GetCard ok",
		zap.String("id", dto.ID),
		zap.String("title", dto.Title),
		zap.String("knowledgeId", dto.KnowledgeID),
		zap.String("status", dto.Status),
		zap.Int("html_len", len(dto.HTMLContent)),
	)
	return dto, nil
}

func (d *Desktop) GetCardDetail(id string) (*CardDetailDTO, error) {
	ctx := d.actionContext()
	log := zaplog.L(ctx).Named("card.api")
	out, err := d.actions.Card.Get(ctx, id)
	if err != nil {
		log.Info("GetCardDetail failed", zap.String("id", id), zap.String("phase", "get"), zap.Error(err))
		return nil, err
	}
	tags, err := d.actions.Card.GetTags(ctx, id)
	if err != nil {
		log.Info("GetCardDetail failed", zap.String("id", id), zap.String("phase", "tags"), zap.Error(err))
		return nil, err
	}
	detail := &CardDetailDTO{CardDTO: *toCardDTO(out.Card, tags)}
	detail.HTMLContent = out.HTMLContent

	if out.Card.ParentID != nil {
		parent, err := d.actions.Card.Get(ctx, *out.Card.ParentID)
		if err != nil {
			log.Info("GetCardDetail failed", zap.String("id", id), zap.String("phase", "parent"), zap.Error(err))
			return nil, err
		}
		detail.Parent = toCardSummaryDTO(parent.Card)
	}

	children, _, err := d.actions.Card.List(ctx, repository.ListCardOptions{
		ParentID:  &id,
		Preload:   []string{"Knowledge", "SRS"},
		OrderBy:   "sort_order",
		OrderDesc: false,
	})
	if err != nil {
		log.Info("GetCardDetail failed", zap.String("id", id), zap.String("phase", "children"), zap.Error(err))
		return nil, err
	}
	detail.Children = toCardSummaryDTOs(children)
	parentID := ""
	if detail.Parent != nil {
		parentID = detail.Parent.ID
	}
	log.Info("GetCardDetail ok",
		zap.String("id", detail.ID),
		zap.String("title", detail.Title),
		zap.Int("tag_count", len(detail.Tags)),
		zap.Int("children_count", len(detail.Children)),
		zap.String("parent_id", parentID),
		zap.Int("html_len", len(detail.HTMLContent)),
	)
	return detail, nil
}

func (d *Desktop) ListCards(filters CardFilters) (*ListCardsResult, error) {
	ctx := d.actionContext()
	log := zaplog.L(ctx).Named("card.api")
	items, total, err := d.actions.Card.List(ctx, repository.ListCardOptions{
		KnowledgeID: filters.KnowledgeID,
		ParentID:    filters.ParentID,
		CardType:    filters.CardType,
		Status:      filters.Status,
		TagIDs:      filters.TagIDs,
		Keyword:     filters.Keyword,
		IsRoot:      filters.IsRoot,
		OrderBy:     filters.OrderBy,
		OrderDesc:   filters.OrderDesc,
		Limit:       filters.Limit,
		Offset:      filters.Offset,
		Preload:     []string{"SRS", "Knowledge"},
	})
	if err != nil {
		log.Info("ListCards failed", zap.Error(err))
		return nil, err
	}
	result := make([]*CardDTO, 0, len(items))
	for _, item := range items {
		tags, _ := d.actions.Card.GetTags(ctx, item.ID)
		result = append(result, toCardDTO(item, tags))
	}
	sample := make([]string, 0, min(5, len(result)))
	for i := range result {
		if i >= 5 {
			break
		}
		sample = append(sample, result[i].ID)
	}
	log.Info("ListCards ok",
		zap.Int64("total", total),
		zap.Int("returned", len(result)),
		zapOptionalString("knowledgeId", filters.KnowledgeID),
		zapOptionalString("parentId", filters.ParentID),
		zap.String("keyword", truncateRunes(filters.Keyword, 40)),
		zap.Int("limit", filters.Limit),
		zap.Int("offset", filters.Offset),
		zap.Strings("sample_card_ids", sample),
	)
	return &ListCardsResult{Items: result, Total: total}, nil
}

func (d *Desktop) GetCardChildren(parentID string) ([]*CardDTO, error) {
	ctx := d.actionContext()
	log := zaplog.L(ctx).Named("card.api")
	items, _, err := d.actions.Card.List(ctx, repository.ListCardOptions{
		ParentID:  &parentID,
		Preload:   []string{"SRS", "Knowledge"},
		OrderBy:   "sort_order",
		OrderDesc: false,
	})
	if err != nil {
		log.Info("GetCardChildren failed", zap.String("parentId", parentID), zap.Error(err))
		return nil, err
	}
	result := make([]*CardDTO, 0, len(items))
	for _, item := range items {
		tags, _ := d.actions.Card.GetTags(ctx, item.ID)
		result = append(result, toCardDTO(item, tags))
	}
	sample := make([]string, 0, min(5, len(result)))
	for i := range result {
		if i >= 5 {
			break
		}
		sample = append(sample, result[i].ID)
	}
	log.Info("GetCardChildren ok", zap.String("parentId", parentID), zap.Int("count", len(result)), zap.Strings("sample_card_ids", sample))
	return result, nil
}

func (d *Desktop) UpdateCard(id string, req UpdateCardRequest) error {
	ctx := d.actionContext()
	log := zaplog.L(ctx).Named("card.api")
	if strings.TrimSpace(req.Title) == "" {
		log.Info("UpdateCard rejected", zap.String("id", id), zap.String("reason", "empty_title"))
		return repository.ErrInvalidInput
	}
	err := d.actions.Card.Update(ctx, id, card.UpdateInput{
		Title:       strings.TrimSpace(req.Title),
		HTMLContent: req.HTMLContent,
		Status:      strings.TrimSpace(req.Status),
	})
	if err != nil {
		log.Info("UpdateCard failed", zap.String("id", id), zap.Error(err))
		return err
	}
	log.Info("UpdateCard ok",
		zap.String("id", id),
		zap.String("title", strings.TrimSpace(req.Title)),
		zap.String("status", strings.TrimSpace(req.Status)),
		zap.Int("html_len", len(req.HTMLContent)),
	)
	return nil
}

func (d *Desktop) DeleteCard(id string) error {
	ctx := d.actionContext()
	log := zaplog.L(ctx).Named("card.api")
	err := d.actions.Card.Delete(ctx, id)
	if err != nil {
		log.Info("DeleteCard failed", zap.String("id", id), zap.Error(err))
		return err
	}
	log.Info("DeleteCard ok", zap.String("id", id))
	return nil
}

func (d *Desktop) MoveCard(req MoveCardRequest) error {
	ctx := d.actionContext()
	log := zaplog.L(ctx).Named("card.api")
	if strings.TrimSpace(req.CardID) == "" || req.TargetIndex < 0 {
		log.Info("MoveCard rejected", zap.String("reason", "invalid_input"))
		return repository.ErrInvalidInput
	}
	err := d.actions.Card.Move(ctx, card.MoveInput{
		CardID:         strings.TrimSpace(req.CardID),
		TargetParentID: req.TargetParentID,
		TargetIndex:    req.TargetIndex,
	})
	if err != nil {
		log.Info("MoveCard failed", zap.String("cardId", req.CardID), zap.Error(err))
		return err
	}
	log.Info("MoveCard ok",
		zap.String("cardId", strings.TrimSpace(req.CardID)),
		zapOptionalString("targetParentId", req.TargetParentID),
		zap.Int("targetIndex", req.TargetIndex),
	)
	return nil
}

func (d *Desktop) ReorderCardChildren(req ReorderCardChildrenRequest) error {
	ctx := d.actionContext()
	log := zaplog.L(ctx).Named("card.api")
	if strings.TrimSpace(req.KnowledgeID) == "" || len(req.OrderedChildIDs) == 0 {
		log.Info("ReorderCardChildren rejected", zap.String("reason", "invalid_input"))
		return repository.ErrInvalidInput
	}
	err := d.actions.Card.ReorderChildren(ctx, card.ReorderChildrenInput{
		KnowledgeID:     strings.TrimSpace(req.KnowledgeID),
		ParentID:        req.ParentID,
		OrderedChildIDs: req.OrderedChildIDs,
	})
	if err != nil {
		log.Info("ReorderCardChildren failed", zap.String("knowledgeId", req.KnowledgeID), zap.Error(err))
		return err
	}
	log.Info("ReorderCardChildren ok",
		zap.String("knowledgeId", strings.TrimSpace(req.KnowledgeID)),
		zapOptionalString("parentId", req.ParentID),
		zap.Int("child_count", len(req.OrderedChildIDs)),
		zap.Strings("ordered_head", firstStrings(req.OrderedChildIDs, 5)),
	)
	return nil
}

func (d *Desktop) AddCardTags(cardID string, tagIDs []string) error {
	ctx := d.actionContext()
	log := zaplog.L(ctx).Named("card.api")
	err := d.actions.Card.AddTags(ctx, cardID, tagIDs)
	if err != nil {
		log.Info("AddCardTags failed", zap.String("cardId", cardID), zap.Error(err))
		return err
	}
	log.Info("AddCardTags ok", zap.String("cardId", cardID), zap.Int("tag_count", len(tagIDs)), zap.Strings("sample_tag_ids", firstStrings(tagIDs, 5)))
	return nil
}

func (d *Desktop) RemoveCardTags(cardID string, tagIDs []string) error {
	ctx := d.actionContext()
	log := zaplog.L(ctx).Named("card.api")
	err := d.actions.Card.RemoveTags(ctx, cardID, tagIDs)
	if err != nil {
		log.Info("RemoveCardTags failed", zap.String("cardId", cardID), zap.Error(err))
		return err
	}
	log.Info("RemoveCardTags ok", zap.String("cardId", cardID), zap.Int("tag_count", len(tagIDs)), zap.Strings("sample_tag_ids", firstStrings(tagIDs, 5)))
	return nil
}

func (d *Desktop) GetCardTags(cardID string) ([]*TagDTO, error) {
	ctx := d.actionContext()
	log := zaplog.L(ctx).Named("card.api")
	items, err := d.actions.Card.GetTags(ctx, cardID)
	if err != nil {
		log.Info("GetCardTags failed", zap.String("cardId", cardID), zap.Error(err))
		return nil, err
	}
	out := toTagDTOs(items)
	sample := make([]string, 0, min(5, len(out)))
	for i := range out {
		if i >= 5 {
			break
		}
		sample = append(sample, out[i].ID+":"+out[i].Name)
	}
	log.Info("GetCardTags ok", zap.String("cardId", cardID), zap.Int("count", len(out)), zap.Strings("sample", sample))
	return out, nil
}

func (d *Desktop) SuspendCard(cardID string) error {
	ctx := d.actionContext()
	log := zaplog.L(ctx).Named("card.api")
	err := d.actions.Card.Suspend(ctx, cardID)
	if err != nil {
		log.Info("SuspendCard failed", zap.String("cardId", cardID), zap.Error(err))
		return err
	}
	log.Info("SuspendCard ok", zap.String("cardId", cardID))
	return nil
}

func (d *Desktop) ResumeCard(cardID string) error {
	ctx := d.actionContext()
	log := zaplog.L(ctx).Named("card.api")
	err := d.actions.Card.Resume(ctx, cardID)
	if err != nil {
		log.Info("ResumeCard failed", zap.String("cardId", cardID), zap.Error(err))
		return err
	}
	log.Info("ResumeCard ok", zap.String("cardId", cardID))
	return nil
}

func toCardDTO(model *models.Card, tags []*models.Tag) *CardDTO {
	if model == nil {
		return nil
	}
	dto := &CardDTO{
		ID:               model.ID,
		KnowledgeID:      model.KnowledgeID,
		SourceDocumentID: model.SourceDocumentID,
		ParentID:         model.ParentID,
		Title:            model.Title,
		CardType:         model.CardType,
		HTMLPath:         model.Path,
		Status:           model.Status,
		CreatedAt:        model.CreatedAt,
		UpdatedAt:        model.UpdatedAt,
	}
	if model.Knowledge != nil {
		dto.KnowledgeName = model.Knowledge.Name
	}
	if model.SRS != nil {
		dto.SRS = toSRSDTO(model.SRS)
	}
	if tags != nil {
		dto.Tags = toTagDTOs(tags)
	}
	return dto
}

func toCardSummaryDTO(model *models.Card) *CardSummaryDTO {
	if model == nil {
		return nil
	}
	return &CardSummaryDTO{
		ID:        model.ID,
		ParentID:  model.ParentID,
		Title:     model.Title,
		CardType:  model.CardType,
		Status:    model.Status,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}
}

func toCardSummaryDTOs(items []*models.Card) []*CardSummaryDTO {
	result := make([]*CardSummaryDTO, 0, len(items))
	for _, item := range items {
		result = append(result, toCardSummaryDTO(item))
	}
	return result
}
