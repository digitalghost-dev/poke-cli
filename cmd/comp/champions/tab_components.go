// Tab rendering helpers for the Champions dashboard.
// The Bubble Tea lifecycle stays in dashboard.go; this file builds the per-tab views.

package champions

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/table"
	"github.com/digitalghost-dev/poke-cli/cmd/comp/shell"
)

func newTeamsTable(teams []teamRow, width, height int) table.Model {
	columns := teamColumns(width)
	rows := make([]table.Row, 0, len(teams))

	for _, team := range teams {
		rows = append(rows, table.Row{
			fmt.Sprintf("%d", team.Rank),
			team.Author,
			team.Record,
			team.Tournament,
			joinOrDash(team.Archetypes),
			joinOrDash(team.Pokemon),
		})
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(max(height-14, 5)),
		table.WithWidth(tableWidth(columns)),
	)
	t.SetStyles(shell.TableStyles())
	return t
}

func renderTeamsTable(table table.Model) string {
	if len(table.Rows()) == 0 {
		return "No data available"
	}

	return table.View()
}

func teamColumns(width int) []table.Column {
	separatorWidth := 6 * 2
	availableWidth := max(width-separatorWidth-4-8, 40)

	authorWidth := min(max(availableWidth/5, 12), 20)
	tournamentWidth := min(max(availableWidth/4, 16), 30)
	archetypesWidth := min(max(availableWidth/5, 12), 20)
	teamWidth := max(availableWidth-authorWidth-tournamentWidth-archetypesWidth, 18)

	return []table.Column{
		{Title: "Rank", Width: 4},
		{Title: "Author", Width: authorWidth},
		{Title: "Record", Width: 8},
		{Title: "Tournament", Width: tournamentWidth},
		{Title: "Archetypes", Width: archetypesWidth},
		{Title: "Team", Width: teamWidth},
	}
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
