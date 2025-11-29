// ABOUTME: PNG generator that converts ANSI-colored text to PNG images.
// ABOUTME: Parses ANSI escape sequences and renders characters at 16x scale with transparency.

package export

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"regexp"
	"strconv"
	"unicode/utf8"
)

// CellSize defines the default pixel dimensions per terminal character cell.
// 4x scale per dimension = 16x16 pixels per cell.
const CellSize = 16

// Unicode block characters
const (
	FullBlock       = '█' // U+2588
	UpperHalfBlock  = '▀' // U+2580
	LowerHalfBlock  = '▄' // U+2584
	LightShade      = '░' // U+2591
	MediumShade     = '▒' // U+2592
	DarkShade       = '▓' // U+2593
)

// Alpha values for shade characters (out of 255)
const (
	LightShadeAlpha  = 64  // ~25%
	MediumShadeAlpha = 128 // ~50%
	DarkShadeAlpha   = 191 // ~75%
)

// Regex patterns for parsing ANSI escape sequences
var (
	// Matches 24-bit foreground color: ESC[38;2;R;G;Bm
	colorRegex = regexp.MustCompile(`\x1b\[38;2;(\d+);(\d+);(\d+)m`)
	// Matches any ANSI escape sequence (for stripping)
	ansiStripRegex = regexp.MustCompile(`\x1b\[[0-9;]*m`)
)

// PNGOptions contains configuration for PNG generation
type PNGOptions struct {
	CellWidth  int // Pixels per character cell width (default: CellSize)
	CellHeight int // Pixels per character cell height (default: CellSize)
}

// DefaultPNGOptions returns default PNG generation options (16x16 per cell)
func DefaultPNGOptions() PNGOptions {
	return PNGOptions{
		CellWidth:  CellSize,
		CellHeight: CellSize,
	}
}

// GeneratePNG creates a PNG image from rendered ANSI lines.
// Returns PNG data as bytes or error.
func GeneratePNG(lines []string, options PNGOptions) ([]byte, error) {
	if len(lines) == 0 {
		return nil, fmt.Errorf("no content to export")
	}

	// Use defaults if zero values provided
	if options.CellWidth == 0 {
		options.CellWidth = CellSize
	}
	if options.CellHeight == 0 {
		options.CellHeight = CellSize
	}

	// Calculate image dimensions based on max line width
	maxWidth := 0
	for _, line := range lines {
		lineWidth := countVisibleChars(line)
		if lineWidth > maxWidth {
			maxWidth = lineWidth
		}
	}

	// Handle edge case of all-empty lines
	if maxWidth == 0 {
		maxWidth = 1
	}

	imgWidth := maxWidth * options.CellWidth
	imgHeight := len(lines) * options.CellHeight

	// Create RGBA image with transparent background (zero-initialized = transparent)
	img := image.NewRGBA(image.Rect(0, 0, imgWidth, imgHeight))

	// Render each line
	for lineIdx, line := range lines {
		renderLineToImage(img, line, lineIdx, options)
	}

	// Encode to PNG
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil, fmt.Errorf("failed to encode PNG: %v", err)
	}

	return buf.Bytes(), nil
}

// renderLineToImage renders a single line of ANSI text to the image
func renderLineToImage(img *image.RGBA, line string, lineIdx int, options PNGOptions) {
	// Default color (white, fully opaque)
	currentColor := color.RGBA{R: 255, G: 255, B: 255, A: 255}
	charIdx := 0

	// Process line character by character, tracking ANSI state
	i := 0
	lineBytes := []byte(line)

	for i < len(lineBytes) {
		// Check for ANSI escape sequence (starts with ESC [)
		if lineBytes[i] == 0x1b && i+1 < len(lineBytes) && lineBytes[i+1] == '[' {
			// Find end of escape sequence
			seqStart := i
			i += 2 // Skip ESC [

			// Scan for terminator (letter)
			for i < len(lineBytes) && !isAnsiTerminator(lineBytes[i]) {
				i++
			}
			if i < len(lineBytes) {
				i++ // Include terminator
			}

			// Parse the sequence
			seq := string(lineBytes[seqStart:i])
			if matches := colorRegex.FindStringSubmatch(seq); matches != nil {
				r, _ := strconv.Atoi(matches[1])
				g, _ := strconv.Atoi(matches[2])
				b, _ := strconv.Atoi(matches[3])
				currentColor = color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: 255}
			}
			// Reset codes (\x1b[0m) are handled implicitly - color stays until changed
			continue
		}

		// Decode UTF-8 character
		r, size := utf8.DecodeRune(lineBytes[i:])
		if r == utf8.RuneError && size == 1 {
			// Invalid UTF-8, skip byte
			i++
			continue
		}

		// Render the character
		drawCell(img, charIdx, lineIdx, r, currentColor, options)
		charIdx++
		i += size
	}
}

// drawCell draws a single character cell to the image
func drawCell(img *image.RGBA, x, y int, char rune, c color.RGBA, options PNGOptions) {
	cellX := x * options.CellWidth
	cellY := y * options.CellHeight
	halfHeight := options.CellHeight / 2

	switch char {
	case FullBlock:
		// Fill entire cell
		fillRect(img, cellX, cellY, options.CellWidth, options.CellHeight, c)

	case UpperHalfBlock:
		// Fill top half only
		fillRect(img, cellX, cellY, options.CellWidth, halfHeight, c)

	case LowerHalfBlock:
		// Fill bottom half only
		fillRect(img, cellX, cellY+halfHeight, options.CellWidth, halfHeight, c)

	case LightShade:
		// Full cell with low alpha
		shadeColor := color.RGBA{R: c.R, G: c.G, B: c.B, A: LightShadeAlpha}
		fillRect(img, cellX, cellY, options.CellWidth, options.CellHeight, shadeColor)

	case MediumShade:
		// Full cell with medium alpha
		shadeColor := color.RGBA{R: c.R, G: c.G, B: c.B, A: MediumShadeAlpha}
		fillRect(img, cellX, cellY, options.CellWidth, options.CellHeight, shadeColor)

	case DarkShade:
		// Full cell with high alpha
		shadeColor := color.RGBA{R: c.R, G: c.G, B: c.B, A: DarkShadeAlpha}
		fillRect(img, cellX, cellY, options.CellWidth, options.CellHeight, shadeColor)

	case ' ':
		// Space - leave transparent (do nothing)

	default:
		// For any other printable character, fill as full block
		// This handles edge cases where other characters might be used
		if char > 32 { // Printable ASCII/Unicode
			fillRect(img, cellX, cellY, options.CellWidth, options.CellHeight, c)
		}
		// Non-printable or control chars: leave transparent
	}
}

// fillRect fills a rectangle in the image with the given color
func fillRect(img *image.RGBA, x, y, width, height int, c color.RGBA) {
	bounds := img.Bounds()
	for dy := 0; dy < height; dy++ {
		for dx := 0; dx < width; dx++ {
			px, py := x+dx, y+dy
			if px >= bounds.Min.X && px < bounds.Max.X && py >= bounds.Min.Y && py < bounds.Max.Y {
				img.SetRGBA(px, py, c)
			}
		}
	}
}

// countVisibleChars counts visible (non-ANSI) characters in a line
func countVisibleChars(line string) int {
	stripped := ansiStripRegex.ReplaceAllString(line, "")
	return utf8.RuneCountInString(stripped)
}

// isAnsiTerminator checks if a byte terminates an ANSI escape sequence
func isAnsiTerminator(b byte) bool {
	return (b >= 'A' && b <= 'Z') || (b >= 'a' && b <= 'z')
}
