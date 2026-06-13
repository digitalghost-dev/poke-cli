package shell

import (
	"fmt"

	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"
	"github.com/digitalghost-dev/poke-cli/cmd/comp/web"
)

type dashboardModel struct {
	spec       Spec
	conn       ConnFunc
	styles     *Styles
	activeTab  int
	width      int
	height     int
	tournament string
	decoded    *Decoded
	table      table.Model
	goBack     bool
	err        error
}

type dataMsg struct {
	decoded Decoded
	err     error
}

func newDashboard(spec Spec, conn ConnFunc, location string) dashboardModel {
	return dashboardModel{
		spec:       spec,
		conn:       conn,
		styles:     NewStyles(),
		tournament: location,
	}
}

func fetchDashboard(spec Spec, location string, conn ConnFunc) tea.Cmd {
	return func() tea.Msg {
		body, err := conn(spec.DashboardURL(location))
		if err != nil {
			return dataMsg{err: err}
		}
		decoded, err := spec.Decode(body)
		if err != nil {
			return dataMsg{err: err}
		}
		return dataMsg{decoded: decoded}
	}
}

func newTable(spec Spec, rows []table.Row, width, height int) table.Model {
	columns := spec.Columns(width - 8)
	tableWidth := len(columns) * 2
	for _, c := range columns {
		tableWidth += c.Width
	}
	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(max(height-14, 5)),
		table.WithWidth(tableWidth),
	)
	t.SetStyles(TableStyles())
	return t
}

func (m dashboardModel) Init() tea.Cmd {
	return fetchDashboard(m.spec, m.tournament, m.conn)
}

func (m dashboardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "b":
			m.goBack = true
			return m, tea.Quit
		case "w":
			return m, web.Open("https://web.poke-cli.com/")
		case "right", "l", "n", "tab":
			m.activeTab = min(m.activeTab+1, len(m.spec.Tabs)-1)
			return m, nil
		case "left", "h", "p", "shift+tab":
			m.activeTab = max(m.activeTab-1, 0)
			return m, nil
		}
		if m.activeTab == 1 && m.decoded != nil {
			var cmd tea.Cmd
			m.table, cmd = m.table.Update(msg)
			return m, cmd
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		if m.decoded != nil {
			m.table = newTable(m.spec, m.decoded.TableRows, m.width, m.height)
		}
		return m, nil

	case dataMsg:
		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}
		d := msg.decoded
		m.decoded = &d
		m.table = newTable(m.spec, d.TableRows, m.width, m.height)
	}

	return m, nil
}

func (m dashboardModel) renderTab(contentWidth int) string {
	if m.err != nil {
		return fmt.Sprintf("fetch error: %v", m.err)
	}
	if m.decoded == nil {
		return "  Loading..."
	}
	switch m.activeTab {
	case 0:
		return m.decoded.Overview(contentWidth, m.styles.HighlightColor)
	case 1:
		return m.table.View()
	case 2:
		return m.decoded.ExtraTab(contentWidth)
	case 3:
		return BarChart(m.decoded.Countries, contentWidth, 20)
	}
	return ""
}

func (m dashboardModel) View() tea.View {
	if m.styles == nil {
		return tea.NewView("")
	}

	body := m.styles.Render(m.spec.Tabs, m.activeTab, m.width, m.renderTab)

	v := tea.NewView(body)
	v.AltScreen = true
	return v
}
