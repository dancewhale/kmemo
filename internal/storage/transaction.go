package storage

import (
	"context"

	"gorm.io/gorm"
)

// Transaction 在单连接上执行回调；callback 内请使用传入的 tx，勿混用外层 DB 以免死锁或未提交读。
func (s *Storage) Transaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
	return s.db.WithContext(ctx).Transaction(fn)
}
