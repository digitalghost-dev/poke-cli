package champions

import (
	"encoding/json"

	tea "charm.land/bubbletea/v2"
	"github.com/digitalghost-dev/poke-cli/cmd/comp/shell"
)

const topTeamsURL = "https://uoddayfnfkebrijlpfbh.supabase.co/rest/v1/pikalytics_top_teams?select=rank,author,record,tournament,archetypes,pokemon,web_url&order=rank"

type dashboardData struct {
	Teams []teamRow
}

type dataMsg struct {
	data *dashboardData
	err  error
}

type teamRow struct {
	Rank       int      `json:"rank"`
	Author     string   `json:"author"`
	Record     string   `json:"record"`
	Tournament string   `json:"tournament"`
	Archetypes []string `json:"archetypes"`
	Pokemon    []string `json:"pokemon"`
	WebURL     string   `json:"web_url"`
}

func fetchDashboardData(conn shell.ConnFunc) tea.Cmd {
	return func() tea.Msg {
		body, err := conn(topTeamsURL)
		if err != nil {
			return dataMsg{err: err}
		}

		var teams []teamRow
		if err := json.Unmarshal(body, &teams); err != nil {
			return dataMsg{err: err}
		}

		return dataMsg{
			data: &dashboardData{
				Teams: teams,
			},
		}
	}
}
