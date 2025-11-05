// Package ansifonts provides a standalone library for rendering text using ANSI art fonts.
//
// The library dynamically discovers ANSI fonts from a fonts directory and provides
// simple APIs for loading fonts and rendering text with various formatting options.
//
// Basic usage:
//
//	font, _ := ansifonts.LoadFont("dogica")
//	lines := ansifonts.RenderText("Hello", font)
//	for _, line := range lines {
//	    fmt.Println(line)
//	}
//
// Advanced usage with options:
//
//	options := ansifonts.RenderOptions{
//	    CharSpacing: 2,
//	    TextColor: "#FF0000",
//	    UseGradient: true,
//	    GradientColor: "#0000FF",
//	}
//	lines := ansifonts.RenderTextWithOptions("Hello", font, options)
package ansifonts

import "fmt"

// FontData represents the overall structure of our .bit font file (JSON format)
type FontData struct {
	Name       string              `json:"name"`
	Author     string              `json:"author"`
	License    string              `json:"license"`
	Characters map[string][]string `json:"characters"`
}

// Font represents a loaded font with its metadata
type Font struct {
	Name     string
	FontData FontData
}

// TextAlignment represents text alignment options
type TextAlignment int

const (
	LeftAlign TextAlignment = iota
	CenterAlign
	RightAlign
)

// GradientDirection represents gradient direction options
type GradientDirection int

const (
	UpDown GradientDirection = iota
	DownUp
	LeftRight
	RightLeft
)

// ColorMode represents different color application modes
type ColorMode int

const (
	SingleColor ColorMode = iota
	Gradient
	Rainbow
)

// ShadowStyle represents shadow style options
type ShadowStyle int

const (
	LightShade ShadowStyle = iota
	MediumShade
	DarkShade
)

// RenderOptions contains all the options for rendering text
type RenderOptions struct {
	// Spacing options
	CharSpacing int // Character spacing (0 to 10)
	WordSpacing int // Word spacing (0 to 20)
	LineSpacing int // Line spacing (0 to 10)

	// Text alignment
	Alignment TextAlignment

	// Text color options
	TextColor         string // Hex color code (e.g., "#FFFFFF")
	GradientColor     string // Hex color code for gradient end color
	GradientDirection GradientDirection
	UseGradient       bool

	// Rainbow effect options
	ColorMode      ColorMode // SingleColor, Gradient, or Rainbow
	RainbowColors  []string  // Custom rainbow colors (hex codes), defaults to standard rainbow if empty
	RainbowFrame   int       // Animation frame for rainbow cycling (default: 0)
	RainbowSpeed   int       // How many frames before color shifts (default: 5)

	// Text scale
	ScaleFactor float64 // 0.5: half size, 1.0: normal, 2.0: double, 4.0: quadruple

	// Shadow options
	ShadowEnabled          bool
	ShadowHorizontalOffset int // -5 to 5
	ShadowVerticalOffset   int // -5 to 5
	ShadowStyle            ShadowStyle

	// Multi-line text
	TextLines []string
}

// DefaultRenderOptions returns RenderOptions with default values
func DefaultRenderOptions() RenderOptions {
	return RenderOptions{
		CharSpacing:            1,
		WordSpacing:            2,
		LineSpacing:            1,
		Alignment:              CenterAlign,
		TextColor:              "#FFFFFF",
		UseGradient:            false,
		ColorMode:              SingleColor,
		RainbowColors:          []string{}, // Will use default rainbow colors if needed
		RainbowFrame:           0,
		RainbowSpeed:           5,
		ScaleFactor:            1.0,
		ShadowEnabled:          false,
		ShadowHorizontalOffset: 0,
		ShadowVerticalOffset:   0,
		ShadowStyle:            LightShade,
		TextLines:              []string{},
	}
}

// Validation constants for RenderOptions
const (
	MinCharSpacing  = 0
	MaxCharSpacing  = 10
	MinWordSpacing  = 0
	MaxWordSpacing  = 20
	MinLineSpacing  = 0
	MaxLineSpacing  = 10
	MinScaleFactor  = 0.5 // 0.5x
	MaxScaleFactor  = 4.0 // 4x
	MinShadowOffset = -5
	MaxShadowOffset = 5
)

// Validate checks if the RenderOptions are valid and returns an error if not
func (opts *RenderOptions) Validate() error {
	// Validate spacing
	if opts.CharSpacing < MinCharSpacing || opts.CharSpacing > MaxCharSpacing {
		return &ValidationError{Field: "CharSpacing", Value: opts.CharSpacing, Min: MinCharSpacing, Max: MaxCharSpacing}
	}
	if opts.WordSpacing < MinWordSpacing || opts.WordSpacing > MaxWordSpacing {
		return &ValidationError{Field: "WordSpacing", Value: opts.WordSpacing, Min: MinWordSpacing, Max: MaxWordSpacing}
	}
	if opts.LineSpacing < MinLineSpacing || opts.LineSpacing > MaxLineSpacing {
		return &ValidationError{Field: "LineSpacing", Value: opts.LineSpacing, Min: MinLineSpacing, Max: MaxLineSpacing}
	}

	// Validate scale factor
	if opts.ScaleFactor < MinScaleFactor || opts.ScaleFactor > MaxScaleFactor {
		return &ScaleValidationError{Field: "ScaleFactor", Value: opts.ScaleFactor, Min: MinScaleFactor, Max: MaxScaleFactor}
	}

	// Validate shadow offsets
	if opts.ShadowHorizontalOffset < MinShadowOffset || opts.ShadowHorizontalOffset > MaxShadowOffset {
		return &ValidationError{Field: "ShadowHorizontalOffset", Value: opts.ShadowHorizontalOffset, Min: MinShadowOffset, Max: MaxShadowOffset}
	}
	if opts.ShadowVerticalOffset < MinShadowOffset || opts.ShadowVerticalOffset > MaxShadowOffset {
		return &ValidationError{Field: "ShadowVerticalOffset", Value: opts.ShadowVerticalOffset, Min: MinShadowOffset, Max: MaxShadowOffset}
	}

	// Validate alignment
	if opts.Alignment < LeftAlign || opts.Alignment > RightAlign {
		return &ValidationError{Field: "Alignment", Value: int(opts.Alignment), Min: int(LeftAlign), Max: int(RightAlign)}
	}

	// Validate gradient direction
	if opts.GradientDirection < UpDown || opts.GradientDirection > RightLeft {
		return &ValidationError{Field: "GradientDirection", Value: int(opts.GradientDirection), Min: int(UpDown), Max: int(RightLeft)}
	}

	// Validate shadow style
	if opts.ShadowStyle < LightShade || opts.ShadowStyle > DarkShade {
		return &ValidationError{Field: "ShadowStyle", Value: int(opts.ShadowStyle), Min: int(LightShade), Max: int(DarkShade)}
	}

	// Validate color mode
	if opts.ColorMode < SingleColor || opts.ColorMode > Rainbow {
		return &ValidationError{Field: "ColorMode", Value: int(opts.ColorMode), Min: int(SingleColor), Max: int(Rainbow)}
	}

	// Validate color format (basic hex color validation)
	if !isValidHexColor(opts.TextColor) {
		return &ColorValidationError{Field: "TextColor", Value: opts.TextColor}
	}
	if opts.UseGradient && !isValidHexColor(opts.GradientColor) {
		return &ColorValidationError{Field: "GradientColor", Value: opts.GradientColor}
	}

	// Validate rainbow colors if provided
	if opts.ColorMode == Rainbow && len(opts.RainbowColors) > 0 {
		for i, color := range opts.RainbowColors {
			if !isValidHexColor(color) {
				return &ColorValidationError{Field: fmt.Sprintf("RainbowColors[%d]", i), Value: color}
			}
		}
	}

	return nil
}

// ValidationError represents a validation error for RenderOptions
type ValidationError struct {
	Field string
	Value int
	Min   int
	Max   int
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("invalid %s: %d (must be between %d and %d)", e.Field, e.Value, e.Min, e.Max)
}

// ScaleValidationError represents a scale validation error
type ScaleValidationError struct {
	Field string
	Value float64
	Min   float64
	Max   float64
}

func (e *ScaleValidationError) Error() string {
	return fmt.Sprintf("invalid %s: %.1f (must be between %.1f and %.1f)", e.Field, e.Value, e.Min, e.Max)
}

// ColorValidationError represents a color validation error
type ColorValidationError struct {
	Field string
	Value string
}

func (e *ColorValidationError) Error() string {
	return fmt.Sprintf("invalid %s: %s (must be a valid hex color like #FFFFFF)", e.Field, e.Value)
}

// isValidHexColor checks if a string is a valid hex color
func isValidHexColor(color string) bool {
	if len(color) != 7 || color[0] != '#' {
		return false
	}
	for _, c := range color[1:] {
		if !((c >= '0' && c <= '9') || (c >= 'A' && c <= 'F') || (c >= 'a' && c <= 'f')) {
			return false
		}
	}
	return true
}

// ShadowStyleOption represents shadow style options
type ShadowStyleOption struct {
	Name string
	Char rune
	Hex  string
}

// Default shadow style options
var shadowStyleOptions = []ShadowStyleOption{
	{"Light Shade", '░', ""},  // U+2591 LIGHT SHADE - Uses main text color
	{"Medium Shade", '▒', ""}, // U+2592 MEDIUM SHADE - Uses main text color
	{"Dark Shade", '▓', ""},   // U+2593 DARK SHADE - Uses main text color
}

// pixelCoord represents a coordinate on the character grid, with support for half-pixels
type pixelCoord struct {
	x      float64 // Use float64 for half-pixel precision
	y      int
	isHalf bool // Flag to indicate if this is a half-pixel position
}

// DescenderInfo holds information about a character's descender properties
type DescenderInfo struct {
	HasDescender    bool
	BaselineHeight  int // Height of the main character body (excluding descender)
	DescenderHeight int // Height of the descender part
	TotalHeight     int // Total character height
	VerticalOffset  int // How much to offset this character vertically
}
