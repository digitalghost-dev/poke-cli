package types

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/digitalghost-dev/poke-cli/cmd/utils"
	"github.com/digitalghost-dev/poke-cli/styling"
)

func TypesCommand() (string, error) {
	var output strings.Builder

	flag.Usage = func() {
		helpMessage := styling.HelpBorder.Render(
			"Get details about a specific typing.\n\n",
			styling.StyleBold.Render("USAGE:"),
			fmt.Sprintf("\n\t%s %s %s", "poke-cli", styling.StyleBold.Render("types"), "[flag]"),
			"\n\n",
			styling.StyleBold.Render("FLAGS:"),
			fmt.Sprintf("\n\t%-30s %s", "-h, --help", "Prints out the help menu"),
		)
		output.WriteString(helpMessage)
	}

	flag.Parse()

	// Handle help flag
	if utils.CheckHelpFlag(&output, flag.Usage) {
		return output.String(), nil
	}

	// Validate arguments
	if err := utils.ValidateTypesArgs(os.Args); err != nil {
		output.WriteString(err.Error())
		return output.String(), err
	}

	endpoint := strings.ToLower(os.Args[1])[0:4]
	if err := runTypeSelectionTable(endpoint); err != nil {
		output.WriteString(err.Error())
		return output.String(), err
	}

	return output.String(), nil
}

type model struct {
	quitting       bool
	table          table.Model
	selectedOption string
}

// Init initializes the model
func (m model) Init() tea.Cmd {
	return nil
}

// Update handles user input and updates the model state
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var bubbleCmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		case "enter":
			// User selected a type
			m.selectedOption = m.table.SelectedRow()[0]
			return m, tea.Quit
		}
	}

	// Handle other updates (like navigation)
	m.table, bubbleCmd = m.table.Update(msg)
	return m, bubbleCmd
}

// View renders the current UI
func (m model) View() string {
	if m.quitting {
		return "\n  Goodbye! \n"
	}

	// Don't render anything if a selection has been made
	if m.selectedOption != "" {
		return ""
	}

	// Render the type selection table with instructions
	return fmt.Sprintf("Select a type!\n%s\n%s",
		styling.TypesTableBorder.Render(m.table.View()),
		styling.KeyMenu.Render("↑ (move up) • ↓ (move down)\nenter (select) • ctrl+c | esc (quit)"))
}

func createTypeSelectionTable() model {
	types := []string{"Normal", "Fire", "Water", "Electric", "Grass", "Ice",
		"Fighting", "Poison", "Ground", "Flying", "Psychic", "Bug", "Dark",
		"Rock", "Ghost", "Dragon", "Steel", "Fairy"}

	rows := make([]table.Row, len(types))
	for i, t := range types {
		rows[i] = []string{t}
	}

	tbl := table.New(
		table.WithColumns([]table.Column{{Title: "Type", Width: 16}}),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("#FFCC00")).
		BorderBottom(true)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("#000")).
		Background(lipgloss.Color("#FFCC00"))
	tbl.SetStyles(s)

	return model{table: tbl}
}

func runTypeSelectionTable(endpoint string) error {
	m := createTypeSelectionTable()

	programModel, err := tea.NewProgram(m).Run()
	if err != nil {
		return fmt.Errorf("error running program: %w", err)
	}

	if finalModel, ok := programModel.(model); ok && finalModel.selectedOption != "" {
		DamageTable(strings.ToLower(finalModel.selectedOption), endpoint)
	}

	return nil
}
