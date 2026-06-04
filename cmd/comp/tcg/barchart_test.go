package tcg

import (
	"strings"
	"testing"
)

func TestBarChart_Empty(t *testing.T) {
	result := barChart([]barChartItem{}, 80, 20)
	if result != "" {
		t.Errorf("expected empty string for empty input, got %q", result)
	}
}

func TestBarChart_AllZeroTotals(t *testing.T) {
	items := []barChartItem{
		{Label: "USA", Total: 0},
		{Label: "Japan", Total: 0},
	}
	// should not panic (division by zero)
	result := barChart(items, 80, 20)
	if result == "" {
		t.Error("expected non-empty output for non-empty input")
	}
	if !strings.Contains(result, "USA") {
		t.Error("expected output to contain label")
	}
}

func TestBarChart_SingleEntry(t *testing.T) {
	items := []barChartItem{
		{Label: "USA", Total: 10},
	}
	result := barChart(items, 80, 20)
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
	items := []barChartItem{
		{Label: "France", Total: 5},
		{Label: "USA", Total: 20},
		{Label: "Japan", Total: 10},
	}
	result := barChart(items, 80, 20)
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
	items := []barChartItem{
		{Label: "A", Total: 100},
		{Label: "B", Total: 90},
		{Label: "C", Total: 80},
		{Label: "D", Total: 70},
		{Label: "E", Total: 60},
		{Label: "F", Total: 50},
		{Label: "G", Total: 40},
		{Label: "H", Total: 30},
		{Label: "I", Total: 20},
		{Label: "J", Total: 10},
		{Label: "K", Total: 5},
	}
	result := barChart(items, 80, 20)
	lines := strings.Split(strings.TrimRight(result, "\n"), "\n")

	if len(lines) != 10 {
		t.Fatalf("expected 10 lines (9 + Other), got %d", len(lines))
	}
	if !strings.Contains(result, "Other") {
		t.Error("expected 'Other' entry for entries beyond top 9")
	}
}

func TestBarChart_ExactlyNine(t *testing.T) {
	items := make([]barChartItem, 9)
	for i := range items {
		items[i] = barChartItem{Label: "X", Total: i + 1}
	}
	result := barChart(items, 80, 20)
	lines := strings.Split(strings.TrimRight(result, "\n"), "\n")

	if len(lines) != 9 {
		t.Fatalf("expected 9 lines, got %d", len(lines))
	}
	if strings.Contains(result, "Other") {
		t.Error("should not have 'Other' entry when exactly 9 items")
	}
}

func TestBarChart_DoesNotMutateInput(t *testing.T) {
	items := []barChartItem{
		{Label: "France", Total: 5},
		{Label: "USA", Total: 20},
		{Label: "Japan", Total: 10},
	}
	original := make([]barChartItem, len(items))
	copy(original, items)

	barChart(items, 80, 20)

	for i, s := range items {
		if s != original[i] {
			t.Errorf("input was mutated at index %d: got %v, want %v", i, s, original[i])
		}
	}
}

func TestBarChart_NarrowWidth(t *testing.T) {
	items := []barChartItem{
		{Label: "USA", Total: 10},
	}
	result := barChart(items, 5, 20)
	if result == "" {
		t.Error("expected non-empty output even for very narrow width")
	}
}

func TestBarChart_OtherExceedsTopEntry(t *testing.T) {
	// Regression: "Other" total can exceed any individual entry's total,
	// which previously caused a negative strings.Repeat count panic.
	items := []barChartItem{
		{Label: "A", Total: 10},
		{Label: "B", Total: 9},
		{Label: "C", Total: 8},
		{Label: "D", Total: 7},
		{Label: "E", Total: 6},
		{Label: "F", Total: 5},
		{Label: "G", Total: 4},
		{Label: "H", Total: 3},
		{Label: "I", Total: 2},
		{Label: "J", Total: 50},
		{Label: "K", Total: 50},
	}
	result := barChart(items, 80, 20)
	if result == "" {
		t.Error("expected non-empty output")
	}
}

func TestBarChart_MinOneBlock(t *testing.T) {
	// Small totals relative to maxVal should still render at least one block.
	items := []barChartItem{
		{Label: "Big", Total: 428},
		{Label: "Tiny", Total: 1},
	}
	result := barChart(items, 80, 20)
	lines := strings.Split(strings.TrimRight(result, "\n"), "\n")
	tinyLine := lines[1]
	if !strings.Contains(tinyLine, "█") {
		t.Error("expected at least one block for non-zero small total")
	}
}
