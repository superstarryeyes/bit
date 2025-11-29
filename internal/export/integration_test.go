// ABOUTME: Integration tests for PNG export functionality.
// ABOUTME: Tests the full pipeline from ANSI input to PNG file output.

//go:build integration
// +build integration

package export

import (
	"bytes"
	"image/png"
	"os"
	"path/filepath"
	"testing"
)

// TestPNGExportIntegration tests the full export pipeline
func TestPNGExportIntegration(t *testing.T) {
	// Sample ANSI-colored output (red "H" in block characters)
	lines := []string{
		"\x1b[38;2;255;0;0m█\x1b[0m  \x1b[38;2;255;0;0m█\x1b[0m",
		"\x1b[38;2;255;0;0m█\x1b[0m  \x1b[38;2;255;0;0m█\x1b[0m",
		"\x1b[38;2;255;0;0m████\x1b[0m",
		"\x1b[38;2;255;0;0m█\x1b[0m  \x1b[38;2;255;0;0m█\x1b[0m",
		"\x1b[38;2;255;0;0m█\x1b[0m  \x1b[38;2;255;0;0m█\x1b[0m",
	}

	// Generate PNG
	data, err := GeneratePNG(lines, DefaultPNGOptions())
	if err != nil {
		t.Fatalf("GeneratePNG failed: %v", err)
	}

	// Verify PNG is valid
	img, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("PNG decode failed: %v", err)
	}

	// Verify dimensions (4 chars wide x 5 lines = 64x80 pixels at 16x scale)
	bounds := img.Bounds()
	expectedWidth := 4 * CellSize  // 64
	expectedHeight := 5 * CellSize // 80
	if bounds.Dx() != expectedWidth || bounds.Dy() != expectedHeight {
		t.Errorf("expected %dx%d, got %dx%d", expectedWidth, expectedHeight, bounds.Dx(), bounds.Dy())
	}

	// Write to temp file for manual inspection
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test_output.png")
	if err := os.WriteFile(tmpFile, data, 0644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	t.Logf("PNG written to: %s (%d bytes)", tmpFile, len(data))

	// Verify file was created and has content
	info, err := os.Stat(tmpFile)
	if err != nil {
		t.Fatalf("failed to stat temp file: %v", err)
	}
	if info.Size() == 0 {
		t.Error("PNG file is empty")
	}
}

// TestExportManagerBinaryExport tests the ExportManager.ExportBinary method
func TestExportManagerBinaryExport(t *testing.T) {
	em := NewExportManager()

	// Verify PNG is recognized as binary
	if !em.IsBinaryFormat("PNG") {
		t.Error("PNG should be identified as binary format")
	}

	// Verify TXT is not binary
	if em.IsBinaryFormat("TXT") {
		t.Error("TXT should not be identified as binary format")
	}
}
