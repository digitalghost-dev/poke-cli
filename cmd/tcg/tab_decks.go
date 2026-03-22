package tcg

type deckStats struct {
	Deck  string
	Total int
}

func decksContent(s []deckStats, width int) string {
	items := make([]barChartItem, len(s))
	for i, c := range s {
		items[i] = barChartItem{Label: c.Deck, Total: c.Total}
	}

	return barChart(items, width, 30)
}
