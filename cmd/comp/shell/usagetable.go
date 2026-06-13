package shell

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"charm.land/bubbles/v2/table"
)

func newUsageTable(f Frequency, total, width, height int) table.Model {
	avail := width - 8
	const rankW, countW, shareW, barWidth = 4, 8, 18, 11
	nameW := min(max(avail-rankW-countW-shareW-8, 16), 40)

	columns := []table.Column{
		{Title: "#", Width: rankW},
		{Title: f.NameHeader, Width: nameW},
		{Title: f.CountHeader, Width: countW},
		{Title: "Share", Width: shareW},
	}

	items := make([]Tally, len(f.Items))
	copy(items, f.Items)
	sort.Slice(items, func(i, j int) bool { return items[i].Count > items[j].Count })

	rows := make([]table.Row, len(items))
	for i, it := range items {
		rows[i] = table.Row{
			strconv.Itoa(i + 1),
			it.Label,
			strconv.Itoa(it.Count),
			shareCell(it.Count, total, barWidth),
		}
	}

	tableHeight := height - 14
	if f.Caption != "" {
		tableHeight -= 2
	}
	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(max(tableHeight, 5)),
		table.WithWidth(rankW+nameW+countW+shareW+4*2),
	)
	t.SetStyles(TableStyles())
	return t
}

func shareCell(count, total, barWidth int) string {
	filled, pct := 0, 0
	if total > 0 {
		filled = min(count*barWidth/total, barWidth)
		if filled == 0 && count > 0 {
			filled = 1
		}
		pct = count * 100 / total
	}
	bar := strings.Repeat("█", filled) + strings.Repeat("░", barWidth-filled)
	return fmt.Sprintf("%s %3d%%", bar, pct)
}
