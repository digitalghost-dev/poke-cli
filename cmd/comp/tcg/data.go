package tcg

import (
	"encoding/json"
	"image/color"
	"strconv"

	"charm.land/bubbles/v2/table"
	"github.com/digitalghost-dev/poke-cli/cmd/comp/shell"
)

type standingRow struct {
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
	PlayerQty     int    `json:"player_quantity"`
}

func decode(body []byte) (shell.Decoded, error) {
	var rows []standingRow
	if err := json.Unmarshal(body, &rows); err != nil {
		return shell.Decoded{}, err
	}

	d := shell.Decoded{
		TableRows: make([]table.Row, len(rows)),
		Countries: countryItems(rows),
	}
	for i, r := range rows {
		d.TableRows[i] = table.Row{
			strconv.Itoa(r.Rank), r.Name, strconv.Itoa(r.Points), r.Record,
			r.OppWinPct, r.OppOppWinPct, r.Deck, r.PlayerCountry,
		}
	}

	d.Extra = shell.Frequency{
		NameHeader:  "Deck",
		CountHeader: "Players",
		Caption:     "Based on the top 256 players per event.",
		Items:       deckItems(rows),
	}

	var tournament, tType, date, winner, winningDeck string
	var total int
	if len(rows) > 0 {
		first := rows[0]
		tournament, tType, date = first.Location, first.Type, first.TextDate
		winner, winningDeck, total = first.Name, first.Deck, first.PlayerQty
	}
	d.Overview = func(contentWidth int, hc color.Color) string {
		return overviewContent(tournament, tType, date, winner, winningDeck, total, contentWidth, hc)
	}
	return d, nil
}

func countryItems(rows []standingRow) []shell.Tally {
	counts := map[string]int{}
	for _, r := range rows {
		if r.PlayerCountry != "" {
			counts[r.PlayerCountry]++
		}
	}
	items := make([]shell.Tally, 0, len(counts))
	for country, n := range counts {
		items = append(items, shell.Tally{Label: country, Count: n})
	}
	return items
}

func deckItems(rows []standingRow) []shell.Tally {
	counts := map[string]int{}
	for _, r := range rows {
		counts[r.Deck]++
	}
	items := make([]shell.Tally, 0, len(counts))
	for deck, n := range counts {
		items = append(items, shell.Tally{Label: deck, Count: n})
	}
	return items
}
