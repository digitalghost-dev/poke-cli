package tcg

type countryStats struct {
	Country string
	Total   int
}

func countriesContent(s []countryStats, width int) string {
	items := make([]barChartItem, len(s))
	for i, c := range s {
		items[i] = barChartItem{Label: c.Country, Total: c.Total}
	}

	return barChart(items, width, 20)
}
