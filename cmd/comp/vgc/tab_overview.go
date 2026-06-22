package vgc

import (
	"fmt"
	"image/color"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/digitalghost-dev/poke-cli/cmd/comp/shell"
)

func baseName(name string) string {
	if i := strings.IndexByte(name, '['); i > 0 {
		return strings.TrimSpace(name[:i])
	}
	return name
}

func teamGrid(team []string) string {
	if len(team) == 0 {
		return "—"
	}

	bullets := make([]string, len(team))
	for i, name := range team {
		bullets[i] = "• " + name
	}

	rows := (len(bullets) + 1) / 2
	left := bullets[:rows]
	right := bullets[rows:]

	cellWidth := func(items []string) int {
		w := 0
		for _, s := range items {
			if lw := lipgloss.Width(s); lw > w {
				w = lw
			}
		}
		return w
	}
	leftStyle := lipgloss.NewStyle().Width(cellWidth(left))
	rightStyle := lipgloss.NewStyle().Width(cellWidth(right))

	lines := make([]string, rows)
	for i := range rows {
		r := ""
		if i < len(right) {
			r = right[i]
		}
		lines[i] = leftStyle.Render(left[i]) + "   " + rightStyle.Render(r)
	}
	return strings.Join(lines, "\n")
}

func overviewContent(tournament, tournamentType, tournamentDate, winner string, winnerTeam []string, totalPlayers, contentWidth int, highlightColor color.Color) string {
	header := fmt.Sprintf("%s · %s · %s", tournament, tournamentType, tournamentDate)

	statBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(highlightColor).
		Padding(1, 2).
		Width(26).
		Align(lipgloss.Center)

	totalBox := statBox.Render("Total Players\n\n" + shell.FormatInt(totalPlayers))
	winnerBox := statBox.Render("Winner\n\n" + winner)

	teamBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(highlightColor).
		Padding(1, 2).
		Align(lipgloss.Center).
		Render("Winning Team\n\n" + teamGrid(winnerTeam))

	boxes := lipgloss.JoinHorizontal(lipgloss.Top, totalBox, "  ", winnerBox, "  ", teamBox)

	content := header + "\n\n" + boxes
	return lipgloss.NewStyle().Width(contentWidth).Align(lipgloss.Center).Render(content)
}
