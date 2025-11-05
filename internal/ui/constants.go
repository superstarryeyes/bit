package ui

// Panel identifiers using iota for type safety
type FocusedPanel int

const (
	TextInputPanel FocusedPanel = iota
	FontPanel
	SpacingPanel
	ColorPanel
	ScalePanel
	ShadowPanel
	BackgroundPanel
	AnimationPanel
	TotalPanels // Used for modulo operations
)

// Spacing modes for the spacing panel
type SpacingMode int

const (
	CharacterSpacingMode SpacingMode = iota
	WordSpacingMode
	LineSpacingMode
	TotalSpacingModes
)

// Text input modes for the text input panel
type TextInputMode int

const (
	TextEntryMode TextInputMode = iota
	TextAlignmentMode
	TotalTextInputModes
)

// Color sub-modes for the color panel
type ColorSubMode int

const (
	TextColorMode ColorSubMode = iota
	GradientColorMode
	GradientDirectionMode
	RainbowMode
	TotalColorSubModes
)

// Shadow sub-modes for the shadow panel
type ShadowSubMode int

const (
	HorizontalShadowMode ShadowSubMode = iota
	VerticalShadowMode
	ShadowStyleMode
	TotalShadowSubModes
)

// Background sub-modes for the background panel
type BackgroundSubMode int

const (
	BackgroundTypeMode BackgroundSubMode = iota
	TotalBackgroundSubModes
)

// Background types
type BackgroundType int

const (
	BackgroundNone BackgroundType = iota
	BackgroundLavaLamp
	BackgroundWavyGrid
	BackgroundTicker
	BackgroundStarfield
	TotalBackgroundTypes
)

// Text alignment options
type TextAlignment int

const (
	LeftAlignment TextAlignment = iota
	CenterAlignment
	RightAlignment
	TotalAlignments
)

// Text scale options
type TextScale int

const (
	ScaleHalf TextScale = -1 // 0.5x
	ScaleOne  TextScale = 0  // 1x
	ScaleTwo  TextScale = 1  // 2x
	ScaleFour TextScale = 2  // 4x
	MinScale            = ScaleHalf
	MaxScale            = ScaleFour
)

// Gradient direction indices
type GradientDirection int

const (
	GradientUpDown GradientDirection = iota
	GradientDownUp
	GradientLeftRight
	GradientRightLeft
	TotalGradientDirections
)

// Cursor mode constants
const (
	CursorBlink = 1
)

// Shadow pixel range constants
const (
	MinShadowPixels         = -5
	MaxShadowPixels         = 5
	DefaultShadowPixels     = 5 // Index for "Off" in shadowPixelOptions
	MinVerticalShadowPixels = -5
	MaxVerticalShadowPixels = 5
)

// Spacing range constants
const (
	MinCharSpacing = 0
	MaxCharSpacing = 10
	MinWordSpacing = 0
	MaxWordSpacing = 20
	MinLineSpacing = 0
	MaxLineSpacing = 10
)

// Default values
const (
	DefaultCharSpacing = 2
	DefaultWordSpacing = 2
	DefaultLineSpacing = 1
	DefaultTextScale   = ScaleOne
)

// Text input constraints
const (
	TextInputCharLimit     = 50
	TextInputMinWidth      = 17
	TextInputMaxWidth      = 50
	FilenameInputCharLimit = 50
	FilenameInputWidth     = 40
	MaxFilenameLength      = 200 // Maximum filename length before extension
)

// Layout thresholds
const (
	MinWidthSingleRow         = 65
	ComfortableWidthSingleRow = 80
	LayoutReservedMargin      = 12 // Fixed margin for borders and spacing
	LayoutMinPanelWidth       = 8  // Absolute minimum panel width
	LayoutSpacerWidth         = 1  // Fixed spacer between panels
)

// Color constants
const (
	MaxRGBValue          = 255 // Maximum RGB color value
	MaxShadowRepeatCount = 20  // Maximum repetition for shadow characters
)

// Background animation constants
const (
	RainbowAnimationSpeed = 5   // Frames per color shift
	TickerSpeed           = 3   // Frames per scroll step
	BlobCount             = 4   // Number of blobs in lava lamp
	DefaultGridSize       = 10  // Grid spacing for wavy grid
	StarCount             = 50  // Number of stars in starfield
	StarSpeed             = 1.5 // Speed of stars moving toward viewer
)

// Animation sub-modes for the animation panel
type AnimationSubMode int

const (
	AnimationTypeMode AnimationSubMode = iota
	AnimationSpeedMode
	TotalAnimationSubModes
)

// Animation types for text
type AnimationType int

const (
	AnimationNone AnimationType = iota
	AnimationScrollLeft
	AnimationScrollRight
	TotalAnimationTypes
)

// Animation speeds
type AnimationSpeed int

const (
	AnimationSlow AnimationSpeed = iota
	AnimationMedium
	AnimationFast
	TotalAnimationSpeeds
)

// Animation speed values (frames per scroll step)
const (
	ScrollSpeedSlow   = 8  // Scroll every 8 frames (slower)
	ScrollSpeedMedium = 4  // Scroll every 4 frames
	ScrollSpeedFast   = 2  // Scroll every 2 frames (faster)
)
