package repository

import (
	"gorm.io/gorm"
)

// RepositoryFactory 仓储工厂
type RepositoryFactory interface {
	Knowledge() KnowledgeRepository
	Card() CardRepository
	Tag() TagRepository
	SRS() SRSRepository
	ReviewLog() ReviewLogRepository
	FSRSParameter() FSRSParameterRepository
	Asset() AssetRepository
	SourceDocument() SourceDocumentRepository
}

type repositoryFactory struct {
	db *gorm.DB
}

// NewRepositoryFactory 创建仓储工厂
func NewRepositoryFactory(db *gorm.DB) RepositoryFactory {
	return &repositoryFactory{db: db}
}

func (f *repositoryFactory) Knowledge() KnowledgeRepository {
	return NewKnowledgeRepository(f.db)
}

func (f *repositoryFactory) Card() CardRepository {
	return NewCardRepository(f.db)
}

func (f *repositoryFactory) Tag() TagRepository {
	return NewTagRepository(f.db)
}

func (f *repositoryFactory) SRS() SRSRepository {
	return NewSRSRepository(f.db)
}

func (f *repositoryFactory) ReviewLog() ReviewLogRepository {
	return NewReviewLogRepository(f.db)
}

func (f *repositoryFactory) FSRSParameter() FSRSParameterRepository {
	return NewFSRSParameterRepository(f.db)
}

func (f *repositoryFactory) Asset() AssetRepository {
	return NewAssetRepository(f.db)
}

func (f *repositoryFactory) SourceDocument() SourceDocumentRepository {
	return NewSourceDocumentRepository(f.db)
}
