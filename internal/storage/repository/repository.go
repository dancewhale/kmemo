// Package repository 提供领域对象的持久化操作接口
package repository

import (
	"gorm.io/gorm"
)

// Repository 所有仓储的基础接口
type Repository interface {
	// WithTx 返回使用指定事务的仓储实例
	WithTx(tx *gorm.DB) Repository
}

// Transaction 事务接口
type Transaction interface {
	Commit() error
	Rollback() error
}
