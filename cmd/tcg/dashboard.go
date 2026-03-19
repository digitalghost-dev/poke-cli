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

type DeckStats struct {
	Deck  string
	Total int
}

type model struct {
	tabs           []string
	styles         *styles
	activeTab      int
	width          int
	tournament     string
	tournamentDate string
	tournamentType string
	standings      []standingRow
	countryStats   []CountryStats
	deckStats      []DeckStats
	totalPlayers   int
	winner         string
	winningDeck    string
	flag           string
	goBack         bool
	err            error
}


func overviewView(m model, contentWidth int) string {
	if len(m.standings) == 0 {
		return "  Loading..."
	}
	return OverviewContent(m.flag, m.tournament, m.tournamentType, m.tournamentDate, m.winner, m.winningDeck, m.totalPlayers, contentWidth)
}

func countriesView(s []CountryStats, width int) string {
	return CountryBarChart(s, width)
}

func (m model) Init() tea.Cmd {
	return fetchData(m.tournament)
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

	case tea.WindowSizeMsg:
		m.width = msg.Width
		return m, nil

	case standingsDataMsg:
		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}
		m.standings = msg.items
		if len(msg.items) > 0 {
			first := msg.items[0]
			m.totalPlayers = first.PlayerQty
			m.winner = first.Name
			m.winningDeck = first.Deck
			m.flag = countryFlag(first.ISOCode)
			m.tournamentDate = first.TextDate
			m.tournamentType = first.Type
		}
		countryCounts := map[string]int{}
		for _, row := range msg.items {
			countryCounts[row.PlayerCountry]++
		}
		for country, count := range countryCounts {
			m.countryStats = append(m.countryStats, CountryStats{Country: country, Total: count})
		}
		deckCounts := map[string]int{}
		for _, row := range msg.items {
			deckCounts[row.Deck]++
		}
		for deck, count := range deckCounts {
			m.deckStats = append(m.deckStats, DeckStats{Deck: deck, Total: count})
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
			border.BottomRight = "└"
		} else if isLast && !isActive {
			border.BottomRight = "┴"
		}
		style = style.Border(border)
		renderedTabs = append(renderedTabs, style.Render(t))
	}

	row := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)

	// Use terminal width if available, otherwise fall back to tab row width.
	// doc has Padding(1,2,1,2) = 4 horizontal chars; window border = 2 chars.
	windowWidth := max(m.width-8, lipgloss.Width(row)-2)
	contentWidth := windowWidth - 2

	// Fill the gap between the tab row and the window's right edge so the top
	// border line stretches the full width of the window.
	highlightColor := lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	fillWidth := (windowWidth + 2) - lipgloss.Width(row)
	if fillWidth > 0 {
		fill := lipgloss.NewStyle().Foreground(highlightColor).
			Render(strings.Repeat("─", fillWidth-1) + "┐")
		row = row + fill
	}

	var content string
	switch m.activeTab {
	case 0:
		if m.err != nil {
			content = fmt.Sprintf("fetch error: %v", m.err)
		} else {
			content = overviewView(m, contentWidth)
		}
	case 3:
		if m.err != nil {
			content = fmt.Sprintf("fetch error: %v", m.err)
		} else {
			content = countriesView(m.countryStats, contentWidth)
		}
	}

	doc.WriteString(row)
	doc.WriteString("\n")
	doc.WriteString(s.window.Width(windowWidth).Render(content))
	doc.WriteString("\n")
	doc.WriteString(styling.KeyMenu.Render("← → (switch tab) • b (back) • ctrl+c | esc (quit)"))

	return s.doc.Render(doc.String())
}
