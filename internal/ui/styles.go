package ui

import (
	"github.com/charmbracelet/lipgloss"
)

// colorOptions references the centralized color palette
var colorOptions = GetANSIColorOptions()

// Shadow pixel options for shadow extension
type ShadowPixelOption struct {
	Name   string
	Pixels int
}

var shadowPixelOptions = []ShadowPixelOption{
	{"5 pixels ←", -5},
	{"4 pixels ←", -4},
	{"3 pixels ←", -3},
	{"2 pixels ←", -2},
	{"1 pixel ←", -1},
	{"Off", 0},
	{"1 pixel →", 1},
	{"2 pixels →", 2},
	{"3 pixels →", 3},
	{"4 pixels →", 4},
	{"5 pixels →", 5},
}

// Vertical shadow pixel options
var verticalShadowPixelOptions = []ShadowPixelOption{
	{"5 pixels ↑", -5},
	{"4 pixels ↑", -4},
	{"3 pixels ↑", -3},
	{"2 pixels ↑", -2},
	{"1 pixel ↑", -1},
	{"Off", 0},
	{"1 pixel ↓", 1},
	{"2 pixels ↓", 2},
	{"3 pixels ↓", 3},
	{"4 pixels ↓", 4},
	{"5 pixels ↓", 5},
}

// Shadow style options using ANSI block characters
type ShadowStyleOption struct {
	Name string
	Char rune
	Hex  string
}

var shadowStyleOptions = []ShadowStyleOption{
	{"Light Shade", '░', ""},  // U+2591 LIGHT SHADE - Uses main text color
	{"Medium Shade", '▒', ""}, // U+2592 MEDIUM SHADE - Uses main text color
	{"Dark Shade", '▓', ""},   // U+2593 DARK SHADE - Uses main text color
}

// Gradient direction options
var gradientDirectionOptions = []GradientDirectionOption{
	{"Up-Down"},
	{"Down-Up"},
	{"Left-Right"},
	{"Right-Left"},
}

// Color variables - now referencing the centralized color palette
var (
	ColorWhite     = ColorPalette["White"]
	ColorRed       = ColorPalette["PureRed"]
	ColorTextInput = ColorPalette["TextInput"]
	ColorExport    = ColorPalette["Export"]
	ColorGray      = ColorPalette["FaintGray"]
	ColorFaint     = ColorPalette["VeryFaint"]
)

// Base styles for the application
var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(ColorPalette["TitleFG"])).
			Background(lipgloss.Color(ColorPalette["TitleBG"])).
			Padding(0, 1)

	// Text input styles
	textInputCursorStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(ColorWhite))
	textInputTextStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(ColorWhite)).
				Background(lipgloss.Color(ColorTextInput))
	textInputPlaceholderStyle = lipgloss.NewStyle().
					Foreground(lipgloss.Color(ColorWhite)).
					Background(lipgloss.Color(ColorTextInput)).
					Faint(true)

	// Filename input styles
	filenameInputTextStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(ColorWhite)).
				Background(lipgloss.Color(ColorExport))
	filenameInputPlaceholderStyle = lipgloss.NewStyle().
					Foreground(lipgloss.Color(ColorWhite)).
					Background(lipgloss.Color(ColorExport)).
					Faint(true)

	// Warning style
	warningStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorWhite)).
			Background(lipgloss.Color(ColorRed)).
			Bold(true).
			Padding(0, 1)
)

// LabelStyles holds all label styles for different panel types
type LabelStyles struct {
	TextInput   lipgloss.Style
	Font        lipgloss.Style
	CharSpacing lipgloss.Style
	WordSpacing lipgloss.Style
	LineSpacing lipgloss.Style
	Color       lipgloss.Style
	Scale       lipgloss.Style
}

// createLabelStyles creates and returns all label styles
func createLabelStyles() LabelStyles {
	return LabelStyles{
		TextInput: lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorPalette["TextInput"])).
			Bold(true),
		Font: lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorPalette["FontPanel"])).
			Bold(true),
		CharSpacing: lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorPalette["CharSpacing"])).
			Bold(true),
		WordSpacing: lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorPalette["WordSpacing"])).
			Bold(true),
		LineSpacing: lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorPalette["LineSpacing"])).
			Bold(true),
		Color: lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorPalette["ColorPanel"])).
			Bold(true),
		Scale: lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorPalette["ScalePanel"])).
			Bold(true),
	}
}

// Panel styles factory functions for dynamic sizing
func createPanelStyles(panelWidth int) (map[string]lipgloss.Style, map[string]lipgloss.Style) {
	normalStyles := map[string]lipgloss.Style{
		"textInput": lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(ColorPalette["TextInput"])).
			Padding(0, 1).
			Width(panelWidth).
			Height(1),
		"font": lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(ColorPalette["FontPanel"])).
			Padding(0, 1).
			Width(panelWidth).
			Height(1),
		"charSpacing": lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(ColorPalette["CharSpacing"])).
			Padding(0, 1).
			Width(panelWidth).
			Height(1),
		"wordSpacing": lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(ColorPalette["WordSpacing"])).
			Padding(0, 1).
			Width(panelWidth).
			Height(1),
		"lineSpacing": lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(ColorPalette["LineSpacing"])).
			Padding(0, 1).
			Width(panelWidth).
			Height(1),
		"color": lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(ColorPalette["ColorPanel"])).
			Padding(0, 1).
			Width(panelWidth).
			Height(1),
		"scale": lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(ColorPalette["ScalePanel"])).
			Padding(0, 1).
			Width(panelWidth).
			Height(1),
		"shadow": lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(ColorPalette["Shadow"])).
			Padding(0, 1).
			Width(panelWidth).
			Height(1),
		"verticalShadow": lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(ColorPalette["Shadow"])).
			Padding(0, 1).
			Width(panelWidth).
			Height(1),
		"background": lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(ColorPalette["Background"])).
			Padding(0, 1).
			Width(panelWidth).
			Height(1),
		"animation": lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(ColorPalette["Animation"])).
			Padding(0, 1).
			Width(panelWidth).
			Height(1),
	}

	selectedStyles := map[string]lipgloss.Style{
		"textInput": lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(ColorPalette["TextInput"])).
			Background(lipgloss.Color(ColorPalette["TextInput"])).
			Foreground(lipgloss.Color(ColorPalette["White"])).
			Padding(0, 1).
			Width(panelWidth).
			Height(1),
		"font": lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(ColorPalette["FontPanel"])).
			Background(lipgloss.Color(ColorPalette["FontPanel"])).
			Foreground(lipgloss.Color(ColorPalette["White"])).
			Padding(0, 1).
			Width(panelWidth).
			Height(1),
		"charSpacing": lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(ColorPalette["CharSpacing"])).
			Background(lipgloss.Color(ColorPalette["CharSpacing"])).
			Foreground(lipgloss.Color(ColorPalette["White"])).
			Padding(0, 1).
			Width(panelWidth).
			Height(1),
		"wordSpacing": lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(ColorPalette["WordSpacing"])).
			Background(lipgloss.Color(ColorPalette["WordSpacing"])).
			Foreground(lipgloss.Color(ColorPalette["White"])).
			Padding(0, 1).
			Width(panelWidth).
			Height(1),
		"lineSpacing": lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(ColorPalette["LineSpacing"])).
			Background(lipgloss.Color(ColorPalette["LineSpacing"])).
			Foreground(lipgloss.Color(ColorPalette["White"])).
			Padding(0, 1).
			Width(panelWidth).
			Height(1),
		"color": lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(ColorPalette["ColorPanel"])).
			Background(lipgloss.Color(ColorPalette["ColorPanel"])).
			Foreground(lipgloss.Color(ColorPalette["Black"])).
			Padding(0, 1).
			Width(panelWidth).
			Height(1),
		"scale": lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(ColorPalette["ScalePanel"])).
			Background(lipgloss.Color(ColorPalette["ScalePanel"])).
			Foreground(lipgloss.Color(ColorPalette["Black"])).
			Padding(0, 1).
			Width(panelWidth).
			Height(1),
		"shadow": lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(ColorPalette["Shadow"])).
			Background(lipgloss.Color(ColorPalette["Shadow"])).
			Foreground(lipgloss.Color(ColorPalette["White"])).
			Padding(0, 1).
			Width(panelWidth).
			Height(1),
		"verticalShadow": lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(ColorPalette["Shadow"])).
			Background(lipgloss.Color(ColorPalette["Shadow"])).
			Foreground(lipgloss.Color(ColorPalette["White"])).
			Padding(0, 1).
			Width(panelWidth).
			Height(1),
		"background": lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(ColorPalette["Background"])).
			Background(lipgloss.Color(ColorPalette["Background"])).
			Foreground(lipgloss.Color(ColorPalette["White"])).
			Padding(0, 1).
			Width(panelWidth).
			Height(1),
		"animation": lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(ColorPalette["Animation"])).
			Background(lipgloss.Color(ColorPalette["Animation"])).
			Foreground(lipgloss.Color(ColorPalette["White"])).
			Padding(0, 1).
			Width(panelWidth).
			Height(1),
	}

	return normalStyles, selectedStyles
}

// Create fixed text display style with dynamic sizing
func createFixedTextDisplayStyle(width, height int) lipgloss.Style {
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(ColorPalette["TextDisplay"])).
		PaddingTop(0).
		PaddingBottom(0).
		PaddingLeft(0).
		PaddingRight(0).
		Align(lipgloss.Center).
		AlignVertical(lipgloss.Center).
		Width(width).
		Height(height)
}
