package tcg

import (
	"encoding/json"
	"fmt"
	"net/url"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type standingRow struct {
	Rank          int    `json:"rank"`
	Name          string `json:"name"`
	Points        int    `json:"points"`
	Record        string `json:"record"`
	Deck          string `json:"deck"`
	PlayerCountry string `json:"player_country"`
	CountryCode   string `json:"country_code"`
	Location      string `json:"location"`
	TextDate      string `json:"text_date"`
	Type          string `json:"type"`
	ISOCode       string `json:"iso_code"`
	PlayerQty     int    `json:"player_quantity"`
}

type standingsDataMsg struct {
	items []standingRow
	err   error
}

func fetchStandings(tournament string) tea.Cmd {
	return func() tea.Msg {
		cols := "rank,name,points,record,deck,player_country,country_code,location,text_date,type,iso_code,player_quantity"
		endpoint := "https://uoddayfnfkebrijlpfbh.supabase.co/rest/v1/standings?select=" + cols + "&location=eq." + url.QueryEscape(tournament) + "&order=rank"
		body, err := supabaseConn(endpoint)
		if err != nil {
			return standingsDataMsg{err: err}
		}

		var rows []standingRow
		if err = json.Unmarshal(body, &rows); err != nil {
			return standingsDataMsg{err: err}
		}

		return standingsDataMsg{items: rows}
	}
}

func countryFlag(isoCode string) string {
	code := strings.ToUpper(isoCode)
	if len(code) != 2 {
		return ""
	}
	return string(rune(0x1F1E6+(rune(code[0])-'A'))) + string(rune(0x1F1E6+(rune(code[1])-'A')))
}

func OverviewContent(flag, tournament, tournamentType, tournamentDate, winner, winningDeck string, totalPlayers, contentWidth int) string {
	highlightColor := lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	header := fmt.Sprintf("%s  %s · %s · %s", flag, tournament, tournamentType, tournamentDate)

	statBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(highlightColor).
		Padding(1, 2).
		Width(26).
		Align(lipgloss.Center)

	totalBox := statBox.Render(fmt.Sprintf("Total Players\n\n%s", formatInt(totalPlayers)))
	winnerBox := statBox.Render(fmt.Sprintf("Winner\n\n%s", winner))
	deckBox := statBox.Render(fmt.Sprintf("Winning Deck\n\n%s", winningDeck))

	boxes := lipgloss.JoinHorizontal(lipgloss.Top, totalBox, "  ", winnerBox, "  ", deckBox)

	content := header + "\n\n" + boxes
	return lipgloss.NewStyle().Width(contentWidth).Align(lipgloss.Center).Render(content)
}

func CountryBarChart(s []CountryStats, width int) string {
	if len(s) == 0 {
		return ""
	}

	sorted := make([]CountryStats, len(s))
	copy(sorted, s)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Total > sorted[j].Total
	})

	display := sorted
	if len(sorted) > 9 {
		other := 0
		for _, stat := range sorted[9:] {
			other += stat.Total
		}
		display = append(sorted[:9], CountryStats{Country: "Other", Total: other})
	}

	const labelWidth = 16
	const countWidth = 5
	maxBarWidth := width - labelWidth - countWidth - 4
	if maxBarWidth < 10 {
		maxBarWidth = 10
	}

	maxVal := 0
	for _, stat := range display {
		if stat.Total > maxVal {
			maxVal = stat.Total
		}
	}

	var sb strings.Builder
	for _, stat := range display {
		barWidth := stat.Total * maxBarWidth / maxVal
		bar := strings.Repeat("█", barWidth) + strings.Repeat(" ", maxBarWidth-barWidth)
		sb.WriteString(fmt.Sprintf("%-*s %s %*d\n", labelWidth, stat.Country, bar, countWidth, stat.Total))
	}
	return sb.String()
}
