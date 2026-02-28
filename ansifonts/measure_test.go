package ansifonts

import "testing"

func TestStripANSI(t *testing.T) {
	input := "A\x1b[31mB\x1b[0mC"
	got := StripANSI(input)
	if got != "ABC" {
		t.Fatalf("StripANSI() = %q, want %q", got, "ABC")
	}
}

func TestMeasureBlockWithANSI(t *testing.T) {
	lines := []string{
		"hi",
		"",
		"\x1b[31mcolor\x1b[0m",
	}

	w, h := MeasureBlock(lines)
	if w != 5 {
		t.Fatalf("MeasureBlock() width = %d, want %d", w, 5)
	}
	if h != 3 {
		t.Fatalf("MeasureBlock() height = %d, want %d", h, 3)
	}
}
