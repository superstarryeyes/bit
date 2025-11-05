package ui

import (
	"github.com/superstarryeyes/bit/ansifonts"
)

// Centralized color palette for consistent color usage across the application
var ColorPalette = map[string]string{
	"White":       "#FFFFFF",
	"Red":         "#FF5555",
	"Green":       "#50FA7B",
	"Blue":        "#8BE9FD",
	"Yellow":      "#F1FA8C",
	"Magenta":     "#FF79C6",
	"Cyan":        "#8BE9FD",
	"Gray":        "#6272A4",
	"PureRed":     "#FF0000",
	"TextInput":   "#FF6B6B",
	"Export":      "#4A90E2",
	"FaintGray":   "#626264",
	"VeryFaint":   "#626262",
	"TitleFG":     "#FAFAFA",
	"TitleBG":     "#7D56F4",
	"PanelBorder": "#874BFD",
	"Selected":    "#F25D94",
	"TextDisplay": "#04B575",
	"Shadow":      "#A0A0A0",
	"FontPanel":   "#4ECDC4",
	"CharSpacing": "#45B7D1",
	"WordSpacing": "#96CEB4",
	"LineSpacing": "#9B59B6",
	"ColorPanel":  "#FECA57",
	"ScalePanel":  "#FF9FF3",
	"Black":       "#000000",
	"Background":  "#B19CD9",
	"Animation":   "#FF6B9D",
}

// ANSIColorMap provides mappings from ANSI color codes to hex values
// This references the canonical color map from the ansifonts package
var ANSIColorMap = ansifonts.ANSIColorMap

// GetANSIColorOptions returns a comprehensive list of color options
// specifically for text rendering, sourced from standard ANSI colors
func GetANSIColorOptions() []ColorOption {
	// Define the ANSI color mappings directly from the library
	ansiColors := []struct {
		Code string
		Name string
	}{
		{"31", "Red"},
		{"32", "Green"},
		{"33", "Yellow"},
		{"34", "Blue"},
		{"35", "Magenta"},
		{"36", "Cyan"},
		{"37", "White"},
		{"91", "Bright Red"},
		{"92", "Bright Green"},
		{"93", "Bright Yellow"},
		{"94", "Bright Blue"},
		{"95", "Bright Magenta"},
		{"96", "Bright Cyan"},
		{"97", "Bright White"},
	}

	// Create options slice with initial capacity
	options := make([]ColorOption, 0, len(ansiColors))

	// Add all ANSI colors from the library
	for _, ansiColor := range ansiColors {
		if hex, exists := ANSIColorMap[ansiColor.Code]; exists {
			options = append(options, ColorOption{
				Name: ansiColor.Name,
				Hex:  hex,
			})
		}
	}

	return options
}
