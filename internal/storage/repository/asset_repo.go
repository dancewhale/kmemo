package repository

import (
	"context"
	"kmemo/internal/storage/dao"
	"kmemo/internal/storage/models"
	"time"

	"gorm.io/gorm"
)

// AssetRepository Asset 仓储
type AssetRepository interface {
	Repository

	Create(ctx context.Context, a *models.Asset) error
	GetByID(ctx context.Context, id string) (*models.Asset, error)
	GetByCardID(ctx context.Context, cardID string, kind ...string) ([]*models.Asset, error)
	GetByChecksum(ctx context.Context, checksum string) (*models.Asset, error)
	Update(ctx context.Context, a *models.Asset) error
	Delete(ctx context.Context, id string) error
	MarkMissing(ctx context.Context, id string) error
}

type assetRepo struct {
	db *gorm.DB
}

func NewAssetRepository(db *gorm.DB) AssetRepository {
	return &assetRepo{db: db}
}

func (r *assetRepo) WithTx(tx *gorm.DB) Repository {
	return &assetRepo{db: tx}
}

func (r *assetRepo) Create(ctx context.Context, a *models.Asset) error {
	return dao.Use(r.db).Asset.WithContext(ctx).Create(a)
}

func (r *assetRepo) GetByID(ctx context.Context, id string) (*models.Asset, error) {
	a, err := dao.Use(r.db).Asset.WithContext(ctx).Where(dao.Asset.ID.Eq(id)).First()
	return a, convertError(err)
}

func (r *assetRepo) GetByCardID(ctx context.Context, cardID string, kind ...string) ([]*models.Asset, error) {
	q := dao.Use(r.db).Asset.WithContext(ctx).Where(dao.Asset.CardID.Eq(cardID))
	if len(kind) > 0 {
		q = q.Where(dao.Asset.Kind.In(kind...))
	}
	list, err := q.Find()
	return list, convertError(err)
}

func (r *assetRepo) GetByChecksum(ctx context.Context, checksum string) (*models.Asset, error) {
	a, err := dao.Use(r.db).Asset.WithContext(ctx).Where(dao.Asset.Checksum.Eq(checksum)).First()
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return a, convertError(err)
}

func (r *assetRepo) Update(ctx context.Context, a *models.Asset) error {
	a.UpdatedAt = time.Now()
	_, err := dao.Use(r.db).Asset.WithContext(ctx).Where(dao.Asset.ID.Eq(a.ID)).Updates(a)
	return convertError(err)
}

func (r *assetRepo) Delete(ctx context.Context, id string) error {
	now := time.Now()
	_, err := dao.Use(r.db).Asset.WithContext(ctx).Where(dao.Asset.ID.Eq(id)).Update(dao.Asset.DeletedAt, &now)
	return convertError(err)
}

func (r *assetRepo) MarkMissing(ctx context.Context, id string) error {
	_, err := dao.Use(r.db).Asset.WithContext(ctx).Where(dao.Asset.ID.Eq(id)).Update(dao.Asset.Status, "missing")
	return convertError(err)
}
