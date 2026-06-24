package types

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/digitalghost-dev/poke-cli/cmd/utils"
	"github.com/digitalghost-dev/poke-cli/styling"
)

func TypesCommand(args []string) (string, error) {
	var output strings.Builder

	usage := func() {
		output.WriteString(
			utils.GenerateHelpMessage(
				utils.HelpConfig{
					Description: "Get details about a specific typing.",
					CmdName:     "types",
				},
			),
		)
	}

	if utils.CheckHelpFlag(args, usage) {
		return output.String(), nil
	}

	// Validate arguments
	if err := utils.ValidateArgs(
		args,
		utils.Validator{MaxArgs: 2, CmdName: "types", RequireName: false, HasFlags: false},
	); err != nil {
		output.WriteString(err.Error())
		return output.String(), err
	}

	const endpoint = "type"
	chart, err := runTypeSelectionTable(endpoint)
	if err != nil {
		output.WriteString(err.Error())
		return output.String(), err
	}
	output.WriteString(chart)

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
	case tea.KeyPressMsg:
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
func (m model) View() tea.View {
	if m.quitting {
		return tea.NewView("\n  Goodbye! \n")
	}

	// Don't render anything if a selection has been made
	if m.selectedOption != "" {
		return tea.NewView("")
	}

	// Render the type selection table with instructions
	return tea.NewView(fmt.Sprintf("Select a type!\n%s\n%s",
		styling.TypesTableBorder.Render(m.table.View()),
		styling.KeyMenu.Render("↑ (move up) • ↓ (move down)\nenter (select) • ctrl+c | esc (quit)")))
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
		table.WithWidth(16),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(styling.ThemeColor).
		BorderBottom(true)
	s.Selected = s.Selected.
		Foreground(styling.ContrastText(styling.ThemeColor)).
		Background(styling.ThemeColor)
	tbl.SetStyles(s)

	return model{table: tbl}
}

func runTypeSelectionTable(endpoint string) (string, error) {
	m := createTypeSelectionTable()

	programModel, err := tea.NewProgram(m).Run()
	if err != nil {
		return "", fmt.Errorf("error running program: %w", err)
	}

	if finalModel, ok := programModel.(model); ok && finalModel.selectedOption != "" {
		return DamageTable(strings.ToLower(finalModel.selectedOption), endpoint)
	}

	return "", nil
}
