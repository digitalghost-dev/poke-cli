package tcg

import (
	"fmt"
	"sort"
	"strings"
)

type barChartItem struct {
	Label string
	Total int
}

func barChart(s []barChartItem, width, labelWidth int) string {
	if len(s) == 0 {
		return ""
	}

	sorted := make([]barChartItem, len(s))
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
		display = append(sorted[:9], barChartItem{Label: "Other", Total: other})
	}

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
		barWidth := 0
		if maxVal > 0 {
			barWidth = stat.Total * maxBarWidth / maxVal
			if barWidth == 0 && stat.Total > 0 {
				barWidth = 1
			}
		}
		bar := strings.Repeat("█", barWidth) + strings.Repeat(" ", maxBarWidth-barWidth)
		sb.WriteString(fmt.Sprintf("%-*s %s %*d\n", labelWidth, stat.Label, bar, countWidth, stat.Total))
	}
	return sb.String()
}
