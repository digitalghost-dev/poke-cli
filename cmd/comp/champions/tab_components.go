// Tab rendering helpers for the Champions dashboard.
// The Bubble Tea lifecycle stays in dashboard.go; this file builds the per-tab views.

package champions

import (
	"fmt"
	"strconv"
	"strings"

	"charm.land/bubbles/v2/table"
	"charm.land/lipgloss/v2"
	"github.com/digitalghost-dev/poke-cli/cmd/comp/shell"
	"github.com/digitalghost-dev/poke-cli/styling"
)

var captionStyle = lipgloss.NewStyle().Foreground(styling.Gray).Italic(true)

// Overview tab
func newOverviewTable(rows []compInfoRow, height int) table.Model {
	const nameWidth = 22
	columns := []table.Column{{Title: "Pokémon", Width: nameWidth}}

	trows := make([]table.Row, 0, len(rows))
	for _, row := range rows {
		trows = append(trows, table.Row{row.Pokemon})
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(trows),
		table.WithFocused(true),
		table.WithHeight(max(height-12, 5)),
		table.WithWidth(nameWidth+4),
	)
	t.SetStyles(shell.TableStyles())
	return t
}

func renderOverview(pokemonTable table.Model, rows []compInfoRow, width int) string {
	if len(rows) == 0 {
		return "No data available"
	}

	caption := captionStyle.Render("Select a Pokémon to see its most common moves, items, abilities, and teammates from recent Champions events.")

	detailWidth := max(width-pokemonTable.Width()-4, 40)
	detail := renderPokemonDetail(selectedCompInfo(pokemonTable, rows), detailWidth)
	body := lipgloss.JoinHorizontal(lipgloss.Top, pokemonTable.View(), "  ", detail)
	return caption + "\n\n" + body
}

func selectedCompInfo(pokemonTable table.Model, rows []compInfoRow) compInfoRow {
	if len(rows) == 0 {
		return compInfoRow{}
	}

	idx := min(max(pokemonTable.Cursor(), 0), len(rows)-1)
	return rows[idx]
}

func renderPokemonDetail(row compInfoRow, width int) string {
	colWidth := min(max((width-3)/2, 18), 34)

	moves := renderStatColumn("Common Moves", row.CommonMoves, colWidth)
	items := renderStatColumn("Common Items", row.CommonItems, colWidth)
	abilities := renderStatColumn("Common Abilities", row.CommonAbilities, colWidth)
	teammates := renderStatColumn("Common Teammates", row.CommonTeammates, colWidth)

	var b strings.Builder
	b.WriteString(styling.Yellow.Render(row.Pokemon))
	b.WriteString("\n\n")
	b.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, moves, "   ", items))
	b.WriteString("\n\n")
	b.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, abilities, "   ", teammates))
	if row.WebURL != "" {
		b.WriteString("\n\n")
		b.WriteString(detailLine("Link", row.WebURL, width))
	}
	return b.String()
}

func renderStatColumn(title string, stats []commonStat, width int) string {
	var b strings.Builder
	b.WriteString(styling.StyleBold.Render(title))
	b.WriteString("\n")
	if len(stats) == 0 {
		b.WriteString("-")
	} else {
		for i, stat := range stats {
			if i > 0 {
				b.WriteString("\n")
			}
			b.WriteString(statLine(stat, width))
		}
	}
	return lipgloss.NewStyle().Width(width).Render(b.String())
}

func statLine(stat commonStat, width int) string {
	const pctWidth = 6
	nameWidth := min(max(width-pctWidth-1, 6), 20)
	name := lipgloss.NewStyle().Width(nameWidth).Render(truncateName(stat.Name, nameWidth))
	return fmt.Sprintf("%s %*.1f%%", name, pctWidth-1, stat.UsagePercent)
}

func truncateName(name string, width int) string {
	if lipgloss.Width(name) <= width {
		return name
	}
	runes := []rune(name)
	if width <= 1 {
		return string(runes[:max(width, 0)])
	}
	return string(runes[:width-1]) + "…"
}

// Usage tab
func newUsageTable(rows []usageRow, height int) table.Model {
	const barWidth = 22
	columns := []table.Column{
		{Title: "#", Width: 4},
		{Title: "Pokémon", Width: 22},
		{Title: "Usage", Width: 7},
		{Title: "Share", Width: barWidth},
	}

	trows := make([]table.Row, 0, len(rows))
	for _, row := range rows {
		trows = append(trows, table.Row{
			strconv.Itoa(row.Rank),
			row.Pokemon,
			fmt.Sprintf("%.1f%%", row.UsagePercent),
			usageBar(row.UsagePercent, barWidth),
		})
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(trows),
		table.WithFocused(true),
		table.WithHeight(max(height-12, 5)),
		table.WithWidth(tableWidth(columns)),
	)
	t.SetStyles(shell.TableStyles())
	return t
}

func renderUsage(usageTable table.Model, rows []usageRow) string {
	if len(rows) == 0 {
		return "No data available"
	}

	caption := captionStyle.Render("Share of teams at recent Champions events that used each Pokémon.")
	return caption + "\n\n" + usageTable.View()
}

func usageBar(pct float64, width int) string {
	filled := min(int(pct*float64(width)/100), width)
	if filled == 0 && pct > 0 {
		filled = 1
	}
	return strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
}

// Top Teams tab
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

// Speed Tiers tab
func newSpeedTable(rows []speedTierRow, height int) table.Model {
	columns := []table.Column{
		{Title: "#", Width: 4},
		{Title: "Pokémon", Width: 20},
		{Title: "Base", Width: 5},
		{Title: "Min", Width: 5},
		{Title: "Max", Width: 5},
		{Title: "Scarf", Width: 6},
	}

	trows := make([]table.Row, 0, len(rows))
	for _, row := range rows {
		trows = append(trows, table.Row{
			strconv.Itoa(row.Rank),
			row.Pokemon,
			strconv.Itoa(row.BaseSpe),
			strconv.Itoa(row.NegMin),
			strconv.Itoa(row.Max),
			strconv.Itoa(row.MaxScarf),
		})
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(trows),
		table.WithFocused(true),
		table.WithHeight(max(height-12, 5)),
		table.WithWidth(tableWidth(columns)),
	)
	t.SetStyles(shell.TableStyles())
	return t
}

func renderSpeedTiers(speedTable table.Model, rows []speedTierRow) string {
	if len(rows) == 0 {
		return "No data available"
	}

	caption := captionStyle.Render("Speed stats at level 50. Min = 0 EVs, negative nature. Max = 252 EVs, positive nature.")
	detail := renderSpeedDetail(selectedSpeedTier(speedTable, rows))
	body := lipgloss.JoinHorizontal(lipgloss.Top, speedTable.View(), "  ", detail)
	return caption + "\n\n" + body
}

func selectedSpeedTier(speedTable table.Model, rows []speedTierRow) speedTierRow {
	if len(rows) == 0 {
		return speedTierRow{}
	}

	idx := min(max(speedTable.Cursor(), 0), len(rows)-1)
	return rows[idx]
}

func renderSpeedDetail(row speedTierRow) string {
	var b strings.Builder
	b.WriteString(styling.Yellow.Render("Selected Pokémon"))
	b.WriteString("\n")

	b.WriteString(styling.StyleBold.Render(row.Pokemon))
	b.WriteString("\n\n")

	stats := []struct {
		label string
		value int
	}{
		{"Base Speed", row.BaseSpe},
		{"Min (0 EV -Spe)", row.NegMin},
		{"Neutral (0 EV)", row.Neutral0},
		{"Neutral (252 EV)", row.Neutral252},
		{"Max (252 EV +Spe)", row.Max},
		{"Neutral + Scarf", row.NeutralScarf},
		{"Max + Scarf", row.MaxScarf},
	}

	for i, s := range stats {
		if i > 0 {
			b.WriteString("\n")
		}
		b.WriteString(speedStatLine(s.label, s.value))
	}
	return b.String()
}

func speedStatLine(label string, value int) string {
	const labelWidth = 19
	padded := lipgloss.NewStyle().Width(labelWidth).Render(label)
	return padded + styling.StyleBold.Render(strconv.Itoa(value))
}
