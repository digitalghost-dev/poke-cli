package champions

import (
	"strings"
	"testing"
)

func TestJoinOrDash(t *testing.T) {
	tests := []struct {
		name string
		in   []string
		want string
	}{
		{"empty", nil, "-"},
		{"single", []string{"Miraidon"}, "Miraidon"},
		{"multiple", []string{"Miraidon", "Flutter Mane"}, "Miraidon, Flutter Mane"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := joinOrDash(tt.in); got != tt.want {
				t.Errorf("joinOrDash(%v) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestTeamCore(t *testing.T) {
	tests := []struct {
		name string
		in   []string
		want string
	}{
		{"empty", nil, "-"},
		{"three or fewer", []string{"A", "B", "C"}, "A, B, C"},
		{"more than three", []string{"A", "B", "C", "D", "E"}, "A, B, C +2"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := teamCore(tt.in); got != tt.want {
				t.Errorf("teamCore(%v) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestNewOverviewTable(t *testing.T) {
	rows := testCompInfo()
	tbl := newOverviewTable(rows, 40)
	if len(tbl.Rows()) != len(rows) {
		t.Fatalf("expected %d rows, got %d", len(rows), len(tbl.Rows()))
	}
	if tbl.Rows()[0][0] != "Miraidon" {
		t.Errorf("expected first row Miraidon, got %q", tbl.Rows()[0][0])
	}
}

func TestSelectedCompInfo(t *testing.T) {
	rows := testCompInfo()

	if got := selectedCompInfo(newOverviewTable(nil, 40), nil); got.Pokemon != "" {
		t.Errorf("expected zero compInfoRow for empty rows, got %+v", got)
	}

	tbl := newOverviewTable(rows, 40)
	if got := selectedCompInfo(tbl, rows); got.Pokemon != "Miraidon" {
		t.Errorf("expected first Pokémon selected, got %q", got.Pokemon)
	}
}

func TestRenderOverview(t *testing.T) {
	if got := renderOverview(newOverviewTable(nil, 40), nil, 120); got != "No data available" {
		t.Errorf("expected empty-state message, got %q", got)
	}

	rows := testCompInfo()
	out := renderOverview(newOverviewTable(rows, 40), rows, 120)
	for _, want := range []string{"Miraidon", "Common Moves", "Common Items", "Common Abilities", "Common Teammates", "Protect", "Flutter Mane"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected overview to contain %q", want)
		}
	}
}

func TestRenderPokemonDetail(t *testing.T) {
	withLink := testCompInfo()[0]
	out := renderPokemonDetail(withLink, 90)
	for _, want := range []string{"Miraidon", "Protect", "Hadron Engine", "Choice Specs", "Flutter Mane", "Link", "https://example.com/pokemon/1", "90.5%", "100.0%"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected detail to contain %q, got:\n%s", want, out)
		}
	}

	noLink := testCompInfo()[1]
	out2 := renderPokemonDetail(noLink, 90)
	if strings.Contains(out2, "Link") {
		t.Error("expected no Link line when WebURL is empty")
	}
	if !strings.Contains(out2, "Common Items") {
		t.Error("expected Common Items heading even with no items")
	}
}

func TestRenderStatColumn(t *testing.T) {
	empty := renderStatColumn("Common Items", nil, 30)
	if !strings.Contains(empty, "Common Items") || !strings.Contains(empty, "-") {
		t.Errorf("expected title and dash placeholder, got %q", empty)
	}

	stats := []commonStat{{Name: "Protect", UsagePercent: 90.5}, {Name: "Tailwind", UsagePercent: 10.25}}
	out := renderStatColumn("Common Moves", stats, 30)
	for _, want := range []string{"Common Moves", "Protect", "90.5%", "Tailwind", "10.2%"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected stat column to contain %q, got:\n%s", want, out)
		}
	}
}

func TestStatLine(t *testing.T) {
	line := statLine(commonStat{Name: "Shadow Sneak", UsagePercent: 89.758}, 24)
	if !strings.Contains(line, "Shadow Sneak") || !strings.Contains(line, "89.8%") {
		t.Errorf("unexpected stat line: %q", line)
	}
}

func TestTruncateName(t *testing.T) {
	tests := []struct {
		name  string
		width int
		want  string
	}{
		{"Protect", 10, "Protect"},
		{"King's Shield", 6, "King'…"},
		{"Charizard-Mega-Y", 1, "C"},
		{"", 5, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := truncateName(tt.name, tt.width); got != tt.want {
				t.Errorf("truncateName(%q, %d) = %q, want %q", tt.name, tt.width, got, tt.want)
			}
		})
	}
}

func TestTeamColumns(t *testing.T) {
	for _, width := range []int{200, 120, 40, 10} {
		cols := teamColumns(width)
		if len(cols) != 5 {
			t.Fatalf("width %d: expected 5 columns, got %d", width, len(cols))
		}
		titles := []string{"Player", "Record", "Tournament", "Archetypes", "Core"}
		for i, c := range cols {
			if c.Title != titles[i] {
				t.Errorf("width %d: column %d title = %q, want %q", width, i, c.Title, titles[i])
			}
			if c.Width <= 0 {
				t.Errorf("width %d: column %q has non-positive width %d", width, c.Title, c.Width)
			}
		}
	}
}

func TestTableWidth(t *testing.T) {
	cols := teamColumns(120)
	got := tableWidth(cols)
	want := len(cols) * 2
	for _, c := range cols {
		want += c.Width
	}
	if got != want {
		t.Errorf("tableWidth = %d, want %d", got, want)
	}
}

func TestSelectedTeam(t *testing.T) {
	teams := testTeams()

	if got := selectedTeam(newTeamsTable(nil, 120, 40), nil); got.Player != "" {
		t.Errorf("expected zero teamRow for empty teams, got %+v", got)
	}

	tbl := newTeamsTable(teams, 120, 40)
	if got := selectedTeam(tbl, teams); got.Player != "Alice" {
		t.Errorf("expected first team selected, got %q", got.Player)
	}
}

func TestNewTeamsTable(t *testing.T) {
	teams := testTeams()
	tbl := newTeamsTable(teams, 120, 40)
	if len(tbl.Rows()) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(tbl.Rows()))
	}
	first := tbl.Rows()[0]
	if first[0] != "Alice" || first[1] != "7-1" {
		t.Errorf("unexpected first row: %v", first)
	}
	if !strings.Contains(first[4], "+3") {
		t.Errorf("expected core column to truncate a 6-mon team, got %q", first[4])
	}
}

func TestRenderTeamsTable(t *testing.T) {
	if got := renderTeamsTable(newTeamsTable(nil, 120, 40), nil, 120); got != "No data available" {
		t.Errorf("expected empty-state message, got %q", got)
	}

	teams := testTeams()
	out := renderTeamsTable(newTeamsTable(teams, 120, 40), teams, 120)
	for _, want := range []string{"Alice", "Selected Team", "Worlds 2026"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected rendered table to contain %q", want)
		}
	}
}

func TestRenderTeamDetail(t *testing.T) {
	withLink := testTeams()[0]
	out := renderTeamDetail(withLink, 120)
	for _, want := range []string{"Alice", "7-1", "Worlds 2026", "Hyper Offense", "Miraidon", "Link", "https://example.com/team/1"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected detail to contain %q, got:\n%s", want, out)
		}
	}

	noLink := testTeams()[1]
	out2 := renderTeamDetail(noLink, 120)
	if strings.Contains(out2, "Link") {
		t.Error("expected no Link line when WebURL is empty")
	}
	if !strings.Contains(out2, "-") {
		t.Error("expected dash placeholder for empty archetypes")
	}
}

func TestDetailLine(t *testing.T) {
	line := detailLine("Tournament", "Worlds 2026", 80)
	if !strings.Contains(line, "Tournament") || !strings.Contains(line, "Worlds 2026") {
		t.Errorf("unexpected detail line: %q", line)
	}

	empty := detailLine("Archetypes", "", 80)
	if !strings.HasSuffix(empty, "-") {
		t.Errorf("expected dash for empty value, got %q", empty)
	}

	wrapped := detailLine("Team", strings.Repeat("Pokemon ", 30), 40)
	if !strings.Contains(wrapped, "\n") {
		t.Error("expected long value to wrap onto multiple lines")
	}
}

func TestWrapWords(t *testing.T) {
	if wrapWords("", 10) != nil {
		t.Error("expected nil for empty input")
	}

	lines := wrapWords("one two three four five", 9)
	if len(lines) < 2 {
		t.Errorf("expected wrapping into multiple lines, got %v", lines)
	}
	for _, line := range lines {
		if strings.HasPrefix(line, " ") || strings.HasSuffix(line, " ") {
			t.Errorf("expected trimmed lines, got %q", line)
		}
	}

	long := wrapWords("supercalifragilisticexpialidocious", 5)
	if len(long) < 2 {
		t.Errorf("expected a long single word to be split, got %v", long)
	}
}

func TestSplitLongWord(t *testing.T) {
	if got := splitLongWord("word", 0); len(got) != 1 || got[0] != "word" {
		t.Errorf("expected non-positive width to return the word unchanged, got %v", got)
	}

	got := splitLongWord("abcdefg", 3)
	want := []string{"abc", "def", "g"}
	if len(got) != len(want) {
		t.Fatalf("splitLongWord = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("chunk %d = %q, want %q", i, got[i], want[i])
		}
	}
}

func TestNewSpeedTable(t *testing.T) {
	rows := testSpeedTiers()
	tbl := newSpeedTable(rows, 40)
	if len(tbl.Rows()) != len(rows) {
		t.Fatalf("expected %d rows, got %d", len(rows), len(tbl.Rows()))
	}
	first := tbl.Rows()[0]
	if first[0] != "1" || first[1] != "Mega Aerodactyl" || first[2] != "150" || first[5] != "333" {
		t.Errorf("unexpected first row: %v", first)
	}
}

func TestSelectedSpeedTier(t *testing.T) {
	rows := testSpeedTiers()

	if got := selectedSpeedTier(newSpeedTable(nil, 40), nil); got.Pokemon != "" {
		t.Errorf("expected zero speedTierRow for empty rows, got %+v", got)
	}

	tbl := newSpeedTable(rows, 40)
	if got := selectedSpeedTier(tbl, rows); got.Pokemon != "Mega Aerodactyl" {
		t.Errorf("expected first row selected, got %q", got.Pokemon)
	}
}

func TestRenderSpeedTiers(t *testing.T) {
	if got := renderSpeedTiers(newSpeedTable(nil, 40), nil); got != "No data available" {
		t.Errorf("expected empty-state message, got %q", got)
	}

	rows := testSpeedTiers()
	out := renderSpeedTiers(newSpeedTable(rows, 40), rows)
	for _, want := range []string{"level 50", "Mega Aerodactyl", "Base Speed", "Max + Scarf", "150", "333"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected speed tiers view to contain %q", want)
		}
	}
}

func TestRenderSpeedDetail(t *testing.T) {
	out := renderSpeedDetail(testSpeedTiers()[0])
	for _, want := range []string{"Selected Pokémon", "Mega Aerodactyl", "Base Speed", "150", "Min (0 EV -Spe)", "153", "Max (252 EV +Spe)", "222", "Max + Scarf", "333"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected detail to contain %q, got:\n%s", want, out)
		}
	}
}

func TestSpeedStatLine(t *testing.T) {
	line := speedStatLine("Base Speed", 150)
	if !strings.Contains(line, "Base Speed") || !strings.Contains(line, "150") {
		t.Errorf("unexpected stat line: %q", line)
	}
}
