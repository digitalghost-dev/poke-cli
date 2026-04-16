package search

import (
	"flag"
	"os"
	"strings"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
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

// model structure
type model struct {
	Choice         int
	Chosen         bool
	Quitting       bool
	TextInput      textinput.Model
	ShowResults    bool
	SearchResults  string
	WarningMessage string
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "type name..."
	ti.CharLimit = 20
	ti.SetWidth(20)

	return model{
		TextInput: ti,
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
func (m model) View() tea.View {
	if m.Quitting {
		return tea.NewView("\n  Quitting search...\n\n")
	}
	if m.ShowResults {
		resultsView, _ := RenderInput(m) // Fetch results view
		return tea.NewView(resultsView)
	}
	if !m.Chosen {
		return tea.NewView(RenderSelection(m))
	}
	inputView, _ := RenderInput(m)
	return tea.NewView(inputView)
}
