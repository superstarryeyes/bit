package ui

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/superstarryeyes/bit/ansifonts"
)

// loadFontList loads only the font metadata without loading the actual font data
// This provides a list of available fonts without consuming memory for all font data
func loadFontList() ([]FontInfo, error) {
	var fonts []FontInfo

	// Read from embedded filesystem in ansifonts package
	entries, err := ansifonts.EmbeddedFonts.ReadDir("fonts")
	if err != nil {
		return nil, fmt.Errorf("failed to read embedded fonts directory: %w", err)
	}

	for _, entry := range entries {
		if filepath.Ext(entry.Name()) == ".bit" {
			fontPath := filepath.Join("fonts", entry.Name())
			if strings.Contains(fontPath, "\\") {
				fontPath = strings.ReplaceAll(fontPath, "\\", "/")
			}
			// For lazy loading, we only load the font name from the file
			// This is much more efficient than loading the entire font data
			fontDataBytes, err := ansifonts.EmbeddedFonts.ReadFile(fontPath)
			if err != nil {
				// Log or track skipped files for debugging
				fmt.Printf("Warning: skipping font file %s: %v\n", entry.Name(), err)
				continue
			}

			var fontMetadata struct {
				Name string `json:"name"`
			}
			err = json.Unmarshal(fontDataBytes, &fontMetadata)
			if err != nil {
				// Log or track invalid JSON files for debugging
				fmt.Printf("Warning: skipping invalid font JSON %s: %v\n", entry.Name(), err)
				continue
			}

			fonts = append(fonts, FontInfo{
				Name:   fontMetadata.Name,
				Path:   fontPath,
				Loaded: false, // Font data not loaded yet
			})
		}
	}

	if len(fonts) == 0 {
		return nil, fmt.Errorf("no valid fonts found in embedded directory - please ensure font files are properly embedded")
	}

	// Sort fonts case-insensitively by name
	sort.Slice(fonts, func(i, j int) bool {
		return strings.ToLower(fonts[i].Name) < strings.ToLower(fonts[j].Name)
	})

	return fonts, nil
}

// loadFontData loads the full font data for a specific font
// This is called when a font is actually needed for rendering
func loadFontData(font *FontInfo) error {
	if font.Loaded {
		// Font already loaded, nothing to do
		return nil
	}

	// Load the full font data
	fontDataBytes, err := ansifonts.EmbeddedFonts.ReadFile(font.Path)
	if err != nil {
		return fmt.Errorf("failed to read font file %s: %w", font.Path, err)
	}

	var fontData FontData
	err = json.Unmarshal(fontDataBytes, &fontData)
	if err != nil {
		return fmt.Errorf("failed to parse font JSON %s: %w", font.Path, err)
	}

	// Update the font with loaded data
	font.FontData = &fontData
	font.Loaded = true

	return nil
}
