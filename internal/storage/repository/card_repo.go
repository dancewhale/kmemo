package repository

import (
	"context"
	"kmemo/internal/storage/dao"
	"kmemo/internal/storage/models"
	"time"

	"gorm.io/gorm"
)

// CardRepository Card 聚合根仓储
type CardRepository interface {
	Repository

	// Create 创建卡片
	Create(ctx context.Context, c *models.Card) error

	// GetByID 根据 ID 获取（可选预加载关联）
	GetByID(ctx context.Context, id string, preload ...string) (*models.Card, error)

	// Update 更新卡片
	Update(ctx context.Context, c *models.Card) error

	// Delete 软删除卡片
	Delete(ctx context.Context, id string) error

	// List 列表查询
	List(ctx context.Context, opts ListCardOptions) ([]*models.Card, int64, error)

	// GetByKnowledge 获取指定知识库下的卡片
	GetByKnowledge(ctx context.Context, knowledgeID string, opts ListCardOptions) ([]*models.Card, int64, error)

	// GetChildren 获取子卡片
	GetChildren(ctx context.Context, parentID string) ([]*models.Card, error)

	// UpdateStatus 批量更新状态
	UpdateStatus(ctx context.Context, ids []string, status string) error

	// AddTags 为卡片添加标签
	AddTags(ctx context.Context, cardID string, tagIDs []string) error

	// RemoveTags 移除卡片标签
	RemoveTags(ctx context.Context, cardID string, tagIDs []string) error

	// GetTags 获取卡片的所有标签
	GetTags(ctx context.Context, cardID string) ([]*models.Tag, error)
}

// ListCardOptions 卡片查询选项
type ListCardOptions struct {
	KnowledgeID      *string  // 知识库过滤
	SourceDocumentID *string  // 来源文档过滤
	ParentID         *string  // 父卡片过滤
	CardType         string   // 卡片类型过滤
	Status           string   // 状态过滤
	TagIDs           []string // 标签过滤（AND 关系）
	Keyword          string   // 标题关键词搜索
	IsRoot           *bool    // 是否根卡片
	Limit            int
	Offset           int
	OrderBy          string   // title, created_at, updated_at, sort_order
	OrderDesc        bool
	Preload          []string // 预加载关联: Tags, SRS, Assets, Knowledge
}

type cardRepo struct {
	db *gorm.DB
}

// NewCardRepository 创建 Card 仓储
func NewCardRepository(db *gorm.DB) CardRepository {
	return &cardRepo{db: db}
}

func (r *cardRepo) WithTx(tx *gorm.DB) Repository {
	return &cardRepo{db: tx}
}

func (r *cardRepo) Create(ctx context.Context, c *models.Card) error {
	return dao.Use(r.db).Card.WithContext(ctx).Create(c)
}

func (r *cardRepo) GetByID(ctx context.Context, id string, preload ...string) (*models.Card, error) {
	q := dao.Use(r.db).Card.WithContext(ctx)
	for _, p := range preload {
		switch p {
		case "Tags":
			q = q.Preload(dao.Card.CardTags)
		case "SRS":
			q = q.Preload(dao.Card.SRS)
		case "Assets":
			q = q.Preload(dao.Card.Assets)
		case "Knowledge":
			q = q.Preload(dao.Card.Knowledge)
		}
	}
	c, err := q.Where(dao.Card.ID.Eq(id)).First()
	return c, convertError(err)
}

func (r *cardRepo) Update(ctx context.Context, c *models.Card) error {
	c.UpdatedAt = time.Now()
	_, err := dao.Use(r.db).Card.WithContext(ctx).Where(dao.Card.ID.Eq(c.ID)).Updates(c)
	return convertError(err)
}

func (r *cardRepo) Delete(ctx context.Context, id string) error {
	now := time.Now()
	_, err := dao.Use(r.db).Card.WithContext(ctx).Where(dao.Card.ID.Eq(id)).Update(dao.Card.DeletedAt, &now)
	return convertError(err)
}

func (r *cardRepo) buildCardQuery(ctx context.Context, opts ListCardOptions) (dao.ICardDo, error) {
	q := dao.Use(r.db).Card.WithContext(ctx)

	if opts.KnowledgeID != nil {
		q = q.Where(dao.Card.KnowledgeID.Eq(*opts.KnowledgeID))
	}
	if opts.SourceDocumentID != nil {
		q = q.Where(dao.Card.SourceDocumentID.Eq(*opts.SourceDocumentID))
	}
	if opts.ParentID != nil {
		q = q.Where(dao.Card.ParentID.Eq(*opts.ParentID))
	}
	if opts.CardType != "" {
		q = q.Where(dao.Card.CardType.Eq(opts.CardType))
	}
	if opts.Status != "" {
		q = q.Where(dao.Card.Status.Eq(opts.Status))
	}
	if len(opts.TagIDs) > 0 {
		cardIDs := make([]string, 0)
		err := r.db.WithContext(ctx).
			Model(&models.CardTag{}).
			Select("card_id").
			Where("tag_id IN ?", opts.TagIDs).
			Group("card_id").
			Having("COUNT(DISTINCT tag_id) = ?", len(opts.TagIDs)).
			Pluck("card_id", &cardIDs).Error
		if err != nil {
			return nil, convertError(err)
		}
		if len(cardIDs) == 0 {
			q = q.Where(dao.Card.ID.Eq("__no_matching_card__"))
		} else {
			q = q.Where(dao.Card.ID.In(cardIDs...))
		}
	}
	if opts.IsRoot != nil {
		q = q.Where(dao.Card.IsRoot.Is(*opts.IsRoot))
	}
	if opts.Keyword != "" {
		q = q.Where(dao.Card.Title.Like("%" + opts.Keyword + "%"))
	}

	for _, p := range opts.Preload {
		switch p {
		case "Tags":
			q = q.Preload(dao.Card.CardTags)
		case "SRS":
			q = q.Preload(dao.Card.SRS)
		case "Assets":
			q = q.Preload(dao.Card.Assets)
		case "Knowledge":
			q = q.Preload(dao.Card.Knowledge)
		}
	}

	return q, nil
}

func (r *cardRepo) List(ctx context.Context, opts ListCardOptions) ([]*models.Card, int64, error) {
	q, err := r.buildCardQuery(ctx, opts)
	if err != nil {
		return nil, 0, err
	}

	count, err := q.Count()
	if err != nil {
		return nil, 0, convertError(err)
	}

	if opts.OrderBy != "" {
		switch opts.OrderBy {
		case "title":
			if opts.OrderDesc {
				q = q.Order(dao.Card.Title.Desc())
			} else {
				q = q.Order(dao.Card.Title)
			}
		case "updated_at":
			if opts.OrderDesc {
				q = q.Order(dao.Card.UpdatedAt.Desc())
			} else {
				q = q.Order(dao.Card.UpdatedAt)
			}
		case "sort_order":
			if opts.OrderDesc {
				q = q.Order(dao.Card.SortOrder.Desc())
			} else {
				q = q.Order(dao.Card.SortOrder)
			}
		default:
			if opts.OrderDesc {
				q = q.Order(dao.Card.CreatedAt.Desc())
			} else {
				q = q.Order(dao.Card.CreatedAt)
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

func (r *cardRepo) GetByKnowledge(ctx context.Context, knowledgeID string, opts ListCardOptions) ([]*models.Card, int64, error) {
	opts.KnowledgeID = &knowledgeID
	return r.List(ctx, opts)
}

func (r *cardRepo) GetChildren(ctx context.Context, parentID string) ([]*models.Card, error) {
	list, err := dao.Use(r.db).Card.WithContext(ctx).Where(dao.Card.ParentID.Eq(parentID)).Find()
	return list, convertError(err)
}

func (r *cardRepo) UpdateStatus(ctx context.Context, ids []string, status string) error {
	_, err := dao.Use(r.db).Card.WithContext(ctx).Where(dao.Card.ID.In(ids...)).Update(dao.Card.Status, status)
	return convertError(err)
}

func (r *cardRepo) AddTags(ctx context.Context, cardID string, tagIDs []string) error {
	for _, tagID := range tagIDs {
		ct := &models.CardTag{
			CardID:    cardID,
			TagID:     tagID,
			CreatedAt: time.Now(),
		}
		if err := dao.Use(r.db).CardTag.WithContext(ctx).Create(ct); err != nil {
			return convertError(err)
		}
	}
	return nil
}

func (r *cardRepo) RemoveTags(ctx context.Context, cardID string, tagIDs []string) error {
	_, err := dao.Use(r.db).CardTag.WithContext(ctx).Where(
		dao.CardTag.CardID.Eq(cardID),
		dao.CardTag.TagID.In(tagIDs...),
	).Delete()
	return convertError(err)
}

func (r *cardRepo) GetTags(ctx context.Context, cardID string) ([]*models.Tag, error) {
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
