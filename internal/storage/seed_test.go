package storage_test

import (
	"path/filepath"
	"testing"

	"kmemo/internal/storage"
	"kmemo/internal/storage/models"
)

func TestSeedDefaultData_idempotent(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	dsn := "file:" + filepath.ToSlash(filepath.Join(dir, "seed.db")) + "?_pragma=foreign_keys(1)"

	s, err := storage.New(storage.Options{DSN: dsn}, true)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()

	var n int64
	if err := s.DB().Model(&models.Knowledge{}).Count(&n).Error; err != nil {
		t.Fatal(err)
	}
	if n != 1 {
		t.Fatalf("expected 1 default knowledge, got %d", n)
	}
	var name string
	if err := s.DB().Model(&models.Knowledge{}).
		Select("name").
		Where("id = ?", storage.DefaultSeedKnowledgeID).
		Scan(&name).Error; err != nil {
		t.Fatal(err)
	}
	if name != "默认知识库" {
		t.Fatalf("unexpected default name %q", name)
	}

	if err := storage.SeedDefaultData(s.DB()); err != nil {
		t.Fatal(err)
	}
	if err := s.DB().Model(&models.Knowledge{}).Count(&n).Error; err != nil {
		t.Fatal(err)
	}
	if n != 1 {
		t.Fatalf("after second seed, expected 1 row, got %d", n)
	}
}
