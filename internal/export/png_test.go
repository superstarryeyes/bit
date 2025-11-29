// ABOUTME: Tests for PNG generation from ANSI-colored text output.
// ABOUTME: Verifies correct parsing of ANSI codes and pixel rendering at 16x scale.

package export

import (
	"bytes"
	"image"
	"image/png"
	"testing"
)

func TestGeneratePNG_EmptyInput(t *testing.T) {
	_, err := GeneratePNG([]string{}, DefaultPNGOptions())
	if err == nil {
		t.Error("expected error for empty input, got nil")
	}
}

func TestGeneratePNG_NilInput(t *testing.T) {
	_, err := GeneratePNG(nil, DefaultPNGOptions())
	if err == nil {
		t.Error("expected error for nil input, got nil")
	}
}

func TestGeneratePNG_OnlySpaces(t *testing.T) {
	// A line with only spaces should produce a valid (but transparent) PNG
	lines := []string{"   "}
	data, err := GeneratePNG(lines, DefaultPNGOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should still produce valid PNG
	if !bytes.HasPrefix(data, []byte{0x89, 'P', 'N', 'G'}) {
		t.Error("invalid PNG header")
	}
}

func TestGeneratePNG_ValidPNGHeader(t *testing.T) {
	// Simple red full block character
	lines := []string{"\x1b[38;2;255;0;0m█\x1b[0m"}

	data, err := GeneratePNG(lines, DefaultPNGOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify PNG magic bytes
	if !bytes.HasPrefix(data, []byte{0x89, 'P', 'N', 'G'}) {
		t.Error("invalid PNG header - missing PNG magic bytes")
	}
}

func TestGeneratePNG_SingleCharDimensions(t *testing.T) {
	// One character should produce 16x16 pixels
	lines := []string{"\x1b[38;2;255;0;0m█\x1b[0m"}

	data, err := GeneratePNG(lines, DefaultPNGOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	img, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("failed to decode PNG: %v", err)
	}

	bounds := img.Bounds()
	if bounds.Dx() != CellSize || bounds.Dy() != CellSize {
		t.Errorf("expected %dx%d, got %dx%d", CellSize, CellSize, bounds.Dx(), bounds.Dy())
	}
}

func TestGeneratePNG_MultipleCharDimensions(t *testing.T) {
	// Three characters on one line should produce 48x16 pixels
	lines := []string{"\x1b[38;2;255;0;0m███\x1b[0m"}

	data, err := GeneratePNG(lines, DefaultPNGOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	img, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("failed to decode PNG: %v", err)
	}

	bounds := img.Bounds()
	expectedWidth := 3 * CellSize
	if bounds.Dx() != expectedWidth || bounds.Dy() != CellSize {
		t.Errorf("expected %dx%d, got %dx%d", expectedWidth, CellSize, bounds.Dx(), bounds.Dy())
	}
}

func TestGeneratePNG_MultipleLinesimensions(t *testing.T) {
	// Two lines of one character each should produce 16x32 pixels
	lines := []string{
		"\x1b[38;2;255;0;0m█\x1b[0m",
		"\x1b[38;2;0;255;0m█\x1b[0m",
	}

	data, err := GeneratePNG(lines, DefaultPNGOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	img, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("failed to decode PNG: %v", err)
	}

	bounds := img.Bounds()
	expectedHeight := 2 * CellSize
	if bounds.Dx() != CellSize || bounds.Dy() != expectedHeight {
		t.Errorf("expected %dx%d, got %dx%d", CellSize, expectedHeight, bounds.Dx(), bounds.Dy())
	}
}

func TestGeneratePNG_TransparentBackground(t *testing.T) {
	// Space character should result in transparent pixel
	lines := []string{" "}

	data, err := GeneratePNG(lines, DefaultPNGOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	img, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("failed to decode PNG: %v", err)
	}

	// Check center pixel is transparent
	rgba, ok := img.(*image.NRGBA)
	if !ok {
		// Try RGBA
		rgbaImg, ok := img.(*image.RGBA)
		if !ok {
			t.Fatalf("unexpected image type: %T", img)
		}
		pixel := rgbaImg.RGBAAt(CellSize/2, CellSize/2)
		if pixel.A != 0 {
			t.Errorf("expected transparent pixel (A=0), got A=%d", pixel.A)
		}
		return
	}
	pixel := rgba.NRGBAAt(CellSize/2, CellSize/2)
	if pixel.A != 0 {
		t.Errorf("expected transparent pixel (A=0), got A=%d", pixel.A)
	}
}

func TestGeneratePNG_ColorParsing(t *testing.T) {
	// Red full block
	lines := []string{"\x1b[38;2;255;0;0m█\x1b[0m"}

	data, err := GeneratePNG(lines, DefaultPNGOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	img, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("failed to decode PNG: %v", err)
	}

	// Check center pixel is red
	r, g, b, a := img.At(CellSize/2, CellSize/2).RGBA()
	// RGBA returns values in [0, 65535] range
	r8, g8, b8, a8 := uint8(r>>8), uint8(g>>8), uint8(b>>8), uint8(a>>8)

	if r8 != 255 || g8 != 0 || b8 != 0 {
		t.Errorf("expected red pixel (255,0,0), got (%d,%d,%d)", r8, g8, b8)
	}
	if a8 != 255 {
		t.Errorf("expected opaque pixel (A=255), got A=%d", a8)
	}
}

func TestGeneratePNG_UpperHalfBlock(t *testing.T) {
	// Upper half block should only fill top 8 rows
	lines := []string{"\x1b[38;2;0;255;0m▀\x1b[0m"}

	data, err := GeneratePNG(lines, DefaultPNGOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	img, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("failed to decode PNG: %v", err)
	}

	// Top half should be green (row 4)
	_, _, _, topA := img.At(CellSize/2, CellSize/4).RGBA()
	if uint8(topA>>8) == 0 {
		t.Error("expected top half to be opaque, got transparent")
	}

	// Bottom half should be transparent (row 12)
	_, _, _, bottomA := img.At(CellSize/2, CellSize*3/4).RGBA()
	if uint8(bottomA>>8) != 0 {
		t.Errorf("expected bottom half to be transparent, got A=%d", uint8(bottomA>>8))
	}
}

func TestGeneratePNG_LowerHalfBlock(t *testing.T) {
	// Lower half block should only fill bottom 8 rows
	lines := []string{"\x1b[38;2;0;0;255m▄\x1b[0m"}

	data, err := GeneratePNG(lines, DefaultPNGOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	img, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("failed to decode PNG: %v", err)
	}

	// Top half should be transparent (row 4)
	_, _, _, topA := img.At(CellSize/2, CellSize/4).RGBA()
	if uint8(topA>>8) != 0 {
		t.Errorf("expected top half to be transparent, got A=%d", uint8(topA>>8))
	}

	// Bottom half should be blue (row 12)
	_, _, _, bottomA := img.At(CellSize/2, CellSize*3/4).RGBA()
	if uint8(bottomA>>8) == 0 {
		t.Error("expected bottom half to be opaque, got transparent")
	}
}

func TestGeneratePNG_ShadeCharacterLightShade(t *testing.T) {
	// Light shade (░) should have low alpha
	lines := []string{"\x1b[38;2;255;255;255m░\x1b[0m"}

	data, err := GeneratePNG(lines, DefaultPNGOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	img, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("failed to decode PNG: %v", err)
	}

	_, _, _, a := img.At(CellSize/2, CellSize/2).RGBA()
	alpha := uint8(a >> 8)
	// Light shade should have alpha around 64 (25%)
	if alpha < 32 || alpha > 96 {
		t.Errorf("expected light shade alpha ~64, got %d", alpha)
	}
}

func TestGeneratePNG_CustomCellSize(t *testing.T) {
	lines := []string{"\x1b[38;2;255;0;0m█\x1b[0m"}

	opts := PNGOptions{
		CellWidth:  8,
		CellHeight: 8,
	}
	data, err := GeneratePNG(lines, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	img, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("failed to decode PNG: %v", err)
	}

	bounds := img.Bounds()
	if bounds.Dx() != 8 || bounds.Dy() != 8 {
		t.Errorf("expected 8x8 with custom options, got %dx%d", bounds.Dx(), bounds.Dy())
	}
}

func TestDefaultPNGOptions(t *testing.T) {
	opts := DefaultPNGOptions()

	if opts.CellWidth != CellSize {
		t.Errorf("expected default CellWidth=%d, got %d", CellSize, opts.CellWidth)
	}
	if opts.CellHeight != CellSize {
		t.Errorf("expected default CellHeight=%d, got %d", CellSize, opts.CellHeight)
	}
}

func TestTerminalAspectRatioPNGOptions(t *testing.T) {
	opts := TerminalAspectRatioPNGOptions()

	// Terminal characters are ~2:1 height:width, so cells should be twice as tall as wide.
	if opts.CellHeight != opts.CellWidth*2 {
		t.Errorf("expected 2:1 height:width ratio, got width=%d height=%d",
			opts.CellWidth, opts.CellHeight)
	}
}

func TestGeneratePNG_TerminalAspectRatioDimensions(t *testing.T) {
	// Three characters on one line with terminal aspect ratio
	lines := []string{"\x1b[38;2;255;0;0m███\x1b[0m"}

	opts := TerminalAspectRatioPNGOptions()
	data, err := GeneratePNG(lines, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	img, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("failed to decode PNG: %v", err)
	}

	bounds := img.Bounds()
	expectedWidth := 3 * opts.CellWidth
	expectedHeight := opts.CellHeight
	if bounds.Dx() != expectedWidth || bounds.Dy() != expectedHeight {
		t.Errorf("expected %dx%d, got %dx%d", expectedWidth, expectedHeight, bounds.Dx(), bounds.Dy())
	}
}

func TestGeneratePNG_ColorConsistency(t *testing.T) {
	// Test that all characters with the same color code produce uniform color
	// This mimics a single-color (non-gradient) rendering where every character
	// has the same \x1b[38;2;R;G;Bm sequence
	cyan := "\x1b[38;2;0;191;255m" // Deep Sky Blue
	reset := "\x1b[0m"

	// 4 cyan full blocks in a row
	lines := []string{
		cyan + "█" + reset + cyan + "█" + reset + cyan + "█" + reset + cyan + "█" + reset,
	}

	opts := TerminalAspectRatioPNGOptions()
	data, err := GeneratePNG(lines, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	img, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("failed to decode PNG: %v", err)
	}

	// Sample color from center of each character cell
	expectedR, expectedG, expectedB := uint32(0), uint32(191), uint32(255)
	halfW, halfH := opts.CellWidth/2, opts.CellHeight/2

	for charIdx := 0; charIdx < 4; charIdx++ {
		x := charIdx*opts.CellWidth + halfW
		y := halfH
		r, g, b, _ := img.At(x, y).RGBA()
		// RGBA returns 16-bit values, shift to 8-bit
		r8, g8, b8 := r>>8, g>>8, b>>8

		if r8 != expectedR || g8 != expectedG || b8 != expectedB {
			t.Errorf("char %d: expected RGB(%d,%d,%d), got RGB(%d,%d,%d)",
				charIdx, expectedR, expectedG, expectedB, r8, g8, b8)
		}
	}
}

func TestCountVisibleChars(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{"plain text", "hello", 5},
		{"with ANSI color", "\x1b[38;2;255;0;0mhello\x1b[0m", 5},
		{"only ANSI codes", "\x1b[38;2;255;0;0m\x1b[0m", 0},
		{"mixed", "\x1b[38;2;255;0;0ma\x1b[0mb", 2},
		{"empty", "", 0},
		{"unicode blocks", "█▀▄", 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := countVisibleChars(tt.input)
			if got != tt.expected {
				t.Errorf("countVisibleChars(%q) = %d, want %d", tt.input, got, tt.expected)
			}
		})
	}
}
