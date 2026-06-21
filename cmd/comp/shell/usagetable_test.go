package shell

import (
	"strings"
	"testing"
)

func TestNewUsageTable_SortedRanked(t *testing.T) {
	f := Frequency{NameHeader: "Pokémon", CountHeader: "Teams", Items: []Tally{
		{Label: "Incineroar", Count: 5},
		{Label: "Flutter Mane", Count: 8},
		{Label: "Rillaboom", Count: 2},
	}}
	tbl := newUsageTable(f, 10, 120, 40)
	rows := tbl.Rows()

	if len(rows) != 3 {
		t.Fatalf("expected 3 rows, got %d", len(rows))
	}
	if rows[0][0] != "1" || rows[0][1] != "Flutter Mane" || rows[0][2] != "8" {
		t.Errorf("expected rank-1 row to be Flutter Mane/8, got %v", rows[0])
	}
	if rows[2][1] != "Rillaboom" {
		t.Errorf("expected lowest count last, got %q", rows[2][1])
	}
	if !strings.Contains(rows[0][3], "80%") || !strings.Contains(rows[0][3], "█") {
		t.Errorf("expected rank-1 share to show a full-ish bar and 80%%, got %q", rows[0][3])
	}
}

func TestShareCell(t *testing.T) {
	full := shareCell(20, 20, 11)
	if strings.Count(full, "█") != 11 || !strings.Contains(full, "100%") {
		t.Errorf("expected a full 11-block bar at 100%%, got %q", full)
	}

	half := shareCell(10, 20, 11)
	if strings.Count(half, "█") != 5 || !strings.Contains(half, "50%") {
		t.Errorf("expected ~half bar at 50%% (absolute), got %q", half)
	}

	none := shareCell(0, 20, 11)
	if strings.Contains(none, "█") || !strings.Contains(none, "0%") {
		t.Errorf("expected no filled blocks at 0%%, got %q", none)
	}

	over := shareCell(30, 20, 11)
	if strings.Count(over, "█") != 11 {
		t.Errorf("expected filled to clamp at barWidth when count>total, got %q", over)
	}
	if !strings.Contains(over, "100%") {
		t.Errorf("expected percentage to clamp at 100%% when count>total, got %q", over)
	}

	_ = shareCell(5, 0, 11)
}
