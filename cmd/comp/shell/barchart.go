package shell

import (
	"fmt"
	"sort"
	"strings"
)

type Tally struct {
	Label string
	Count int
}

func BarChart(s []Tally, width, labelWidth int) string {
	if len(s) == 0 {
		return ""
	}

	sorted := make([]Tally, len(s))
	copy(sorted, s)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Count > sorted[j].Count
	})

	display := sorted
	if len(sorted) > 9 {
		other := 0
		for _, stat := range sorted[9:] {
			other += stat.Count
		}
		display = append(sorted[:9], Tally{Label: "Other", Count: other})
	}

	const countWidth = 5
	maxBarWidth := max(width-labelWidth-countWidth-4, 10)

	maxVal := 0
	for _, stat := range display {
		if stat.Count > maxVal {
			maxVal = stat.Count
		}
	}

	var sb strings.Builder
	for _, stat := range display {
		barWidth := 0
		if maxVal > 0 {
			barWidth = stat.Count * maxBarWidth / maxVal
			if barWidth == 0 && stat.Count > 0 {
				barWidth = 1
			}
		}
		bar := strings.Repeat("█", barWidth) + strings.Repeat(" ", maxBarWidth-barWidth)
		fmt.Fprintf(&sb, "%-*s %s %*d\n", labelWidth, stat.Label, bar, countWidth, stat.Count)
	}
	return sb.String()
}
