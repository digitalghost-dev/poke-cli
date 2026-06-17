// Tab rendering helpers for the Champions dashboard.
// The Bubble Tea lifecycle stays in dashboard.go; this file builds the per-tab views.

package champions

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/table"
	"github.com/digitalghost-dev/poke-cli/cmd/comp/shell"
	"github.com/digitalghost-dev/poke-cli/styling"
)

// Top Teams section
func newTeamsTable(teams []teamRow, width, height int) table.Model {
	columns := teamColumns(width)
	rows := make([]table.Row, 0, len(teams))

	for _, team := range teams {
		rows = append(rows, table.Row{
			team.Player,
			team.Record,
			team.Tournament,
			joinOrDash(team.Archetypes),
			teamCore(team.Pokemon),
		})
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(max(height-22, 5)),
		table.WithWidth(tableWidth(columns)),
	)
	t.SetStyles(shell.TableStyles())
	return t
}

func renderTeamsTable(teamsTable table.Model, teams []teamRow, width int) string {
	if len(teamsTable.Rows()) == 0 {
		return "No data available"
	}

	detail := renderTeamDetail(selectedTeam(teamsTable, teams), width)
	return teamsTable.View() + "\n\n" + detail
}

func teamColumns(width int) []table.Column {
	const recordWidth = 8

	separatorWidth := 5 * 2
	availableWidth := max(width-separatorWidth-recordWidth, 40)

	playerWidth := min(max(availableWidth*20/100, 8), 24)
	tournamentWidth := min(max(availableWidth*28/100, 12), 40)
	archetypesWidth := min(max(availableWidth*20/100, 8), 24)
	teamWidth := max(availableWidth-playerWidth-tournamentWidth-archetypesWidth, 8)

	return []table.Column{
		{Title: "Player", Width: playerWidth},
		{Title: "Record", Width: recordWidth},
		{Title: "Tournament", Width: tournamentWidth},
		{Title: "Archetypes", Width: archetypesWidth},
		{Title: "Core", Width: teamWidth},
	}
}

func selectedTeam(teamsTable table.Model, teams []teamRow) teamRow {
	if len(teams) == 0 {
		return teamRow{}
	}

	idx := min(max(teamsTable.Cursor(), 0), len(teams)-1)
	return teams[idx]
}

func renderTeamDetail(team teamRow, width int) string {
	var b strings.Builder

	title := fmt.Sprintf("%s (%s)", team.Player, team.Record)
	b.WriteString(styling.Yellow.Render("Selected Team"))
	b.WriteString("\n")
	b.WriteString(styling.StyleBold.Render(title))
	b.WriteString("\n")
	b.WriteString(detailLine("Tournament", team.Tournament, width))
	b.WriteString("\n")
	b.WriteString(detailLine("Archetypes", joinOrDash(team.Archetypes), width))
	b.WriteString("\n")
	b.WriteString(detailLine("Team", joinOrDash(team.Pokemon), width))
	if team.WebURL != "" {
		b.WriteString("\n")
		b.WriteString(detailLine("Link", team.WebURL, width))
	}

	return b.String()
}

func tableWidth(columns []table.Column) int {
	width := len(columns) * 2
	for _, c := range columns {
		width += c.Width
	}
	return width
}

func joinOrDash(values []string) string {
	if len(values) == 0 {
		return "-"
	}
	return strings.Join(values, ", ")
}

func teamCore(pokemon []string) string {
	if len(pokemon) <= 3 {
		return joinOrDash(pokemon)
	}
	return strings.Join(pokemon[:3], ", ") + fmt.Sprintf(" +%d", len(pokemon)-3)
}

func detailLine(label, value string, width int) string {
	prefix := styling.StyleBold.Render(label + ": ")
	plainPrefixWidth := len(label) + 2
	lineWidth := max(width-plainPrefixWidth, 20)
	lines := wrapWords(value, lineWidth)
	if len(lines) == 0 {
		return prefix + "-"
	}

	var b strings.Builder
	b.WriteString(prefix)
	b.WriteString(lines[0])
	for _, line := range lines[1:] {
		b.WriteString("\n")
		b.WriteString(strings.Repeat(" ", plainPrefixWidth))
		b.WriteString(line)
	}
	return b.String()
}

func wrapWords(value string, width int) []string {
	words := strings.Fields(value)
	if len(words) == 0 {
		return nil
	}

	lines := make([]string, 0, 2)
	current := words[0]
	for _, word := range words[1:] {
		if len(current) > width {
			lines = append(lines, splitLongWord(current, width)...)
			current = word
			continue
		}
		if len(current)+1+len(word) > width {
			lines = append(lines, current)
			current = word
			continue
		}
		current += " " + word
	}
	if len(current) > width {
		lines = append(lines, splitLongWord(current, width)...)
	} else {
		lines = append(lines, current)
	}
	return lines
}

func splitLongWord(word string, width int) []string {
	if width <= 0 {
		return []string{word}
	}

	var lines []string
	for len(word) > width {
		lines = append(lines, word[:width])
		word = word[width:]
	}
	if word != "" {
		lines = append(lines, word)
	}
	return lines
}

// Speed Tiers section
