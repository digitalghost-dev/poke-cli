package search

import (
	"flag"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/digitalghost-dev/poke-cli/cmd/utils"
)

func SearchCommand() (string, error) {
	var output strings.Builder

	flag.Usage = func() {
		output.WriteString(
			utils.GenerateHelpMessage(
				utils.HelpConfig{
					Description: "Search for a resource by name or partial match.",
					CmdName:     "search",
				},
			),
		)
	}

	if utils.CheckHelpFlag(&output, flag.Usage) {
		return output.String(), nil
	}

	flag.Parse()

	if err := utils.ValidateArgs(os.Args, utils.Validator{MaxArgs: 3, CmdName: "search", RequireName: false, HasFlags: false}); err != nil {
		output.WriteString(err.Error())
		return output.String(), err
	}

	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		return "", err
	}

	return output.String(), nil
}

// Model structure
type Model struct {
	Choice         int
	Chosen         bool
	Quitting       bool
	TextInput      textinput.Model
	ShowResults    bool
	SearchResults  string
	WarningMessage string
}

func initialModel() Model {
	ti := textinput.New()
	ti.Placeholder = "type name..."
	ti.CharLimit = 20
	ti.Width = 20

	return Model{
		TextInput: ti,
	}
}

// Init initializes the program.
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles keypresses and updates the state.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
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
