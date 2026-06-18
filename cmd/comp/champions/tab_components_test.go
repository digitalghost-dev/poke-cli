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
