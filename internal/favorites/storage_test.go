// ABOUTME: Tests for favorites file storage operations.
// ABOUTME: Validates loading, saving, and config directory handling.

package favorites

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestGetConfigDir_CreatesDirectory(t *testing.T) {
	// Use temp dir as home
	tmpHome := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpHome)
	defer os.Setenv("HOME", originalHome)

	dir, err := GetConfigDir()
	if err != nil {
		t.Fatalf("GetConfigDir failed: %v", err)
	}

	expectedDir := filepath.Join(tmpHome, ".config", "bit")
	if dir != expectedDir {
		t.Errorf("unexpected dir: got %q, want %q", dir, expectedDir)
	}

	// Verify directory exists
	info, err := os.Stat(dir)
	if err != nil {
		t.Fatalf("config dir does not exist: %v", err)
	}
	if !info.IsDir() {
		t.Error("config path is not a directory")
	}
}

func TestLoad_MissingFile_ReturnsEmptyStore(t *testing.T) {
	tmpHome := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpHome)
	defer os.Setenv("HOME", originalHome)

	store, err := Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if store.Favorites == nil {
		t.Error("Favorites should not be nil")
	}
	if len(store.Favorites) != 0 {
		t.Errorf("expected empty favorites, got %d", len(store.Favorites))
	}
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	tmpHome := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpHome)
	defer os.Setenv("HOME", originalHome)

	original := &FavoritesStore{
		Favorites: []Favorite{
			{
				ID:          "test-1",
				Name:        "Test Favorite",
				CreatedAt:   time.Now().UTC().Truncate(time.Second),
				Text:        "Hello World",
				FontName:    "BlockFont",
				CharSpacing: 2,
			},
		},
	}

	// Save
	err := Save(original)
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Load
	loaded, err := Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if len(loaded.Favorites) != 1 {
		t.Fatalf("expected 1 favorite, got %d", len(loaded.Favorites))
	}

	fav := loaded.Favorites[0]
	if fav.ID != "test-1" {
		t.Errorf("ID mismatch: got %q", fav.ID)
	}
	if fav.Name != "Test Favorite" {
		t.Errorf("Name mismatch: got %q", fav.Name)
	}
	if fav.Text != "Hello World" {
		t.Errorf("Text mismatch: got %q", fav.Text)
	}
}

func TestLoad_InvalidJSON_ReturnsError(t *testing.T) {
	tmpHome := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpHome)
	defer os.Setenv("HOME", originalHome)

	// Create config dir and write invalid JSON
	configDir := filepath.Join(tmpHome, ".config", "bit")
	err := os.MkdirAll(configDir, 0755)
	if err != nil {
		t.Fatalf("failed to create config dir: %v", err)
	}

	filePath := filepath.Join(configDir, "favorites.json")
	err = os.WriteFile(filePath, []byte("not valid json{"), 0644)
	if err != nil {
		t.Fatalf("failed to write invalid json: %v", err)
	}

	_, err = Load()
	if err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
}

func TestGetFavoritesFilePath(t *testing.T) {
	tmpHome := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpHome)
	defer os.Setenv("HOME", originalHome)

	path, err := GetFavoritesFilePath()
	if err != nil {
		t.Fatalf("GetFavoritesFilePath failed: %v", err)
	}

	expected := filepath.Join(tmpHome, ".config", "bit", "favorites.json")
	if path != expected {
		t.Errorf("path mismatch: got %q, want %q", path, expected)
	}
}
