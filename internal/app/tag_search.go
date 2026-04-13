package app

import (
	"strings"
	"time"

	"github.com/google/uuid"

	"kmemo/internal/storage/models"
	"kmemo/internal/storage/repository"
)

type TagDTO struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Color       string    `json:"color"`
	Icon        string    `json:"icon"`
	Description string    `json:"description"`
	CardCount   int       `json:"cardCount"`
	CreatedAt   time.Time `json:"createdAt"`
}

type CreateTagRequest struct {
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Color       string `json:"color"`
	Icon        string `json:"icon"`
	Description string `json:"description"`
}

type UpdateTagRequest struct {
	Name        string `json:"name"`
	Color       string `json:"color"`
	Icon        string `json:"icon"`
	Description string `json:"description"`
}

func (d *Desktop) CreateTag(req CreateTagRequest) (string, error) {
	if strings.TrimSpace(req.Name) == "" || strings.TrimSpace(req.Slug) == "" {
		return "", repository.ErrInvalidInput
	}
	ctx := d.actionContext()
	now := time.Now().UTC()
	tag := &models.Tag{ID: uuid.NewString(), Name: strings.TrimSpace(req.Name), Slug: strings.TrimSpace(req.Slug), Color: req.Color, Icon: req.Icon, Description: req.Description, CreatedAt: now, UpdatedAt: now}
	if err := d.actions.Tag.Create(ctx, tag); err != nil {
		return "", err
	}
	return tag.ID, nil
}

func (d *Desktop) GetTag(id string) (*TagDTO, error) {
	ctx := d.actionContext()
	tag, err := d.actions.Tag.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return toTagDTO(tag), nil
}

func (d *Desktop) ListTags() ([]*TagDTO, error) {
	ctx := d.actionContext()
	items, err := d.actions.Tag.List(ctx)
	if err != nil {
		return nil, err
	}
	return toTagDTOs(items), nil
}

func (d *Desktop) UpdateTag(id string, req UpdateTagRequest) error {
	ctx := d.actionContext()
	return d.actions.Tag.Update(ctx, &models.Tag{ID: id, Name: strings.TrimSpace(req.Name), Color: req.Color, Icon: req.Icon, Description: req.Description})
}

func (d *Desktop) DeleteTag(id string) error {
	ctx := d.actionContext()
	return d.actions.Tag.Delete(ctx, id)
}

func (d *Desktop) SearchCardsByTags(tagIDs []string) ([]*CardDTO, error) {
	ctx := d.actionContext()
	items, err := d.actions.Search.SearchCardsByTags(ctx, tagIDs)
	if err != nil {
		return nil, err
	}
	result := make([]*CardDTO, 0, len(items))
	for _, item := range items {
		result = append(result, toCardDTO(item, nil))
	}
	return result, nil
}

func toTagDTO(model *models.Tag) *TagDTO {
	if model == nil {
		return nil
	}
	return &TagDTO{ID: model.ID, Name: model.Name, Slug: model.Slug, Color: model.Color, Icon: model.Icon, Description: model.Description, CardCount: model.CardCount, CreatedAt: model.CreatedAt}
}

func toTagDTOs(items []*models.Tag) []*TagDTO {
	result := make([]*TagDTO, 0, len(items))
	for _, item := range items {
		result = append(result, toTagDTO(item))
	}
	return result
}
