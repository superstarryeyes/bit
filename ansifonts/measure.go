package ansifonts

import (
	"regexp"
	"unicode/utf8"
)

// ANSI escape sequence regex for accurate stripping
var ansiRegex = regexp.MustCompile(`\x1b\[[0-9;]*m`)

// StripANSI removes ANSI escape sequences for accurate width calculation
func StripANSI(s string) string {
	return ansiRegex.ReplaceAllString(s, "")
}

// MeasureBlock returns the maximum width and total height of a text block.
// Width is calculated using rune count after ANSI sequences are stripped.
func MeasureBlock(lines []string) (w int, h int) {
	h = len(lines)
	for _, line := range lines {
		lineWidth := utf8.RuneCountInString(StripANSI(line))
		if lineWidth > w {
			w = lineWidth
		}
	}
	return w, h
}
