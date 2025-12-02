package ui

import (
	"fmt"

	"github.com/superstarryeyes/bit/ansifonts"
)

// loadFontList loads only the font metadata without loading the actual font data
// This provides a list of available fonts without consuming memory for all font data
func loadFontList() ([]FontInfo, error) {
	// Use the unified ListFonts function which includes both custom and embedded fonts
	fontNames, err := ansifonts.ListFonts()
	if err != nil {
		return nil, fmt.Errorf("failed to list fonts: %w", err)
	}

	if len(fontNames) == 0 {
		return nil, fmt.Errorf("no fonts available - please ensure font files are properly embedded or loaded")
	}

	var fonts []FontInfo
	for _, fontName := range fontNames {
		fonts = append(fonts, FontInfo{
			Name:   fontName,
			Path:   "", // Path not relevant when using unified loader
			Loaded: false, // Font data not loaded yet
		})
	}

	return fonts, nil
}

// loadFontData loads the full font data for a specific font
// This is called when a font is actually needed for rendering
func loadFontData(font *FontInfo) error {
	if font.Loaded {
		// Font already loaded, nothing to do
		return nil
	}

	// Use the unified LoadFont function which handles both custom and embedded fonts
	loadedFont, err := ansifonts.LoadFont(font.Name)
	if err != nil {
		return fmt.Errorf("failed to load font %s: %w", font.Name, err)
	}

	// Update the font with loaded data - convert ansifonts.FontData to ui.FontData
	font.FontData = &FontData{
		Name:       loadedFont.FontData.Name,
		Author:     loadedFont.FontData.Author,
		License:    loadedFont.FontData.License,
		Characters: loadedFont.FontData.Characters,
	}
	font.Loaded = true

	return nil
}
