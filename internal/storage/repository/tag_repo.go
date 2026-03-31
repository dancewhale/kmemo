package repository

import (
	"context"
	"kmemo/internal/storage/dao"
	"kmemo/internal/storage/models"
	"time"

	"gorm.io/gorm"
)

// TagRepository Tag 仓储
type TagRepository interface {
	Repository

	Create(ctx context.Context, t *models.Tag) error
	GetByID(ctx context.Context, id string) (*models.Tag, error)
	GetBySlug(ctx context.Context, slug string) (*models.Tag, error)
	Update(ctx context.Context, t *models.Tag) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, opts ListTagOptions) ([]*models.Tag, int64, error)
	GetByCardID(ctx context.Context, cardID string) ([]*models.Tag, error)
	UpdateCardCount(ctx context.Context, tagID string) error
}

type ListTagOptions struct {
	Keyword   string
	Limit     int
	Offset    int
	OrderBy   string
	OrderDesc bool
}

type tagRepo struct {
	db *gorm.DB
}

func NewTagRepository(db *gorm.DB) TagRepository {
	return &tagRepo{db: db}
}

func (r *tagRepo) WithTx(tx *gorm.DB) Repository {
	return &tagRepo{db: tx}
}

func (r *tagRepo) Create(ctx context.Context, t *models.Tag) error {
	return dao.Use(r.db).Tag.WithContext(ctx).Create(t)
}

func (r *tagRepo) GetByID(ctx context.Context, id string) (*models.Tag, error) {
	t, err := dao.Use(r.db).Tag.WithContext(ctx).Where(dao.Tag.ID.Eq(id)).First()
	return t, convertError(err)
}

func (r *tagRepo) GetBySlug(ctx context.Context, slug string) (*models.Tag, error) {
	t, err := dao.Use(r.db).Tag.WithContext(ctx).Where(dao.Tag.Slug.Eq(slug)).First()
	return t, convertError(err)
}

func (r *tagRepo) Update(ctx context.Context, t *models.Tag) error {
	t.UpdatedAt = time.Now()
	_, err := dao.Use(r.db).Tag.WithContext(ctx).Where(dao.Tag.ID.Eq(t.ID)).Updates(t)
	return convertError(err)
}

func (r *tagRepo) Delete(ctx context.Context, id string) error {
	now := time.Now()
	_, err := dao.Use(r.db).Tag.WithContext(ctx).Where(dao.Tag.ID.Eq(id)).Update(dao.Tag.DeletedAt, &now)
	return convertError(err)
}

func (r *tagRepo) List(ctx context.Context, opts ListTagOptions) ([]*models.Tag, int64, error) {
	q := dao.Use(r.db).Tag.WithContext(ctx)

	if opts.Keyword != "" {
		pattern := "%" + opts.Keyword + "%"
		q = q.Where(dao.Tag.Name.Like(pattern))
	}

	count, err := q.Count()
	if err != nil {
		return nil, 0, convertError(err)
	}

	if opts.OrderBy != "" {
		switch opts.OrderBy {
		case "sort_order":
			if opts.OrderDesc {
				q = q.Order(dao.Tag.SortOrder.Desc())
			} else {
				q = q.Order(dao.Tag.SortOrder)
			}
		case "card_count":
			if opts.OrderDesc {
				q = q.Order(dao.Tag.CardCount.Desc())
			} else {
				q = q.Order(dao.Tag.CardCount)
			}
		case "created_at":
			if opts.OrderDesc {
				q = q.Order(dao.Tag.CreatedAt.Desc())
			} else {
				q = q.Order(dao.Tag.CreatedAt)
			}
		default:
			if opts.OrderDesc {
				q = q.Order(dao.Tag.Name.Desc())
			} else {
				q = q.Order(dao.Tag.Name)
			}
		}
	}

	if opts.Limit > 0 {
		q = q.Limit(opts.Limit)
	}
	if opts.Offset > 0 {
		q = q.Offset(opts.Offset)
	}

	list, err := q.Find()
	return list, count, convertError(err)
}

func (r *tagRepo) GetByCardID(ctx context.Context, cardID string) ([]*models.Tag, error) {
	cardTags, err := dao.Use(r.db).CardTag.WithContext(ctx).Preload(dao.CardTag.Tag).Where(
		dao.CardTag.CardID.Eq(cardID),
	).Find()
	if err != nil {
		return nil, convertError(err)
	}

	tags := make([]*models.Tag, 0, len(cardTags))
	for _, ct := range cardTags {
		if ct.Tag != nil {
			tags = append(tags, ct.Tag)
		}
	}
	return tags, nil
}

func (r *tagRepo) UpdateCardCount(ctx context.Context, tagID string) error {
	count, err := dao.Use(r.db).CardTag.WithContext(ctx).Where(dao.CardTag.TagID.Eq(tagID)).Count()
	if err != nil {
		return convertError(err)
	}
	_, err = dao.Use(r.db).Tag.WithContext(ctx).Where(dao.Tag.ID.Eq(tagID)).Update(dao.Tag.CardCount, int(count))
	return convertError(err)
}
