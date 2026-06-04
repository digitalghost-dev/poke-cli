package tcg

import (
	"strings"
	"testing"
)

func sampleRows() []standingRows {
	return []standingRows{
		{Rank: 1, Name: "Ash Ketchum", Points: 47, Record: "15 - 1 - 2", OppWinPct: "58.10%", OppOppWinPct: "60.56%", Deck: "gardevoir", PlayerCountry: "United States", ISOCode: "US"},
		{Rank: 2, Name: "Misty Williams", Points: 44, Record: "14 - 2 - 2", OppWinPct: "68.56%", OppOppWinPct: "61.67%", Deck: "dragapult/dusknoir", PlayerCountry: "Japan", ISOCode: "JP"},
		{Rank: 3, Name: "Brock Harrison", Points: 41, Record: "13 - 2 - 2", OppWinPct: "69.01%", OppOppWinPct: "63.78%", Deck: "dragapult", PlayerCountry: "United Kingdom", ISOCode: "GB"},
	}
}

func TestStandingsTable_ContainsHeaders(t *testing.T) {
	m := standingsTable(sampleRows(), 120, 40)
	view := m.View()
	for _, header := range []string{"Rank", "Name", "Points", "Record", "OPW%", "OOPW%", "Deck", "Country"} {
		if !strings.Contains(view, header) {
			t.Errorf("expected table to contain header %q", header)
		}
	}
}

func TestStandingsTable_ContainsRowData(t *testing.T) {
	m := standingsTable(sampleRows(), 120, 40)
	view := m.View()
	for _, s := range []string{"Ash Ketchum", "gardevoir", "United States", "47"} {
		if !strings.Contains(view, s) {
			t.Errorf("expected table to contain %q", s)
		}
	}
}

func TestStandingsTable_EmptyRows(t *testing.T) {
	m := standingsTable([]standingRows{}, 120, 40)
	view := m.View()
	if view == "" {
		t.Fatal("expected non-empty view even with no rows")
	}
}

func TestStandingsTable_NarrowWidth(t *testing.T) {
	m := standingsTable(sampleRows(), 10, 40)
	if m.View() == "" {
		t.Fatal("expected non-empty view for narrow width")
	}
}

func TestStandingsTable_ShortHeight(t *testing.T) {
	m := standingsTable(sampleRows(), 120, 5)
	if m.View() == "" {
		t.Fatal("expected non-empty view for short height")
	}
}
