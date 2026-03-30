package storage

// Schema 演进建议：
//  1. 在 models/ 中增改字段与索引标签，调用 AutoMigrate（幂等，适合多数迭代）。
//  2. 模型结构体变更后运行 go run ./internal/storage/gen（或 task db:gen）以更新 dao。
//  3. 复杂数据回填、多步变更可引入 gormigrate / atlas / 手写 SQL，与 AutoMigrate 组合使用。
import (
	"fmt"

	"gorm.io/gorm"

	"kmemo/internal/storage/models"
)

// migrateModels 返回 AutoMigrate 的模型列表；顺序尽量满足外键依赖（GORM 会多次尝试）。
func migrateModels() []any {
	return []any{
		&models.FSRSParameter{},
		&models.Knowledge{},
		&models.SourceDocument{},
		&models.Card{},
		&models.Asset{},
		&models.Tag{},
		&models.CardTag{},
		&models.SearchIndexState{},
		&models.CardSRS{},
		&models.ReviewLog{},
	}
}

// AutoMigrate 对当前连接执行 schema 迁移；幂等，可重复调用。
func AutoMigrate(db *gorm.DB) error {
	if err := db.AutoMigrate(migrateModels()...); err != nil {
		return fmt.Errorf("automigrate: %w", err)
	}
	return nil
}
