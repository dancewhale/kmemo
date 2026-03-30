package storage_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"kmemo/internal/storage"
)

func ExampleNew() {
	dir, err := os.MkdirTemp("", "kmemo-storage-*")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dir)

	dbPath := filepath.Join(dir, "test.db")
	dsn := "file:" + filepath.ToSlash(dbPath) + "?_pragma=foreign_keys(1)"

	s, err := storage.New(storage.Options{
		Driver: "sqlite",
		DSN:    dsn,
	}, true)
	if err != nil {
		panic(err)
	}
	defer s.Close()

	// 业务代码中可在此使用 dao 包：dao.Card.WithContext(ctx).Where(...).First()
	_ = s.DB()
	fmt.Println("ok")

	// Output:
	// ok
}

func TestAutoMigrate_idempotent(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "t.db")
	dsn := "file:" + filepath.ToSlash(dbPath) + "?_pragma=foreign_keys(1)"

	s, err := storage.New(storage.Options{DSN: dsn}, true)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	if err := s.AutoMigrate(); err != nil {
		t.Fatal(err)
	}
}
