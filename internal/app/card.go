package app

import (
	"strings"
	"time"

	"kmemo/internal/actions/card"
	"kmemo/internal/storage/models"
	"kmemo/internal/storage/repository"
)

type CardDTO struct {
	ID               string     `json:"id"`
	KnowledgeID      string     `json:"knowledgeId"`
	KnowledgeName    string     `json:"knowledgeName"`
	SourceDocumentID *string    `json:"sourceDocumentId"`
	ParentID         *string    `json:"parentId"`
	Title            string     `json:"title"`
	CardType         string     `json:"cardType"`
	HTMLPath         string     `json:"htmlPath"`
	HTMLContent      string     `json:"htmlContent"`
	Status           string     `json:"status"`
	Tags             []*TagDTO  `json:"tags"`
	SRS              *SRSDTO    `json:"srs"`
	CreatedAt        time.Time  `json:"createdAt"`
	UpdatedAt        time.Time  `json:"updatedAt"`
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

type CardFilters struct {
	KnowledgeID *string  `json:"knowledgeId"`
	CardType    string   `json:"cardType"`
	Status      string   `json:"status"`
	TagIDs      []string `json:"tagIds"`
	Keyword     string   `json:"keyword"`
	ParentID    *string  `json:"parentId"`
	IsRoot      *bool    `json:"isRoot"`
	OrderBy     string   `json:"orderBy"`   // title, created_at, updated_at, sort_order；空则按仓储默认（无显式排序）
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

func (d *Desktop) CreateCard(req CreateCardRequest) (string, error) {
	if strings.TrimSpace(req.KnowledgeID) == "" || strings.TrimSpace(req.Title) == "" || strings.TrimSpace(req.CardType) == "" {
		return "", repository.ErrInvalidInput
	}
	ctx := d.actionContext()
	return d.actions.Card.Create(ctx, card.CreateInput{
		KnowledgeID:      strings.TrimSpace(req.KnowledgeID),
		SourceDocumentID: req.SourceDocumentID,
		ParentID:         req.ParentID,
		Title:            strings.TrimSpace(req.Title),
		CardType:         strings.TrimSpace(req.CardType),
		HTMLContent:      req.HTMLContent,
		TagIDs:           req.TagIDs,
	})
}

func (d *Desktop) GetCard(id string) (*CardDTO, error) {
	ctx := d.actionContext()
	out, err := d.actions.Card.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	dto := toCardDTO(out.Card, nil)
	dto.HTMLContent = out.HTMLContent
	return dto, nil
}

func (d *Desktop) GetCardDetail(id string) (*CardDetailDTO, error) {
	ctx := d.actionContext()
	out, err := d.actions.Card.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	tags, err := d.actions.Card.GetTags(ctx, id)
	if err != nil {
		return nil, err
	}
	detail := &CardDetailDTO{CardDTO: *toCardDTO(out.Card, tags)}
	detail.HTMLContent = out.HTMLContent

	if out.Card.ParentID != nil {
		parent, err := d.actions.Card.Get(ctx, *out.Card.ParentID)
		if err != nil {
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
		return nil, err
	}
	detail.Children = toCardSummaryDTOs(children)
	return detail, nil
}

func (d *Desktop) ListCards(filters CardFilters) ([]*CardDTO, int64, error) {
	ctx := d.actionContext()
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
		return nil, 0, err
	}
	result := make([]*CardDTO, 0, len(items))
	for _, item := range items {
		tags, _ := d.actions.Card.GetTags(ctx, item.ID)
		result = append(result, toCardDTO(item, tags))
	}
	return result, total, nil
}

func (d *Desktop) GetCardChildren(parentID string) ([]*CardDTO, error) {
	ctx := d.actionContext()
	items, _, err := d.actions.Card.List(ctx, repository.ListCardOptions{
		ParentID:  &parentID,
		Preload:   []string{"SRS", "Knowledge"},
		OrderBy:   "sort_order",
		OrderDesc: false,
	})
	if err != nil {
		return nil, err
	}
	result := make([]*CardDTO, 0, len(items))
	for _, item := range items {
		tags, _ := d.actions.Card.GetTags(ctx, item.ID)
		result = append(result, toCardDTO(item, tags))
	}
	return result, nil
}

func (d *Desktop) UpdateCard(id string, req UpdateCardRequest) error {
	if strings.TrimSpace(req.Title) == "" {
		return repository.ErrInvalidInput
	}
	ctx := d.actionContext()
	return d.actions.Card.Update(ctx, id, card.UpdateInput{
		Title:       strings.TrimSpace(req.Title),
		HTMLContent: req.HTMLContent,
		Status:      strings.TrimSpace(req.Status),
	})
}

func (d *Desktop) DeleteCard(id string) error {
	ctx := d.actionContext()
	return d.actions.Card.Delete(ctx, id)
}

func (d *Desktop) AddCardTags(cardID string, tagIDs []string) error {
	ctx := d.actionContext()
	return d.actions.Card.AddTags(ctx, cardID, tagIDs)
}

func (d *Desktop) RemoveCardTags(cardID string, tagIDs []string) error {
	ctx := d.actionContext()
	return d.actions.Card.RemoveTags(ctx, cardID, tagIDs)
}

func (d *Desktop) GetCardTags(cardID string) ([]*TagDTO, error) {
	ctx := d.actionContext()
	items, err := d.actions.Card.GetTags(ctx, cardID)
	if err != nil {
		return nil, err
	}
	return toTagDTOs(items), nil
}

func (d *Desktop) SuspendCard(cardID string) error {
	ctx := d.actionContext()
	return d.actions.Card.Suspend(ctx, cardID)
}

func (d *Desktop) ResumeCard(cardID string) error {
	ctx := d.actionContext()
	return d.actions.Card.Resume(ctx, cardID)
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
