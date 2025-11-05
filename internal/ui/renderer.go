package ui

import (
	"fmt"

	"github.com/superstarryeyes/bit/ansifonts"
)

func (m *model) renderText() {
	if len(m.font.fonts) == 0 || m.textInput.currentText == "" {
		m.uiState.renderedLines = []string{"No text or fonts available"}
		return
	}

	// Load font data if not already loaded (lazy loading)
	selectedFont := &m.font.fonts[m.font.selectedFont]
	if !selectedFont.Loaded {
		err := loadFontData(selectedFont)
		if err != nil {
			m.uiState.renderedLines = []string{fmt.Sprintf("Error loading font: %v", err)}
			return
		}
	}

	font := *selectedFont.FontData

	// Convert internal FontData to ansifonts.FontData
	ansiFontData := ansifonts.FontData{
		Name:       font.Name,
		Author:     font.Author,
		License:    font.License,
		Characters: font.Characters,
	}

	// Determine color mode
	var colorMode ansifonts.ColorMode
	if m.color.rainbowEnabled {
		colorMode = ansifonts.Rainbow
	} else if m.color.gradientEnabled && m.color.gradientColor != m.color.textColor {
		colorMode = ansifonts.Gradient
	} else {
		colorMode = ansifonts.SingleColor
	}

	// Create render options from model settings
	options := ansifonts.RenderOptions{
		CharSpacing:            m.spacing.charSpacing,
		WordSpacing:            m.spacing.wordSpacing,
		LineSpacing:            m.spacing.lineSpacing,
		Alignment:              ansifonts.TextAlignment(m.textInput.alignment),
		TextColor:              colorOptions[m.color.textColor].Hex,
		GradientColor:          colorOptions[m.color.gradientColor].Hex,
		GradientDirection:      ansifonts.GradientDirection(m.color.gradientDirection),
		UseGradient:            m.color.gradientEnabled && m.color.gradientColor != m.color.textColor,
		ColorMode:              colorMode,
		RainbowFrame:           m.background.frame, // Use background frame for rainbow animation
		RainbowSpeed:           RainbowAnimationSpeed,
		ScaleFactor:            m.getScaleFactorFloat(),
		ShadowEnabled:          m.shadow.enabled,
		ShadowHorizontalOffset: m.shadow.horizontalOffset,
		ShadowVerticalOffset:   m.shadow.verticalOffset,
		ShadowStyle:            ansifonts.ShadowStyle(m.shadow.style),
	}

	// Check for half-pixel usage to show warning in UI
	// The ansifonts library will automatically disable shadows if needed
	hasHalfPixels := ansifonts.DetectHalfPixelUsage(m.textInput.currentText, ansiFontData, m.getScaleFactorFloat())
	m.shadow.showWarning = hasHalfPixels && m.shadow.enabled && (m.shadow.horizontalOffset != 0 || m.shadow.verticalOffset != 0)

	// Clear previous rendered lines to prevent memory leak
	m.uiState.renderedLines = nil

	// Render using the ansifonts library - all rendering logic is centralized there
	m.uiState.renderedLines = ansifonts.RenderTextWithFont(m.textInput.currentText, ansiFontData, options)
}

// getScaleFactorFloat converts the UI scale enum to a float64 scale factor
func (m *model) getScaleFactorFloat() float64 {
	switch m.scale.scale {
	case ScaleHalf:
		return 0.5
	case ScaleOne:
		return 1.0
	case ScaleTwo:
		return 2.0
	case ScaleFour:
		return 4.0
	default:
		return 1.0
	}
}
