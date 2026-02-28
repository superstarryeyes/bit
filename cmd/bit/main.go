package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/superstarryeyes/bit/ansifonts"
	"github.com/superstarryeyes/bit/internal/fit"
	"github.com/superstarryeyes/bit/internal/ui"
)

func main() {
	// Define CLI flags
	var fontName string
	var textColor string
	var gradientColor string
	var gradientDirection string
	var charSpacing int
	var wordSpacing int
	var lineSpacing int
	var scaleInt int
	var shadowEnabled bool
	var shadowH int
	var shadowV int
	var shadowStyle int
	var alignment string
	var list bool
	var version bool
	var loadFontPath string
	var fitFonts string
	var fitScales string
	var targetW int
	var targetH int
	var priority string
	var limit int
	var showDims bool
	var noPreview bool

	flag.StringVar(&fontName, "font", "", "Font name to use (default: first available font)")
	flag.StringVar(&textColor, "color", "", "Text color: ANSI code (31) or hex (#FF0000)")
	flag.StringVar(&gradientColor, "gradient", "", "Gradient end color: ANSI code (34) or hex (#0000FF)")
	flag.StringVar(&gradientDirection, "direction", "down", "Gradient direction: down, up, right, left")
	flag.IntVar(&charSpacing, "char-spacing", 2, "Character spacing (0 to 10)")
	flag.IntVar(&wordSpacing, "word-spacing", 2, "Word spacing (0 to 20)")
	flag.IntVar(&lineSpacing, "line-spacing", 1, "Line spacing (0 to 10)")
	flag.IntVar(&scaleInt, "scale", 0, "Text scale: -1 (0.5x), 0 (1x), 1 (2x), 2 (4x)")
	flag.BoolVar(&shadowEnabled, "shadow", false, "Enable shadow effect")
	flag.IntVar(&shadowH, "shadow-h", 1, "Shadow horizontal offset (-5 to 5)")
	flag.IntVar(&shadowV, "shadow-v", 1, "Shadow vertical offset (-5 to 5)")
	flag.IntVar(&shadowStyle, "shadow-style", 1, "Shadow style: 0 (light), 1 (medium), 2 (dark)")
	flag.StringVar(&alignment, "align", "center", "Text alignment: left, center, right")
	flag.BoolVar(&list, "list", false, "List all available fonts")
	flag.BoolVar(&version, "version", false, "Show version information")
	flag.StringVar(&loadFontPath, "load", "", "Path to a custom font file (.bit) OR a directory of fonts")
	flag.StringVar(&fitFonts, "fonts", "", "CSV list of fonts to test in fit mode")
	flag.StringVar(&fitScales, "scales", "", "CSV list of scale factors to test in fit mode")
	flag.IntVar(&targetW, "target-w", 0, "Target width for fit mode")
	flag.IntVar(&targetH, "target-h", 0, "Target height for fit mode")
	flag.StringVar(&priority, "priority", "width", "Fit priority: width or height")
	flag.IntVar(&limit, "limit", 10, "Max results to show in fit mode")
	flag.BoolVar(&showDims, "show-dims", false, "Show measured dimensions in fit mode")
	flag.BoolVar(&noPreview, "no-preview", false, "Disable preview render in fit mode")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Bit - Terminal ANSI Logo Designer & Font Library\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "  bit                    Start interactive UI\n")
		fmt.Fprintf(os.Stderr, "  bit [options] <text>   Render text with CLI options\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nColor Codes:\n")
		fmt.Fprintf(os.Stderr, "  30=black         31=red              32=green          33=yellow\n")
		fmt.Fprintf(os.Stderr, "  34=blue          35=magenta          36=cyan           37=white\n")
		fmt.Fprintf(os.Stderr, "  90=gray          91=bright-red       92=bright-green   93=bright-yellow\n")
		fmt.Fprintf(os.Stderr, "  94=bright-blue   95=bright-magenta   96=bright-cyan\n")
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  bit                                                    # Start interactive UI\n")
		fmt.Fprintf(os.Stderr, "  bit -list                                              # List all fonts\n")
		fmt.Fprintf(os.Stderr, "  bit \"Hello World\"                                      # Quick render\n")
		fmt.Fprintf(os.Stderr, "  bit -font ithaca -color 31 \"Red\"                       # With font and color\n")
		fmt.Fprintf(os.Stderr, "  bit -font ithaca -color \"#FF0000\" \"Red Hex\"            # Hex color\n")
		fmt.Fprintf(os.Stderr, "  bit -font dogica -color 31 -gradient 34 \"Gradient\"     # Gradient\n")
		fmt.Fprintf(os.Stderr, "  bit -font pressstart -color 32 -shadow \"Shadow\"        # With shadow\n")
		fmt.Fprintf(os.Stderr, "  bit -load ./myfont.bit \"Custom\"                        # Load custom font file\n")
		fmt.Fprintf(os.Stderr, "  bit -load ./fonts/ -list                               # Load custom font directory\n")
		fmt.Fprintf(os.Stderr, "  bit -target-w 40 -target-h 10 \"Hello\"                  # Fit mode\n")
	}

	flag.Parse()

	fitMode := targetW > 0 || targetH > 0

	// Process custom font loading BEFORE other operations
	if loadFontPath != "" {
		loadedFonts, err := ansifonts.RegisterCustomPath(loadFontPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading custom fonts from '%s': %v\n", loadFontPath, err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, "Loaded %d custom fonts: %v\n", len(loadedFonts), loadedFonts)
	}

	// Show version
	if version {
		fmt.Println("Bit - Terminal ANSI Logo Designer & Font Library")
		fmt.Println("Version: 0.3.0")
		fmt.Println("https://github.com/superstarryeyes/bit")
		return
	}

	// List fonts
	if list {
		fonts, err := ansifonts.ListFonts()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error listing fonts: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Available fonts:")
		for _, font := range fonts {
			fmt.Printf("  %s\n", font)
		}
		return
	}

	// If no arguments provided, start interactive UI
	if flag.NArg() == 0 && !list && !version && !fitMode {
		m, err := ui.InitialModel()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error initializing application: %v\n", err)
			os.Exit(1)
		}

		p := tea.NewProgram(m, tea.WithAltScreen())
		if _, err := p.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error running application: %v\n", err)
			os.Exit(1)
		}
		return
	}

	// CLI mode: render text
	text := strings.Join(flag.Args(), " ")
	if text == "" {
		text = "Hello"
	}

	// Replace literal \n with actual newlines
	text = strings.ReplaceAll(text, "\\n", "\n")

	// Helper function to parse color (ANSI code or hex)
	parseColor := func(colorInput string, defaultColor string) string {
		if colorInput == "" {
			return defaultColor
		}
		// Check if it's a hex color (starts with #)
		if strings.HasPrefix(colorInput, "#") {
			// Validate hex format
			if len(colorInput) == 7 {
				return colorInput
			}
			fmt.Fprintf(os.Stderr, "Warning: Invalid hex color '%s', using default\n", colorInput)
			return defaultColor
		}
		// Try ANSI code mapping using centralized color map
		if color, ok := ansifonts.ANSIColorMap[colorInput]; ok {
			return color
		}
		fmt.Fprintf(os.Stderr, "Warning: Unknown color code '%s', using default\n", colorInput)
		return defaultColor
	}

	// Convert scaleInt to actual scale factor
	var scale float64
	switch scaleInt {
	case -1:
		scale = 0.5 // 0.5x
	case 0:
		scale = 1.0 // 1x
	case 1:
		scale = 2.0 // 2x
	case 2:
		scale = 4.0 // 4x
	default:
		scale = 1.0 // Default to 1x if invalid value provided
		fmt.Fprintf(os.Stderr, "Warning: Invalid scale value '%d', using default scale (1x)\n", scaleInt)
	}

	// Build render options
	options := ansifonts.RenderOptions{
		CharSpacing: charSpacing,
		WordSpacing: wordSpacing,
		LineSpacing: lineSpacing,
		ScaleFactor: scale,
	}

	// Set alignment
	switch alignment {
	case "left":
		options.Alignment = ansifonts.LeftAlign
	case "center":
		options.Alignment = ansifonts.CenterAlign
	case "right":
		options.Alignment = ansifonts.RightAlign
	default:
		options.Alignment = ansifonts.CenterAlign
	}

	// Set colors (supports both ANSI codes and hex)
	options.TextColor = parseColor(textColor, "#FFFFFF")

	// Set gradient
	if gradientColor != "" {
		options.GradientColor = parseColor(gradientColor, options.TextColor)

		// Set gradient direction
		switch gradientDirection {
		case "down":
			options.GradientDirection = ansifonts.UpDown
		case "up":
			options.GradientDirection = ansifonts.DownUp
		case "right":
			options.GradientDirection = ansifonts.LeftRight
		case "left":
			options.GradientDirection = ansifonts.RightLeft
		default:
			options.GradientDirection = ansifonts.UpDown
		}

		options.UseGradient = true
	}

	// Set shadow
	if shadowEnabled {
		options.ShadowEnabled = true
		options.ShadowHorizontalOffset = shadowH
		options.ShadowVerticalOffset = shadowV

		// Validate shadow style
		if shadowStyle < 0 || shadowStyle > 2 {
			shadowStyle = 1
		}
		options.ShadowStyle = ansifonts.ShadowStyle(shadowStyle)
	}

	if fitMode {
		fontsToUse, err := resolveFitFonts(fitFonts, fontName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error resolving fonts: %v\n", err)
			os.Exit(1)
		}

		scalesToUse, err := resolveFitScales(fitScales, scale)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing scales: %v\n", err)
			os.Exit(1)
		}

		candidates, err := fit.FindBest(text, fontsToUse, scalesToUse, options, targetW, targetH, fit.Priority(priority))
		if err != nil {
			if missingErr, ok := err.(*fit.MissingFontError); ok {
				fmt.Fprintf(os.Stderr, "Warning: skipped missing fonts: %s\n", strings.Join(missingErr.Missing, ", "))
			} else {
				fmt.Fprintf(os.Stderr, "Error running fit: %v\n", err)
				os.Exit(1)
			}
		}

		if len(candidates) == 0 {
			fmt.Fprintln(os.Stderr, "No fit candidates found")
			os.Exit(1)
		}

		if limit <= 0 || limit > len(candidates) {
			limit = len(candidates)
		}

		fmt.Println("Fit results:")
		for i := 0; i < limit; i++ {
			candidate := candidates[i]
			if showDims {
				fmt.Printf("%2d. %s scale=%.2f %dx%d dw=%d dh=%d\n", i+1, candidate.Font, candidate.Scale, candidate.W, candidate.H, candidate.DW, candidate.DH)
				continue
			}
			fmt.Printf("%2d. %s scale=%.2f dw=%d dh=%d\n", i+1, candidate.Font, candidate.Scale, candidate.DW, candidate.DH)
		}

		if noPreview {
			return
		}

		best := candidates[0]
		bestFont, err := ansifonts.LoadFont(best.Font)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading font '%s': %v\n", best.Font, err)
			os.Exit(1)
		}
		options.ScaleFactor = best.Scale
		fmt.Println()
		preview := ansifonts.RenderTextWithOptions(text, bestFont, options)
		for _, line := range preview {
			fmt.Println(line)
		}
		return
	}

	// If no font specified, use the first available font
	if fontName == "" {
		fonts, err := ansifonts.ListFonts()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error listing fonts: %v\n", err)
			os.Exit(1)
		}
		if len(fonts) == 0 {
			fmt.Fprintf(os.Stderr, "No fonts available\n")
			os.Exit(1)
		}
		fontName = fonts[0]
	}

	// Load the font
	font, err := ansifonts.LoadFont(fontName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading font '%s': %v\n", fontName, err)
		os.Exit(1)
	}

	// Render and print
	rendered := ansifonts.RenderTextWithOptions(text, font, options)
	for _, line := range rendered {
		fmt.Println(line)
	}
}

func resolveFitFonts(fontsCSV string, fallbackFont string) ([]string, error) {
	if fontsCSV != "" {
		return splitCSV(fontsCSV), nil
	}
	if fallbackFont != "" {
		return []string{fallbackFont}, nil
	}
	return ansifonts.ListFonts()
}

func resolveFitScales(scalesCSV string, fallbackScale float64) ([]float64, error) {
	if scalesCSV == "" {
		return []float64{fallbackScale}, nil
	}

	parts := splitCSV(scalesCSV)
	if len(parts) == 0 {
		return nil, fmt.Errorf("no scales provided")
	}

	values := make([]float64, 0, len(parts))
	for _, raw := range parts {
		value, err := strconv.ParseFloat(raw, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid scale %q", raw)
		}
		values = append(values, value)
	}
	return values, nil
}

func splitCSV(value string) []string {
	parts := strings.Split(value, ",")
	results := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed == "" {
			continue
		}
		results = append(results, trimmed)
	}
	return results
}
