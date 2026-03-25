package tcg

import (
	"strings"
	"testing"
)

func TestOverviewContent(t *testing.T) {
	tests := []struct {
		name         string
		flag         string
		tournament   string
		tType        string
		tDate        string
		winner       string
		winningDeck  string
		totalPlayers int
		contentWidth int
		contains     []string
	}{
		{
			name:         "all fields present in output",
			flag:         "🇺🇸",
			tournament:   "Dallas",
			tType:        "Regional",
			tDate:        "January 10-12, 2025",
			winner:       "Ash Ketchum",
			winningDeck:  "gardevoir",
			totalPlayers: 1024,
			contentWidth: 80,
			contains:     []string{"Dallas", "Regional", "January 10-12, 2025", "Ash Ketchum", "gardevoir", "1,024", "Total Players", "Winner", "Winning Deck"},
		},
		{
			name:         "large player count formatted with commas",
			totalPlayers: 1000000,
			contentWidth: 80,
			contains:     []string{"1,000,000"},
		},
		{
			name:         "empty strings do not panic",
			contentWidth: 80,
			contains:     []string{"Total Players", "Winner", "Winning Deck"},
		},
		{
			name:         "narrow content width does not panic",
			tournament:   "Sydney",
			winner:       "Trainer Red",
			totalPlayers: 500,
			contentWidth: 10,
			contains:     []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := overviewContent(tt.flag, tt.tournament, tt.tType, tt.tDate, tt.winner, tt.winningDeck, tt.totalPlayers, tt.contentWidth)
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

func TestFormatInt(t *testing.T) {
	tests := []struct {
		n    int
		want string
	}{
		{0, "0"},
		{999, "999"},
		{1000, "1,000"},
		{4010, "4,010"},
		{10000, "10,000"},
		{1000000, "1,000,000"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got := formatInt(tt.n)
			if got != tt.want {
				t.Errorf("formatInt(%d) = %q, want %q", tt.n, got, tt.want)
			}
		})
	}
}
