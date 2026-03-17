package card

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/digitalghost-dev/poke-cli/styling"
)

var seriesIDMap = map[string]string{
	"Mega Evolution":   "me",
	"Scarlet & Violet": "sv",
	"Sword & Shield":   "swsh",
	"Sun & Moon":       "sm",
}

type SeriesModel struct {
	List     list.Model
	Choice   string
	SeriesID string
	Quitting bool
}

func (m SeriesModel) Init() tea.Cmd {
	return nil
}

func (m SeriesModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.Quitting = true
			return m, tea.Quit
		case "enter":
			i, ok := m.List.SelectedItem().(styling.Item)
			if ok {
				m.Choice = string(i)
				m.SeriesID = seriesIDMap[string(i)]
			}
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.List.SetWidth(msg.Width)
		return m, nil
	}

	var cmd tea.Cmd
	m.List, cmd = m.List.Update(msg)
	return m, cmd
}

func (m SeriesModel) View() string {
	if m.Quitting {
		return "\n  Quitting card search...\n\n"
	}
	if m.Choice != "" {
		return styling.QuitTextStyle.Render("Series selected:", m.Choice)
	}

	return "\n" + m.List.View()
}

func SeriesList() SeriesModel {
	items := []list.Item{
		styling.Item("Mega Evolution"),
		styling.Item("Scarlet & Violet"),
		styling.Item("Sword & Shield"),
		styling.Item("Sun & Moon"),
	}

	const listWidth = 20
	const listHeight = 12

	l := list.New(items, styling.ItemDelegate{}, listWidth, listHeight)
	l.Title = "First, pick a series"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = styling.TitleStyle
	l.Styles.PaginationStyle = styling.PaginationStyle
	l.Styles.HelpStyle = styling.HelpStyle

	return SeriesModel{List: l}
}
