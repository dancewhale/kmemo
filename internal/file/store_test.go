package file

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func setupTestStore(t *testing.T) (FileStore, string) {
	// Create temporary directory
	tmpDir, err := os.MkdirTemp("", "filestore-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}

	// Create file store
	config := DefaultConfig(tmpDir)
	store, err := NewFileStore(config)
	if err != nil {
		os.RemoveAll(tmpDir)
		t.Fatalf("failed to create file store: %v", err)
	}

	return store, tmpDir
}

func cleanupTestStore(tmpDir string) {
	os.RemoveAll(tmpDir)
}

func TestCreateFileObject(t *testing.T) {
	store, tmpDir := setupTestStore(t)
	defer cleanupTestStore(tmpDir)

	ctx := context.Background()

	// Test creating a card with slug
	input := CreateFileObjectInput{
		Kind: FileObjectKindCard,
		ID:   "0195f3f0-a8b7-7c8d-b1b9-d45c8f23a1ee",
		Slug: "test-card",
		Ext:  "html",
		Data: []byte("<h1>Test Card</h1>"),
	}

	loc, err := store.CreateFileObject(ctx, input)
	if err != nil {
		t.Fatalf("CreateFileObject failed: %v", err)
	}

	// Verify location
	if loc.Ref.Kind != FileObjectKindCard {
		t.Errorf("expected kind %s, got %s", FileObjectKindCard, loc.Ref.Kind)
	}
	if loc.Ref.ID != "0195f3f0-a8b7-7c8d-b1b9-d45c8f23a1ee" {
		t.Errorf("expected id %s, got %s", "0195f3f0-a8b7-7c8d-b1b9-d45c8f23a1ee", loc.Ref.ID)
	}
	if loc.Scope != FileObjectScopeActive {
		t.Errorf("expected scope %s, got %s", FileObjectScopeActive, loc.Scope)
	}
	if loc.Slug != "test-card" {
		t.Errorf("expected slug %s, got %s", "test-card", loc.Slug)
	}
	if loc.Ext != "html" {
		t.Errorf("expected ext %s, got %s", "html", loc.Ext)
	}

	// Verify file exists
	if _, err := os.Stat(loc.Path); os.IsNotExist(err) {
		t.Errorf("file does not exist at %s", loc.Path)
	}

	// Verify file content
	content, err := os.ReadFile(loc.Path)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}
	if string(content) != "<h1>Test Card</h1>" {
		t.Errorf("expected content %s, got %s", "<h1>Test Card</h1>", string(content))
	}

	// Test creating duplicate should fail
	_, err = store.CreateFileObject(ctx, input)
	if err == nil {
		t.Error("expected error when creating duplicate, got nil")
	}
}

func TestCreateFileObjectAsset(t *testing.T) {
	store, tmpDir := setupTestStore(t)
	defer cleanupTestStore(tmpDir)

	ctx := context.Background()

	// Test creating an asset without slug
	input := CreateFileObjectInput{
		Kind: FileObjectKindAsset,
		ID:   "0195f3f0-a8b7-7c8d-b1b9-d45c8f23a1ee",
		Slug: "", // Assets don't use slug
		Ext:  "png",
		Data: []byte("fake-image-data"),
	}

	loc, err := store.CreateFileObject(ctx, input)
	if err != nil {
		t.Fatalf("CreateFileObject failed: %v", err)
	}

	// Verify file name doesn't contain slug
	expectedFileName := "0195f3f0-a8b7-7c8d-b1b9-d45c8f23a1ee.png"
	if loc.FileName != expectedFileName {
		t.Errorf("expected filename %s, got %s", expectedFileName, loc.FileName)
	}
}

func TestFindFileObject(t *testing.T) {
	store, tmpDir := setupTestStore(t)
	defer cleanupTestStore(tmpDir)

	ctx := context.Background()

	// Create a test object
	input := CreateFileObjectInput{
		Kind: FileObjectKindCard,
		ID:   "0195f3f0-a8b7-7c8d-b1b9-d45c8f23a1ee",
		Slug: "test-card",
		Ext:  "html",
		Data: []byte("<h1>Test</h1>"),
	}
	_, err := store.CreateFileObject(ctx, input)
	if err != nil {
		t.Fatalf("CreateFileObject failed: %v", err)
	}

	// Test finding by kind + id
	lookup := FileObjectLookup{
		Ref: FileObjectRef{
			Kind: FileObjectKindCard,
			ID:   "0195f3f0-a8b7-7c8d-b1b9-d45c8f23a1ee",
		},
		Name: nil,
	}

	loc, err := store.FindFileObject(ctx, lookup, FileObjectScopeActive)
	if err != nil {
		t.Fatalf("FindFileObject failed: %v", err)
	}
	if loc == nil {
		t.Fatal("expected to find object, got nil")
	}
	if loc.Slug != "test-card" {
		t.Errorf("expected slug %s, got %s", "test-card", loc.Slug)
	}

	// Test finding by kind + id + slug + ext
	lookupWithHint := FileObjectLookup{
		Ref: FileObjectRef{
			Kind: FileObjectKindCard,
			ID:   "0195f3f0-a8b7-7c8d-b1b9-d45c8f23a1ee",
		},
		Name: &FileObjectNameHint{
			Slug: "test-card",
			Ext:  "html",
		},
	}

	loc2, err := store.FindFileObject(ctx, lookupWithHint, FileObjectScopeActive)
	if err != nil {
		t.Fatalf("FindFileObject with hint failed: %v", err)
	}
	if loc2 == nil {
		t.Fatal("expected to find object with hint, got nil")
	}
	if loc2.Path != loc.Path {
		t.Errorf("expected same path, got different paths")
	}

	// Test finding non-existent object
	lookupMissing := FileObjectLookup{
		Ref: FileObjectRef{
			Kind: FileObjectKindCard,
			ID:   "0195f3f0-0000-7000-8000-000000000000",
		},
		Name: nil,
	}

	loc3, err := store.FindFileObject(ctx, lookupMissing, FileObjectScopeActive)
	if err != nil {
		t.Fatalf("FindFileObject for missing object failed: %v", err)
	}
	if loc3 != nil {
		t.Error("expected nil for missing object, got location")
	}
}

func TestReadFileObject(t *testing.T) {
	store, tmpDir := setupTestStore(t)
	defer cleanupTestStore(tmpDir)

	ctx := context.Background()

	// Create a test object
	testContent := []byte("<h1>Test Content</h1>")
	input := CreateFileObjectInput{
		Kind: FileObjectKindCard,
		ID:   "0195f3f0-a8b7-7c8d-b1b9-d45c8f23a1ee",
		Slug: "test-card",
		Ext:  "html",
		Data: testContent,
	}
	_, err := store.CreateFileObject(ctx, input)
	if err != nil {
		t.Fatalf("CreateFileObject failed: %v", err)
	}

	// Read the object
	lookup := FileObjectLookup{
		Ref: FileObjectRef{
			Kind: FileObjectKindCard,
			ID:   "0195f3f0-a8b7-7c8d-b1b9-d45c8f23a1ee",
		},
		Name: nil,
	}

	content, loc, err := store.ReadFileObject(ctx, lookup, FileObjectScopeActive)
	if err != nil {
		t.Fatalf("ReadFileObject failed: %v", err)
	}
	if string(content) != string(testContent) {
		t.Errorf("expected content %s, got %s", string(testContent), string(content))
	}
	if loc.Slug != "test-card" {
		t.Errorf("expected slug %s, got %s", "test-card", loc.Slug)
	}
}

func TestOverwriteFileObject(t *testing.T) {
	store, tmpDir := setupTestStore(t)
	defer cleanupTestStore(tmpDir)

	ctx := context.Background()

	// Create initial object
	input := CreateFileObjectInput{
		Kind: FileObjectKindCard,
		ID:   "0195f3f0-a8b7-7c8d-b1b9-d45c8f23a1ee",
		Slug: "original-slug",
		Ext:  "html",
		Data: []byte("<h1>Original</h1>"),
	}
	_, err := store.CreateFileObject(ctx, input)
	if err != nil {
		t.Fatalf("CreateFileObject failed: %v", err)
	}

	// Overwrite with new content and slug
	overwriteInput := OverwriteFileObjectInput{
		Kind: FileObjectKindCard,
		ID:   "0195f3f0-a8b7-7c8d-b1b9-d45c8f23a1ee",
		Slug: "updated-slug",
		Ext:  "html",
		Data: []byte("<h1>Updated</h1>"),
	}

	loc, err := store.OverwriteFileObject(ctx, overwriteInput)
	if err != nil {
		t.Fatalf("OverwriteFileObject failed: %v", err)
	}

	// Verify slug changed
	if loc.Slug != "updated-slug" {
		t.Errorf("expected slug %s, got %s", "updated-slug", loc.Slug)
	}

	// Verify content changed
	content, err := os.ReadFile(loc.Path)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}
	if string(content) != "<h1>Updated</h1>" {
		t.Errorf("expected content %s, got %s", "<h1>Updated</h1>", string(content))
	}

	// Verify old file doesn't exist
	oldPath := filepath.Join(filepath.Dir(loc.Path), "0195f3f0-a8b7-7c8d-b1b9-d45c8f23a1ee--original-slug.html")
	if _, err := os.Stat(oldPath); !os.IsNotExist(err) {
		t.Error("old file should not exist after slug change")
	}
}

func TestMoveToTrashAndRestore(t *testing.T) {
	store, tmpDir := setupTestStore(t)
	defer cleanupTestStore(tmpDir)

	ctx := context.Background()

	// Create object
	input := CreateFileObjectInput{
		Kind: FileObjectKindCard,
		ID:   "0195f3f0-a8b7-7c8d-b1b9-d45c8f23a1ee",
		Slug: "test-card",
		Ext:  "html",
		Data: []byte("<h1>Test</h1>"),
	}
	createLoc, err := store.CreateFileObject(ctx, input)
	if err != nil {
		t.Fatalf("CreateFileObject failed: %v", err)
	}

	// Move to trash
	ref := FileObjectRef{
		Kind: FileObjectKindCard,
		ID:   "0195f3f0-a8b7-7c8d-b1b9-d45c8f23a1ee",
	}

	trashLoc, err := store.MoveFileObjectToTrash(ctx, ref)
	if err != nil {
		t.Fatalf("MoveFileObjectToTrash failed: %v", err)
	}

	// Verify scope changed
	if trashLoc.Scope != FileObjectScopeTrash {
		t.Errorf("expected scope %s, got %s", FileObjectScopeTrash, trashLoc.Scope)
	}

	// Verify original file doesn't exist
	if _, err := os.Stat(createLoc.Path); !os.IsNotExist(err) {
		t.Error("original file should not exist after moving to trash")
	}

	// Verify trash file exists
	if _, err := os.Stat(trashLoc.Path); os.IsNotExist(err) {
		t.Error("trash file should exist")
	}

	// Restore from trash
	restoreLoc, err := store.RestoreFileObjectFromTrash(ctx, ref)
	if err != nil {
		t.Fatalf("RestoreFileObjectFromTrash failed: %v", err)
	}

	// Verify scope changed back
	if restoreLoc.Scope != FileObjectScopeActive {
		t.Errorf("expected scope %s, got %s", FileObjectScopeActive, restoreLoc.Scope)
	}

	// Verify trash file doesn't exist
	if _, err := os.Stat(trashLoc.Path); !os.IsNotExist(err) {
		t.Error("trash file should not exist after restore")
	}

	// Verify active file exists
	if _, err := os.Stat(restoreLoc.Path); os.IsNotExist(err) {
		t.Error("active file should exist after restore")
	}
}

func TestPermanentlyDelete(t *testing.T) {
	store, tmpDir := setupTestStore(t)
	defer cleanupTestStore(tmpDir)

	ctx := context.Background()

	// Create and move to trash
	input := CreateFileObjectInput{
		Kind: FileObjectKindCard,
		ID:   "0195f3f0-a8b7-7c8d-b1b9-d45c8f23a1ee",
		Slug: "test-card",
		Ext:  "html",
		Data: []byte("<h1>Test</h1>"),
	}
	_, err := store.CreateFileObject(ctx, input)
	if err != nil {
		t.Fatalf("CreateFileObject failed: %v", err)
	}

	ref := FileObjectRef{
		Kind: FileObjectKindCard,
		ID:   "0195f3f0-a8b7-7c8d-b1b9-d45c8f23a1ee",
	}

	trashLoc, err := store.MoveFileObjectToTrash(ctx, ref)
	if err != nil {
		t.Fatalf("MoveFileObjectToTrash failed: %v", err)
	}

	// Permanently delete
	err = store.PermanentlyDeleteFileObject(ctx, ref)
	if err != nil {
		t.Fatalf("PermanentlyDeleteFileObject failed: %v", err)
	}

	// Verify file doesn't exist
	if _, err := os.Stat(trashLoc.Path); !os.IsNotExist(err) {
		t.Error("file should not exist after permanent deletion")
	}
}

func TestRenameSlug(t *testing.T) {
	store, tmpDir := setupTestStore(t)
	defer cleanupTestStore(tmpDir)

	ctx := context.Background()

	// Create object
	input := CreateFileObjectInput{
		Kind: FileObjectKindCard,
		ID:   "0195f3f0-a8b7-7c8d-b1b9-d45c8f23a1ee",
		Slug: "original-slug",
		Ext:  "html",
		Data: []byte("<h1>Test</h1>"),
	}
	originalLoc, err := store.CreateFileObject(ctx, input)
	if err != nil {
		t.Fatalf("CreateFileObject failed: %v", err)
	}

	// Rename slug
	ref := FileObjectRef{
		Kind: FileObjectKindCard,
		ID:   "0195f3f0-a8b7-7c8d-b1b9-d45c8f23a1ee",
	}

	newLoc, err := store.RenameFileObjectSlug(ctx, ref, "new-slug")
	if err != nil {
		t.Fatalf("RenameFileObjectSlug failed: %v", err)
	}

	// Verify slug changed
	if newLoc.Slug != "new-slug" {
		t.Errorf("expected slug %s, got %s", "new-slug", newLoc.Slug)
	}

	// Verify old file doesn't exist
	if _, err := os.Stat(originalLoc.Path); !os.IsNotExist(err) {
		t.Error("old file should not exist after rename")
	}

	// Verify new file exists
	if _, err := os.Stat(newLoc.Path); os.IsNotExist(err) {
		t.Error("new file should exist after rename")
	}
}

func TestScanFileObjects(t *testing.T) {
	store, tmpDir := setupTestStore(t)
	defer cleanupTestStore(tmpDir)

	ctx := context.Background()

	// Create multiple objects
	for i := 0; i < 3; i++ {
		input := CreateFileObjectInput{
			Kind: FileObjectKindCard,
			ID:   FileObjectID(fmt.Sprintf("0195f3f0-a8b7-7c8d-b1b9-d45c8f23a1e%d", i)),
			Slug: fmt.Sprintf("card-%d", i),
			Ext:  "html",
			Data: []byte(fmt.Sprintf("<h1>Card %d</h1>", i)),
		}
		_, err := store.CreateFileObject(ctx, input)
		if err != nil {
			t.Fatalf("CreateFileObject failed: %v", err)
		}
	}

	// Scan all cards
	kind := FileObjectKindCard
	options := ScanFileObjectsOptions{
		Kind:  &kind,
		Scope: FileObjectScopeActive,
	}

	locationCh, errorCh := store.ScanFileObjects(ctx, options)

	// Collect results
	var locations []FileObjectLocation
	for loc := range locationCh {
		locations = append(locations, loc)
	}

	// Check for errors
	for err := range errorCh {
		t.Fatalf("ScanFileObjects error: %v", err)
	}

	// Verify count
	if len(locations) != 3 {
		t.Errorf("expected 3 locations, got %d", len(locations))
	}
}

func TestBucketPath(t *testing.T) {
	tests := []struct {
		id       FileObjectID
		expected string
	}{
		{"0195f3f0-a8b7-7c8d-b1b9-d45c8f23a1ee", "019/5/f/3"},
		{"01234567-89ab-7cde-8f01-234567890abc", "012/3/4/5"},
	}

	for _, tt := range tests {
		result := BuildBucketPath(tt.id)
		if result != tt.expected {
			t.Errorf("BuildBucketPath(%s) = %s, expected %s", tt.id, result, tt.expected)
		}
	}
}

func TestSlugNormalization(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Hello World", "hello-world"},
		{"Test  Multiple   Spaces", "test-multiple-spaces"},
		{"CamelCase", "camelcase"},
		{"with-hyphens", "with-hyphens"},
		{"with_underscores", "with_underscores"},
		{"with.dots", "with.dots"},
		{"", "untitled"},
		{"   ", "untitled"},
		{"---test---", "test"},
		{"test@#$%special", "testspecial"},
	}

	for _, tt := range tests {
		result := NormalizeSlug(tt.input)
		if result != tt.expected {
			t.Errorf("NormalizeSlug(%q) = %q, expected %q", tt.input, result, tt.expected)
		}
	}
}

func TestUUIDValidation(t *testing.T) {
	tests := []struct {
		id    string
		valid bool
	}{
		{"0195f3f0-a8b7-7c8d-b1b9-d45c8f23a1ee", true},
		{"0195f3f0-a8b7-7000-8000-d45c8f23a1ee", true},
		{"0195f3f0-a8b7-6c8d-b1b9-d45c8f23a1ee", false}, // version 6, not 7
		{"0195f3f0-a8b7-7c8d-c1b9-d45c8f23a1ee", false}, // invalid variant
		{"not-a-uuid", false},
		{"", false},
	}

	for _, tt := range tests {
		result := ValidateUUIDv7(tt.id)
		if result != tt.valid {
			t.Errorf("ValidateUUIDv7(%q) = %v, expected %v", tt.id, result, tt.valid)
		}
	}
}
