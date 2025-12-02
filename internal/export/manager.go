// ABOUTME: Export manager handles saving rendered ANSI art to various file formats.
// ABOUTME: Supports text formats (TXT, GO, JS, PY, RS, SH) and binary formats (PNG).

package export

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
)

// MaxFilenameLength is the maximum filename length before extension
const MaxFilenameLength = 200

// ansiRegex is compiled once at package level for efficiency
var ansiRegex = regexp.MustCompile(`\x1b\[[0-9;]*m`)

// stripANSI removes ANSI escape sequences from text
func stripANSI(s string) string {
	return ansiRegex.ReplaceAllString(s, "")
}

// ExportManager handles exporting text in various formats
type ExportManager struct {
	formats []ExportFormat
	basePath string // Base directory for exports (defaults to Desktop)
}

// NewExportManager creates a new export manager with supported formats
func NewExportManager() *ExportManager {
	return &ExportManager{
		formats:  SupportedFormats,
		basePath: getDesktopPath(),
	}
}

// getDesktopPath returns the user's Desktop directory path
func getDesktopPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		// Fallback to current directory if we can't get home
		cwd, _ := os.Getwd()
		return cwd
	}
	return filepath.Join(home, "Desktop")
}

// GetSupportedFormats returns the list of supported export formats
func (em *ExportManager) GetSupportedFormats() []ExportFormat {
	return em.formats
}

// Export saves the content to a file in the specified format
func (em *ExportManager) Export(content, filename, formatName string) error {
	// Find the format
	var format *ExportFormat
	for _, f := range em.formats {
		if f.Name == formatName {
			format = &f
			break
		}
	}

	if format == nil {
		return fmt.Errorf("unsupported format: %s", formatName)
	}

	// Sanitize filename to prevent path traversal attacks
	filename = SanitizeFilename(filename)
	if filename == "" {
		return fmt.Errorf("invalid filename")
	}

	// Ensure filename has the correct extension
	if !strings.HasSuffix(filename, format.Extension) {
		filename += format.Extension
	}

	// Create full file path using filepath.Join for safety
	filePath := filepath.Join(em.basePath, filepath.Base(filename))

	// Write content to file
	err := os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("failed to write file: %v", err)
	}

	return nil
}

// ExportBinary saves binary content (like PNG) to a file in the specified format
func (em *ExportManager) ExportBinary(content []byte, filename, formatName string) error {
	// Find the format
	var format *ExportFormat
	for _, f := range em.formats {
		if f.Name == formatName {
			format = &f
			break
		}
	}

	if format == nil {
		return fmt.Errorf("unsupported format: %s", formatName)
	}

	if !format.IsBinary {
		return fmt.Errorf("format %s is not a binary format, use Export() instead", formatName)
	}

	// Sanitize filename to prevent path traversal attacks
	filename = SanitizeFilename(filename)
	if filename == "" {
		return fmt.Errorf("invalid filename")
	}

	// Ensure filename has the correct extension
	if !strings.HasSuffix(filename, format.Extension) {
		filename += format.Extension
	}

	// Create full file path using filepath.Join for safety
	filePath := filepath.Join(em.basePath, filepath.Base(filename))

	// Write binary content to file
	err := os.WriteFile(filePath, content, 0644)
	if err != nil {
		return fmt.Errorf("failed to write file: %v", err)
	}

	return nil
}

// IsBinaryFormat returns true if the format requires binary export
func (em *ExportManager) IsBinaryFormat(name string) bool {
	format := em.GetFormatByName(name)
	if format != nil {
		return format.IsBinary
	}
	return false
}

// CheckFileExists checks if a file already exists at the given path
func (em *ExportManager) CheckFileExists(filename, formatName string) (bool, string, error) {
	// Find the format
	var format *ExportFormat
	for _, f := range em.formats {
		if f.Name == formatName {
			format = &f
			break
		}
	}

	if format == nil {
		return false, "", fmt.Errorf("unsupported format: %s", formatName)
	}

	// Sanitize filename
	filename = SanitizeFilename(filename)
	if filename == "" {
		return false, "", fmt.Errorf("invalid filename")
	}

	// Ensure filename has the correct extension
	if !strings.HasSuffix(filename, format.Extension) {
		filename += format.Extension
	}

	// Create full file path using filepath.Join for safety
	filePath := filepath.Join(em.basePath, filepath.Base(filename))

	// Check if file exists
	if _, err := os.Stat(filePath); err == nil {
		return true, filename, nil
	}

	return false, filename, nil
}

// SanitizeFilename removes dangerous characters from filenames
func SanitizeFilename(filename string) string {
	// Handle empty filename
	if filename == "" {
		return ""
	}

	// Extract base filename to prevent path traversal
	// First normalize path separators
	normalized := strings.ReplaceAll(filename, "\\", "/")
	base := filepath.Base(normalized)

	// Remove path separators and other problematic characters (but keep dots)
	invalidChars := []string{"/", ":", "*", "?", "\"", "<", ">", "|"}
	for _, char := range invalidChars {
		base = strings.ReplaceAll(base, char, "")
	}

	// Trim whitespace from start/end
	base = strings.Trim(base, " ")

	// Check for reserved Windows filenames (case-insensitive)
	// Check the base name before adding the extension
	ext := filepath.Ext(base)
	nameWithoutExt := strings.TrimSuffix(base, ext)

	reservedNames := []string{"CON", "PRN", "AUX", "NUL", "COM1", "COM2", "COM3", "COM4", "COM5",
		"COM6", "COM7", "COM8", "COM9", "LPT1", "LPT2", "LPT3", "LPT4", "LPT5", "LPT6", "LPT7", "LPT8", "LPT9"}
	upperName := strings.ToUpper(nameWithoutExt)
	if slices.Contains(reservedNames, upperName) {
		base = "_" + base
	}

	// Ensure filename is not empty after sanitization
	if base == "" {
		return ""
	}

	// Limit filename length (well under filesystem limit, leaving room for extension)
	// Use rune-based truncation to properly handle UTF-8 characters
	baseRunes := []rune(base)
	if len(baseRunes) > MaxFilenameLength {
		// Preserve extension when truncating if it exists
		ext := filepath.Ext(base)
		extRunes := []rune(ext)

		if len(extRunes) >= MaxFilenameLength {
			// Extreme case: extension is longer than our limit, just truncate the whole thing
			base = string(baseRunes[:MaxFilenameLength])
		} else {
			// Normal case: truncate name part, preserving extension
			maxNameLength := MaxFilenameLength - len(extRunes)
			if maxNameLength <= 0 {
				// Should not happen given the check above, but be defensive
				base = string(baseRunes[:MaxFilenameLength])
			} else {
				nameRunes := []rune(nameWithoutExt)
				if len(nameRunes) > maxNameLength {
					base = string(nameRunes[:maxNameLength]) + ext
				} else {
					// This case means the name part is already within limits
					// but the total is over due to the extension, so we truncate the whole thing
					base = string(baseRunes[:MaxFilenameLength])
				}
			}
		}
	}

	return base
}

// GetFormatByName returns the export format with the given name
func (em *ExportManager) GetFormatByName(name string) *ExportFormat {
	for _, format := range em.formats {
		if format.Name == name {
			return &format
		}
	}
	return nil
}

// GetFormatNames returns a slice of all format names in order
func (em *ExportManager) GetFormatNames() []string {
	names := make([]string, len(em.formats))
	for i, format := range em.formats {
		names[i] = format.Name
	}
	return names
}

// GetFormatDescription returns the description for a given format name
func (em *ExportManager) GetFormatDescription(name string) string {
	format := em.GetFormatByName(name)
	if format != nil {
		return format.Description
	}
	return "Unknown Format"
}

// GetFormatExtension returns the file extension for a given format name
func (em *ExportManager) GetFormatExtension(name string) string {
	format := em.GetFormatByName(name)
	if format != nil {
		return format.Extension
	}
	return ".txt"
}

// GetDefaultFormat returns the default export format (first in the list)
func (em *ExportManager) GetDefaultFormat() string {
	if len(em.formats) > 0 {
		return em.formats[0].Name
	}
	return "TXT"
}

// GenerateTXTCode creates plain text content by stripping ANSI codes
func GenerateTXTCode(lines []string) string {
	// Join rendered lines with newlines and strip ANSI codes
	rawContent := strings.Join(lines, "\n")
	return stripANSI(rawContent)
}

// GenerateGoCode creates Go source code that reproduces the ANSI art
func GenerateGoCode(lines []string) string {
	var builder strings.Builder

	// Write package and imports (minimal imports for standalone version)
	builder.WriteString("package main\n\n")
	builder.WriteString("import (\n")
	builder.WriteString("\t\"fmt\"\n")
	builder.WriteString(")\n\n")

	// Write main function
	builder.WriteString("func main() {\n")

	// Add the rendered lines with embedded ANSI codes
	builder.WriteString("\tlines := []string{\n")
	for _, line := range lines {
		// Escape quotes in the line
		escapedLine := strings.ReplaceAll(line, "\"", "\\\"")
		builder.WriteString(fmt.Sprintf("\t\t\"%s\",\n", escapedLine))
	}
	builder.WriteString("\t}\n\n")

	// Print each line
	builder.WriteString("\tfor _, line := range lines {\n")
	builder.WriteString("\t\tfmt.Println(line)\n")
	builder.WriteString("\t}\n")
	builder.WriteString("}\n")

	return builder.String()
}

// GenerateJSCode creates JavaScript source code that reproduces the ANSI art
func GenerateJSCode(lines []string) string {
	var builder strings.Builder

	// Write file header
	builder.WriteString("/* Generated JavaScript ANSI Art */\n")
	builder.WriteString("\n")

	// Write the array of lines
	builder.WriteString("const ansiArtLines = [\n")
	for _, line := range lines {
		// Escape quotes and backslashes in the line
		escapedLine := strings.ReplaceAll(line, "\\", "\\\\")
		escapedLine = strings.ReplaceAll(escapedLine, "\"", "\\\"")
		builder.WriteString(fmt.Sprintf("  \"%s\",\n", escapedLine))
	}
	builder.WriteString("];\n\n")

	// Write function to display the art
	builder.WriteString("function displayAnsiArt() {\n")
	builder.WriteString("  ansiArtLines.forEach(function(line) {\n")
	builder.WriteString("    console.log(line);\n")
	builder.WriteString("  });\n")
	builder.WriteString("}\n\n")
	builder.WriteString("displayAnsiArt();\n")

	return builder.String()
}

// GeneratePythonCode creates Python source code that reproduces the ANSI art
func GeneratePythonCode(lines []string) string {
	var builder strings.Builder

	// Write file header
	builder.WriteString("# Generated Python ANSI Art\n")
	builder.WriteString("\n")

	// Write the array of lines
	builder.WriteString("ansi_art_lines = [\n")
	for _, line := range lines {
		// Escape quotes and backslashes in the line
		escapedLine := strings.ReplaceAll(line, "\\", "\\\\")
		escapedLine = strings.ReplaceAll(escapedLine, "\"", "\\\"")
		escapedLine = strings.ReplaceAll(escapedLine, "'", "\\'")
		builder.WriteString(fmt.Sprintf("    \"%s\",\n", escapedLine))
	}
	builder.WriteString("]\n\n")

	// Write function to display the art
	builder.WriteString("def display_ansi_art():\n")
	builder.WriteString("    for line in ansi_art_lines:\n")
	builder.WriteString("        print(line)\n\n")
	builder.WriteString("if __name__ == \"__main__\":\n")
	builder.WriteString("    display_ansi_art()\n")

	return builder.String()
}

// GenerateRustCode creates Rust source code that reproduces the ANSI art
func GenerateRustCode(lines []string) string {
	var builder strings.Builder

	// Write file header
	builder.WriteString("// Generated Rust ANSI Art\n")
	builder.WriteString("fn main() {\n")
	builder.WriteString("    let ansi_art_lines = vec![\n")

	// Write the array of lines
	for _, line := range lines {
		// Escape quotes and backslashes in the line
		escapedLine := strings.ReplaceAll(line, "\\", "\\\\")
		escapedLine = strings.ReplaceAll(escapedLine, "\"", "\\\"")
		builder.WriteString(fmt.Sprintf("        \"%s\",\n", escapedLine))
	}
	builder.WriteString("    ];\n\n")

	// Write code to display the art
	builder.WriteString("    for line in ansi_art_lines {\n")
	builder.WriteString("        println!(\"{}\", line);\n")
	builder.WriteString("    }\n")
	builder.WriteString("}\n")

	return builder.String()
}

// GenerateBashCode creates Bash script that reproduces the ANSI art
func GenerateBashCode(lines []string) string {
	var builder strings.Builder

	// Write file header
	builder.WriteString("#!/bin/bash\n")
	builder.WriteString("# Generated Bash ANSI Art\n")
	builder.WriteString("\n")

	// Write array of lines
	builder.WriteString("ansi_art_lines=(\n")
	for _, line := range lines {
		// Escape quotes and backslashes in the line
		escapedLine := strings.ReplaceAll(line, "\\", "\\\\")
		escapedLine = strings.ReplaceAll(escapedLine, "\"", "\\\"")
		escapedLine = strings.ReplaceAll(escapedLine, "$", "\\$")
		escapedLine = strings.ReplaceAll(escapedLine, "`", "\\`")
		builder.WriteString(fmt.Sprintf("    \"%s\"\n", escapedLine))
	}
	builder.WriteString(")\n\n")

	// Write function to display the art
	builder.WriteString("display_ansi_art() {\n")
	builder.WriteString("    for line in \"${ansi_art_lines[@]}\"; do\n")
	builder.WriteString("        echo -e \"$line\"\n")
	builder.WriteString("    done\n")
	builder.WriteString("}\n\n")
	builder.WriteString("# Call the function\n")
	builder.WriteString("display_ansi_art\n")

	return builder.String()
}
