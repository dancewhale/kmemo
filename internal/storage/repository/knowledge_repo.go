package repository

import (
	"context"
	"kmemo/internal/storage/dao"
	"kmemo/internal/storage/models"
	"time"

	"gorm.io/gorm"
)

// KnowledgeRepository Knowledge 聚合根仓储
type KnowledgeRepository interface {
	Repository

	// Create 创建知识库
	Create(ctx context.Context, k *models.Knowledge) error

	// GetByID 根据 ID 获取
	GetByID(ctx context.Context, id string) (*models.Knowledge, error)

	// Update 更新知识库
	Update(ctx context.Context, k *models.Knowledge) error

	// Delete 软删除知识库
	Delete(ctx context.Context, id string) error

	// List 列表查询
	List(ctx context.Context, opts ListKnowledgeOptions) ([]*models.Knowledge, int64, error)

	// GetTree 获取树形结构（含子节点）
	GetTree(ctx context.Context, rootID string) (*models.Knowledge, error)

	// GetChildren 获取直接子节点
	GetChildren(ctx context.Context, parentID string) ([]*models.Knowledge, error)

	// Move 移动到新父节点
	Move(ctx context.Context, id string, newParentID *string) error

	// Archive 归档
	Archive(ctx context.Context, id string) error

	// Unarchive 取消归档
	Unarchive(ctx context.Context, id string) error
}

// ListKnowledgeOptions 查询选项
type ListKnowledgeOptions struct {
	ParentID        *string // 父节点过滤
	IncludeArchived bool    // 是否包含已归档
	Keyword         string  // 名称/描述关键词搜索
	Limit           int     // 分页大小
	Offset          int     // 分页偏移
	OrderBy         string  // 排序字段: name, created_at, updated_at
	OrderDesc       bool    // 是否降序
}

type knowledgeRepo struct {
	db *gorm.DB
}

// NewKnowledgeRepository 创建 Knowledge 仓储
func NewKnowledgeRepository(db *gorm.DB) KnowledgeRepository {
	return &knowledgeRepo{db: db}
}

func (r *knowledgeRepo) WithTx(tx *gorm.DB) Repository {
	return &knowledgeRepo{db: tx}
}

func (r *knowledgeRepo) Create(ctx context.Context, k *models.Knowledge) error {
	return dao.Use(r.db).Knowledge.WithContext(ctx).Create(k)
}

func (r *knowledgeRepo) GetByID(ctx context.Context, id string) (*models.Knowledge, error) {
	k, err := dao.Use(r.db).Knowledge.WithContext(ctx).Where(
		dao.Knowledge.ID.Eq(id),
	).First()
	if err != nil {
		return nil, convertError(err)
	}
	return k, nil
}

func (r *knowledgeRepo) Update(ctx context.Context, k *models.Knowledge) error {
	k.UpdatedAt = time.Now()
	_, err := dao.Use(r.db).Knowledge.WithContext(ctx).Where(
		dao.Knowledge.ID.Eq(k.ID),
	).Updates(k)
	return convertError(err)
}

func (r *knowledgeRepo) Delete(ctx context.Context, id string) error {
	now := time.Now()
	_, err := dao.Use(r.db).Knowledge.WithContext(ctx).Where(
		dao.Knowledge.ID.Eq(id),
	).Update(dao.Knowledge.ArchivedAt, &now)
	return convertError(err)
}

func (r *knowledgeRepo) List(ctx context.Context, opts ListKnowledgeOptions) ([]*models.Knowledge, int64, error) {
	q := dao.Use(r.db).Knowledge.WithContext(ctx)

	if opts.ParentID != nil {
		q = q.Where(dao.Knowledge.ParentID.Eq(*opts.ParentID))
	}
	if !opts.IncludeArchived {
		q = q.Where(dao.Knowledge.ArchivedAt.IsNull())
	}
	if opts.Keyword != "" {
		pattern := "%" + opts.Keyword + "%"
		q = q.Where(dao.Knowledge.Name.Like(pattern))
	}

	count, err := q.Count()
	if err != nil {
		return nil, 0, convertError(err)
	}

	if opts.OrderBy != "" {
		switch opts.OrderBy {
		case "name":
			if opts.OrderDesc {
				q = q.Order(dao.Knowledge.Name.Desc())
			} else {
				q = q.Order(dao.Knowledge.Name)
			}
		case "updated_at":
			if opts.OrderDesc {
				q = q.Order(dao.Knowledge.UpdatedAt.Desc())
			} else {
				q = q.Order(dao.Knowledge.UpdatedAt)
			}
		default:
			if opts.OrderDesc {
				q = q.Order(dao.Knowledge.CreatedAt.Desc())
			} else {
				q = q.Order(dao.Knowledge.CreatedAt)
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

func (r *knowledgeRepo) GetTree(ctx context.Context, rootID string) (*models.Knowledge, error) {
	k, err := dao.Use(r.db).Knowledge.WithContext(ctx).Preload(dao.Knowledge.Children).Where(
		dao.Knowledge.ID.Eq(rootID),
	).First()
	return k, convertError(err)
}

func (r *knowledgeRepo) GetChildren(ctx context.Context, parentID string) ([]*models.Knowledge, error) {
	list, err := dao.Use(r.db).Knowledge.WithContext(ctx).Where(
		dao.Knowledge.ParentID.Eq(parentID),
	).Find()
	return list, convertError(err)
}

func (r *knowledgeRepo) Move(ctx context.Context, id string, newParentID *string) error {
	_, err := dao.Use(r.db).Knowledge.WithContext(ctx).Where(
		dao.Knowledge.ID.Eq(id),
	).Update(dao.Knowledge.ParentID, newParentID)
	return convertError(err)
}

func (r *knowledgeRepo) Archive(ctx context.Context, id string) error {
	now := time.Now()
	_, err := dao.Use(r.db).Knowledge.WithContext(ctx).Where(
		dao.Knowledge.ID.Eq(id),
	).Update(dao.Knowledge.ArchivedAt, &now)
	return convertError(err)
}

func (r *knowledgeRepo) Unarchive(ctx context.Context, id string) error {
	_, err := dao.Use(r.db).Knowledge.WithContext(ctx).Where(
		dao.Knowledge.ID.Eq(id),
	).Update(dao.Knowledge.ArchivedAt, nil)
	return convertError(err)
}
