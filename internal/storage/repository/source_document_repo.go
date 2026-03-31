package repository

import (
	"context"
	"kmemo/internal/storage/dao"
	"kmemo/internal/storage/models"
	"time"

	"gorm.io/gorm"
)

// SourceDocumentRepository 来源文档仓储
type SourceDocumentRepository interface {
	Repository

	Create(ctx context.Context, sd *models.SourceDocument) error
	GetByID(ctx context.Context, id string) (*models.SourceDocument, error)
	GetByKnowledgeID(ctx context.Context, knowledgeID string) ([]*models.SourceDocument, error)
	Update(ctx context.Context, sd *models.SourceDocument) error
	Delete(ctx context.Context, id string) error
}

type sourceDocumentRepo struct {
	db *gorm.DB
}

func NewSourceDocumentRepository(db *gorm.DB) SourceDocumentRepository {
	return &sourceDocumentRepo{db: db}
}

func (r *sourceDocumentRepo) WithTx(tx *gorm.DB) Repository {
	return &sourceDocumentRepo{db: tx}
}

func (r *sourceDocumentRepo) Create(ctx context.Context, sd *models.SourceDocument) error {
	return dao.Use(r.db).SourceDocument.WithContext(ctx).Create(sd)
}

func (r *sourceDocumentRepo) GetByID(ctx context.Context, id string) (*models.SourceDocument, error) {
	sd, err := dao.Use(r.db).SourceDocument.WithContext(ctx).Where(dao.SourceDocument.ID.Eq(id)).First()
	return sd, convertError(err)
}

func (r *sourceDocumentRepo) GetByKnowledgeID(ctx context.Context, knowledgeID string) ([]*models.SourceDocument, error) {
	list, err := dao.Use(r.db).SourceDocument.WithContext(ctx).Where(dao.SourceDocument.KnowledgeID.Eq(knowledgeID)).Find()
	return list, convertError(err)
}

func (r *sourceDocumentRepo) Update(ctx context.Context, sd *models.SourceDocument) error {
	sd.UpdatedAt = time.Now()
	_, err := dao.Use(r.db).SourceDocument.WithContext(ctx).Where(dao.SourceDocument.ID.Eq(sd.ID)).Updates(sd)
	return convertError(err)
}

func (r *sourceDocumentRepo) Delete(ctx context.Context, id string) error {
	_, err := dao.Use(r.db).SourceDocument.WithContext(ctx).Where(dao.SourceDocument.ID.Eq(id)).Delete()
	return convertError(err)
}
