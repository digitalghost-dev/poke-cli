package shell

import (
	"strings"
	"testing"
)

func TestBarChart_Empty(t *testing.T) {
	result := BarChart([]Tally{}, 80, 20)
	if result != "" {
		t.Errorf("expected empty string for empty input, got %q", result)
	}
}

func TestBarChart_AllZeroTotals(t *testing.T) {
	items := []Tally{
		{Label: "USA", Count: 0},
		{Label: "Japan", Count: 0},
	}

	result := BarChart(items, 80, 20)
	if result == "" {
		t.Error("expected non-empty output for non-empty input")
	}
	if !strings.Contains(result, "USA") {
		t.Error("expected output to contain label")
	}
}

func TestBarChart_SingleEntry(t *testing.T) {
	items := []Tally{
		{Label: "USA", Count: 10},
	}
	result := BarChart(items, 80, 20)
	if result == "" {
		t.Error("expected non-empty output for single entry")
	}
	if !strings.Contains(result, "USA") {
		t.Error("expected output to contain label")
	}
	if !strings.Contains(result, "10") {
		t.Error("expected output to contain count")
	}
}

func TestBarChart_SortsDescending(t *testing.T) {
	items := []Tally{
		{Label: "France", Count: 5},
		{Label: "USA", Count: 20},
		{Label: "Japan", Count: 10},
	}
	result := BarChart(items, 80, 20)
	lines := strings.Split(strings.TrimRight(result, "\n"), "\n")

	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	if !strings.HasPrefix(strings.TrimSpace(lines[0]), "USA") {
		t.Errorf("expected USA first (highest), got %q", lines[0])
	}
	if !strings.HasPrefix(strings.TrimSpace(lines[2]), "France") {
		t.Errorf("expected France last (lowest), got %q", lines[2])
	}
}

func TestBarChart_TopNineWithOther(t *testing.T) {
	items := []Tally{
		{Label: "A", Count: 100},
		{Label: "B", Count: 90},
		{Label: "C", Count: 80},
		{Label: "D", Count: 70},
		{Label: "E", Count: 60},
		{Label: "F", Count: 50},
		{Label: "G", Count: 40},
		{Label: "H", Count: 30},
		{Label: "I", Count: 20},
		{Label: "J", Count: 10},
		{Label: "K", Count: 5},
	}
	result := BarChart(items, 80, 20)
	lines := strings.Split(strings.TrimRight(result, "\n"), "\n")

	if len(lines) != 10 {
		t.Fatalf("expected 10 lines (9 + Other), got %d", len(lines))
	}
	if !strings.Contains(result, "Other") {
		t.Error("expected 'Other' entry for entries beyond top 9")
	}
}

func TestBarChart_ExactlyNine(t *testing.T) {
	items := make([]Tally, 9)
	for i := range items {
		items[i] = Tally{Label: "X", Count: i + 1}
	}
	result := BarChart(items, 80, 20)
	lines := strings.Split(strings.TrimRight(result, "\n"), "\n")

	if len(lines) != 9 {
		t.Fatalf("expected 9 lines, got %d", len(lines))
	}
	if strings.Contains(result, "Other") {
		t.Error("should not have 'Other' entry when exactly 9 items")
	}
}

func TestBarChart_DoesNotMutateInput(t *testing.T) {
	items := []Tally{
		{Label: "France", Count: 5},
		{Label: "USA", Count: 20},
		{Label: "Japan", Count: 10},
	}
	original := make([]Tally, len(items))
	copy(original, items)

	BarChart(items, 80, 20)

	for i, s := range items {
		if s != original[i] {
			t.Errorf("input was mutated at index %d: got %v, want %v", i, s, original[i])
		}
	}
}

func TestBarChart_NarrowWidth(t *testing.T) {
	items := []Tally{
		{Label: "USA", Count: 10},
	}
	result := BarChart(items, 5, 20)
	if result == "" {
		t.Error("expected non-empty output even for very narrow width")
	}
}

func TestBarChart_OtherExceedsTopEntry(t *testing.T) {
	items := []Tally{
		{Label: "A", Count: 10},
		{Label: "B", Count: 9},
		{Label: "C", Count: 8},
		{Label: "D", Count: 7},
		{Label: "E", Count: 6},
		{Label: "F", Count: 5},
		{Label: "G", Count: 4},
		{Label: "H", Count: 3},
		{Label: "I", Count: 2},
		{Label: "J", Count: 50},
		{Label: "K", Count: 50},
	}
	result := BarChart(items, 80, 20)
	if result == "" {
		t.Error("expected non-empty output")
	}
}

func TestBarChart_MinOneBlock(t *testing.T) {
	items := []Tally{
		{Label: "Big", Count: 428},
		{Label: "Tiny", Count: 1},
	}
	result := BarChart(items, 80, 20)
	lines := strings.Split(strings.TrimRight(result, "\n"), "\n")
	tinyLine := lines[1]
	if !strings.Contains(tinyLine, "█") {
		t.Error("expected at least one block for non-zero small total")
	}
}
