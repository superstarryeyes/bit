// ABOUTME: Tests for favorites type definitions and JSON serialization.
// ABOUTME: Validates Favorite struct fields and round-trip JSON encoding/decoding.

package favorites

import (
	"encoding/json"
	"testing"
	"time"
)

func TestFavorite_JSONRoundTrip(t *testing.T) {
	original := Favorite{
		ID:        "test-123",
		Name:      "My Cool Art",
		CreatedAt: time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),

		Text:      "Hello\nWorld",
		FontName:  "BlockFont",
		Alignment: 1,

		CharSpacing: 2,
		WordSpacing: 4,
		LineSpacing: 1,

		TextColor:         3,
		GradientEnabled:   true,
		GradientColor:     5,
		GradientDirection: 2,

		Scale: 1,

		ShadowEnabled: true,
		ShadowHOffset: 2,
		ShadowVOffset: -1,
		ShadowStyle:   1,
	}

	// Marshal to JSON
	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("failed to marshal Favorite: %v", err)
	}

	// Unmarshal back
	var decoded Favorite
	err = json.Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("failed to unmarshal Favorite: %v", err)
	}

	// Verify all fields
	if decoded.ID != original.ID {
		t.Errorf("ID mismatch: got %q, want %q", decoded.ID, original.ID)
	}
	if decoded.Name != original.Name {
		t.Errorf("Name mismatch: got %q, want %q", decoded.Name, original.Name)
	}
	if !decoded.CreatedAt.Equal(original.CreatedAt) {
		t.Errorf("CreatedAt mismatch: got %v, want %v", decoded.CreatedAt, original.CreatedAt)
	}
	if decoded.Text != original.Text {
		t.Errorf("Text mismatch: got %q, want %q", decoded.Text, original.Text)
	}
	if decoded.FontName != original.FontName {
		t.Errorf("FontName mismatch: got %q, want %q", decoded.FontName, original.FontName)
	}
	if decoded.Alignment != original.Alignment {
		t.Errorf("Alignment mismatch: got %d, want %d", decoded.Alignment, original.Alignment)
	}
	if decoded.CharSpacing != original.CharSpacing {
		t.Errorf("CharSpacing mismatch: got %d, want %d", decoded.CharSpacing, original.CharSpacing)
	}
	if decoded.WordSpacing != original.WordSpacing {
		t.Errorf("WordSpacing mismatch: got %d, want %d", decoded.WordSpacing, original.WordSpacing)
	}
	if decoded.LineSpacing != original.LineSpacing {
		t.Errorf("LineSpacing mismatch: got %d, want %d", decoded.LineSpacing, original.LineSpacing)
	}
	if decoded.TextColor != original.TextColor {
		t.Errorf("TextColor mismatch: got %d, want %d", decoded.TextColor, original.TextColor)
	}
	if decoded.GradientEnabled != original.GradientEnabled {
		t.Errorf("GradientEnabled mismatch: got %v, want %v", decoded.GradientEnabled, original.GradientEnabled)
	}
	if decoded.GradientColor != original.GradientColor {
		t.Errorf("GradientColor mismatch: got %d, want %d", decoded.GradientColor, original.GradientColor)
	}
	if decoded.GradientDirection != original.GradientDirection {
		t.Errorf("GradientDirection mismatch: got %d, want %d", decoded.GradientDirection, original.GradientDirection)
	}
	if decoded.Scale != original.Scale {
		t.Errorf("Scale mismatch: got %d, want %d", decoded.Scale, original.Scale)
	}
	if decoded.ShadowEnabled != original.ShadowEnabled {
		t.Errorf("ShadowEnabled mismatch: got %v, want %v", decoded.ShadowEnabled, original.ShadowEnabled)
	}
	if decoded.ShadowHOffset != original.ShadowHOffset {
		t.Errorf("ShadowHOffset mismatch: got %d, want %d", decoded.ShadowHOffset, original.ShadowHOffset)
	}
	if decoded.ShadowVOffset != original.ShadowVOffset {
		t.Errorf("ShadowVOffset mismatch: got %d, want %d", decoded.ShadowVOffset, original.ShadowVOffset)
	}
	if decoded.ShadowStyle != original.ShadowStyle {
		t.Errorf("ShadowStyle mismatch: got %d, want %d", decoded.ShadowStyle, original.ShadowStyle)
	}
}

func TestFavoritesStore_JSONRoundTrip(t *testing.T) {
	original := FavoritesStore{
		Favorites: []Favorite{
			{ID: "fav-1", Name: "First", Text: "Hello"},
			{ID: "fav-2", Name: "Second", Text: "World"},
		},
	}

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("failed to marshal FavoritesStore: %v", err)
	}

	var decoded FavoritesStore
	err = json.Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("failed to unmarshal FavoritesStore: %v", err)
	}

	if len(decoded.Favorites) != len(original.Favorites) {
		t.Fatalf("Favorites count mismatch: got %d, want %d", len(decoded.Favorites), len(original.Favorites))
	}

	for i, fav := range decoded.Favorites {
		if fav.ID != original.Favorites[i].ID {
			t.Errorf("Favorite[%d].ID mismatch: got %q, want %q", i, fav.ID, original.Favorites[i].ID)
		}
		if fav.Name != original.Favorites[i].Name {
			t.Errorf("Favorite[%d].Name mismatch: got %q, want %q", i, fav.Name, original.Favorites[i].Name)
		}
	}
}

func TestFavoritesStore_EmptyList(t *testing.T) {
	store := FavoritesStore{Favorites: []Favorite{}}

	data, err := json.Marshal(store)
	if err != nil {
		t.Fatalf("failed to marshal empty FavoritesStore: %v", err)
	}

	var decoded FavoritesStore
	err = json.Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("failed to unmarshal empty FavoritesStore: %v", err)
	}

	if decoded.Favorites == nil {
		t.Error("Favorites should not be nil after unmarshal")
	}
	if len(decoded.Favorites) != 0 {
		t.Errorf("Favorites should be empty, got %d items", len(decoded.Favorites))
	}
}
