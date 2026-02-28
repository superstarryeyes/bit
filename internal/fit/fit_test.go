package fit

import (
	"reflect"
	"testing"
)

func TestSortCandidatesWidthPriority(t *testing.T) {
	candidates := []Candidate{
		{Font: "b", Scale: 1, DW: 1, DH: 4},
		{Font: "a", Scale: 2, DW: 1, DH: 2},
		{Font: "a", Scale: 1, DW: 1, DH: 2},
		{Font: "c", Scale: 1, DW: 0, DH: 7},
	}

	sortCandidates(candidates, PriorityWidth)

	expected := []Candidate{
		{Font: "c", Scale: 1, DW: 0, DH: 7},
		{Font: "a", Scale: 1, DW: 1, DH: 2},
		{Font: "a", Scale: 2, DW: 1, DH: 2},
		{Font: "b", Scale: 1, DW: 1, DH: 4},
	}

	if !reflect.DeepEqual(candidates, expected) {
		t.Fatalf("width priority order mismatch: got %#v", candidates)
	}
}

func TestSortCandidatesHeightPriority(t *testing.T) {
	candidates := []Candidate{
		{Font: "b", Scale: 1, DW: 4, DH: 1},
		{Font: "a", Scale: 2, DW: 2, DH: 1},
		{Font: "a", Scale: 1, DW: 2, DH: 1},
		{Font: "c", Scale: 1, DW: 7, DH: 0},
	}

	sortCandidates(candidates, PriorityHeight)

	expected := []Candidate{
		{Font: "c", Scale: 1, DW: 7, DH: 0},
		{Font: "a", Scale: 1, DW: 2, DH: 1},
		{Font: "a", Scale: 2, DW: 2, DH: 1},
		{Font: "b", Scale: 1, DW: 4, DH: 1},
	}

	if !reflect.DeepEqual(candidates, expected) {
		t.Fatalf("height priority order mismatch: got %#v", candidates)
	}
}
