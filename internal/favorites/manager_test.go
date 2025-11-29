// ABOUTME: Tests for favorites manager business logic.
// ABOUTME: Validates add, remove, list, and get operations.

package favorites

import (
	"os"
	"testing"
)

func setupTestEnv(t *testing.T) func() {
	tmpHome := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpHome)
	return func() {
		os.Setenv("HOME", originalHome)
	}
}

func TestNewManager_LoadsExistingFavorites(t *testing.T) {
	cleanup := setupTestEnv(t)
	defer cleanup()

	// Pre-save some favorites
	store := &FavoritesStore{
		Favorites: []Favorite{
			{ID: "existing-1", Name: "Existing"},
		},
	}
	err := Save(store)
	if err != nil {
		t.Fatalf("failed to save initial store: %v", err)
	}

	// Create manager - should load existing
	mgr, err := NewManager()
	if err != nil {
		t.Fatalf("NewManager failed: %v", err)
	}

	favorites := mgr.List()
	if len(favorites) != 1 {
		t.Errorf("expected 1 favorite, got %d", len(favorites))
	}
}

func TestManager_Add(t *testing.T) {
	cleanup := setupTestEnv(t)
	defer cleanup()

	mgr, err := NewManager()
	if err != nil {
		t.Fatalf("NewManager failed: %v", err)
	}

	fav := Favorite{
		Name:     "Test Art",
		Text:     "Hello",
		FontName: "BlockFont",
	}

	id, err := mgr.Add(fav)
	if err != nil {
		t.Fatalf("Add failed: %v", err)
	}

	if id == "" {
		t.Error("Add should return non-empty ID")
	}

	// Verify it's in the list
	favorites := mgr.List()
	if len(favorites) != 1 {
		t.Fatalf("expected 1 favorite, got %d", len(favorites))
	}

	if favorites[0].ID != id {
		t.Errorf("ID mismatch: got %q, want %q", favorites[0].ID, id)
	}
	if favorites[0].Name != "Test Art" {
		t.Errorf("Name mismatch: got %q", favorites[0].Name)
	}
}

func TestManager_Add_PersistsToDisk(t *testing.T) {
	cleanup := setupTestEnv(t)
	defer cleanup()

	mgr, err := NewManager()
	if err != nil {
		t.Fatalf("NewManager failed: %v", err)
	}

	_, err = mgr.Add(Favorite{Name: "Persisted", Text: "Test"})
	if err != nil {
		t.Fatalf("Add failed: %v", err)
	}

	// Create new manager - should load persisted favorite
	mgr2, err := NewManager()
	if err != nil {
		t.Fatalf("NewManager failed: %v", err)
	}

	favorites := mgr2.List()
	if len(favorites) != 1 {
		t.Fatalf("expected 1 persisted favorite, got %d", len(favorites))
	}
	if favorites[0].Name != "Persisted" {
		t.Errorf("Name mismatch: got %q", favorites[0].Name)
	}
}

func TestManager_Get(t *testing.T) {
	cleanup := setupTestEnv(t)
	defer cleanup()

	mgr, err := NewManager()
	if err != nil {
		t.Fatalf("NewManager failed: %v", err)
	}

	id, err := mgr.Add(Favorite{Name: "GetTest", Text: "Hello"})
	if err != nil {
		t.Fatalf("Add failed: %v", err)
	}

	fav, err := mgr.Get(id)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if fav.Name != "GetTest" {
		t.Errorf("Name mismatch: got %q", fav.Name)
	}
}

func TestManager_Get_NotFound(t *testing.T) {
	cleanup := setupTestEnv(t)
	defer cleanup()

	mgr, err := NewManager()
	if err != nil {
		t.Fatalf("NewManager failed: %v", err)
	}

	_, err = mgr.Get("nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent ID")
	}
}

func TestManager_Remove(t *testing.T) {
	cleanup := setupTestEnv(t)
	defer cleanup()

	mgr, err := NewManager()
	if err != nil {
		t.Fatalf("NewManager failed: %v", err)
	}

	id, err := mgr.Add(Favorite{Name: "ToRemove"})
	if err != nil {
		t.Fatalf("Add failed: %v", err)
	}

	err = mgr.Remove(id)
	if err != nil {
		t.Fatalf("Remove failed: %v", err)
	}

	favorites := mgr.List()
	if len(favorites) != 0 {
		t.Errorf("expected 0 favorites after remove, got %d", len(favorites))
	}
}

func TestManager_Remove_NotFound(t *testing.T) {
	cleanup := setupTestEnv(t)
	defer cleanup()

	mgr, err := NewManager()
	if err != nil {
		t.Fatalf("NewManager failed: %v", err)
	}

	err = mgr.Remove("nonexistent")
	if err == nil {
		t.Error("expected error for removing nonexistent ID")
	}
}

func TestManager_Remove_PersistsToDisk(t *testing.T) {
	cleanup := setupTestEnv(t)
	defer cleanup()

	mgr, err := NewManager()
	if err != nil {
		t.Fatalf("NewManager failed: %v", err)
	}

	id, _ := mgr.Add(Favorite{Name: "ToRemove"})
	mgr.Remove(id)

	// Create new manager - should reflect removal
	mgr2, err := NewManager()
	if err != nil {
		t.Fatalf("NewManager failed: %v", err)
	}

	if len(mgr2.List()) != 0 {
		t.Error("removal was not persisted")
	}
}

func TestManager_List_ReturnsInOrder(t *testing.T) {
	cleanup := setupTestEnv(t)
	defer cleanup()

	mgr, err := NewManager()
	if err != nil {
		t.Fatalf("NewManager failed: %v", err)
	}

	mgr.Add(Favorite{Name: "First"})
	mgr.Add(Favorite{Name: "Second"})
	mgr.Add(Favorite{Name: "Third"})

	favorites := mgr.List()
	if len(favorites) != 3 {
		t.Fatalf("expected 3 favorites, got %d", len(favorites))
	}

	// Should be in insertion order
	if favorites[0].Name != "First" {
		t.Errorf("expected First, got %q", favorites[0].Name)
	}
	if favorites[1].Name != "Second" {
		t.Errorf("expected Second, got %q", favorites[1].Name)
	}
	if favorites[2].Name != "Third" {
		t.Errorf("expected Third, got %q", favorites[2].Name)
	}
}
