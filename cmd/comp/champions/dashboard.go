package champions

import (
	"fmt"

	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"
	"github.com/digitalghost-dev/poke-cli/cmd/comp/shell"
	"github.com/digitalghost-dev/poke-cli/cmd/comp/web"
)

var tabs = []string{"Pokémon Overview", "Usage", "Top Teams", "Speed Tiers"}

type dashboardModel struct {
	activeTab int
	conn      shell.ConnFunc
	data      *dashboardData
	err       error
	goBack    bool
	height    int
	quit      bool
	overview  table.Model
	speed     table.Model
	styles    *shell.Styles
	teams     table.Model
	usage     table.Model
	width     int
}

func newDashboard(conn shell.ConnFunc) dashboardModel {
	return dashboardModel{
		conn:   conn,
		styles: shell.NewStyles(),
	}
}

func (m dashboardModel) renderTab(contentWidth int) string {
	if m.err != nil {
		return fmt.Sprintf("fetch error: %v", m.err)
	}

	if m.data == nil {
		return "  Loading..."
	}

	switch m.activeTab {
	case 0:
		return renderOverview(m.overview, m.data.CompInfo, contentWidth)
	case 1:
		return renderUsage(m.usage, m.data.Usage)
	case 2:
		return renderTeamsTable(m.teams, m.data.Teams, contentWidth)
	case 3:
		return renderSpeedTiers(m.speed, m.data.SpeedTiers)
	default:
		return ""
	}
}

func (m dashboardModel) Init() tea.Cmd {
	return fetchDashboardData(m.conn)
}

func (m dashboardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.quit = true
			return m, tea.Quit
		case "b":
			m.goBack = true
			return m, tea.Quit
		case "w":
			return m, web.Open("https://web.poke-cli.com/")
		case "right", "l", "tab":
			m.activeTab = min(m.activeTab+1, len(tabs)-1)
			return m, nil
		case "left", "h", "shift+tab":
			m.activeTab = max(m.activeTab-1, 0)
			return m, nil
		}
		if m.data != nil {
			switch m.activeTab {
			case 0:
				var cmd tea.Cmd
				m.overview, cmd = m.overview.Update(msg)
				return m, cmd
			case 1:
				var cmd tea.Cmd
				m.usage, cmd = m.usage.Update(msg)
				return m, cmd
			case 2:
				var cmd tea.Cmd
				m.teams, cmd = m.teams.Update(msg)
				return m, cmd
			case 3:
				var cmd tea.Cmd
				m.speed, cmd = m.speed.Update(msg)
				return m, cmd
			}
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		if m.data != nil {
			m.overview = newOverviewTable(m.data.CompInfo, m.height)
			m.teams = newTeamsTable(m.data.Teams, contentWidth(m.width), m.height)
			m.usage = newUsageTable(m.data.Usage, m.height)
			m.speed = newSpeedTable(m.data.SpeedTiers, m.height)
		}
		return m, nil
	case dataMsg:
		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}
		m.data = msg.data
		m.overview = newOverviewTable(m.data.CompInfo, m.height)
		m.teams = newTeamsTable(m.data.Teams, contentWidth(m.width), m.height)
		m.usage = newUsageTable(m.data.Usage, m.height)
		m.speed = newSpeedTable(m.data.SpeedTiers, m.height)
		return m, nil
	}

	return m, nil
}

func (m dashboardModel) View() tea.View {
	if m.quit {
		return tea.NewView("\n Goodbye! \n")
	}

	if m.styles == nil {
		return tea.NewView("")
	}

	body := m.styles.Render(tabs, m.activeTab, m.width, m.renderTab)

	v := tea.NewView(body)
	v.AltScreen = true
	return v
}

func contentWidth(width int) int {
	return max(width-10, 40)
}
