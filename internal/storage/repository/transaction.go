package repository

import (
	"context"

	"gorm.io/gorm"
)

// TransactionManager 事务管理器
type TransactionManager interface {
	// BeginTx 开启事务
	BeginTx(ctx context.Context) (*gorm.DB, error)

	// WithTx 在事务中执行函数
	WithTx(ctx context.Context, fn func(tx *gorm.DB) error) error
}

type transactionManager struct {
	db *gorm.DB
}

// NewTransactionManager 创建事务管理器
func NewTransactionManager(db *gorm.DB) TransactionManager {
	return &transactionManager{db: db}
}

func (tm *transactionManager) BeginTx(ctx context.Context) (*gorm.DB, error) {
	return tm.db.WithContext(ctx).Begin(), nil
}

func (tm *transactionManager) WithTx(ctx context.Context, fn func(tx *gorm.DB) error) error {
	return tm.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(tx)
	})
}
