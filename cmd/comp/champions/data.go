package champions

import (
	"encoding/json"

	tea "charm.land/bubbletea/v2"
	"github.com/digitalghost-dev/poke-cli/cmd/comp/shell"
)

const (
	compInfoURL   = "https://uoddayfnfkebrijlpfbh.supabase.co/rest/v1/pikalytics_pokemon_comp_info?select=pokemon,web_url,common_moves,common_abilities,common_items,common_teammates&order=pokemon"
	topTeamsURL   = "https://uoddayfnfkebrijlpfbh.supabase.co/rest/v1/pikalytics_top_teams?select=author,record,tournament,archetypes,pokemon,web_url&order=rank"
	speedTiersURL = "https://uoddayfnfkebrijlpfbh.supabase.co/rest/v1/pikalytics_speed_tiers?select=rank,pokemon,base_spe,neutral_0_sp,neutral_32_sp,neg_spe_0_sp,max_speed,max_scarf,neutral_32_scarf&order=rank"
)

type dashboardData struct {
	CompInfo   []compInfoRow
	Teams      []teamRow
	SpeedTiers []speedTierRow
}

type dataMsg struct {
	data *dashboardData
	err  error
}

type teamRow struct {
	Player     string   `json:"author"`
	Record     string   `json:"record"`
	Tournament string   `json:"tournament"`
	Archetypes []string `json:"archetypes"`
	Pokemon    []string `json:"pokemon"`
	WebURL     string   `json:"web_url"`
}

type compInfoRow struct {
	Pokemon         string       `json:"pokemon"`
	WebURL          string       `json:"web_url"`
	CommonMoves     []commonStat `json:"common_moves"`
	CommonAbilities []commonStat `json:"common_abilities"`
	CommonItems     []commonStat `json:"common_items"`
	CommonTeammates []commonStat `json:"common_teammates"`
}

type commonStat struct {
	Name         string  `json:"name"`
	UsagePercent float64 `json:"usage_percent"`
}

type speedTierRow struct {
	Rank         int    `json:"rank"`
	Pokemon      string `json:"pokemon"`
	BaseSpe      int    `json:"base_spe"`
	Neutral0     int    `json:"neutral_0_sp"`
	Neutral252   int    `json:"neutral_32_sp"`
	NegMin       int    `json:"neg_spe_0_sp"`
	Max          int    `json:"max_speed"`
	MaxScarf     int    `json:"max_scarf"`
	NeutralScarf int    `json:"neutral_32_scarf"`
}

func fetchDashboardData(conn shell.ConnFunc) tea.Cmd {
	return func() tea.Msg {
		compInfo, err := fetchCompInfo(conn)
		if err != nil {
			return dataMsg{err: err}
		}

		teams, err := fetchTopTeams(conn)
		if err != nil {
			return dataMsg{err: err}
		}

		speedTiers, err := fetchSpeedTiers(conn)
		if err != nil {
			return dataMsg{err: err}
		}

		return dataMsg{
			data: &dashboardData{
				CompInfo:   compInfo,
				Teams:      teams,
				SpeedTiers: speedTiers,
			},
		}
	}
}

func fetchCompInfo(conn shell.ConnFunc) ([]compInfoRow, error) {
	body, err := conn(compInfoURL)
	if err != nil {
		return nil, err
	}

	var rows []compInfoRow
	if err := json.Unmarshal(body, &rows); err != nil {
		return nil, err
	}
	return rows, nil
}

func fetchTopTeams(conn shell.ConnFunc) ([]teamRow, error) {
	body, err := conn(topTeamsURL)
	if err != nil {
		return nil, err
	}

	var teams []teamRow
	if err := json.Unmarshal(body, &teams); err != nil {
		return nil, err
	}
	return teams, nil
}

func fetchSpeedTiers(conn shell.ConnFunc) ([]speedTierRow, error) {
	body, err := conn(speedTiersURL)
	if err != nil {
		return nil, err
	}

	var rows []speedTierRow
	if err := json.Unmarshal(body, &rows); err != nil {
		return nil, err
	}
	return rows, nil
}
