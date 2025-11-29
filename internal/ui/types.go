// ABOUTME: Type definitions for the TUI application models and state.
// ABOUTME: Contains all sub-models for text, font, spacing, color, scale, shadow, and export.

package ui

import (
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/superstarryeyes/bit/internal/export"
	"github.com/superstarryeyes/bit/internal/favorites"
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

// exportModel handles export functionality
type exportModel struct {
	active               bool                  // Whether we're in export mode
	format               string                // Selected export format
	filenameInput        textinput.Model       // Text input for filename
	showConfirmation     bool                  // Whether to show export confirmation in header
	confirmationText     string                // The confirmation text to display
	showOverwritePrompt  bool                  // Whether to show overwrite confirmation
	overwriteFilename    string                // Filename that would be overwritten
	overwriteContent     string                // Content to write if user confirms (text formats)
	overwriteBinaryContent []byte              // Content to write if user confirms (binary formats like PNG)
	overwriteFormat      string                // Format for the overwrite
	selectedButton       int                   // 0 = Yes, 1 = No
	manager              *export.ExportManager // Export manager for format information
}

// uiStateModel handles general UI state
type uiStateModel struct {
	focusedPanel  FocusedPanel // Currently focused panel
	width         int
	height        int
	renderedLines []string // Rendered text cache
	usesTwoRows   bool     // Cache the layout decision to prevent flickering
}

// favoritesModel handles favorites functionality
type favoritesModel struct {
	manager          *favorites.Manager // Favorites manager for persistence
	active           bool               // Whether favorites view is open
	selectedIndex    int                // Currently selected favorite in list
	nameInput        textinput.Model    // Text input for naming new favorites
	showNamePrompt   bool               // Whether showing the name input prompt
	showConfirmation bool               // Whether to show confirmation message
	confirmationText string             // Confirmation text to display
}

// model is the main application model composed of sub-models
type model struct {
	textInput textInputModel
	font      fontModel
	spacing   spacingModel
	color     colorModel
	scale     scaleModel
	shadow    shadowModel
	export    exportModel
	favorites favoritesModel
	uiState   uiStateModel
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
