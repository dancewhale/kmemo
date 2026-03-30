package storage

import (
	"fmt"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// openDB 根据 Options.Driver 选择 dialector；切换 MySQL/PostgreSQL 时在此增加分支并共用下方 gorm.Config 组装逻辑。
func openDB(opts Options) (*gorm.DB, error) {
	opts.normalize()
	switch opts.Driver {
	case "sqlite", "":
		return openSQLite(opts)
	default:
		return nil, fmt.Errorf("storage: unsupported driver %q (当前实现为 sqlite)", opts.Driver)
	}
}

// openSQLite 打开 SQLite 并应用 GORM 配置。
func openSQLite(opts Options) (*gorm.DB, error) {
	cfg := &gorm.Config{
		// 显式命名：表名单数、列名 snake_case，贴近本项目的 DDL 风格。
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "",
			SingularTable: true,
		},
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
		// SQLite 下 card.cover_asset_id 与 asset.card_id 形成环依赖，迁移时跳过后端 FK 约束创建；
		// 运行期仍通过 PRAGMA foreign_keys=ON 校验（若列级 FK 已写入 schema）。
		DisableForeignKeyConstraintWhenMigrating: true,
	}

	logLevel := logger.Warn
	switch opts.LogLevel {
	case "silent":
		logLevel = logger.Silent
	case "error":
		logLevel = logger.Error
	case "warn":
		logLevel = logger.Warn
	case "info":
		logLevel = logger.Info
	}

	cfg.Logger = logger.Default.LogMode(logLevel)

	dialector := sqlite.Open(opts.DSN)
	db, err := gorm.Open(dialector, cfg)
	if err != nil {
		return nil, fmt.Errorf("gorm open sqlite: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("gorm sql db: %w", err)
	}

	if opts.MaxOpenConns > 0 {
		sqlDB.SetMaxOpenConns(opts.MaxOpenConns)
	}
	if opts.MaxIdleConns > 0 {
		sqlDB.SetMaxIdleConns(opts.MaxIdleConns)
	}
	if opts.ConnMaxLifetime > 0 {
		sqlDB.SetConnMaxLifetime(opts.ConnMaxLifetime)
	}

	// 运行期强制开启外键（与迁移期是否禁用无关）。
	if err := db.Exec("PRAGMA foreign_keys = ON").Error; err != nil {
		return nil, fmt.Errorf("sqlite pragma foreign_keys: %w", err)
	}
	// WAL 提升桌面端并发读；路径由 DSN file: 指定。
	if err := db.Exec("PRAGMA journal_mode = WAL").Error; err != nil {
		return nil, fmt.Errorf("sqlite pragma journal_mode: %w", err)
	}
	if err := db.Exec("PRAGMA busy_timeout = 5000").Error; err != nil {
		return nil, fmt.Errorf("sqlite pragma busy_timeout: %w", err)
	}

	return db, nil
}
