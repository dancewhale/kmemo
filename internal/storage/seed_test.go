package storage_test

import (
	"path/filepath"
	"testing"

	"github.com/google/uuid"
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
	var row models.Knowledge
	if err := s.DB().Model(&models.Knowledge{}).
		Select("id", "name").
		Limit(1).
		Take(&row).Error; err != nil {
		t.Fatal(err)
	}
	if _, err := uuid.Parse(row.ID); err != nil {
		t.Fatalf("seed knowledge id is not a valid UUID: %q", row.ID)
	}
	if got := uuid.MustParse(row.ID).Version(); got != 7 {
		t.Fatalf("seed knowledge id UUID version = %d, want 7", got)
	}
	if row.Name != "默认知识库" {
		t.Fatalf("unexpected default name %q", row.Name)
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
