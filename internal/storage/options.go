package storage

import "time"

// Options 配置存储层；DSN 使用各 driver 标准格式，便于切换 MySQL / PostgreSQL。
//
// SQLite 示例：
//   file:/path/app.db?_pragma=foreign_keys(1)&_pragma=busy_timeout(5000)&_pragma=journal_mode(WAL)
// MySQL（未来）：
//   user:pass@tcp(127.0.0.1:3306)/kmemo?charset=utf8mb4&parseTime=True&loc=Local
// PostgreSQL（未来）：
//   host=localhost user=kmemo password=secret dbname=kmemo sslmode=disable TimeZone=UTC
type Options struct {
	// Driver 默认 sqlite；切换数据库时改为 mysql / postgres 并配合 DSN。
	Driver string
	// DSN 数据源名称，由 gorm.io/driver/* 解析。
	DSN string

	// LogLevel GORM 日志级别：Silent / Error / Warn / Info
	LogLevel string

	// SlowThreshold 慢查询阈值；0 表示使用 GORM 默认。
	SlowThreshold time.Duration

	// MaxOpenConns / MaxIdleConns / ConnMaxLifetime 透传给 sql.DB；SQLite 通常 MaxOpenConns=1。
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

func (o *Options) normalize() {
	if o.Driver == "" {
		o.Driver = "sqlite"
	}
	if o.LogLevel == "" {
		o.LogLevel = "warn"
	}
	if o.MaxOpenConns == 0 && o.Driver == "sqlite" {
		o.MaxOpenConns = 1
	}
}
