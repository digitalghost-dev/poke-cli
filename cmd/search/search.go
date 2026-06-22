package search

import (
	"strings"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"github.com/digitalghost-dev/poke-cli/cmd/utils"
)

func SearchCommand(args []string) (string, error) {
	var output strings.Builder

	usage := func() {
		output.WriteString(
			utils.GenerateHelpMessage(
				utils.HelpConfig{
					Description: "Search for a resource by name or partial match.",
					CmdName:     "search",
				},
			),
		)
	}

	if utils.CheckHelpFlag(args, usage) {
		return output.String(), nil
	}

	if err := utils.ValidateArgs(
		args,
		utils.Validator{MaxArgs: 2, CmdName: "search", RequireName: false, HasFlags: false},
	); err != nil {
		output.WriteString(err.Error())
		return output.String(), err
	}

	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		return "", err
	}

	return output.String(), nil
}

// model structure
type model struct {
	choice         int
	chosen         bool
	quitting       bool
	textInput      textinput.Model
	showResults    bool
	searchResults  string
	warningMessage string
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "type name..."
	ti.CharLimit = 20
	ti.SetWidth(20)

	return model{
		textInput: ti,
	}
}

// Init initializes the program.
func (m model) Init() tea.Cmd {
	return nil
}

// Update handles keypresses and updates the state.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.quitting = true
			return m, tea.Quit
		}
	}

	if !m.chosen {
		return UpdateSelection(msg, m)
	}
	return UpdateInput(msg, m)
}

// View renders the correct UI screen.
func (m model) View() tea.View {
	if m.quitting {
		return tea.NewView("\n  Quitting search...\n\n")
	}
	if m.showResults {
		resultsView, _ := RenderInput(m) // Fetch results view
		return tea.NewView(resultsView)
	}
	if !m.chosen {
		return tea.NewView(RenderSelection(m))
	}
	inputView, _ := RenderInput(m)
	return tea.NewView(inputView)
}
