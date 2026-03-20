package tcg

import (
	"encoding/json"
	"net/url"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type standingRows struct {
	Rank          int    `json:"rank"`
	Name          string `json:"name"`
	Points        int    `json:"points"`
	Record        string `json:"record"`
	OppWinPct     string `json:"opp_win_percent"`
	OppOppWinPct  string `json:"opp_opp_win_percent"`
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
	items []standingRows
	err   error
}

func fetchData(tournament string) tea.Cmd {
	return func() tea.Msg {
		cols := "rank,name,points,record,opp_win_percent,opp_opp_win_percent,deck,player_country,country_code,location,text_date,type,iso_code,player_quantity"
		endpoint := "https://uoddayfnfkebrijlpfbh.supabase.co/rest/v1/standings?select=" + cols + "&location=eq." + url.QueryEscape(tournament) + "&order=rank"
		body, err := supabaseConn(endpoint)
		if err != nil {
			return standingsDataMsg{err: err}
		}

		var rows []standingRows
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

