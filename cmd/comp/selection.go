package comp

import (
	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"github.com/digitalghost-dev/poke-cli/styling"
)

var compIDMap = map[string]string{
	"TCG Competition Data":   "tcg",
	"VGC Competition Data":   "vgc",
	"Pokémon Champions Data": "champions",
}

type pickerModel struct {
	list     list.Model
	choice   string
	compID   string
	quitting bool
}

func (m pickerModel) Init() tea.Cmd {
	return nil
}

func (m pickerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.quitting = true
			return m, tea.Quit
		case "enter":
			i, ok := m.list.SelectedItem().(styling.Item)
			if ok {
				m.choice = string(i)
				m.compID = compIDMap[string(i)]
			}
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m pickerModel) View() tea.View {
	var content string
	if m.quitting {
		content = "\n  Quitting comp command...\n\n"
	} else if m.choice != "" {
		content = styling.QuitTextStyle.Render("Viewing data for:", m.choice)
	} else {
		content = "\n" + m.list.View()
	}

	v := tea.NewView(content)
	v.AltScreen = true
	return v
}

func CompList() pickerModel {
	items := []list.Item{
		styling.Item("TCG Competition Data"),
		styling.Item("VGC Competition Data"),
		styling.Item("Pokémon Champions Data"),
	}

	const listWidth = 24
	const listHeight = 12

	l := list.New(items, styling.ItemDelegate{}, listWidth, listHeight)
	l.Title = "Pick a competition type"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = styling.TitleStyle
	l.Styles.PaginationStyle = styling.PaginationStyle
	l.Styles.HelpStyle = styling.HelpStyle

	return pickerModel{list: l}
}
