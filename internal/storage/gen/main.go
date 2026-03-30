// 生成器入口：go run ./internal/storage/gen（或由 Task db:gen 调用）。
// 依赖 models 包中的结构体，输出到 internal/storage/dao。
package main

import (
	"path/filepath"
	"runtime"

	"gorm.io/driver/sqlite"
	"gorm.io/gen"
	"gorm.io/gorm"

	"kmemo/internal/storage/models"
)

func main() {
	_, file, _, _ := runtime.Caller(0)
	genDir := filepath.Dir(file)
	repoRoot := filepath.Clean(filepath.Join(genDir, "..", "..", ".."))
	out := filepath.Join(repoRoot, "internal", "storage", "dao")

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	g := gen.NewGenerator(gen.Config{
		OutPath: out,
		// 与 go.mod module 路径一致，供生成代码 import models。
		ModelPkgPath: "kmemo/internal/storage/models",
		// 默认 Query + 接口；未设置 WithoutContext 时 Gen 会生成带 context 的查询方法。
		Mode: gen.WithDefaultQuery | gen.WithQueryInterface,
	})

	g.UseDB(db)

	g.ApplyBasic(
		models.Knowledge{},
		models.SourceDocument{},
		models.Card{},
		models.Asset{},
		models.Tag{},
		models.CardTag{},
		models.SearchIndexState{},
		models.CardSRS{},
		models.FSRSParameter{},
		models.ReviewLog{},
	)

	g.Execute()
}
