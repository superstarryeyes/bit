package ansifonts

import (
	"strings"
	"unicode/utf8"
)

// maxRowLen is a helper to find the maximum rune length among all rows of a glyph,
// effectively determining its bounding box width.
func maxRowLen(rows []string) int {
	maxLen := 0
	for _, r := range rows {
		maxLen = max(maxLen, utf8.RuneCountInString(r))
	}
	return maxLen
}

// normalizeGlyph pads rows of a glyph to a consistent height with empty strings
// to simplify comparison between glyphs of different heights.
func normalizeGlyph(glyph []string, height int) []string {
	out := make([]string, height)
	maxLen := 0
	for _, r := range glyph {
		maxLen = max(maxLen, utf8.RuneCountInString(r))
	}

	for i := range height {
		if i < len(glyph) {
			// Ensure each row has consistent width by padding with spaces if needed
			row := glyph[i]
			if utf8.RuneCountInString(row) < maxLen {
				row += strings.Repeat(" ", maxLen-utf8.RuneCountInString(row))
			}
			out[i] = row
		} else {
			// Pad with spaces (not empty strings) to maintain consistent width
			out[i] = strings.Repeat(" ", maxLen)
		}
	}
	return out
}

// computeKerning calculates the horizontal offset needed to align glyphB relative to glyphA
// such that the minimum distance between their visible pixels is exactly 1 (touching).
// This allows the renderer to add precise character spacing on top of this baseline.
func computeKerning(glyphA, glyphB []string) int {
	if len(glyphA) == 0 || len(glyphB) == 0 {
		return 0
	}

	// Normalize heights for line-by-line comparison
	h := max(len(glyphA), len(glyphB))
	a := normalizeGlyph(glyphA, h)
	b := normalizeGlyph(glyphB, h)

	widthA := maxRowLen(a)
	if widthA == 0 {
		return 0
	}

	minDist := 1000 // effectively infinity
	hasOverlap := false

	// Global bounds for fallback (handling non-vertically overlapping characters)
	maxAGlobal := -1
	minBGlobal := 1000

	for y := 0; y < h; y++ {
		rowA := []rune(a[y])
		rowB := []rune(b[y])

		// Find rightmost pixel in A on this line
		maxA := -1
		for x := len(rowA) - 1; x >= 0; x-- {
			if rowA[x] != ' ' && rowA[x] != 0 {
				maxA = x
				break
			}
		}

		// Find leftmost pixel in B on this line
		minB := -1
		for x := 0; x < len(rowB); x++ {
			if rowB[x] != ' ' && rowB[x] != 0 {
				minB = x
				break
			}
		}

		// Update global bounds
		if maxA != -1 {
			if maxA > maxAGlobal {
				maxAGlobal = maxA
			}
		}
		if minB != -1 {
			if minB < minBGlobal {
				minBGlobal = minB
			}
		}

		// If both have pixels on this line, calculate the horizontal distance
		if maxA != -1 && minB != -1 {
			// Calculate distance: (Start of B - Start of A)
			// We imagine B starts at WidthA.
			// PosA = maxA
			// PosB = WidthA + minB
			// Distance = PosB - PosA
			dist := (widthA + minB) - maxA
			if dist < minDist {
				minDist = dist
			}
			hasOverlap = true
		}
	}

	// If no lines overlap (e.g., punctuation marks at different heights like ' and .),
	// fall back to the bounding box horizontal distance.
	if !hasOverlap {
		if maxAGlobal != -1 && minBGlobal != 1000 {
			minDist = (widthA + minBGlobal) - maxAGlobal
		} else {
			// One or both glyphs are empty space
			return 0
		}
	}

	// We want the closest pixels to be adjacent (distance 1 pixel).
	// This creates a visual gap of 0.
	// The renderer will add the user's specific spacing on top of this.
	// Adjustment = DesiredDistance(1) - ActualDistance(minDist)
	return 1 - minDist
}
