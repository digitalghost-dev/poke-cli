package tcg

import (
	"strings"
	"testing"
)

func TestCountryBarChart_Empty(t *testing.T) {
	result := CountryBarChart([]CountryStats{}, 80)
	if result != "" {
		t.Errorf("expected empty string for empty input, got %q", result)
	}
}

func TestCountryBarChart_SingleEntry(t *testing.T) {
	stats := []CountryStats{
		{Country: "USA", Total: 10},
	}
	result := CountryBarChart(stats, 80)
	if result == "" {
		t.Error("expected non-empty output for single entry")
	}
	if !strings.Contains(result, "USA") {
		t.Error("expected output to contain country name")
	}
	if !strings.Contains(result, "10") {
		t.Error("expected output to contain count")
	}
}

func TestCountryBarChart_SortsDescending(t *testing.T) {
	stats := []CountryStats{
		{Country: "France", Total: 5},
		{Country: "USA", Total: 20},
		{Country: "Japan", Total: 10},
	}
	result := CountryBarChart(stats, 80)
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

func TestCountryBarChart_TopNineWithOther(t *testing.T) {
	stats := []CountryStats{
		{Country: "A", Total: 100},
		{Country: "B", Total: 90},
		{Country: "C", Total: 80},
		{Country: "D", Total: 70},
		{Country: "E", Total: 60},
		{Country: "F", Total: 50},
		{Country: "G", Total: 40},
		{Country: "H", Total: 30},
		{Country: "I", Total: 20},
		{Country: "J", Total: 10},
		{Country: "K", Total: 5},
	}
	result := CountryBarChart(stats, 80)
	lines := strings.Split(strings.TrimRight(result, "\n"), "\n")

	if len(lines) != 10 {
		t.Fatalf("expected 10 lines (9 + Other), got %d", len(lines))
	}
	if !strings.Contains(result, "Other") {
		t.Error("expected 'Other' entry for entries beyond top 9")
	}
}

func TestCountryBarChart_ExactlyNine(t *testing.T) {
	stats := make([]CountryStats, 9)
	for i := range stats {
		stats[i] = CountryStats{Country: "X", Total: i + 1}
	}
	result := CountryBarChart(stats, 80)
	lines := strings.Split(strings.TrimRight(result, "\n"), "\n")

	if len(lines) != 9 {
		t.Fatalf("expected 9 lines, got %d", len(lines))
	}
	if strings.Contains(result, "Other") {
		t.Error("should not have 'Other' entry when exactly 9 countries")
	}
}

func TestCountryBarChart_DoesNotMutateInput(t *testing.T) {
	stats := []CountryStats{
		{Country: "France", Total: 5},
		{Country: "USA", Total: 20},
		{Country: "Japan", Total: 10},
	}
	original := make([]CountryStats, len(stats))
	copy(original, stats)

	CountryBarChart(stats, 80)

	for i, s := range stats {
		if s != original[i] {
			t.Errorf("input was mutated at index %d: got %v, want %v", i, s, original[i])
		}
	}
}

func TestCountryBarChart_NarrowWidth(t *testing.T) {
	stats := []CountryStats{
		{Country: "USA", Total: 10},
	}
	result := CountryBarChart(stats, 5)
	if result == "" {
		t.Error("expected non-empty output even for very narrow width")
	}
}

func TestCountryBarChart_OtherExceedsTopCountry(t *testing.T) {
	// Regression: "Other" total can exceed any individual country's total,
	// which previously caused a negative strings.Repeat count panic.
	stats := []CountryStats{
		{Country: "A", Total: 10},
		{Country: "B", Total: 9},
		{Country: "C", Total: 8},
		{Country: "D", Total: 7},
		{Country: "E", Total: 6},
		{Country: "F", Total: 5},
		{Country: "G", Total: 4},
		{Country: "H", Total: 3},
		{Country: "I", Total: 2},
		{Country: "J", Total: 50},
		{Country: "K", Total: 50},
	}
	result := CountryBarChart(stats, 80)
	if result == "" {
		t.Error("expected non-empty output")
	}
}
