package search

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// MainExecution is the main function for the "search" command.
func MainExecution() {
	initialModel := Model{
		Choice:    0,
		Chosen:    false,
		Quitting:  false,
		TextInput: textinput.New(),
	}
	initialModel.TextInput.Placeholder = "type here"
	initialModel.TextInput.CharLimit = 20
	initialModel.TextInput.Width = 20

	p := tea.NewProgram(initialModel)
	if _, err := p.Run(); err != nil {
		fmt.Println("could not start program:", err)
	}
}

// Model structure
type Model struct {
	Choice         int             // Index of the selected search category (e.g., 0 = Pokémon, 1 = Ability)
	Chosen         bool            // Whether a category has been chosen yet
	Quitting       bool            // Flag to indicate if the program is quitting
	TextInput      textinput.Model // The text input field used for search input
	ShowResults    bool            // Whether to display the search results screen
	SearchResults  string          // The formatted search results to be displayed
	WarningMessage string          // A warning message to show (e.g., empty input)
}

// Init initializes the program.
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles keypresses and updates the state.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		k := msg.String()
		if k == "esc" || k == "ctrl+c" {
			m.Quitting = true
			return m, tea.Quit
		}
	}

	if !m.Chosen {
		return UpdateSelection(msg, m)
	}
	return UpdateInput(msg, m)
}

// View renders the correct UI screen.
func (m Model) View() string {
	if m.Quitting {
		return "\n  Quitting search...\n\n"
	}
	if m.ShowResults {
		resultsView, _ := RenderInput(m) // Fetch results view
		return resultsView
	}
	if !m.Chosen {
		return RenderSelection(m)
	}
	inputView, _ := RenderInput(m)
	return inputView
}
