// ABOUTME: Keyboard event handlers for the TUI application.
// ABOUTME: Handles export mode, text input, font selection, and all panel interactions.

package ui

import (
	"math/rand/v2"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/superstarryeyes/bit/ansifonts"
	"github.com/superstarryeyes/bit/internal/favorites"
)

// handleWindowResize handles terminal window resize events
func (m *model) handleWindowResize(msg tea.WindowSizeMsg) tea.Cmd {
	m.uiState.width = msg.Width
	m.uiState.height = msg.Height

	// Update layout parameters based on new window size
	_, _, _, _, availableWidth := m.calculateLayoutParameters()

	// Set filename input width to be wide enough for the placeholder
	if m.export.filenameInput.Width < FilenameInputWidth {
		m.export.filenameInput.Width = FilenameInputWidth
	}

	// Determine layout based on window size
	m.updateLayoutMode(availableWidth)

	// Re-render text to ensure proper layout
	m.renderText()
	return nil
}

// updateLayoutMode determines whether to use single or two-row layout
func (m *model) updateLayoutMode(availableWidth int) {
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
}

// handleExportModeKeys handles keyboard input when in export mode
func (m *model) handleExportModeKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	// Handle overwrite prompt separately
	if m.export.showOverwritePrompt {
		return m.handleOverwritePromptKeys(msg)
	}

	switch msg.String() {
	case "esc":
		m.export.active = false
		m.export.filenameInput.Blur()
		return m, nil
	case "enter":
		if m.export.filenameInput.Value() != "" {
			m.exportText()
			// Don't close export mode yet - let overwrite prompt handle it
			if !m.export.showOverwritePrompt {
				m.export.active = false
				m.export.filenameInput.Blur()
			}
		}
		return m, nil
	case "shift+tab":
		m.cycleExportFormat(-1)
		return m, nil
	case "tab":
		m.cycleExportFormat(1)
		return m, nil
	default:
		m.export.filenameInput, cmd = m.export.filenameInput.Update(msg)
		return m, cmd
	}
}

// handleOverwritePromptKeys handles keyboard input for the overwrite confirmation prompt
func (m *model) handleOverwritePromptKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "left", "right", "h", "l":
		// Toggle between Yes (0) and No (1)
		m.export.selectedButton = 1 - m.export.selectedButton
		return m, nil
	case "enter":
		if m.export.selectedButton == 0 {
			// Yes - proceed with overwrite
			// Check if this is binary content (PNG) or text content
			if len(m.export.overwriteBinaryContent) > 0 {
				m.performBinaryExport(m.export.overwriteBinaryContent, m.export.overwriteFilename, m.export.overwriteFormat)
			} else {
				m.performExport(m.export.overwriteContent, m.export.overwriteFilename, m.export.overwriteFormat)
			}
		}
		// Close overwrite prompt and export mode, clear overwrite data
		m.export.showOverwritePrompt = false
		m.export.active = false
		m.export.filenameInput.Blur()
		m.export.overwriteBinaryContent = nil // Clear binary content
		m.export.overwriteContent = ""        // Clear text content
		return m, nil
	case "esc":
		// Cancel overwrite, clear overwrite data
		m.export.showOverwritePrompt = false
		m.export.overwriteBinaryContent = nil // Clear binary content
		m.export.overwriteContent = ""        // Clear text content
		return m, nil
	}
	return m, nil
}

// cycleExportFormat cycles through export formats
func (m *model) cycleExportFormat(direction int) {
	// Get format names from export manager
	formatNames := m.export.manager.GetFormatNames()

	// Defensive programming: check if formatNames is empty to prevent division by zero
	if len(formatNames) == 0 {
		return
	}

	currentIndex := 0
	for i, f := range formatNames {
		if f == m.export.format {
			currentIndex = i
			break
		}
	}

	if direction > 0 {
		currentIndex = (currentIndex + 1) % len(formatNames)
	} else {
		currentIndex = (currentIndex - 1 + len(formatNames)) % len(formatNames)
	}

	m.export.format = formatNames[currentIndex]
}

// handleTextPanelUpdate handles updates for the text input panel
func (m *model) handleTextPanelUpdate(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg.String() {
	case "tab":
		m.textInput.mode = TextInputMode((int(m.textInput.mode) + 1) % int(TotalTextInputModes))
		m.textInput.input.Blur()
	case "up", "down":
		if m.textInput.mode == TextEntryMode && m.textInput.input.Focused() {
			m.handleMultiRowNavigation(msg.String())
		} else if m.textInput.mode == TextAlignmentMode {
			m.handleTextAlignment(msg.String())
		}
	case "enter":
		if m.textInput.mode == TextEntryMode {
			m.handleTextInputToggle()
		}
	case "e", "r":
		// Check if text input is focused before handling special keys
		if m.textInput.input.Focused() {
			// Let text input process the key
			m.textInput.input, cmd = m.textInput.input.Update(msg)
			return m, cmd
		}
		// Not focused, handle as special command
		return m, nil
	case "k", "j":
		if m.textInput.mode == TextEntryMode && m.textInput.input.Focused() {
			// Let text input process the key
			m.textInput.input, cmd = m.textInput.input.Update(msg)
			return m, cmd
		} else if m.textInput.mode == TextAlignmentMode {
			m.handleTextAlignment(msg.String())
		}
	default:
		if m.textInput.input.Focused() {
			m.textInput.input, cmd = m.textInput.input.Update(msg)
			return m, cmd
		}
	}

	return m, cmd
}

// handleMultiRowNavigation handles up/down navigation in multi-row text input
func (m *model) handleMultiRowNavigation(direction string) {
	// Save current cursor position before moving
	if m.textInput.currentRow < len(m.textInput.rowCursors) {
		m.textInput.rowCursors[m.textInput.currentRow] = m.textInput.input.Position()
	} else {
		// Extend rowCursors slice if needed
		for i := len(m.textInput.rowCursors); i <= m.textInput.currentRow; i++ {
			m.textInput.rowCursors = append(m.textInput.rowCursors, 0)
		}
		m.textInput.rowCursors[m.textInput.currentRow] = m.textInput.input.Position()
	}

	if isUpKey(direction) {
		if m.textInput.currentRow > 0 {
			m.textInput.textRows[m.textInput.currentRow] = m.textInput.input.Value()
			m.textInput.currentRow--
			m.textInput.input.SetValue(m.textInput.textRows[m.textInput.currentRow])

			// Restore cursor position for this row
			if m.textInput.currentRow < len(m.textInput.rowCursors) {
				cursorPos := m.textInput.rowCursors[m.textInput.currentRow]
				// Ensure cursor position doesn't exceed text length
				textLen := len(m.textInput.textRows[m.textInput.currentRow])
				if cursorPos > textLen {
					cursorPos = textLen
				}
				m.textInput.input.SetCursor(cursorPos)
			} else {
				// Extend rowCursors slice if needed
				for i := len(m.textInput.rowCursors); i <= m.textInput.currentRow; i++ {
					m.textInput.rowCursors = append(m.textInput.rowCursors, 0)
				}
				m.textInput.input.SetCursor(0)
			}
		}
	} else if isDownKey(direction) {
		m.textInput.textRows[m.textInput.currentRow] = m.textInput.input.Value()
		if m.textInput.currentRow < len(m.textInput.textRows)-1 {
			m.textInput.currentRow++
			m.textInput.input.SetValue(m.textInput.textRows[m.textInput.currentRow])

			// Restore cursor position for this row
			if m.textInput.currentRow < len(m.textInput.rowCursors) {
				cursorPos := m.textInput.rowCursors[m.textInput.currentRow]
				// Ensure cursor position doesn't exceed text length
				textLen := len(m.textInput.textRows[m.textInput.currentRow])
				if cursorPos > textLen {
					cursorPos = textLen
				}
				m.textInput.input.SetCursor(cursorPos)
			} else {
				// Extend rowCursors slice if needed
				for i := len(m.textInput.rowCursors); i <= m.textInput.currentRow; i++ {
					m.textInput.rowCursors = append(m.textInput.rowCursors, 0)
				}
				m.textInput.input.SetCursor(0)
			}
		} else {
			m.textInput.textRows = append(m.textInput.textRows, "")
			// Extend rowCursors slice
			m.textInput.rowCursors = append(m.textInput.rowCursors, 0)
			m.textInput.currentRow++
			m.textInput.input.SetValue("")
			m.textInput.input.SetCursor(0)
		}
	}
	m.updateCurrentTextFromRows()
	m.renderText()
}

// handleTextAlignment handles text alignment changes
func (m *model) handleTextAlignment(direction string) {
	if isUpKey(direction) {
		m.textInput.alignment = TextAlignment((int(m.textInput.alignment) - 1 + int(TotalAlignments)) % int(TotalAlignments))
	} else {
		m.textInput.alignment = TextAlignment((int(m.textInput.alignment) + 1) % int(TotalAlignments))
	}
	m.renderText()
}

// handleTextInputToggle toggles text input focus
func (m *model) handleTextInputToggle() {
	if m.textInput.input.Focused() {
		// Save cursor position before blurring
		if m.textInput.currentRow < len(m.textInput.rowCursors) {
			m.textInput.rowCursors[m.textInput.currentRow] = m.textInput.input.Position()
		} else {
			// Extend rowCursors slice if needed
			for i := len(m.textInput.rowCursors); i <= m.textInput.currentRow; i++ {
				m.textInput.rowCursors = append(m.textInput.rowCursors, 0)
			}
			m.textInput.rowCursors[m.textInput.currentRow] = m.textInput.input.Position()
		}
		m.textInput.textRows[m.textInput.currentRow] = m.textInput.input.Value()
		m.textInput.input.Blur()
		m.updateCurrentTextFromRows()
		m.renderText()
	} else {
		m.textInput.input.Focus()
		if m.textInput.currentRow < len(m.textInput.textRows) {
			currentRowText := m.textInput.textRows[m.textInput.currentRow]
			// Check if this is the initial placeholder text
			isInitialPlaceholder := len(m.textInput.textRows) == 1 &&
				m.textInput.textRows[0] == "Hello" &&
				m.textInput.currentRow == 0

			if isInitialPlaceholder {
				m.textInput.input.SetValue("")
			} else {
				m.textInput.input.SetValue(currentRowText)
			}

			// Restore cursor position for this row
			if m.textInput.currentRow < len(m.textInput.rowCursors) {
				cursorPos := m.textInput.rowCursors[m.textInput.currentRow]
				// Ensure cursor position doesn't exceed text length
				textLen := len(currentRowText)
				if cursorPos > textLen {
					cursorPos = textLen
				}
				m.textInput.input.SetCursor(cursorPos)
			} else {
				// Extend rowCursors slice if needed
				for i := len(m.textInput.rowCursors); i <= m.textInput.currentRow; i++ {
					m.textInput.rowCursors = append(m.textInput.rowCursors, 0)
				}
				m.textInput.input.SetCursor(0)
			}
		}
	}
}

// handleFontPanelUpdate handles updates for the font selection panel (panel 1)
func (m *model) handleFontPanelUpdate(msg tea.KeyMsg) {
	switch {
	case isUpKey(msg.String()):
		m.font.selectedFont = (m.font.selectedFont - 1 + len(m.font.fonts)) % len(m.font.fonts)
		m.renderText()
	case isDownKey(msg.String()):
		m.font.selectedFont = (m.font.selectedFont + 1) % len(m.font.fonts)
		m.renderText()
	}
}

// handleSpacingPanelUpdate handles updates for the spacing panel
func (m *model) handleSpacingPanelUpdate(msg tea.KeyMsg) {
	switch msg.String() {
	case "tab":
		m.spacing.mode = SpacingMode((int(m.spacing.mode) + 1) % int(TotalSpacingModes))
	case "up", "down", "k", "j":
		direction := 1
		if isUpKey(msg.String()) {
			direction = 1
		} else {
			direction = -1
		}

		switch m.spacing.mode {
		case CharacterSpacingMode:
			m.adjustCharSpacing(direction)
		case WordSpacingMode:
			m.adjustWordSpacing(direction)
		case LineSpacingMode:
			m.adjustLineSpacing(direction)
		}
		m.renderText()
	}
}

// adjustCharSpacing adjusts character spacing within bounds
func (m *model) adjustCharSpacing(direction int) {
	m.spacing.charSpacing += direction
	if m.spacing.charSpacing > MaxCharSpacing {
		m.spacing.charSpacing = MaxCharSpacing
	} else if m.spacing.charSpacing < MinCharSpacing {
		m.spacing.charSpacing = MinCharSpacing
	}
}

// adjustWordSpacing adjusts word spacing within bounds
func (m *model) adjustWordSpacing(direction int) {
	m.spacing.wordSpacing += direction
	if m.spacing.wordSpacing > MaxWordSpacing {
		m.spacing.wordSpacing = MaxWordSpacing
	} else if m.spacing.wordSpacing < MinWordSpacing {
		m.spacing.wordSpacing = MinWordSpacing
	}
}

// adjustLineSpacing adjusts line spacing within bounds
func (m *model) adjustLineSpacing(direction int) {
	m.spacing.lineSpacing += direction
	if m.spacing.lineSpacing > MaxLineSpacing {
		m.spacing.lineSpacing = MaxLineSpacing
	} else if m.spacing.lineSpacing < MinLineSpacing {
		m.spacing.lineSpacing = MinLineSpacing
	}
}

// handleColorPanelUpdate handles updates for the color panel
func (m *model) handleColorPanelUpdate(msg tea.KeyMsg) {
	switch {
	case msg.String() == "tab":
		m.color.subMode = ColorSubMode((int(m.color.subMode) + 1) % int(TotalColorSubModes))
	case isUpKey(msg.String()), isDownKey(msg.String()):
		direction := -1
		if isDownKey(msg.String()) {
			direction = 1
		}

		switch m.color.subMode {
		case TextColorMode:
			m.color.textColor = (m.color.textColor + direction + len(colorOptions)) % len(colorOptions)
		case GradientColorMode:
			m.color.gradientColor = (m.color.gradientColor + direction + len(colorOptions)) % len(colorOptions)
			m.color.gradientEnabled = (m.color.gradientColor != m.color.textColor)
		case GradientDirectionMode:
			// Ensure we don't get negative indices even with unusual values
			newIndex := (int(m.color.gradientDirection) + direction) % int(TotalGradientDirections)
			// Handle negative modulo result
			if newIndex < 0 {
				newIndex += int(TotalGradientDirections)
			}
			m.color.gradientDirection = GradientDirection(newIndex)
		}
		m.renderText()
	}
}

// handleScalePanelUpdate handles updates for the scale panel
func (m *model) handleScalePanelUpdate(msg tea.KeyMsg) {
	switch {
	case isUpKey(msg.String()):
		if m.scale.scale < MaxScale {
			m.scale.scale++
			m.renderText()
		}
	case isDownKey(msg.String()):
		if m.scale.scale > MinScale {
			m.scale.scale--
			m.renderText()
		}
	}

	// Update shadow warning based on actual half-pixel condition
	m.updateShadowWarning()
}

// updateShadowWarning checks if half-pixels are detected and updates the warning state accordingly
func (m *model) updateShadowWarning() {
	if len(m.font.fonts) == 0 || m.textInput.currentText == "" {
		m.shadow.showWarning = false
		return
	}

	// Load font data if not already loaded (lazy loading)
	selectedFont := &m.font.fonts[m.font.selectedFont]
	if !selectedFont.Loaded {
		err := loadFontData(selectedFont)
		if err != nil {
			m.shadow.showWarning = false
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

	// Check for half-pixel usage, which affects shadow compatibility
	hasHalfPixels := ansifonts.DetectHalfPixelUsage(m.textInput.currentText, ansiFontData, m.getScaleFactorFloat())

	// Update warning based on actual half-pixel condition and shadow settings
	// Use canonical offset values instead of UI indices
	m.shadow.showWarning = hasHalfPixels && m.shadow.enabled && (m.shadow.horizontalOffset != 0 || m.shadow.verticalOffset != 0)
}

// handleShadowPanelUpdate handles updates for the shadow panel
func (m *model) handleShadowPanelUpdate(msg tea.KeyMsg) {
	switch {
	case msg.String() == "tab":
		m.shadow.subMode = ShadowSubMode((int(m.shadow.subMode) + 1) % int(TotalShadowSubModes))
	case isUpKey(msg.String()), isDownKey(msg.String()):
		switch m.shadow.subMode {
		case HorizontalShadowMode:
			m.handleHorizontalShadow(msg.String())
		case VerticalShadowMode:
			m.handleVerticalShadow(msg.String())
		case ShadowStyleMode:
			m.handleShadowStyle(msg.String())
		}
		m.renderText()

		// Update shadow warning based on current settings
		m.updateShadowWarning()
	}
}

// handleHorizontalShadow handles horizontal shadow adjustments
func (m *model) handleHorizontalShadow(direction string) {
	if isUpKey(direction) {
		m.shadow.horizontalIndex = (m.shadow.horizontalIndex + 1) % len(shadowPixelOptions)
	} else {
		m.shadow.horizontalIndex = (m.shadow.horizontalIndex - 1 + len(shadowPixelOptions)) % len(shadowPixelOptions)
	}

	// Update the canonical offset value from the UI index
	m.shadow.horizontalOffset = shadowPixelOptions[m.shadow.horizontalIndex].Pixels
	m.shadow.enabled = (m.shadow.horizontalOffset != 0 || m.shadow.verticalOffset != 0)
}

// handleVerticalShadow handles vertical shadow adjustments
func (m *model) handleVerticalShadow(direction string) {
	if isUpKey(direction) {
		m.shadow.verticalIndex = (m.shadow.verticalIndex + 1) % len(verticalShadowPixelOptions)
	} else {
		m.shadow.verticalIndex = (m.shadow.verticalIndex - 1 + len(verticalShadowPixelOptions)) % len(verticalShadowPixelOptions)
	}

	// Update the canonical offset value from the UI index
	m.shadow.verticalOffset = verticalShadowPixelOptions[m.shadow.verticalIndex].Pixels
	m.shadow.enabled = (m.shadow.horizontalOffset != 0 || m.shadow.verticalOffset != 0)
}

// handleShadowStyle handles shadow style changes
func (m *model) handleShadowStyle(direction string) {
	if isUpKey(direction) {
		m.shadow.style = (m.shadow.style + 1) % len(shadowStyleOptions)
	} else {
		m.shadow.style = (m.shadow.style - 1 + len(shadowStyleOptions)) % len(shadowStyleOptions)
	}
}

// handleRandomize randomizes font and color settings
func (m *model) handleRandomize() {
	m.font.selectedFont = rand.IntN(len(m.font.fonts))
	m.color.textColor = rand.IntN(len(colorOptions))
	m.color.gradientColor = rand.IntN(len(colorOptions))
	m.color.gradientEnabled = (m.color.gradientColor != m.color.textColor)
	m.color.gradientDirection = GradientDirection(rand.IntN(int(TotalGradientDirections)))

	// Update shadow warning after randomization
	m.updateShadowWarning()
	m.renderText()
}

// handleEnterExportMode enters export mode
func (m *model) handleEnterExportMode() {
	m.export.active = true
	m.export.filenameInput.Focus()
	m.export.filenameInput.SetValue("")
}

// resetConfirmations resets all confirmation and warning messages
func (m *model) resetConfirmations() {
	m.export.showConfirmation = false
	m.export.confirmationText = ""
	m.shadow.showWarning = false
}

// handlePanelNavigation handles left/right panel navigation
func (m *model) handlePanelNavigation(direction int) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	// Check if text input is focused and should handle the key
	if m.uiState.focusedPanel == TextInputPanel && m.textInput.input.Focused() {
		// Let text input handle the key
		var msg tea.KeyMsg
		if direction > 0 {
			msg = tea.KeyMsg{Type: tea.KeyRight}
		} else {
			msg = tea.KeyMsg{Type: tea.KeyLeft}
		}
		m.textInput.input, cmd = m.textInput.input.Update(msg)
		return m, cmd
	}

	// Navigate panels
	if direction > 0 {
		m.uiState.focusedPanel = FocusedPanel((int(m.uiState.focusedPanel) + 1) % int(TotalPanels))
	} else {
		m.uiState.focusedPanel = FocusedPanel((int(m.uiState.focusedPanel) - 1 + int(TotalPanels)) % int(TotalPanels))
	}
	m.textInput.input.Blur()

	return m, cmd
}

// handleEnterFavoritesMode enters favorites mode
func (m *model) handleEnterFavoritesMode() {
	m.favorites.active = true
	m.favorites.showNamePrompt = false
	m.favorites.selectedIndex = 0
	m.favorites.nameInput.Blur()
}

// handleFavoritesModeKeys handles keyboard input when in favorites mode
func (m *model) handleFavoritesModeKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	// Handle name input prompt
	if m.favorites.showNamePrompt {
		return m.handleFavoritesNamePromptKeys(msg)
	}

	favList := m.favorites.manager.List()

	switch msg.String() {
	case "esc":
		m.favorites.active = false
		m.favorites.showConfirmation = false
		return m, nil

	case "s":
		// Save current art as favorite - show name prompt
		m.favorites.showNamePrompt = true
		m.favorites.nameInput.Focus()
		m.favorites.nameInput.SetValue("")
		return m, nil

	case "up":
		if len(favList) > 0 && m.favorites.selectedIndex > 0 {
			m.favorites.selectedIndex--
		}
		return m, nil

	case "down":
		if len(favList) > 0 && m.favorites.selectedIndex < len(favList)-1 {
			m.favorites.selectedIndex++
		}
		return m, nil

	case "enter":
		// Load selected favorite
		if len(favList) > 0 && m.favorites.selectedIndex < len(favList) {
			m.loadFavorite(&favList[m.favorites.selectedIndex])
			m.favorites.active = false
		}
		return m, nil

	case "d", "backspace", "delete":
		// Delete selected favorite
		if len(favList) > 0 && m.favorites.selectedIndex < len(favList) {
			id := favList[m.favorites.selectedIndex].ID
			err := m.favorites.manager.Remove(id)
			if err == nil {
				// Adjust selection if needed
				newList := m.favorites.manager.List()
				if m.favorites.selectedIndex >= len(newList) && m.favorites.selectedIndex > 0 {
					m.favorites.selectedIndex--
				}
			}
		}
		return m, nil

	default:
		return m, cmd
	}
}

// handleFavoritesNamePromptKeys handles keyboard input for the favorites name prompt
func (m *model) handleFavoritesNamePromptKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg.String() {
	case "esc":
		m.favorites.showNamePrompt = false
		m.favorites.nameInput.Blur()
		return m, nil

	case "enter":
		name := m.favorites.nameInput.Value()
		if name == "" {
			// Use first 20 chars of text as default name
			name = m.textInput.currentText
			if len(name) > 20 {
				name = name[:20]
			}
			name = strings.ReplaceAll(name, "\n", " ")
		}

		// Create favorite from current state
		fav := m.createFavoriteFromCurrentState(name)
		_, err := m.favorites.manager.Add(fav)
		if err == nil {
			m.favorites.showConfirmation = true
			m.favorites.confirmationText = "Saved: " + name
		}

		m.favorites.showNamePrompt = false
		m.favorites.nameInput.Blur()
		return m, nil

	default:
		m.favorites.nameInput, cmd = m.favorites.nameInput.Update(msg)
		return m, cmd
	}
}

// createFavoriteFromCurrentState creates a Favorite from current model state
func (m *model) createFavoriteFromCurrentState(name string) favorites.Favorite {
	fontName := ""
	if len(m.font.fonts) > 0 && m.font.selectedFont < len(m.font.fonts) {
		fontName = m.font.fonts[m.font.selectedFont].Name
	}

	return favorites.Favorite{
		Name:      name,
		Text:      m.textInput.currentText,
		FontName:  fontName,
		Alignment: int(m.textInput.alignment),

		CharSpacing: m.spacing.charSpacing,
		WordSpacing: m.spacing.wordSpacing,
		LineSpacing: m.spacing.lineSpacing,

		TextColor:         m.color.textColor,
		GradientEnabled:   m.color.gradientEnabled,
		GradientColor:     m.color.gradientColor,
		GradientDirection: int(m.color.gradientDirection),

		Scale: int(m.scale.scale),

		ShadowEnabled: m.shadow.enabled,
		ShadowHOffset: m.shadow.horizontalOffset,
		ShadowVOffset: m.shadow.verticalOffset,
		ShadowStyle:   m.shadow.style,
	}
}

// loadFavorite restores model state from a Favorite
func (m *model) loadFavorite(fav *favorites.Favorite) {
	// Restore text
	m.textInput.currentText = fav.Text
	m.textInput.textRows = strings.Split(fav.Text, "\n")
	m.textInput.rowCursors = make([]int, len(m.textInput.textRows))
	m.textInput.currentRow = 0
	m.textInput.alignment = TextAlignment(fav.Alignment)

	// Find font by name
	fontFound := false
	for i, font := range m.font.fonts {
		if font.Name == fav.FontName {
			m.font.selectedFont = i
			fontFound = true
			break
		}
	}
	if !fontFound && len(m.font.fonts) > 0 {
		m.font.selectedFont = 0 // Fallback to first font
	}

	// Restore spacing
	m.spacing.charSpacing = fav.CharSpacing
	m.spacing.wordSpacing = fav.WordSpacing
	m.spacing.lineSpacing = fav.LineSpacing

	// Restore color
	m.color.textColor = fav.TextColor
	m.color.gradientEnabled = fav.GradientEnabled
	m.color.gradientColor = fav.GradientColor
	m.color.gradientDirection = GradientDirection(fav.GradientDirection)

	// Restore scale
	m.scale.scale = TextScale(fav.Scale)

	// Restore shadow
	m.shadow.enabled = fav.ShadowEnabled
	m.shadow.horizontalOffset = fav.ShadowHOffset
	m.shadow.verticalOffset = fav.ShadowVOffset
	m.shadow.style = fav.ShadowStyle

	// Update shadow UI indices from canonical values
	m.shadow.horizontalIndex = m.findShadowPixelIndex(fav.ShadowHOffset, shadowPixelOptions)
	m.shadow.verticalIndex = m.findShadowPixelIndex(fav.ShadowVOffset, verticalShadowPixelOptions)

	// Re-render
	m.renderText()
}

// findShadowPixelIndex finds the index for a shadow pixel value
func (m *model) findShadowPixelIndex(value int, options []ShadowPixelOption) int {
	for i, opt := range options {
		if opt.Pixels == value {
			return i
		}
	}
	return 0 // Default to first option
}
