package tcg

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/digitalghost-dev/poke-cli/styling"
)

type styles struct {
	doc         lipgloss.Style
	inactiveTab lipgloss.Style
	activeTab   lipgloss.Style
	window      lipgloss.Style
}

func tabBorderWithBottom(left, middle, right string) lipgloss.Border {
	border := lipgloss.RoundedBorder()
	border.BottomLeft = left
	border.Bottom = middle
	border.BottomRight = right
	return border
}

func newStyles() *styles {
	inactiveTabBorder := tabBorderWithBottom("┴", "─", "┴")
	activeTabBorder := tabBorderWithBottom("┘", " ", "└")
	highlightColor := lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}

	s := new(styles)
	s.doc = lipgloss.NewStyle().
		Padding(1, 2, 1, 2)
	s.inactiveTab = lipgloss.NewStyle().
		Border(inactiveTabBorder, true).
		BorderForeground(highlightColor).
		Padding(0, 1)
	s.activeTab = s.inactiveTab.
		Border(activeTabBorder, true)
	s.window = lipgloss.NewStyle().
		BorderForeground(highlightColor).
		Padding(2, 0).
		Border(lipgloss.NormalBorder()).
		UnsetBorderTop()
	return s
}

type CountryStats struct {
	Country string
	Total   int
}

type model struct {
	tabs         []string
	styles       *styles
	activeTab    int
	tournament   string
	countryStats []CountryStats
	goBack       bool
	err          error
}

func countriesView(s []CountryStats, width int) string {
	return CountryBarChart(s, width)
}

func (m model) Init() tea.Cmd {
	return fetchMetrics(m.tournament)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "b":
			m.goBack = true
			return m, tea.Quit
		case "right", "l", "n", "tab":
			m.activeTab = min(m.activeTab+1, len(m.tabs)-1)
			return m, nil
		case "left", "h", "p", "shift+tab":
			m.activeTab = max(m.activeTab-1, 0)
			return m, nil
		}

	case metricsDataMsg:
		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}
		counts := map[string]int{}
		for _, row := range msg.items {
			counts[row.PlayerCountry]++
		}
		for country, count := range counts {
			m.countryStats = append(m.countryStats, CountryStats{Country: country, Total: count})
		}
	}

	return m, nil
}

func (m model) View() string {
	if m.styles == nil {
		return ""
	}

	doc := strings.Builder{}
	s := m.styles

	var renderedTabs []string

	for i, t := range m.tabs {
		var style lipgloss.Style
		isFirst, isLast, isActive := i == 0, i == len(m.tabs)-1, i == m.activeTab
		if isActive {
			style = s.activeTab
		} else {
			style = s.inactiveTab
		}
		border, _, _, _, _ := style.GetBorder()
		if isFirst && isActive {
			border.BottomLeft = "│"
		} else if isFirst && !isActive {
			border.BottomLeft = "├"
		} else if isLast && isActive {
			border.BottomRight = "│"
		} else if isLast && !isActive {
			border.BottomRight = "┤"
		}
		style = style.Border(border)
		renderedTabs = append(renderedTabs, style.Render(t))
	}

	row := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
	contentWidth := lipgloss.Width(row) - 4 // subtract window borders

	var content string
	switch m.activeTab {
	case 3:
		if m.err != nil {
			content = fmt.Sprintf("fetch error: %v", m.err)
		} else {
			content = countriesView(m.countryStats, contentWidth)
		}
	}

	doc.WriteString(row)
	doc.WriteString("\n")
	doc.WriteString(s.window.Width(lipgloss.Width(row) - 2).Render(content))
	doc.WriteString("\n")
	doc.WriteString(styling.KeyMenu.Render("← → (switch tab) • b (back) • ctrl+c | esc (quit)"))

	return s.doc.Render(doc.String())
}
