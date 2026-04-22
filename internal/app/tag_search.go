package app

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"kmemo/internal/storage/models"
	"kmemo/internal/storage/repository"
	"kmemo/internal/zaplog"
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
	ctx := d.actionContext()
	log := zaplog.L(ctx).Named("tag.api")
	if strings.TrimSpace(req.Name) == "" || strings.TrimSpace(req.Slug) == "" {
		log.Info("CreateTag rejected", zap.String("reason", "empty_name_or_slug"))
		return "", repository.ErrInvalidInput
	}
	now := time.Now().UTC()
	tag := &models.Tag{ID: uuid.NewString(), Name: strings.TrimSpace(req.Name), Slug: strings.TrimSpace(req.Slug), Color: req.Color, Icon: req.Icon, Description: req.Description, CreatedAt: now, UpdatedAt: now}
	if err := d.actions.Tag.Create(ctx, tag); err != nil {
		log.Info("CreateTag failed", zap.String("name", req.Name), zap.Error(err))
		return "", err
	}
	log.Info("CreateTag ok", zap.String("id", tag.ID), zap.String("name", tag.Name), zap.String("slug", tag.Slug))
	return tag.ID, nil
}

func (d *Desktop) GetTag(id string) (*TagDTO, error) {
	ctx := d.actionContext()
	log := zaplog.L(ctx).Named("tag.api")
	tag, err := d.actions.Tag.Get(ctx, id)
	if err != nil {
		log.Info("GetTag failed", zap.String("id", id), zap.Error(err))
		return nil, err
	}
	out := toTagDTO(tag)
	log.Info("GetTag ok", zap.String("id", out.ID), zap.String("name", out.Name), zap.String("slug", out.Slug), zap.Int("cardCount", out.CardCount))
	return out, nil
}

func (d *Desktop) ListTags() ([]*TagDTO, error) {
	ctx := d.actionContext()
	log := zaplog.L(ctx).Named("tag.api")
	items, err := d.actions.Tag.List(ctx)
	if err != nil {
		log.Info("ListTags failed", zap.Error(err))
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
	log.Info("ListTags ok", zap.Int("count", len(out)), zap.Strings("sample", sample))
	return out, nil
}

func (d *Desktop) UpdateTag(id string, req UpdateTagRequest) error {
	ctx := d.actionContext()
	log := zaplog.L(ctx).Named("tag.api")
	err := d.actions.Tag.Update(ctx, &models.Tag{ID: id, Name: strings.TrimSpace(req.Name), Color: req.Color, Icon: req.Icon, Description: req.Description})
	if err != nil {
		log.Info("UpdateTag failed", zap.String("id", id), zap.Error(err))
		return err
	}
	log.Info("UpdateTag ok",
		zap.String("id", id),
		zap.String("name", strings.TrimSpace(req.Name)),
		zap.String("description_excerpt", truncateRunes(req.Description, 80)),
	)
	return nil
}

func (d *Desktop) DeleteTag(id string) error {
	ctx := d.actionContext()
	log := zaplog.L(ctx).Named("tag.api")
	err := d.actions.Tag.Delete(ctx, id)
	if err != nil {
		log.Info("DeleteTag failed", zap.String("id", id), zap.Error(err))
		return err
	}
	log.Info("DeleteTag ok", zap.String("id", id))
	return nil
}

func (d *Desktop) SearchCardsByTags(tagIDs []string) ([]*CardDTO, error) {
	ctx := d.actionContext()
	log := zaplog.L(ctx).Named("tag.api")
	items, err := d.actions.Search.SearchCardsByTags(ctx, tagIDs)
	if err != nil {
		log.Info("SearchCardsByTags failed", zap.Strings("tag_ids", firstStrings(tagIDs, 8)), zap.Error(err))
		return nil, err
	}
	result := make([]*CardDTO, 0, len(items))
	for _, item := range items {
		result = append(result, toCardDTO(item, nil))
	}
	sample := make([]string, 0, min(5, len(result)))
	for i := range result {
		if i >= 5 {
			break
		}
		sample = append(sample, result[i].ID)
	}
	log.Info("SearchCardsByTags ok",
		zap.Int("tag_id_count", len(tagIDs)),
		zap.Strings("tag_ids_head", firstStrings(tagIDs, 8)),
		zap.Int("card_count", len(result)),
		zap.Strings("sample_card_ids", sample),
	)
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
