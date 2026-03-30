// 小型 CLI：供 task db:migrate 调用，对本地 SQLite 执行 AutoMigrate。
package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"kmemo/internal/storage"
)

func main() {
	dbPath := os.Getenv("KMEMO_DB_PATH")
	if dbPath == "" {
		dbPath = "kmemo.dev.db"
	}
	if !filepath.IsAbs(dbPath) {
		wd, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		dbPath = filepath.Join(wd, dbPath)
	}

	dsn := fmt.Sprintf(
		"file:%s?_pragma=foreign_keys(1)&_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)",
		filepath.ToSlash(dbPath),
	)

	s, err := storage.New(storage.Options{
		Driver:   "sqlite",
		DSN:      dsn,
		LogLevel: "warn",
	}, true)
	if err != nil {
		log.Fatalf("migrate: %v", err)
	}
	defer s.Close()

	log.Printf("storage migrate ok: %s\n", dbPath)
}
