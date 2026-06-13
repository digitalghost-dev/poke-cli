package vgc

import (
	"strings"
	"testing"

	"charm.land/lipgloss/v2"
)

func TestOverviewContent(t *testing.T) {
	tests := []struct {
		name         string
		tournament   string
		tType        string
		tDate        string
		winner       string
		winnerTeam   []string
		totalPlayers int
		contentWidth int
		contains     []string
	}{
		{
			name:         "all fields present in output",
			tournament:   "Indianapolis",
			tType:        "Regional",
			tDate:        "May 29-31, 2026",
			winner:       "Arsal Puri",
			winnerTeam:   []string{"Venusaur", "Landorus", "Iron Crown"},
			totalPlayers: 1013,
			contentWidth: 110,
			contains:     []string{"Indianapolis", "Regional", "May 29-31, 2026", "Arsal Puri", "Venusaur", "Iron Crown", "1,013", "Total Players", "Winner", "Winning Team"},
		},
		{
			name:         "large player count formatted with commas",
			totalPlayers: 1000000,
			contentWidth: 80,
			contains:     []string{"1,000,000"},
		},
		{
			name:         "empty team does not panic",
			contentWidth: 80,
			contains:     []string{"Total Players", "Winner", "Winning Team"},
		},
		{
			name:         "narrow content width does not panic",
			tournament:   "Turin",
			winner:       "Trainer Red",
			winnerTeam:   []string{"Miraidon"},
			totalPlayers: 500,
			contentWidth: 10,
			contains:     []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := overviewContent(tt.tournament, tt.tType, tt.tDate, tt.winner, tt.winnerTeam, tt.totalPlayers, tt.contentWidth, lipgloss.Color("#7D56F4"))
			if result == "" {
				t.Fatal("expected non-empty output")
			}
			for _, s := range tt.contains {
				if !strings.Contains(result, s) {
					t.Errorf("expected output to contain %q", s)
				}
			}
		})
	}
}

func TestTeamGrid(t *testing.T) {
	team := []string{"Venusaur", "Charizard", "Garchomp", "Incineroar", "Floette", "Sinistcha"}
	grid := teamGrid(team)
	lines := strings.Split(grid, "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 rows for a 6-member team, got %d: %q", len(lines), grid)
	}
	if !strings.Contains(lines[0], "Venusaur") || !strings.Contains(lines[0], "Incineroar") {
		t.Errorf("expected first row to pair Venusaur + Incineroar, got %q", lines[0])
	}
	if !strings.Contains(grid, "• ") {
		t.Error("expected bullets in the grid")
	}

	if teamGrid(nil) != "—" {
		t.Error("expected em dash placeholder for an empty team")
	}
}

func TestBaseName(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"Tornadus [Incarnate Forme]", "Tornadus"},
		{"Urshifu [Rapid Strike Style]", "Urshifu"},
		{"Venusaur", "Venusaur"},
		{"Ogerpon [Hearthflame Mask]", "Ogerpon"},
		{"", ""},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			if got := baseName(tt.in); got != tt.want {
				t.Errorf("baseName(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}
