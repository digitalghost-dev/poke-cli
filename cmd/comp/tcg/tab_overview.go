package tcg

import (
	"fmt"
	"image/color"
	"strconv"
	"strings"

	"charm.land/lipgloss/v2"
)

func formatInt(n int) string {
	s := strconv.Itoa(n)
	var result strings.Builder
	for i, c := range s {
		if i > 0 && (len(s)-i)%3 == 0 {
			result.WriteRune(',')
		}
		result.WriteRune(c)
	}
	return result.String()
}

func overviewContent(flag, tournament, tournamentType, tournamentDate, winner, winningDeck string, totalPlayers, contentWidth int, highlightColor color.Color) string {
	header := fmt.Sprintf("%s  %s · %s · %s", flag, tournament, tournamentType, tournamentDate)

	statBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(highlightColor).
		Padding(1, 2).
		Width(26).
		Align(lipgloss.Center)

	totalBox := statBox.Render("Total Players\n\n" + formatInt(totalPlayers))
	winnerBox := statBox.Render("Winner\n\n" + winner)
	deckBox := statBox.Render("Winning Deck\n\n" + winningDeck)

	boxes := lipgloss.JoinHorizontal(lipgloss.Top, totalBox, "  ", winnerBox, "  ", deckBox)

	content := header + "\n\n" + boxes
	return lipgloss.NewStyle().Width(contentWidth).Align(lipgloss.Center).Render(content)
}
