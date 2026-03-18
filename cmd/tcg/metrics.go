package tcg

import (
	"encoding/json"
	"fmt"
	"net/url"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type metricsData struct {
	PlayerCountry string `json:"player_country"`
}

type metricsDataMsg struct {
	items []metricsData
	err   error
}

func fetchMetrics(tournament string) tea.Cmd {
	return func() tea.Msg {
		endpoint := "https://uoddayfnfkebrijlpfbh.supabase.co/rest/v1/standings?select=player_country&location=eq." + url.QueryEscape(tournament)
		body, err := supabaseConn(endpoint)
		if err != nil {
			return metricsDataMsg{err: err}
		}

		var allMetrics []metricsData
		if err = json.Unmarshal(body, &allMetrics); err != nil {
			return metricsDataMsg{err: err}
		}

		return metricsDataMsg{items: allMetrics}
	}
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
