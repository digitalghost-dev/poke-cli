package vgc

import (
	"net/url"

	"charm.land/bubbles/v2/table"
	"github.com/digitalghost-dev/poke-cli/cmd/comp/shell"
	"github.com/digitalghost-dev/poke-cli/connections"
)

const baseURL = "https://uoddayfnfkebrijlpfbh.supabase.co/rest/v1/comp_vgc_standings_view"

func Run() (back bool, err error) {
	return shell.Run(Spec(), connections.CallTCGData)
}

func Spec() shell.Spec {
	return shell.Spec{
		Tabs:    []string{"Overview", "Standings", "Usage", "Countries"},
		ListURL: baseURL + "?select=location,text_date&rank=eq.1&order=start_date.desc",
		DashboardURL: func(location string) string {
			cols := "rank,name,points,record,opp_win_percent,opp_opp_win_percent,team,player_country,country_code,location,text_date,type,player_quantity"
			return baseURL + "?select=" + cols + "&location=eq." + url.QueryEscape(location) + "&order=rank"
		},
		Columns: standingsColumns,
		Decode:  decode,
	}
}

func standingsColumns(width int) []table.Column {
	fixedWidth := 4 + 6 + 10 + 7 + 7
	separators := 7 * 2
	flexWidth := max(width-fixedWidth-separators, 40)
	nameWidth := min(max(flexWidth*3/5, 20), 32)
	countryWidth := min(max(flexWidth-nameWidth, 16), 24)
	return []table.Column{
		{Title: "Rank", Width: 4},
		{Title: "Name", Width: nameWidth},
		{Title: "Points", Width: 6},
		{Title: "Record", Width: 10},
		{Title: "OPW%", Width: 7},
		{Title: "OOPW%", Width: 7},
		{Title: "Country", Width: countryWidth},
	}
}
