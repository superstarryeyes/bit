package ui

import (
	"math"
	"math/rand"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Random number helpers for background effects
func randomFloat() float64 {
	return rand.Float64()
}

func randomInt(n int) int {
	return rand.Intn(n)
}

// Neon color palette for backgrounds
var neonColors = []lipgloss.Color{
	lipgloss.Color("201"), // Magenta
	lipgloss.Color("51"),  // Cyan
	lipgloss.Color("226"), // Yellow
	lipgloss.Color("129"), // Purple
	lipgloss.Color("46"),  // Green
	lipgloss.Color("214"), // Orange
	lipgloss.Color("196"), // Red
}

// Rainbow colors for cycling animation
var rainbowColors = []lipgloss.Color{
	lipgloss.Color("196"), // Red
	lipgloss.Color("214"), // Orange
	lipgloss.Color("226"), // Yellow
	lipgloss.Color("46"),  // Green
	lipgloss.Color("51"),  // Cyan
	lipgloss.Color("39"),  // Blue
	lipgloss.Color("201"), // Magenta
}

// NewLavaLamp creates a new lava lamp effect with floating blobs
func NewLavaLamp(width, height int) *LavaLamp {
	blobs := []Blob{
		{
			X:          float64(width) / 4,
			Y:          float64(height) / 3,
			VX:         0.3,
			VY:         0.2,
			Radius:     8,
			ColorIndex: 0,
		},
		{
			X:          float64(width) * 3 / 4,
			Y:          float64(height) / 2,
			VX:         -0.25,
			VY:         0.15,
			Radius:     10,
			ColorIndex: 1,
		},
		{
			X:          float64(width) / 2,
			Y:          float64(height) * 2 / 3,
			VX:         0.2,
			VY:         -0.3,
			Radius:     7,
			ColorIndex: 2,
		},
		{
			X:          float64(width) / 3,
			Y:          float64(height) / 4,
			VX:         -0.15,
			VY:         0.25,
			Radius:     9,
			ColorIndex: 3,
		},
	}

	return &LavaLamp{
		Blobs:  blobs,
		Width:  width,
		Height: height,
		Frame:  0,
	}
}

// UpdateLavaLamp moves the blobs and adds organic wobble
func UpdateLavaLamp(l *LavaLamp) {
	l.Frame++

	for i := range l.Blobs {
		blob := &l.Blobs[i]

		// Update position
		blob.X += blob.VX
		blob.Y += blob.VY

		// Bounce off edges
		if blob.X < 0 || blob.X > float64(l.Width) {
			blob.VX = -blob.VX
		}
		if blob.Y < 0 || blob.Y > float64(l.Height) {
			blob.VY = -blob.VY
		}

		// Add organic wobble
		wobbleX := math.Sin(float64(l.Frame+i*37)/30.0) * 0.05
		wobbleY := math.Cos(float64(l.Frame+i*41)/25.0) * 0.05
		blob.VX += wobbleX
		blob.VY += wobbleY

		// Damping to prevent too much speed
		blob.VX *= 0.99
		blob.VY *= 0.99
	}
}

// RenderLavaLamp generates the metaball effect with gradient characters
func RenderLavaLamp(l *LavaLamp) []string {
	// Create field map
	field := make([][]float64, l.Height)
	colorMap := make([][]int, l.Height)

	for y := 0; y < l.Height; y++ {
		field[y] = make([]float64, l.Width)
		colorMap[y] = make([]int, l.Width)

		for x := 0; x < l.Width; x++ {
			// Calculate field strength from all blobs
			totalField := 0.0
			closestBlob := 0
			maxField := 0.0

			for i, blob := range l.Blobs {
				dx := float64(x) - blob.X
				dy := float64(y) - blob.Y
				distance := math.Sqrt(dx*dx + dy*dy)

				// Metaball field: 1/distance^2 * radius^2
				if distance > 0 {
					blobField := (blob.Radius * blob.Radius) / (distance * distance)
					totalField += blobField

					// Track which blob contributes most (for color)
					if blobField > maxField {
						maxField = blobField
						closestBlob = i
					}
				}
			}

			field[y][x] = totalField
			colorMap[y][x] = closestBlob
		}
	}

	// Render with gradient characters
	gradientChars := []string{" ", "░", "▒", "▓", "█"}
	lines := make([]string, l.Height)

	for y := 0; y < l.Height; y++ {
		var b strings.Builder

		for x := 0; x < l.Width; x++ {
			strength := field[y][x]
			blobIndex := colorMap[y][x]

			// Choose character based on field strength
			var char string
			if strength < 0.3 {
				char = " " // Empty space
			} else if strength < 0.8 {
				char = gradientChars[1] // ░
			} else if strength < 1.5 {
				char = gradientChars[2] // ▒
			} else if strength < 2.5 {
				char = gradientChars[3] // ▓
			} else {
				char = gradientChars[4] // █
			}

			// Color from closest blob
			color := neonColors[blobIndex%len(neonColors)]

			styled := lipgloss.NewStyle().
				Foreground(color).
				Render(char)

			b.WriteString(styled)
		}

		lines[y] = b.String()
	}

	return lines
}

// NewWavyGrid creates a new wavy grid background
func NewWavyGrid(width, height int) *WavyGrid {
	return &WavyGrid{
		Width:    width,
		Height:   height,
		Frame:    0,
		GridSize: DefaultGridSize,
	}
}

// UpdateWavyGrid advances the animation frame
func UpdateWavyGrid(g *WavyGrid) {
	g.Frame++
}

// RenderWavyGrid generates the wavy grid as strings
func RenderWavyGrid(g *WavyGrid) []string {
	lines := make([]string, g.Height)

	for y := 0; y < g.Height; y++ {
		var b strings.Builder

		for x := 0; x < g.Width; x++ {
			// Calculate wave offset using sine
			waveX := math.Sin(float64(y)/5.0+float64(g.Frame)/20.0) * 2
			waveY := math.Sin(float64(x)/5.0+float64(g.Frame)/20.0) * 2

			// Determine if this position should be a grid line
			gridX := int(float64(x) + waveX)
			gridY := int(float64(y) + waveY)

			isGridLine := (gridX%g.GridSize == 0) || (gridY%g.GridSize == 0)

			var char string
			var color lipgloss.Color

			if isGridLine {
				// Grid line
				if gridX%g.GridSize == 0 && gridY%g.GridSize == 0 {
					char = "+"                          // Intersection
					color = lipgloss.Color("129") // Purple
				} else if gridX%g.GridSize == 0 {
					char = "│"                         // Vertical line
					color = lipgloss.Color("61") // Dark purple
				} else {
					char = "─"                         // Horizontal line
					color = lipgloss.Color("61")
				}
			} else {
				char = " "
				color = lipgloss.Color("0")
			}

			styled := lipgloss.NewStyle().
				Foreground(color).
				Render(char)

			b.WriteString(styled)
		}

		lines[y] = b.String()
	}

	return lines
}

// NewTicker creates a new ticker/sidescroller effect
func NewTicker(text string) *Ticker {
	return &Ticker{
		Text:   text,
		Offset: 0,
		Speed:  TickerSpeed,
	}
}

// UpdateTicker advances the scroll position
func UpdateTicker(t *Ticker, frame int) {
	if frame%t.Speed == 0 {
		t.Offset++
	}
}

// RenderTicker generates the ticker as repeated scrolling text
func RenderTicker(t *Ticker, width, height int) []string {
	lines := make([]string, height)

	// Create a repeating ticker text
	tickerText := " " + t.Text + " "
	repeatCount := (width / len(tickerText)) + 2

	fullText := strings.Repeat(tickerText, repeatCount)
	offset := t.Offset % len(tickerText)

	for y := 0; y < height; y++ {
		var b strings.Builder

		for x := 0; x < width; x++ {
			idx := (x + offset) % len(fullText)
			char := string(fullText[idx])

			// Rainbow color based on position
			colorIdx := (x + y) % len(rainbowColors)
			color := rainbowColors[colorIdx]

			styled := lipgloss.NewStyle().
				Foreground(color).
				Render(char)

			b.WriteString(styled)
		}

		lines[y] = b.String()
	}

	return lines
}

// NewStarfield creates a new 3D starfield effect with flying stars/icons
func NewStarfield(width, height int) *Starfield {
	sf := &Starfield{
		Width:  width,
		Height: height,
		Stars:  make([]Star, StarCount),
		Frame:  0,
	}

	// Star icons with weighted selection (favoring ASCII characters)
	icons := []string{
		"█", "▓", "▒", "░", "●", "◆", "■", "▲", // Regular ASCII icons (common)
		"✦", "✧", "✨", "⋆", "★", "☆", // Star symbols (less common)
		"•", "∘", "·", // Dots (medium)
	}

	// Initialize stars at random positions in 3D space
	for i := range sf.Stars {
		// Weight icon selection (80% ASCII, 20% fancy)
		iconIndex := 0
		roll := randomFloat()
		if roll < 0.80 {
			// 80% chance: Regular ASCII icons
			iconIndex = randomInt(8)
		} else {
			// 20% chance: Star symbols and dots
			iconIndex = 8 + randomInt(len(icons)-8)
		}

		sf.Stars[i] = Star{
			X:    (randomFloat() - 0.5) * 100,
			Y:    (randomFloat() - 0.5) * 100,
			Z:    randomFloat() * 100,
			Icon: icons[iconIndex],
		}
	}

	return sf
}

// UpdateStarfield moves stars toward the viewer (3D perspective effect)
func UpdateStarfield(sf *Starfield) {
	sf.Frame++

	for i := range sf.Stars {
		star := &sf.Stars[i]

		// Move star forward (toward viewer)
		star.Z -= StarSpeed

		// Reset star if it passes the viewer
		if star.Z <= 0 {
			star.X = (randomFloat() - 0.5) * 100
			star.Y = (randomFloat() - 0.5) * 100
			star.Z = 100
		}
	}
}

// RenderStarfield generates the 3D starfield with perspective projection
func RenderStarfield(sf *Starfield) []string {
	// Create canvas
	canvas := make([][]string, sf.Height)
	for y := 0; y < sf.Height; y++ {
		canvas[y] = make([]string, sf.Width)
		for x := 0; x < sf.Width; x++ {
			canvas[y][x] = " "
		}
	}

	// Project stars onto 2D screen using perspective
	for _, star := range sf.Stars {
		// Perspective projection: scale = focal_distance / z
		scale := 100 / star.Z
		screenX := int(star.X*scale) + sf.Width/2
		screenY := int(star.Y*scale) + sf.Height/2

		// Draw star if on screen
		if screenX >= 0 && screenX < sf.Width && screenY >= 0 && screenY < sf.Height {
			// Calculate depth (0 = far, 1 = close)
			depth := 1.0 - (star.Z / 100)

			var char string
			var color lipgloss.Color

			if depth > 0.7 {
				// Close - show full icon with bright color
				char = star.Icon
				color = neonColors[int(depth*10)%len(neonColors)]
			} else if depth > 0.4 {
				// Medium distance - show bullet
				char = "•"
				color = lipgloss.Color("75") // Medium blue
			} else {
				// Far - show small dot
				char = "·"
				color = lipgloss.Color("60") // Dark blue
			}

			styled := lipgloss.NewStyle().
				Foreground(color).
				Render(char)

			canvas[screenY][screenX] = styled
		}
	}

	// Convert canvas to string lines
	lines := make([]string, sf.Height)
	for y := 0; y < sf.Height; y++ {
		var b strings.Builder
		for x := 0; x < sf.Width; x++ {
			b.WriteString(canvas[y][x])
		}
		lines[y] = b.String()
	}

	return lines
}

// CompositeBackground overlays rendered text on top of a background
func CompositeBackground(background []string, textLines []string, textX, textY, width, height int) []string {
	result := make([]string, height)

	// Initialize with background or empty space
	for y := 0; y < height; y++ {
		if y < len(background) {
			result[y] = background[y]
		} else {
			result[y] = strings.Repeat(" ", width)
		}
	}

	// Overlay text
	for i, line := range textLines {
		y := textY + i
		if y >= 0 && y < height {
			result[y] = overlayString(result[y], line, textX, width)
		}
	}

	return result
}

// overlayString overlays src onto dst at position x, preserving styled strings
func overlayString(dst, src string, x, maxWidth int) string {
	if x < 0 {
		return dst
	}

	// Get visual widths (ignoring ANSI codes)
	dstWidth := lipgloss.Width(dst)
	srcWidth := lipgloss.Width(src)

	// Build result by extracting visible portions
	var result strings.Builder

	// Left side: extract x visible characters from dst
	if x > 0 {
		leftPart := extractVisibleChars(dst, 0, x)
		result.WriteString(leftPart)
	}

	// Middle: the overlay
	result.WriteString(src)

	// Right side: extract remaining characters from dst after the overlay
	rightStart := x + srcWidth
	if rightStart < dstWidth {
		rightPart := extractVisibleChars(dst, rightStart, dstWidth-rightStart)
		result.WriteString(rightPart)
	}

	return result.String()
}

// extractVisibleChars extracts count visible characters starting from position start
// Handles ANSI escape codes properly
func extractVisibleChars(s string, start, count int) string {
	if count <= 0 {
		return ""
	}

	runes := []rune(s)
	visibleCount := 0
	inEscape := false
	startIdx := -1
	endIdx := -1

	for i, r := range runes {
		// Track ANSI escape sequences
		if r == '\x1b' {
			inEscape = true
		} else if inEscape && r == 'm' {
			inEscape = false
			continue
		}

		// Only count visible characters
		if !inEscape {
			if visibleCount == start && startIdx == -1 {
				startIdx = i
			}
			if visibleCount >= start {
				visibleCount++
				if visibleCount-start >= count {
					endIdx = i + 1
					break
				}
			} else {
				visibleCount++
			}
		}
	}

	if startIdx == -1 {
		return ""
	}
	if endIdx == -1 {
		endIdx = len(runes)
	}

	return string(runes[startIdx:endIdx])
}
