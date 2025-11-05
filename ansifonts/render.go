package ansifonts

import (
	"fmt"
	"math"
	"regexp"
	"strings"
	"unicode/utf8"
)

// DefaultRainbowColors provides the standard rainbow color palette
var DefaultRainbowColors = []string{
	"#FF0000", // Red
	"#FF7F00", // Orange
	"#FFFF00", // Yellow
	"#00FF00", // Green
	"#00FFFF", // Cyan
	"#0000FF", // Blue
	"#8B00FF", // Violet
}

// DetectHalfPixelUsage checks if the current text rendering would use half-pixels
// that would interfere with shadow rendering. This function should only return true
// when shadows would actually cause visual artifacts.
func DetectHalfPixelUsage(text string, fontData FontData, scaleFactor float64) bool {
	if text == "" {
		return false
	}

	runes := []rune(text)
	if len(runes) == 0 {
		return false
	}

	// Only check: Look for half-pixel ANSI characters (▀, ▄) in the SCALED font definition
	// These are the only characters that actually cause issues with shadow rendering
	halfPixelChars := []rune{'▀', '▄'}
	for _, r := range runes {
		charStr := string(r)
		if bitmapLines, ok := fontData.Characters[charStr]; ok {
			scaledBitmapLines := scaleCharacter(bitmapLines, scaleFactor)
			for _, line := range scaledBitmapLines {
				for _, halfPixelChar := range halfPixelChars {
					if strings.ContainsRune(line, halfPixelChar) {
						return true
					}
				}
			}
		}
	}

	return false
}

// RenderTextWithFont renders text using the specified font with advanced rendering options
func RenderTextWithFont(text string, fontData FontData, options RenderOptions) []string {
	if text == "" {
		return []string{}
	}

	// Validate font data
	if fontData.Name == "" {
		return []string{}
	}

	if fontData.Characters == nil {
		return []string{}
	}

	// Encapsulate shadow compatibility logic within the library
	// If half-pixels are detected and shadows are enabled with non-zero offsets,
	// automatically disable shadows to prevent visual artifacts
	if options.ShadowEnabled && (options.ShadowHorizontalOffset != 0 || options.ShadowVerticalOffset != 0) {
		hasHalfPixels := DetectHalfPixelUsage(text, fontData, options.ScaleFactor)
		if hasHalfPixels {
			options.ShadowEnabled = false
		}
	}

	// Split text into lines to process each one independently
	textLines := strings.Split(text, "\n")
	var allRenderedLines []string
	var renderedTextLines [][]string

	// First pass: render each text line and find the maximum width for alignment
	maxTextLineWidth := 0
	for _, line := range textLines {
		if line == "" {
			renderedTextLines = append(renderedTextLines, []string{""})
			continue
		}

		lineRendered := renderTextWithFont(line, fontData, options.CharSpacing, float64(options.WordSpacing), options.ScaleFactor)
		lineRendered = stripEmptyLines(lineRendered)

		lineWidth := 0
		for _, row := range lineRendered {
			lineWidth = max(lineWidth, utf8.RuneCountInString(stripANSI(row)))
		}

		maxTextLineWidth = max(maxTextLineWidth, lineWidth)
		renderedTextLines = append(renderedTextLines, lineRendered)
	}

	// Second pass: apply alignment, styling, and shadow to each text line's block
	for i, lineRendered := range renderedTextLines {
		if len(lineRendered) == 1 && lineRendered[0] == "" {
			if i > 0 {
				allRenderedLines = append(allRenderedLines, "")
			}
			continue
		}

		// Apply alignment to the current line's rendered block
		alignedBlock := applyAlignmentToTextLine(lineRendered, maxTextLineWidth, options.Alignment)

		// Apply styling and shadow
		finalBlock := applyStylingAndShadow(alignedBlock, options)

		// Add configurable spacing between text lines
		if i > 0 && len(allRenderedLines) > 0 {
			for range options.LineSpacing {
				allRenderedLines = append(allRenderedLines, "")
			}
		}

		allRenderedLines = append(allRenderedLines, finalBlock...)
	}

	// Final pass to ensure all lines have the same width for consistent rendering
	maxWidth := 0
	for _, line := range allRenderedLines {
		maxWidth = max(maxWidth, utf8.RuneCountInString(stripANSI(line)))
	}

	for i, line := range allRenderedLines {
		lineWidth := utf8.RuneCountInString(stripANSI(line))
		if lineWidth < maxWidth {
			allRenderedLines[i] = line + strings.Repeat(" ", maxWidth-lineWidth)
		}
	}

	return allRenderedLines
}

// applyStylingAndShadow provides a unified way to render a text block and its shadow.
// It correctly handles both single colors and independent gradients.
func applyStylingAndShadow(plainBlock []string, options RenderOptions) []string {
	if len(plainBlock) == 0 {
		return plainBlock
	}

	// --- Parameter Setup ---
	var shadowPixels, verticalShadowPixels int
	var shadowChar rune
	if options.ShadowEnabled {
		shadowPixels = options.ShadowHorizontalOffset
		verticalShadowPixels = options.ShadowVerticalOffset
		shadowChar = shadowStyleOptions[options.ShadowStyle].Char
	}

	// Determine color mode and setup colors
	// Priority: ColorMode field takes precedence over legacy UseGradient
	colorMode := options.ColorMode
	if colorMode == SingleColor && options.UseGradient && options.GradientColor != options.TextColor {
		colorMode = Gradient // Support legacy UseGradient flag
	}

	// Setup colors based on mode
	startColorHex := options.TextColor
	var endColorHex string
	var rainbowColors []string

	if colorMode == Gradient {
		endColorHex = options.GradientColor
	} else if colorMode == Rainbow {
		// Use custom rainbow colors if provided, otherwise use defaults
		if len(options.RainbowColors) > 0 {
			rainbowColors = options.RainbowColors
		} else {
			rainbowColors = DefaultRainbowColors
		}
	}

	startR, startG, startB := hexToRGB(startColorHex)
	endR, endG, endB := hexToRGB(endColorHex)

	// Single color setup
	shadowStyleHex := shadowStyleOptions[options.ShadowStyle].Hex
	var shadowColorForStyle string
	if shadowStyleHex != "" {
		shadowColorForStyle = shadowStyleHex
	} else {
		shadowColorForStyle = startColorHex // Shadow inherits main text color by default
	}

	// --- Canvas Calculation ---
	blockHeight := len(plainBlock)
	blockWidth := 0
	for _, line := range plainBlock {
		blockWidth = max(blockWidth, utf8.RuneCountInString(line))
	}

	canvasMinX, canvasMaxX := 0, blockWidth
	canvasMinY, canvasMaxY := 0, blockHeight

	if shadowPixels < 0 {
		canvasMinX = shadowPixels
	} else if shadowPixels > 0 {
		canvasMaxX = blockWidth + shadowPixels
	}
	if verticalShadowPixels < 0 {
		canvasMinY = verticalShadowPixels
	} else if verticalShadowPixels > 0 {
		canvasMaxY = blockHeight + verticalShadowPixels
	}
	canvasWidth := canvasMaxX - canvasMinX
	canvasHeight := canvasMaxY - canvasMinY

	// --- Canvas Creation ---
	type canvasCell struct {
		char    rune
		isMain  bool
		lineIdx int // Original row index for gradient calculation
		charIdx int // Original col index for gradient calculation
	}
	canvas := make([][]canvasCell, canvasHeight)
	for i := range canvas {
		canvas[i] = make([]canvasCell, canvasWidth)
		// All cells are already zero-initialized with char: 0 (rune zero value)
		// Set them to space explicitly
		for j := range canvas[i] {
			canvas[i][j] = canvasCell{char: ' '}
		}
	}

	// --- Render to Canvas (Shadow first, then Main Text) ---
	if options.ShadowEnabled {
		shadowOffsetX := -canvasMinX + shadowPixels
		shadowOffsetY := -canvasMinY + verticalShadowPixels
		for y, line := range plainBlock {
			lineRunes := []rune(line)
			for x, r := range lineRunes {
				if r != ' ' {
					targetX, targetY := shadowOffsetX+x, shadowOffsetY+y
					if targetX >= 0 && targetX < canvasWidth && targetY >= 0 && targetY < canvasHeight {
						canvas[targetY][targetX] = canvasCell{char: shadowChar, isMain: false, lineIdx: y, charIdx: x}
					}
				}
			}
		}
	}

	mainOffsetX := -canvasMinX
	mainOffsetY := -canvasMinY
	for y, line := range plainBlock {
		lineRunes := []rune(line)
		for x, r := range lineRunes {
			if r != ' ' {
				targetX, targetY := mainOffsetX+x, mainOffsetY+y
				if targetX >= 0 && targetX < canvasWidth && targetY >= 0 && targetY < canvasHeight {
					canvas[targetY][targetX] = canvasCell{char: r, isMain: true, lineIdx: y, charIdx: x}
				}
			}
		}
	}

	// --- Convert Canvas to Styled Strings ---
	var result []string
	for y := range canvasHeight {
		var builder strings.Builder
		for x := range canvasWidth {
			cell := canvas[y][x]
			if cell.char == ' ' {
				builder.WriteRune(' ')
				continue
			}

			var cellColorHex string
			if colorMode == Rainbow && cell.isMain {
				// Rainbow mode: cycle through rainbow colors based on character position and animation frame
				// The frame offset creates the animation effect - colors shift as frame increments
				frameOffset := 0
				if options.RainbowSpeed > 0 {
					frameOffset = options.RainbowFrame / options.RainbowSpeed
				}
				colorIdx := (cell.charIdx + cell.lineIdx + frameOffset) % len(rainbowColors)
				cellColorHex = rainbowColors[colorIdx]
			} else if colorMode == Gradient {
				var factor float64
				switch options.GradientDirection {
				case UpDown: // Up-Down
					if blockHeight > 1 {
						factor = float64(cell.lineIdx) / float64(blockHeight-1)
					}
				case DownUp: // Down-Up
					if blockHeight > 1 {
						factor = 1.0 - (float64(cell.lineIdx) / float64(blockHeight-1))
					}
				case LeftRight, RightLeft: // Left-Right, Right-Left
					// For horizontal gradients, calculate factor based on the entire block width
					// rather than individual line widths to ensure consistency across characters
					// with varying heights (ascenders/descenders)
					if canvasWidth > 1 {
						// Calculate the actual x position in the canvas for gradient calculation
						actualX := x
						factor = float64(actualX) / float64(canvasWidth-1)
					}
					if options.GradientDirection == RightLeft {
						factor = 1.0 - factor
					}
				}
				r := int(float64(startR) + factor*float64(endR-startR))
				g := int(float64(startG) + factor*float64(endG-startG))
				b := int(float64(startB) + factor*float64(endB-startB))
				cellColorHex = rgbToHex(clamp(r, 0, 255), clamp(g, 0, 255), clamp(b, 0, 255))
			} else {
				// Single color mode
				if cell.isMain {
					cellColorHex = startColorHex
				} else {
					cellColorHex = shadowColorForStyle
				}
			}
			// Use true color (24-bit RGB) for smoother gradients
			r, g, b := hexToRGB(cellColorHex)
			builder.WriteString(fmt.Sprintf("\x1b[38;2;%d;%d;%dm%s\x1b[0m", r, g, b, string(cell.char)))
		}
		result = append(result, strings.TrimRight(builder.String(), " "))
	}
	return result
}

// applyAlignmentToTextLine applies alignment to a single rendered text line
func applyAlignmentToTextLine(lineRendered []string, maxTextLineWidth int, alignment TextAlignment) []string {
	if len(lineRendered) == 0 {
		return lineRendered
	}

	// Find the actual width of this text line (use the widest row)
	lineWidth := 0
	for _, row := range lineRendered {
		lineWidth = max(lineWidth, utf8.RuneCountInString(stripANSI(row)))
	}

	// If this line is already at max width, no alignment needed
	if lineWidth >= maxTextLineWidth {
		return lineRendered
	}

	// Calculate padding once for the entire text line based on its maximum width
	var leftPadding int
	switch alignment {
	case LeftAlign: // Left alignment
		leftPadding = 0
	case CenterAlign: // Center alignment
		leftPadding = (maxTextLineWidth - lineWidth) / 2
	case RightAlign: // Right alignment
		leftPadding = maxTextLineWidth - lineWidth
	default:
		// Default to left alignment
		leftPadding = 0
	}

	// Apply the same padding to all rows in this text line
	alignedRows := make([]string, len(lineRendered))
	for i, row := range lineRendered {
		rowWidth := utf8.RuneCountInString(stripANSI(row))

		// Add left padding
		alignedRow := strings.Repeat(" ", leftPadding) + row

		// Add right padding to make each row the same total width
		totalPaddingNeeded := maxTextLineWidth - rowWidth
		rightPaddingForThisRow := totalPaddingNeeded - leftPadding
		if rightPaddingForThisRow > 0 {
			alignedRow += strings.Repeat(" ", rightPaddingForThisRow)
		}

		alignedRows[i] = alignedRow
	}

	return alignedRows
}

// ANSI escape sequence regex for accurate stripping
var ansiRegex = regexp.MustCompile(`\x1b\[[0-9;]*m`)

// stripANSI removes ANSI escape sequences for accurate width calculation
func stripANSI(s string) string {
	return ansiRegex.ReplaceAllString(s, "")
}

// stripEmptyLines removes empty lines from both the top and bottom of rendered text
// This ensures consistent spacing behavior regardless of whether characters have descenders
func stripEmptyLines(lines []string) []string {
	if len(lines) == 0 {
		return lines
	}

	// Find first non-empty line from top
	start := 0
	for start < len(lines) {
		// Check if line contains any non-space characters
		trimmed := strings.TrimSpace(lines[start])
		if trimmed != "" {
			break
		}
		start++
	}

	// If no non-empty lines found, return empty slice
	if start >= len(lines) {
		return []string{}
	}

	// Find last non-empty line from bottom
	end := len(lines) - 1
	for end > start {
		// Check if line contains any non-space characters
		trimmed := strings.TrimSpace(lines[end])
		if trimmed != "" {
			break
		}
		end--
	}

	// Return the slice without empty lines at top or bottom
	return lines[start : end+1]
}

// renderTextWithFont renders text using the specified font with proven rendering logic
func renderTextWithFont(text string, fontData FontData, baseCharSpacing int, wordSpacing float64, scaleFactor float64) []string {
	if text == "" {
		return []string{}
	}

	// Analyze descender properties for this font and scale
	descenderInfo := analyzeDescenderProperties(fontData, scaleFactor)

	// Find max character height (accounting for scaling and descender adjustments)
	maxCharHeight := 0
	for _, bitmapLines := range fontData.Characters {
		scaledLines := scaleCharacter(bitmapLines, scaleFactor)
		maxCharHeight = max(maxCharHeight, len(scaledLines))
	}

	// Also consider the height needed for proper descender alignment
	for _, info := range descenderInfo {
		requiredHeight := info.TotalHeight + info.VerticalOffset
		maxCharHeight = max(maxCharHeight, requiredHeight)
	}

	if maxCharHeight == 0 {
		return []string{"Font has no character data"}
	}

	// Determine a default width for missing characters or the 'space' character
	defaultMissingCharWidth := 4
	if spaceBitmapLines, ok := fontData.Characters[" "]; ok && len(spaceBitmapLines) > 0 {
		defaultMissingCharWidth = utf8.RuneCountInString(spaceBitmapLines[0])
	} else {
		// Fallback to average width of common characters if space is not defined
		for _, char := range "xM!" {
			if charBitmapLines, found := fontData.Characters[string(char)]; found && len(charBitmapLines) > 0 {
				defaultMissingCharWidth = utf8.RuneCountInString(charBitmapLines[0])
				break
			}
		}
	}

	// Pre-calculate character visual bounding box widths and cache kerning values.
	charWidths := make(map[string]int)
	charHeights := make(map[string]int)          // Track individual character heights
	charOffsets := make(map[string]int)          // Track vertical offset for proper descender alignment
	adjustedBitmaps := make(map[string][]string) // Cache adjusted character bitmaps
	kerningCache := make(map[[2]string]int)

	runes := []rune(text)
	for i, r := range runes {
		charStr := string(r)
		if _, exists := charWidths[charStr]; !exists {
			if charStr == " " {
				// Handle manual space character as half-pixel (0.5 pixels)
				manualSpaceWidth := 0.5
				charWidths[charStr] = int(math.Ceil(manualSpaceWidth)) // Use ceiling for width calculation (results in 1)
				charHeights[charStr] = maxCharHeight                   // Use max height for consistent synchronization
				charOffsets[charStr] = 0
				adjustedBitmaps[charStr] = []string{strings.Repeat(" ", int(math.Ceil(manualSpaceWidth)))}
			} else if bitmapLines, ok := fontData.Characters[charStr]; ok {
				// Apply scaling to the bitmap
				scaledBitmapLines := scaleCharacter(bitmapLines, scaleFactor)

				// Filter out empty bitmap rows that can cause synchronization issues
				filteredBitmapLines := make([]string, 0, len(scaledBitmapLines))
				for _, line := range scaledBitmapLines {
					if line != "" {
						filteredBitmapLines = append(filteredBitmapLines, line)
					}
				}

				// If all lines were empty (like problematic space definitions), treat as space
				if len(filteredBitmapLines) == 0 {
					manualSpaceWidth := 0.5
					charWidths[charStr] = int(math.Ceil(manualSpaceWidth))
					charHeights[charStr] = maxCharHeight
					charOffsets[charStr] = 0
					adjustedBitmaps[charStr] = []string{strings.Repeat(" ", int(math.Ceil(manualSpaceWidth)))}
				} else {
					// Use descender information for proper alignment
					if info, hasDescenderInfo := descenderInfo[charStr]; hasDescenderInfo {
						// Adjust character bitmap for proper descender alignment
						adjustedBitmap := adjustCharacterForDescenders(filteredBitmapLines, info, maxCharHeight)
						adjustedBitmaps[charStr] = adjustedBitmap
						charWidths[charStr] = maxRowLen(adjustedBitmap)
						charHeights[charStr] = len(adjustedBitmap)
						charOffsets[charStr] = 0 // Offset is already applied in adjustedBitmap
					} else {
						// Fallback to original logic for characters without descender info
						charWidths[charStr] = maxRowLen(filteredBitmapLines)
						charHeights[charStr] = len(filteredBitmapLines)
						adjustedBitmaps[charStr] = filteredBitmapLines

						// Calculate vertical offset to center characters that are shorter than max height
						if len(filteredBitmapLines) < maxCharHeight {
							charOffsets[charStr] = (maxCharHeight - len(filteredBitmapLines)) / 2
						} else {
							charOffsets[charStr] = 0
						}
					}
				}
			} else {
				charWidths[charStr] = defaultMissingCharWidth
				charHeights[charStr] = maxCharHeight // Use max height for missing chars
				charOffsets[charStr] = 0
				adjustedBitmaps[charStr] = []string{strings.Repeat(" ", defaultMissingCharWidth)}
			}
		}

		// Pre-calculate kerning for pairs
		if i < len(runes)-1 {
			nextCharStr := string(runes[i+1])
			pair := [2]string{charStr, nextCharStr}
			if _, exists := kerningCache[pair]; !exists {
				if charStr == " " || nextCharStr == " " {
					kerningCache[pair] = 0
				} else {
					// Use adjusted bitmaps for kerning calculation to account for descender alignment
					leftBitmap, leftExists := adjustedBitmaps[charStr]
					rightBitmap, rightExists := adjustedBitmaps[nextCharStr]

					if !leftExists || !rightExists {
						kerningCache[pair] = 0
					} else {
						kerningCache[pair] = computeKerning(leftBitmap, rightBitmap)
					}
				}
			}
		}
	}

	var result []string

	// Render the text row by row
	for i := range maxCharHeight {
		lineRunes := make([]rune, 0)
		charStartPositions := make([]float64, len(runes)) // Use float64 for half-pixel precision

		if len(runes) > 0 {
			charStartPositions[0] = 0
		}

		// First pass: Calculate the absolute starting X-position for each character
		for idx := range runes {
			charStr := string(runes[idx])
			if idx > 0 {
				prevCharStr := string(runes[idx-1])
				var prevCharTotalAdvance float64 // Use float64 for half-pixel precision

				if prevCharStr == " " {
					// Determine if this space is a word boundary or character-level spacing
					isWordBoundary := isSpaceAtWordBoundary(runes, idx-1)
					if isWordBoundary {
						prevCharTotalAdvance = 0.5 + wordSpacing // Word boundary space gets word spacing
					} else {
						prevCharTotalAdvance = 0.5 // Character-level space remains half-pixel only
					}
				} else {
					prevCharTotalAdvance = float64(charWidths[prevCharStr])
					optimalInterCharSpacing := kerningCache[[2]string{prevCharStr, charStr}]

					// Handle half-pixel spacing for odd/even height differences
					heightDiff := charHeights[prevCharStr] - charHeights[charStr]
					halfPixelAdjustment := 0.0

					// If heights differ by an odd number, adjust by half a pixel
					if heightDiff%2 != 0 && i >= charHeights[charStr] {
						halfPixelAdjustment = 0.5
					}

					if baseCharSpacing == 0 {
						prevCharTotalAdvance += float64(optimalInterCharSpacing) + halfPixelAdjustment
					} else {
						prevCharTotalAdvance += float64(optimalInterCharSpacing) + float64(baseCharSpacing) + halfPixelAdjustment
					}
				}

				charStartPositions[idx] = charStartPositions[idx-1] + prevCharTotalAdvance
			}
		}

		// Second pass: Place each character's fragment onto the lineRunes canvas
		cumulativeError := 0.0 // Track cumulative rounding errors
		for idx := range runes {
			charStr := string(runes[idx])
			currentXOffset := charStartPositions[idx] + cumulativeError
			fragment := ""

			if charStr == " " {
				// Calculate the actual space width based on whether it's a word boundary
				isWordBoundary := isSpaceAtWordBoundary(runes, idx)
				var spaceWidth float64
				if isWordBoundary {
					spaceWidth = 0.5 + wordSpacing
				} else {
					spaceWidth = 0.5
				}
				fragment = strings.Repeat(" ", int(math.Ceil(spaceWidth)))
			} else if adjustedBitmap, ok := adjustedBitmaps[charStr]; ok {
				// Use the pre-calculated adjusted bitmap that already accounts for descender alignment
				if i >= 0 && i < len(adjustedBitmap) {
					fragment = adjustedBitmap[i]
				} else {
					fragment = ""
				}
			} else {
				fragment = strings.Repeat(" ", charWidths[charStr])
			}

			// Convert float64 position to integer for rendering with proper rounding
			renderXOffset := int(math.Round(currentXOffset))

			// Update cumulative error for next character
			cumulativeError += currentXOffset - float64(renderXOffset)

			// Ensure lineRunes has enough capacity for main text
			requiredLength := renderXOffset + utf8.RuneCountInString(fragment)
			for len(lineRunes) < requiredLength {
				lineRunes = append(lineRunes, ' ')
			}

			// Place the fragment into lineRunes at the calculated position
			fragmentRunes := []rune(fragment)
			for fragIdx, fragRune := range fragmentRunes {
				targetPos := renderXOffset + fragIdx
				if targetPos >= 0 && targetPos < len(lineRunes) {
					// Place character, preserving original proven logic
					if fragRune != ' ' || lineRunes[targetPos] == ' ' {
						lineRunes[targetPos] = fragRune
					}
				}
			}
		}

		// Output the line, trimming any trailing spaces
		resultLine := string(lineRunes)
		result = append(result, strings.TrimRight(resultLine, " "))
	}

	return result
}

// isSpaceAtWordBoundary determines if a space character is at a word boundary
// versus being used for character-level spacing.
// A space is a word boundary only if it separates multi-character words.
// Single characters separated by spaces are treated as character-level spacing.
func isSpaceAtWordBoundary(runes []rune, spaceIndex int) bool {
	if spaceIndex < 0 || spaceIndex >= len(runes) || runes[spaceIndex] != ' ' {
		return false
	}

	// Get the word before and after this space (across any number of spaces)
	wordBefore := getWordBeforeSpaceSequence(runes, spaceIndex)
	wordAfter := getWordAfterSpaceSequence(runes, spaceIndex)

	// Only apply word spacing if BOTH words are multi-character
	// This applies regardless of how many spaces are between them
	return len(wordBefore) > 1 && len(wordAfter) > 1
}

// getWordBeforeSpaceSequence extracts the word before a space, skipping over any preceding spaces
func getWordBeforeSpaceSequence(runes []rune, spaceIndex int) string {
	// First, skip backwards over any spaces to find the start of the preceding word
	i := spaceIndex - 1
	for i >= 0 && runes[i] == ' ' {
		i--
	}

	// Now extract the word (non-space characters)
	var word []rune
	for i >= 0 && runes[i] != ' ' {
		word = append([]rune{runes[i]}, word...)
		i--
	}
	return string(word)
}

// getWordAfterSpaceSequence extracts the word after a space, skipping over any following spaces
func getWordAfterSpaceSequence(runes []rune, spaceIndex int) string {
	// First, skip forwards over any spaces to find the start of the next word
	i := spaceIndex + 1
	for i < len(runes) && runes[i] == ' ' {
		i++
	}

	// Now extract the word (non-space characters)
	var word []rune
	for i < len(runes) && runes[i] != ' ' {
		word = append(word, runes[i])
		i++
	}
	return string(word)
}

// hexToRGB converts a hex color string to RGB values (more robustly)
func hexToRGB(hex string) (int, int, int) {
	if hex == "" {
		return 0, 0, 0
	}
	if hex[0] == '#' {
		hex = hex[1:]
	}

	if len(hex) != 6 {
		return 0, 0, 0
	}

	r := clamp(hexCharToInt(hex[0])*16+hexCharToInt(hex[1]), 0, 255)
	g := clamp(hexCharToInt(hex[2])*16+hexCharToInt(hex[3]), 0, 255)
	b := clamp(hexCharToInt(hex[4])*16+hexCharToInt(hex[5]), 0, 255)

	return r, g, b
}

// rgbToHex converts RGB values to a hex color string
func rgbToHex(r, g, b int) string {
	return fmt.Sprintf("#%02X%02X%02X", r, g, b)
}

// clamp ensures a value is within a specified range
func clamp(value, minVal, maxVal int) int {
	return max(minVal, min(value, maxVal))
}

// hexCharToInt converts a single hex character to its integer value
func hexCharToInt(c byte) int {
	switch {
	case '0' <= c && c <= '9':
		return int(c - '0')
	case 'a' <= c && c <= 'f':
		return int(c - 'a' + 10)
	case 'A' <= c && c <= 'F':
		return int(c - 'A' + 10)
	}
	return 0
}
