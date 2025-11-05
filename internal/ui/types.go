package ui

import (
	"github.com/superstarryeyes/bit/internal/export"
	"github.com/charmbracelet/bubbles/textinput"
)

// textInputModel handles text entry and alignment
type textInputModel struct {
	input       textinput.Model
	currentText string
	textRows    []string      // Multiple rows of text
	rowCursors  []int         // Cursor positions for each row
	currentRow  int           // Currently selected row for editing
	alignment   TextAlignment // Text alignment
	mode        TextInputMode // Text input panel sub-mode
}

// fontModel handles font selection
type fontModel struct {
	fonts        []FontInfo
	selectedFont int
}

// spacingModel handles character, word, and line spacing
type spacingModel struct {
	charSpacing int
	wordSpacing int
	lineSpacing int
	mode        SpacingMode // Spacing panel sub-mode
}

// colorModel handles text color and gradient settings
type colorModel struct {
	textColor         int               // Index into colorOptions array
	gradientColor     int               // Index into colorOptions array for gradient end color
	gradientEnabled   bool              // Whether gradient is enabled
	gradientDirection GradientDirection // Gradient direction
	rainbowEnabled    bool              // Whether rainbow mode is enabled
	subMode           ColorSubMode      // Color panel sub-mode
}

// scaleModel handles text scaling
type scaleModel struct {
	scale TextScale // Text scaling factor
}

// shadowModel handles shadow settings
type shadowModel struct {
	enabled          bool          // Whether shadow is enabled
	horizontalOffset int           // Horizontal shadow offset in pixels (canonical value -5..5)
	verticalOffset   int           // Vertical shadow offset in pixels (canonical value -5..5)
	horizontalIndex  int           // Index into shadowPixelOptions array (UI mapping)
	verticalIndex    int           // Index into verticalShadowPixelOptions array (UI mapping)
	style            int           // Index into shadowStyleOptions array (ANSI block styles)
	showWarning      bool          // Whether to show the shadow warning message
	subMode          ShadowSubMode // Shadow panel sub-mode
}

// backgroundModel handles background effects
type backgroundModel struct {
	enabled       bool              // Whether background is enabled
	backgroundType BackgroundType   // Type of background effect
	subMode       BackgroundSubMode // Background panel sub-mode
	lavaLamp      *LavaLamp         // Lava lamp effect engine
	wavyGrid      *WavyGrid         // Wavy grid effect engine
	ticker        *Ticker           // Ticker/sidescroller effect
	starfield     *Starfield        // Starfield effect engine
	frame         int               // Animation frame counter
}

// animationModel handles text animation
type animationModel struct {
	animationType AnimationType      // Type of animation
	speed         AnimationSpeed     // Animation speed
	subMode       AnimationSubMode   // Animation panel sub-mode
	scrollOffset  int                // Current horizontal scroll offset (in characters)
}

// exportModel handles export functionality
type exportModel struct {
	active              bool                  // Whether we're in export mode
	format              string                // Selected export format
	filenameInput       textinput.Model       // Text input for filename
	showConfirmation    bool                  // Whether to show export confirmation in header
	confirmationText    string                // The confirmation text to display
	showOverwritePrompt bool                  // Whether to show overwrite confirmation
	overwriteFilename   string                // Filename that would be overwritten
	overwriteContent    string                // Content to write if user confirms
	overwriteFormat     string                // Format for the overwrite
	selectedButton      int                   // 0 = Yes, 1 = No
	manager             *export.ExportManager // Export manager for format information
}

// uiStateModel handles general UI state
type uiStateModel struct {
	focusedPanel  FocusedPanel // Currently focused panel
	width         int
	height        int
	renderedLines []string // Rendered text cache
	usesTwoRows   bool     // Cache the layout decision to prevent flickering
}

// model is the main application model composed of sub-models
type model struct {
	textInput  textInputModel
	font       fontModel
	spacing    spacingModel
	color      colorModel
	scale      scaleModel
	shadow     shadowModel
	background backgroundModel
	animation  animationModel
	export     exportModel
	uiState    uiStateModel
}

// FontInfo holds information about available fonts
type FontInfo struct {
	Name     string
	Path     string
	FontData *FontData // Pointer to allow for lazy loading (nil when not loaded)
	Loaded   bool      // Flag to indicate if font data is loaded
}

// FontData represents the overall structure of our .bit font file (JSON format)
type FontData struct {
	Name       string              `json:"name"`
	Author     string              `json:"author"`
	License    string              `json:"license"`
	Characters map[string][]string `json:"characters"`
}

// Color options for text with proper hex codes
type ColorOption struct {
	Name string
	Hex  string
}

// Gradient direction options
type GradientDirectionOption struct {
	Name string
}

// Background effect structures
// LavaLamp represents the metaball/lava lamp effect with floating blobs
type LavaLamp struct {
	Blobs  []Blob
	Width  int
	Height int
	Frame  int
}

// Blob represents a single floating blob for the lava lamp effect
type Blob struct {
	X, Y       float64 // Position
	VX, VY     float64 // Velocity
	Radius     float64 // Size
	ColorIndex int     // Index into color palette
}

// WavyGrid represents an animated grid background with sine wave distortion
type WavyGrid struct {
	Width    int
	Height   int
	Frame    int
	GridSize int
}

// Ticker represents a sidescroller ticker effect
type Ticker struct {
	Text   string
	Offset int
	Speed  int // How many frames before scrolling
}

// Starfield represents a 3D starfield effect with perspective projection
type Starfield struct {
	Stars  []Star
	Width  int
	Height int
	Frame  int
}

// Star represents a single star/icon in 3D space
type Star struct {
	X, Y, Z float64 // 3D position (Z is depth)
	Icon    string  // Display character
}
