package fit

import (
	"fmt"
	"sort"
	"strings"

	"github.com/superstarryeyes/bit/ansifonts"
)

// Candidate represents a font and scale candidate with its measured dimensions.
type Candidate struct {
	Font  string
	Scale float64
	W     int
	H     int
	DW    int
	DH    int
}

// Priority indicates which dimension should be prioritized in sorting.
type Priority string

const (
	PriorityWidth  Priority = "width"
	PriorityHeight Priority = "height"
)

// MissingFontError reports fonts that could not be loaded.
type MissingFontError struct {
	Missing []string
}

func (err *MissingFontError) Error() string {
	return fmt.Sprintf("failed to load %d font(s): %s", len(err.Missing), strings.Join(err.Missing, ", "))
}

// FindBest renders text with each font/scale, measures it, and returns sorted candidates.
func FindBest(text string, fonts []string, scales []float64, baseOptions ansifonts.RenderOptions, targetW int, targetH int, priority Priority) ([]Candidate, error) {
	if len(fonts) == 0 {
		return nil, fmt.Errorf("no fonts provided")
	}
	if len(scales) == 0 {
		return nil, fmt.Errorf("no scales provided")
	}

	normalizedPriority, err := normalizePriority(priority)
	if err != nil {
		return nil, err
	}

	fontCache := make(map[string]*ansifonts.Font, len(fonts))
	var candidates []Candidate
	var missingFonts []string

	for _, fontName := range fonts {
		fontObj, ok := fontCache[fontName]
		if !ok {
			loadedFont, loadErr := ansifonts.LoadFont(fontName)
			if loadErr != nil {
				missingFonts = append(missingFonts, fontName)
				continue
			}
			fontCache[fontName] = loadedFont
			fontObj = loadedFont
		}

		for _, scale := range scales {
			options := baseOptions
			options.ScaleFactor = scale
			rendered := ansifonts.RenderTextWithOptions(text, fontObj, options)
			w, h := ansifonts.MeasureBlock(rendered)
			candidate := Candidate{
				Font:  fontObj.Name,
				Scale: scale,
				W:     w,
				H:     h,
				DW:    distance(targetW, w),
				DH:    distance(targetH, h),
			}
			if candidate.DW < 0 || candidate.DH < 0 {
				continue // Skip candidates that exceed target dimensions
			}
			candidates = append(candidates, candidate)
		}
	}

	if len(candidates) == 0 {
		if len(missingFonts) > 0 {
			return nil, &MissingFontError{Missing: missingFonts}
		}
		return nil, fmt.Errorf("no candidates generated")
	}

	sortCandidates(candidates, normalizedPriority)

	if len(missingFonts) > 0 {
		return candidates, &MissingFontError{Missing: missingFonts}
	}

	return candidates, nil
}

func normalizePriority(priority Priority) (Priority, error) {
	switch strings.ToLower(string(priority)) {
	case "", string(PriorityWidth):
		return PriorityWidth, nil
	case string(PriorityHeight):
		return PriorityHeight, nil
	default:
		return "", fmt.Errorf("invalid priority: %s", priority)
	}
}

func sortCandidates(candidates []Candidate, priority Priority) {
	sort.SliceStable(candidates, func(i, j int) bool {
		a := candidates[i]
		b := candidates[j]

		if priority == PriorityHeight {
			if a.DH != b.DH {
				return a.DH < b.DH
			}
			if a.DW != b.DW {
				return a.DW < b.DW
			}
		} else {
			if a.DW != b.DW {
				return a.DW < b.DW
			}
			if a.DH != b.DH {
				return a.DH < b.DH
			}
		}

		if a.Font != b.Font {
			return a.Font < b.Font
		}
		return a.Scale < b.Scale
	})
}

func distance(target int, actual int) int {
	if target == 0 {
		return 0
	}
	return  target -actual 
}
