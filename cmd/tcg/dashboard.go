package tcg

import (
	"fmt"
	"image/color"
	"os"
	"strings"

	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/digitalghost-dev/poke-cli/styling"
)

type styles struct {
	doc            lipgloss.Style
	inactiveTab    lipgloss.Style
	activeTab      lipgloss.Style
	window         lipgloss.Style
	highlightColor color.Color
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
	isDark := lipgloss.HasDarkBackground(os.Stdin, os.Stdout)
	ld := lipgloss.LightDark(isDark)
	highlightColor := ld(lipgloss.Color("#874BFD"), lipgloss.Color("#7D56F4"))

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
	s.highlightColor = highlightColor
	return s
}

type model struct {
	conn           func(string) ([]byte, error)
	tabs           []string
	styles         *styles
	activeTab      int
	width          int
	height         int
	standingsTable table.Model
	tournament     string
	tournamentDate string
	tournamentType string
	standings      []standingRows
	countryStats   []countryStats
	deckStats      []deckStats
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
	return overviewContent(m.flag, m.tournament, m.tournamentType, m.tournamentDate, m.winner, m.winningDeck, m.totalPlayers, contentWidth, m.styles.highlightColor)
}

func countriesView(s []countryStats, width int) string {
	return countriesContent(s, width)
}

func (m model) Init() tea.Cmd {
	return fetchData(m.tournament, m.conn)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
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
		if m.activeTab == 1 {
			var cmd tea.Cmd
			m.standingsTable, cmd = m.standingsTable.Update(msg)
			return m, cmd
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		if len(m.standings) > 0 {
			m.standingsTable = standingsTable(m.standings, m.width-8, m.height)
		}
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
			if row.PlayerCountry != "" {
				countryCounts[row.PlayerCountry]++
			}
		}
		for country, count := range countryCounts {
			m.countryStats = append(m.countryStats, countryStats{Country: country, Total: count})
		}
		deckCounts := map[string]int{}
		for _, row := range msg.items {
			deckCounts[row.Deck]++
		}
		for deck, count := range deckCounts {
			m.deckStats = append(m.deckStats, deckStats{Deck: deck, Total: count})
		}
		m.standingsTable = standingsTable(msg.items, m.width-8, m.height)
	}

	return m, nil
}

func (m model) View() tea.View {
	if m.styles == nil {
		return tea.NewView("")
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
	highlightColor := m.styles.highlightColor
	fillWidth := windowWidth - lipgloss.Width(row)
	if fillWidth > 0 {
		fill := lipgloss.NewStyle().Foreground(highlightColor).
			Render(strings.Repeat("─", fillWidth-1) + "┐")
		row = row + fill
	}

	var content string
	switch m.activeTab {
	// Overview Tab
	case 0:
		if m.err != nil {
			content = fmt.Sprintf("fetch error: %v", m.err)
		} else {
			content = overviewView(m, contentWidth)
		}

	// Standings Tab
	case 1:
		if m.err != nil {
			content = fmt.Sprintf("fetch error: %v", m.err)
		} else if len(m.standings) == 0 {
			content = "  Loading..."
		} else {
			content = m.standingsTable.View()
		}

	// Decks Tab
	case 2:
		if m.err != nil {
			content = fmt.Sprintf("fetch error: %v", m.err)
		} else {
			content = decksContent(m.deckStats, contentWidth)
		}

	// Countries Tab
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

	v := tea.NewView(s.doc.Render(doc.String()))
	v.AltScreen = true
	return v
}
