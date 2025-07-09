package search

import (
	"flag"
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/digitalghost-dev/poke-cli/cmd/utils"
	"github.com/digitalghost-dev/poke-cli/styling"
	"os"
)

func SearchCommand() {
	flag.Usage = func() {
		helpMessage := styling.HelpBorder.Render(
			"Search for a resource by name or partial match.\n\n",
			styling.StyleBold.Render("USAGE:"),
			fmt.Sprintf("\n\t%s %s %s", "poke-cli", styling.StyleBold.Render("search"), "[flag]"),
			"\n\n",
			styling.StyleBold.Render("FLAGS:"),
			fmt.Sprintf("\n\t%-15s %s", "-h, --help", "Prints out the help menu.\n\n"),
			styling.StyleItalic.Render("Supports prefix matching using ^ (example: ^char → charizard)"),
		)
		fmt.Println(helpMessage)
	}

	flag.Parse()

	if len(os.Args) == 3 && (os.Args[2] == "-h" || os.Args[2] == "--help") {
		flag.Usage()
		return
	}

	if err := utils.ValidateSearchArgs(os.Args); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	p := tea.NewProgram(initialModel())
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
