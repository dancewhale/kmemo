package repository

import (
	"context"
	"kmemo/internal/storage/dao"
	"kmemo/internal/storage/models"
	"time"

	"gorm.io/gorm"
)

// FSRSParameterRepository FSRS 参数仓储
type FSRSParameterRepository interface {
	Repository

	Create(ctx context.Context, param *models.FSRSParameter) error
	GetByID(ctx context.Context, id string) (*models.FSRSParameter, error)
	GetByName(ctx context.Context, name string) (*models.FSRSParameter, error)
	Update(ctx context.Context, param *models.FSRSParameter) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, opts ListFSRSParameterOptions) ([]*models.FSRSParameter, int64, error)
	GetDefault(ctx context.Context) (*models.FSRSParameter, error)
}

type ListFSRSParameterOptions struct {
	Keyword   string
	Limit     int
	Offset    int
	OrderBy   string
	OrderDesc bool
}

type fsrsParameterRepo struct {
	db *gorm.DB
}

func NewFSRSParameterRepository(db *gorm.DB) FSRSParameterRepository {
	return &fsrsParameterRepo{db: db}
}

func (r *fsrsParameterRepo) WithTx(tx *gorm.DB) Repository {
	return &fsrsParameterRepo{db: tx}
}

func (r *fsrsParameterRepo) Create(ctx context.Context, param *models.FSRSParameter) error {
	return dao.Use(r.db).FSRSParameter.WithContext(ctx).Create(param)
}

func (r *fsrsParameterRepo) GetByID(ctx context.Context, id string) (*models.FSRSParameter, error) {
	p, err := dao.Use(r.db).FSRSParameter.WithContext(ctx).Where(dao.FSRSParameter.ID.Eq(id)).First()
	return p, convertError(err)
}

func (r *fsrsParameterRepo) GetByName(ctx context.Context, name string) (*models.FSRSParameter, error) {
	p, err := dao.Use(r.db).FSRSParameter.WithContext(ctx).Where(dao.FSRSParameter.Name.Eq(name)).First()
	return p, convertError(err)
}

func (r *fsrsParameterRepo) Update(ctx context.Context, param *models.FSRSParameter) error {
	param.UpdatedAt = time.Now()
	_, err := dao.Use(r.db).FSRSParameter.WithContext(ctx).Where(dao.FSRSParameter.ID.Eq(param.ID)).Updates(param)
	return convertError(err)
}

func (r *fsrsParameterRepo) Delete(ctx context.Context, id string) error {
	_, err := dao.Use(r.db).FSRSParameter.WithContext(ctx).Where(dao.FSRSParameter.ID.Eq(id)).Delete()
	return convertError(err)
}

func (r *fsrsParameterRepo) List(ctx context.Context, opts ListFSRSParameterOptions) ([]*models.FSRSParameter, int64, error) {
	q := dao.Use(r.db).FSRSParameter.WithContext(ctx)

	if opts.Keyword != "" {
		q = q.Where(dao.FSRSParameter.Name.Like("%" + opts.Keyword + "%"))
	}

	count, err := q.Count()
	if err != nil {
		return nil, 0, convertError(err)
	}

	if opts.OrderBy != "" {
		switch opts.OrderBy {
		case "name":
			if opts.OrderDesc {
				q = q.Order(dao.FSRSParameter.Name.Desc())
			} else {
				q = q.Order(dao.FSRSParameter.Name)
			}
		case "created_at":
			if opts.OrderDesc {
				q = q.Order(dao.FSRSParameter.CreatedAt.Desc())
			} else {
				q = q.Order(dao.FSRSParameter.CreatedAt)
			}
		case "updated_at":
			if opts.OrderDesc {
				q = q.Order(dao.FSRSParameter.UpdatedAt.Desc())
			} else {
				q = q.Order(dao.FSRSParameter.UpdatedAt)
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

func (r *fsrsParameterRepo) GetDefault(ctx context.Context) (*models.FSRSParameter, error) {
	p, err := dao.Use(r.db).FSRSParameter.WithContext(ctx).Order(dao.FSRSParameter.CreatedAt).First()
	return p, convertError(err)
}
