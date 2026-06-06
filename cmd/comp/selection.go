package comp

import (
	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"github.com/digitalghost-dev/poke-cli/styling"
)

var compIDMap = map[string]string{
	"TCG": "tcg",
	"VGC": "vgc",
}

type pickerModel struct {
	List     list.Model
	Choice   string
	CompID   string
	Quitting bool
}

func (m pickerModel) Init() tea.Cmd {
	return nil
}

func (m pickerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.Quitting = true
			return m, tea.Quit
		case "enter":
			i, ok := m.List.SelectedItem().(styling.Item)
			if ok {
				m.Choice = string(i)
				m.CompID = compIDMap[string(i)]
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

func (m pickerModel) View() tea.View {
	var content string
	if m.Quitting {
		content = "\n  Quitting comp command...\n\n"
	} else if m.Choice != "" {
		content = styling.QuitTextStyle.Render("Viewing data for:", m.Choice)
	} else {
		content = "\n" + m.List.View()
	}

	v := tea.NewView(content)
	v.AltScreen = true
	return v
}

func CompList() pickerModel {
	items := []list.Item{
		styling.Item("TCG"),
		styling.Item("VGC"),
	}

	const listWidth = 20
	const listHeight = 12

	l := list.New(items, styling.ItemDelegate{}, listWidth, listHeight)
	l.Title = "Pick a competition type"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = styling.TitleStyle
	l.Styles.PaginationStyle = styling.PaginationStyle
	l.Styles.HelpStyle = styling.HelpStyle

	return pickerModel{List: l}
}
