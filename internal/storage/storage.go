// Package storage 提供本地优先的 SQLite 访问、GORM 配置、迁移与 GORM Gen DAO 初始化。
package storage

import (
	"fmt"

	"gorm.io/gorm"

	"kmemo/internal/storage/dao"
)

// Storage 统一持有 *gorm.DB 与生成 DAO 的生命周期。
type Storage struct {
	db *gorm.DB
}

// New 打开数据库、可选执行 AutoMigrate，并调用 dao.SetDefault 注册 Gen 查询入口。
func New(opts Options, withMigrate bool) (*Storage, error) {
	opts.normalize()
	if opts.DSN == "" {
		return nil, fmt.Errorf("storage: DSN is required")
	}

	db, err := openDB(opts)
	if err != nil {
		return nil, err
	}

	s := &Storage{db: db}
	if withMigrate {
		if err := AutoMigrate(db); err != nil {
			_ = s.Close()
			return nil, err
		}
	}
	s.initDAO()
	return s, nil
}

// Open 等价于 New(opts, false)，由调用方自行 AutoMigrate。
func Open(opts Options) (*Storage, error) {
	return New(opts, false)
}

// DB 返回底层 *gorm.DB，供自定义查询或传入 Gen 未封装场景。
func (s *Storage) DB() *gorm.DB {
	return s.db
}

// AutoMigrate 对 schema 做增量迁移。
func (s *Storage) AutoMigrate() error {
	return AutoMigrate(s.db)
}

func (s *Storage) initDAO() {
	dao.SetDefault(s.db)
}

// Close 关闭连接池。
func (s *Storage) Close() error {
	sqlDB, err := s.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
