// ABOUTME: Type definitions for the favorites persistence system.
// ABOUTME: Contains Favorite struct with all ASCII art configuration fields.

package favorites

import "time"

// Favorite represents a saved ASCII art configuration
type Favorite struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`

	// Text content
	Text      string `json:"text"`
	FontName  string `json:"font_name"`
	Alignment int    `json:"alignment"`

	// Spacing
	CharSpacing int `json:"char_spacing"`
	WordSpacing int `json:"word_spacing"`
	LineSpacing int `json:"line_spacing"`

	// Color
	TextColor         int  `json:"text_color"`
	GradientEnabled   bool `json:"gradient_enabled"`
	GradientColor     int  `json:"gradient_color"`
	GradientDirection int  `json:"gradient_direction"`

	// Scale
	Scale int `json:"scale"`

	// Shadow
	ShadowEnabled bool `json:"shadow_enabled"`
	ShadowHOffset int  `json:"shadow_h_offset"`
	ShadowVOffset int  `json:"shadow_v_offset"`
	ShadowStyle   int  `json:"shadow_style"`
}

// FavoritesStore holds all saved favorites
type FavoritesStore struct {
	Favorites []Favorite `json:"favorites"`
}
