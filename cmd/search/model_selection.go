package search

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/digitalghost-dev/poke-cli/styling"
)

// UpdateSelection handles navigation in the selection menu.
func UpdateSelection(msg tea.Msg, m Model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "down":
			m.Choice++
			if m.Choice > 2 {
				m.Choice = 2
			}
		case "up":
			m.Choice--
			if m.Choice < 0 {
				m.Choice = 0
			}
		case "enter":
			m.Chosen = true
			m.TextInput.Focus()
			return m, textinput.Blink
		}
	}
	return m, nil
}

// RenderSelection renders the selection menu.
func RenderSelection(m Model) string {
	c := m.Choice
	greeting := styling.StyleItalic.Render("Search for a resource and return a matching selection table")
	choices := fmt.Sprintf(
		"%s\n%s\n%s",
		checkbox("Ability", c == 0),
		checkbox("Move", c == 1),
		checkbox("Pokémon", c == 2),
	)
	help := styling.KeyMenu.Render("↑ (move up) • ↓ (move down)\nenter (select) • ctrl+c | esc (quit)")

	return greeting + "\n\nWhat would you like to search?\n\n" + choices + "\n\n" + help + "\n"
}

// checkbox renders checkboxes for the selection menu.
func checkbox(label string, checked bool) string {
	if checked {
		return styling.CheckboxStyle.Render("> " + label)
	}
	return label
}
