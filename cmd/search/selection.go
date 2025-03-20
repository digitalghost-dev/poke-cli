package search

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Styling
var (
	keywordStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("211"))
	checkboxStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))
	mainStyle     = lipgloss.NewStyle().MarginLeft(2)
)

// SearchCommand "main" function for the "search" command.
// Passed to "cli.go" as an available command.
func SearchCommand() {
	initialModel := model{
		Choice:    0,
		Chosen:    false,
		Quitting:  false,
		textInput: textinput.New(),
	}
	initialModel.textInput.Placeholder = "Enter name..."
	initialModel.textInput.CharLimit = 20
	initialModel.textInput.Width = 20

	p := tea.NewProgram(initialModel)
	if _, err := p.Run(); err != nil {
		fmt.Println("could not start program:", err)
	}
}

// Model structure
type model struct {
	Choice    int
	Chosen    bool
	Quitting  bool
	textInput textinput.Model
}

// Initialize the program
func (m model) Init() tea.Cmd {
	return nil
}

// Handle updates (key presses and commands)
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		k := msg.String()
		if k == "q" || k == "esc" || k == "ctrl+c" {
			m.Quitting = true
			return m, tea.Quit
		}
	}

	if !m.Chosen {
		return updateChoices(msg, m)
	}
	return updateTextInput(msg, m)
}

// Render the view
func (m model) View() string {
	var s string
	if m.Quitting {
		return "\n  See you later!\n\n"
	}
	if !m.Chosen {
		s = choicesView(m)
	} else {
		s = chosenView(m)
	}
	return mainStyle.Render("\n" + s + "\n\n")
}

// Handle navigation and selection in the choice view
func updateChoices(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "j", "down":
			m.Choice++
			if m.Choice > 1 {
				m.Choice = 1
			}
		case "k", "up":
			m.Choice--
			if m.Choice < 0 {
				m.Choice = 0
			}
		case "enter":
			m.Chosen = true
			m.textInput.Focus()
			return m, textinput.Blink
		}
	}
	return m, nil
}

// Handle text input updates
func updateTextInput(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter: // User presses enter after input
			m.Quitting = true
			fmt.Printf("User searched for: %s\n", m.textInput.Value()) // Optional: Print entered text before quitting
			return m, tea.Quit
		case tea.KeyEsc: // Escape to go back
			m.Chosen = false
			m.textInput.Reset()
			return m, nil
		}
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

// Choice View
func choicesView(m model) string {
	c := m.Choice
	tpl := "What would you like to search?\n\n%s\n\n"
	choices := fmt.Sprintf(
		"%s\n%s",
		checkbox("Pokémon", c == 0),
		checkbox("Ability", c == 1),
	)
	return fmt.Sprintf(tpl, choices)
}

// Text Input View
func chosenView(m model) string {
	var msg string

	switch m.Choice {
	case 0:
		msg = "Enter a Pokémon name:"
	case 1:
		msg = "Enter an Ability name:"
	default:
		msg = "Enter your search query:"
	}

	return fmt.Sprintf(
		"%s\n\n%s\n\n%s",
		msg,
		m.textInput.View(),
		"(Press Enter to confirm)",
	)
}

// Checkbox function to display selection
func checkbox(label string, checked bool) string {
	if checked {
		return checkboxStyle.Render("[x] " + label)
	}
	return fmt.Sprintf("[ ] %s", label)
}
