package ui

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/charmbracelet/lipgloss"
)

// truncateText truncates text to fit within maxWidth, adding "..." if needed
func truncateText(text string, maxWidth int) string {
	if maxWidth <= 0 {
		return ""
	}
	if maxWidth <= 3 {
		return strings.Repeat(".", maxWidth)
	}
	if utf8.RuneCountInString(text) <= maxWidth {
		return text
	}
	// Reserve 3 characters for "..."
	truncated := []rune(text)[:maxWidth-3]
	return string(truncated) + "..."
}

// calculateLayoutParameters calculates the layout parameters for the UI panels
// Uses consistent thresholds with updateLayoutMode
func (m *model) calculateLayoutParameters() (int, int, int, int, int) {
	// Use consistent layout thresholds from constants
	const (
		reservedMargin = LayoutReservedMargin
		minPanelWidth  = LayoutMinPanelWidth
		spacerWidth    = LayoutSpacerWidth
	)

	availableWidth := m.uiState.width - reservedMargin

	// Use the same threshold logic as updateLayoutMode to prevent layout thrashing
	// Check if we're only showing the placeholder text
	isOnlyPlaceholder := true
	if len(m.textInput.textRows) > 0 {
		for _, row := range m.textInput.textRows {
			if strings.TrimSpace(row) != "" {
				isOnlyPlaceholder = false
				break
			}
		}
	}

	// Check if text input is focused
	isTextInputFocused := m.uiState.focusedPanel == TextInputPanel && m.textInput.mode == TextEntryMode && m.textInput.input.Focused()

	if m.uiState.usesTwoRows {
		// Switch back to single row when we have comfortable width
		if availableWidth >= ComfortableWidthSingleRow && !isTextInputFocused {
			m.uiState.usesTwoRows = false
		}
	} else {
		// Switch to two rows only when needed and not focused
		if !isTextInputFocused && availableWidth < MinWidthSingleRow && !isOnlyPlaceholder {
			m.uiState.usesTwoRows = true
		}
	}

	// Calculate panel count based on layout
	totalPanels := 8 // Single row: 8 panels
	if m.uiState.usesTwoRows {
		totalPanels = 4 // Two rows: 4 panels per row
	}

	// Calculate panel width with fixed spacing
	totalSpacerWidth := (totalPanels - 1) * spacerWidth
	panelWidth := (availableWidth - totalSpacerWidth) / totalPanels

	// Enforce minimum panel width
	panelWidth = max(panelWidth, minPanelWidth)

	// Calculate content width (accounting for borders and padding)
	contentWidth := max(panelWidth-4, 3)

	return panelWidth, contentWidth, spacerWidth, totalPanels, availableWidth
}

// createPanelContents creates the content strings for all panels
func (m *model) createPanelContents(contentWidth int) (string, string, string, string, string, string, string, string) {
	// Helper function to count non-empty rows
	countNonEmptyRows := func(rows []string) int {
		count := 0
		for _, row := range rows {
			if strings.TrimSpace(row) != "" {
				count++
			}
		}
		return count
	}

	// Text panel content - depends on current text input mode
	var textPanelContent string
	if m.uiState.focusedPanel == TextInputPanel && m.textInput.mode == TextEntryMode && m.textInput.input.Focused() {
		// When in text input edit mode, show just the textinput component
		textPanelContent = m.textInput.input.View()
	} else if m.uiState.focusedPanel == TextInputPanel && m.textInput.mode == TextAlignmentMode {
		// When in text alignment mode, show current alignment
		alignmentNames := []string{"Left", "Center", "Right"}
		textPanelContent = truncateText(alignmentNames[int(m.textInput.alignment)], contentWidth)
	} else {
		// When not in edit mode, show row count and preview
		nonEmptyRows := countNonEmptyRows(m.textInput.textRows)
		if nonEmptyRows == 0 {
			textPanelContent = truncateText("Enter text...", contentWidth)
		} else if nonEmptyRows == 1 {
			// Find the first non-empty row
			for _, row := range m.textInput.textRows {
				if strings.TrimSpace(row) != "" {
					textPanelContent = truncateText(row, contentWidth)
					break
				}
			}
		} else {
			// Show multi-row indicator with non-empty row count
			firstNonEmptyRow := ""
			for _, row := range m.textInput.textRows {
				if strings.TrimSpace(row) != "" {
					firstNonEmptyRow = row
					break
				}
			}
			preview := truncateText(firstNonEmptyRow, contentWidth-10) // Reserve space for row count
			textPanelContent = fmt.Sprintf("%s (%d rows)", preview, nonEmptyRows)
		}
	}

	var fontPanelContent string
	if len(m.font.fonts) > 0 {
		fontPanelContent = truncateText(m.font.fonts[m.font.selectedFont].Name, contentWidth)
	} else {
		fontPanelContent = truncateText("No fonts", contentWidth)
	}

	// Combined spacing content based on current mode
	var spacingContent string
	if m.spacing.mode == CharacterSpacingMode {
		spacingContent = truncateText(fmt.Sprintf("%d", m.spacing.charSpacing), contentWidth)
	} else if m.spacing.mode == WordSpacingMode {
		spacingContent = truncateText(fmt.Sprintf("%d", m.spacing.wordSpacing), contentWidth)
	} else { // Line spacing
		spacingContent = truncateText(fmt.Sprintf("%d", m.spacing.lineSpacing), contentWidth)
	}

	// Color content based on current sub-mode
	var colorContent string
	if m.color.subMode == TextColorMode {
		colorContent = truncateText(colorOptions[m.color.textColor].Name, contentWidth)
	} else if m.color.subMode == GradientColorMode {
		if m.color.gradientEnabled {
			colorContent = truncateText(colorOptions[m.color.gradientColor].Name, contentWidth)
		} else {
			colorContent = truncateText("None", contentWidth)
		}
	} else if m.color.subMode == GradientDirectionMode {
		colorContent = truncateText(gradientDirectionOptions[int(m.color.gradientDirection)].Name, contentWidth)
	} else { // Rainbow mode
		if m.color.rainbowEnabled {
			colorContent = truncateText("On", contentWidth)
		} else {
			colorContent = truncateText("Off", contentWidth)
		}
	}

	var scaleContent string
	switch m.scale.scale {
	case ScaleHalf:
		scaleContent = truncateText("0.5x", contentWidth)
	case ScaleOne:
		scaleContent = truncateText("1x", contentWidth)
	case ScaleTwo:
		scaleContent = truncateText("2x", contentWidth)
	case ScaleFour:
		scaleContent = truncateText("4x", contentWidth)
	default:
		scaleContent = truncateText("1x", contentWidth)
	}

	// Combined shadow content based on current sub-mode
	var shadowContent string
	if m.shadow.subMode == HorizontalShadowMode {
		shadowContent = truncateText(shadowPixelOptions[m.shadow.horizontalIndex].Name, contentWidth)
	} else if m.shadow.subMode == VerticalShadowMode {
		shadowContent = truncateText(verticalShadowPixelOptions[m.shadow.verticalIndex].Name, contentWidth)
	} else { // Style mode (ANSI character texture)
		// Display the actual ANSI character texture instead of just the name
		styleChar := string(shadowStyleOptions[m.shadow.style].Char)
		// Repeat the character to fill the content width
		if contentWidth > 0 {
			repeatCount := min(contentWidth,
				// Limit the repetition for better visual appearance
				MaxShadowRepeatCount)
			shadowContent = strings.Repeat(styleChar, repeatCount)
		} else {
			shadowContent = styleChar
		}
	}

	// Background content based on current sub-mode
	var backgroundContent string
	if m.background.subMode == BackgroundTypeMode {
		backgroundNames := []string{"None", "Lava Lamp", "Wavy Grid", "Ticker", "Starfield"}
		backgroundContent = truncateText(backgroundNames[int(m.background.backgroundType)], contentWidth)
	}

	// Animation content based on current sub-mode
	var animationContent string
	if m.animation.subMode == AnimationTypeMode {
		animationNames := []string{"None", "Scroll ←", "Scroll →"}
		animationContent = truncateText(animationNames[int(m.animation.animationType)], contentWidth)
	} else if m.animation.subMode == AnimationSpeedMode {
		speedNames := []string{"Slow", "Medium", "Fast"}
		animationContent = truncateText(speedNames[int(m.animation.speed)], contentWidth)
	}

	return textPanelContent, fontPanelContent, spacingContent, colorContent, scaleContent, shadowContent, backgroundContent, animationContent
}

// createStyledPanels creates styled panels with appropriate selection highlighting
func (m *model) createStyledPanels(panelWidth int, textContent, fontContent, spacingContent, colorContent, scaleContent, shadowContent, backgroundContent, animationContent string) (string, string, string, string, string, string, string, string) {
	normalStyles, selectedStyles := createPanelStyles(panelWidth)

	var textPanel, fontPanel, spacingPanel, colorPanel, scalePanel, shadowPanel, backgroundPanel, animationPanel string

	if m.uiState.focusedPanel == TextInputPanel {
		textPanel = selectedStyles["textInput"].Render(textContent)
	} else {
		textPanel = normalStyles["textInput"].Render(textContent)
	}

	if m.uiState.focusedPanel == FontPanel {
		fontPanel = selectedStyles["font"].Render(fontContent)
	} else {
		fontPanel = normalStyles["font"].Render(fontContent)
	}

	if m.uiState.focusedPanel == SpacingPanel {
		if m.spacing.mode == CharacterSpacingMode {
			spacingPanel = selectedStyles["charSpacing"].Render(spacingContent)
		} else if m.spacing.mode == WordSpacingMode {
			spacingPanel = selectedStyles["wordSpacing"].Render(spacingContent)
		} else {
			spacingPanel = selectedStyles["lineSpacing"].Render(spacingContent)
		}
	} else {
		if m.spacing.mode == CharacterSpacingMode {
			spacingPanel = normalStyles["charSpacing"].Render(spacingContent)
		} else if m.spacing.mode == WordSpacingMode {
			spacingPanel = normalStyles["wordSpacing"].Render(spacingContent)
		} else {
			spacingPanel = normalStyles["lineSpacing"].Render(spacingContent)
		}
	}

	if m.uiState.focusedPanel == ColorPanel {
		colorPanel = selectedStyles["color"].Render(colorContent)
	} else {
		colorPanel = normalStyles["color"].Render(colorContent)
	}

	if m.uiState.focusedPanel == ScalePanel {
		scalePanel = selectedStyles["scale"].Render(scaleContent)
	} else {
		scalePanel = normalStyles["scale"].Render(scaleContent)
	}

	if m.uiState.focusedPanel == ShadowPanel {
		// Combined shadow panel styling
		if m.shadow.subMode == HorizontalShadowMode {
			shadowPanel = selectedStyles["shadow"].Render(shadowContent)
		} else { // Vertical shadow
			shadowPanel = selectedStyles["verticalShadow"].Render(shadowContent)
		}
	} else {
		shadowPanel = normalStyles["shadow"].Render(shadowContent)
	}

	if m.uiState.focusedPanel == BackgroundPanel {
		backgroundPanel = selectedStyles["background"].Render(backgroundContent)
	} else {
		backgroundPanel = normalStyles["background"].Render(backgroundContent)
	}

	if m.uiState.focusedPanel == AnimationPanel {
		animationPanel = selectedStyles["animation"].Render(animationContent)
	} else {
		animationPanel = normalStyles["animation"].Render(animationContent)
	}

	return textPanel, fontPanel, spacingPanel, colorPanel, scalePanel, shadowPanel, backgroundPanel, animationPanel
}

// arrangeControlPanels arranges the control panels in either single or double row layout
func (m *model) arrangeControlPanels(spacerWidth int, labeledTextPanel, labeledFontPanel, labeledSpacingPanel, labeledColorPanel, labeledScalePanel, labeledShadowPanel, labeledBackgroundPanel, labeledAnimationPanel string) string {
	// Create spacer with calculated width
	spacer := strings.Repeat(" ", spacerWidth)

	// Calculate the height of labeled panels to ensure consistent layout
	labeledPanelHeight := lipgloss.Height(labeledTextPanel)

	// Arrange labeled control panels based on layout with width validation
	var controlPanelsRow string
	if m.uiState.usesTwoRows {
		// First row: Text, Font, Spacing, Color
		firstRow := lipgloss.JoinHorizontal(lipgloss.Top,
			labeledTextPanel,
			spacer,
			labeledFontPanel,
			spacer,
			labeledSpacingPanel,
			spacer,
			labeledColorPanel,
		)

		// Second row: Scale, Shadow, Background, Animation
		secondRow := lipgloss.JoinHorizontal(lipgloss.Top,
			labeledScalePanel,
			spacer,
			labeledShadowPanel,
			spacer,
			labeledBackgroundPanel,
			spacer,
			labeledAnimationPanel,
		)

		// Combine rows vertically WITHOUT extra spacing to eliminate unnecessary newline
		controlPanelsRow = lipgloss.JoinVertical(lipgloss.Left, firstRow, secondRow)

		// Set a fixed height for the control panels area to prevent jumping
		// Account for both panel rows and label rows
		controlPanelsHeight := labeledPanelHeight * 2 // 2 panel rows
		controlPanelsRow = lipgloss.NewStyle().Height(controlPanelsHeight).Render(controlPanelsRow)
	} else {
		// Single row: all 8 panels - ensure they fit within terminal width
		controlPanelsRow = lipgloss.JoinHorizontal(lipgloss.Top,
			labeledTextPanel,
			spacer,
			labeledFontPanel,
			spacer,
			labeledSpacingPanel,
			spacer,
			labeledColorPanel,
			spacer,
			labeledScalePanel,
			spacer,
			labeledShadowPanel,
			spacer,
			labeledBackgroundPanel,
			spacer,
			labeledAnimationPanel,
		)

		// Set a fixed height for the control panels area
		controlPanelsHeight := labeledPanelHeight
		controlPanelsRow = lipgloss.NewStyle().Height(controlPanelsHeight).Render(controlPanelsRow)
	}

	// Center the control panels row with overflow protection
	controlPanelsWidth := lipgloss.Width(controlPanelsRow)
	maxAllowedWidth := m.uiState.width - 2 // Leave 2 chars margin

	var controlPanels string
	// If the control panels are too wide, don't center them to prevent overflow
	if controlPanelsWidth > maxAllowedWidth {
		controlPanels = controlPanelsRow
	} else {
		controlPanels = lipgloss.NewStyle().
			Width(m.uiState.width).
			Align(lipgloss.Center).
			Render(controlPanelsRow)
	}

	return controlPanels
}
