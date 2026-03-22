package tcg

import (
	"fmt"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
	"github.com/digitalghost-dev/poke-cli/styling"
)

func standingsTable(rows []standingRows, width, height int) table.Model {
	fixedWidth := 4 + 20 + 6 + 10 + 7 + 7 + 18
	separators := 8 * 2
	deckWidth := min(max(width-fixedWidth-separators, 10), 30)

	columns := []table.Column{
		{Title: "Rank", Width: 4},
		{Title: "Name", Width: 20},
		{Title: "Points", Width: 6},
		{Title: "Record", Width: 10},
		{Title: "OPW%", Width: 7},
		{Title: "OOPW%", Width: 7},
		{Title: "Deck", Width: deckWidth},
		{Title: "Country", Width: 18},
	}

	tableRows := make([]table.Row, len(rows))
	for i, r := range rows {
		tableRows[i] = table.Row{
			fmt.Sprintf("%d", r.Rank),
			r.Name,
			fmt.Sprintf("%d", r.Points),
			r.Record,
			r.OppWinPct,
			r.OppOppWinPct,
			r.Deck,
			r.PlayerCountry,
		}
	}

	tableHeight := max(height-14, 5)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(styling.YellowColor).
		BorderBottom(true).
		Bold(true)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("#000")).
		Background(styling.YellowColor)

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(tableRows),
		table.WithFocused(true),
		table.WithHeight(tableHeight),
	)
	t.SetStyles(s)

	return t
}
