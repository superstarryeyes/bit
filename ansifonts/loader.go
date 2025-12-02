// Package ansifonts provides a library for rendering text using ANSI art fonts.
//
// The loader package handles font discovery and loading from the fonts directory.
package ansifonts

import (
	"embed"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
)

//go:embed fonts/*.bit
var EmbeddedFonts embed.FS

// customFontsRegistry holds custom fonts loaded from the filesystem
var customFontsRegistry = make(map[string]FontData)

// validateFontData ensures the JSON has required fields
func validateFontData(fd *FontData) error {
	if fd.Name == "" {
		return fmt.Errorf("font data missing required 'name' field")
	}
	if fd.Characters == nil || len(fd.Characters) == 0 {
		return fmt.Errorf("font data missing required 'characters' field")
	}
	return nil
}

// RegisterFontFile loads a single .bit font file and registers it
func RegisterFontFile(path string) (string, error) {
	// Check file extension (case-insensitive)
	ext := strings.ToLower(filepath.Ext(path))
	if ext != ".bit" {
		return "", fmt.Errorf("file %s does not have .bit extension", path)
	}

	// Read file
	fontBytes, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to read file %s: %w", path, err)
	}

	// Unmarshal and validate
	var fontData FontData
	if err := json.Unmarshal(fontBytes, &fontData); err != nil {
		return "", fmt.Errorf("failed to parse JSON in %s: %w", path, err)
	}

	if err := validateFontData(&fontData); err != nil {
		return "", fmt.Errorf("invalid font data in %s: %w", path, err)
	}

	// Store in registry using lowercase name as key
	key := strings.ToLower(fontData.Name)
	customFontsRegistry[key] = fontData

	return fontData.Name, nil
}

// RegisterFontDirectory loads all .bit font files from a directory
func RegisterFontDirectory(dirPath string) ([]string, error) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory %s: %w", dirPath, err)
	}

	var loadedNames []string
	var errors []string

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		fileName := entry.Name()
		if !strings.HasSuffix(strings.ToLower(fileName), ".bit") {
			continue
		}

		fullPath := filepath.Join(dirPath, fileName)
		fontName, err := RegisterFontFile(fullPath)
		if err != nil {
			errors = append(errors, fmt.Sprintf("failed to load %s: %v", fileName, err))
			continue
		}

		loadedNames = append(loadedNames, fontName)
	}

	// If no fonts could be loaded, return an error
	if len(loadedNames) == 0 {
		if len(errors) > 0 {
			return nil, fmt.Errorf("no fonts could be loaded from directory %s. Errors: %s", dirPath, strings.Join(errors, "; "))
		}
		return nil, fmt.Errorf("no .bit font files found in directory %s", dirPath)
	}

	// Log errors for partially failed loads (but still return success)
	if len(errors) > 0 {
		fmt.Fprintf(os.Stderr, "Warning: Some fonts failed to load: %s\n", strings.Join(errors, "; "))
	}

	return loadedNames, nil
}

// RegisterCustomPath is the smart entry point that handles both files and directories
func RegisterCustomPath(path string) ([]string, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("path %s does not exist: %w", path, err)
	}

	if info.IsDir() {
		return RegisterFontDirectory(path)
	}

	// It's a file
	fontName, err := RegisterFontFile(path)
	if err != nil {
		return nil, err
	}

	return []string{fontName}, nil
}

// LoadFont loads a font by name, checking custom fonts first, then embedded fonts
func LoadFont(name string) (*Font, error) {
	// Check custom fonts registry first (allows overriding embedded fonts)
	key := strings.ToLower(name)
	if fontData, exists := customFontsRegistry[key]; exists {
		return &Font{
			Name:     fontData.Name,
			FontData: fontData,
		}, nil
	}

	// Fall back to embedded fonts
	fontPath := path.Join("fonts", name+".bit")
	fontBytes, err := EmbeddedFonts.ReadFile(fontPath)
	if err != nil {
		return nil, fmt.Errorf("font '%s' not found in custom or embedded fonts", name)
	}

	var fontData FontData
	err = json.Unmarshal(fontBytes, &fontData)
	if err != nil {
		return nil, err
	}

	return &Font{
		Name:     fontData.Name,
		FontData: fontData,
	}, nil
}

// ListFonts returns a list of available font names from both custom and embedded fonts
func ListFonts() ([]string, error) {
	// Get embedded fonts
	entries, err := EmbeddedFonts.ReadDir("fonts")
	if err != nil {
		return nil, err
	}

	fontSet := make(map[string]bool)
	
	// Add embedded fonts
	for _, entry := range entries {
		if !entry.IsDir() && path.Ext(entry.Name()) == ".bit" {
			fontName := strings.TrimSuffix(entry.Name(), ".bit")
			fontSet[fontName] = true
		}
	}

	// Add custom fonts (these may override embedded fonts)
	for _, fontData := range customFontsRegistry {
		fontSet[fontData.Name] = true
	}

	// Convert to sorted slice
	fonts := make([]string, 0, len(fontSet))
	for fontName := range fontSet {
		fonts = append(fonts, fontName)
	}
	sort.Strings(fonts)

	return fonts, nil
}
