package tcg

import (
	"fmt"
	"image/color"

	"charm.land/lipgloss/v2"
	"github.com/digitalghost-dev/poke-cli/cmd/comp/shell"
)

func overviewContent(tournament, tournamentType, tournamentDate, winner, winningDeck string, totalPlayers, contentWidth int, highlightColor color.Color) string {
	header := fmt.Sprintf("%s · %s · %s", tournament, tournamentType, tournamentDate)

	statBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(highlightColor).
		Padding(1, 2).
		Width(26).
		Align(lipgloss.Center)

	totalBox := statBox.Render("Total Players\n\n" + shell.FormatInt(totalPlayers))
	winnerBox := statBox.Render("Winner\n\n" + winner)
	deckBox := statBox.Render("Winning Deck\n\n" + winningDeck)

	boxes := lipgloss.JoinHorizontal(lipgloss.Top, totalBox, "  ", winnerBox, "  ", deckBox)

	content := header + "\n\n" + boxes
	return lipgloss.NewStyle().Width(contentWidth).Align(lipgloss.Center).Render(content)
}
