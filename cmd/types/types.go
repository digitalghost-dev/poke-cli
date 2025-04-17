package types

import (
	"flag"
	"fmt"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/digitalghost-dev/poke-cli/cmd"
	"github.com/digitalghost-dev/poke-cli/styling"
	"os"
	"strings"
)

type model struct {
	table          table.Model
	selectedOption string // Track the selected option
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var bubbleCmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			m.selectedOption = "quit"
			return m, tea.Quit
		case "enter":
			selectedRow := m.table.SelectedRow()
			m.selectedOption = selectedRow[0]
			return m, tea.Batch(
				tea.Quit,
			)
		}
	}
	m.table, bubbleCmd = m.table.Update(msg)
	return m, bubbleCmd
}

func (m model) View() string {
	// When an option is selected, no longer display the table.
	if m.selectedOption != "" {
		return ""
	}
	// Otherwise, display the table
	return "Select a type!\n" +
		styling.TypesTableBorder.Render(m.table.View()) +
		"\n" +
		styling.KeyMenu.Render("↑ (move up) • ↓ (move down)\nenter (select) • ctrl+c | esc (quit)")
}

// Function that generates and handles the type selection table
func tableGeneration(endpoint string) table.Model {
	columns := []table.Column{{Title: "Type", Width: 16}}
	rows := []table.Row{
		{"Normal"}, {"Fire"}, {"Water"}, {"Electric"}, {"Grass"}, {"Ice"},
		{"Fighting"}, {"Poison"}, {"Ground"}, {"Flying"}, {"Psychic"}, {"Bug"},
		{"Rock"}, {"Ghost"}, {"Dragon"}, {"Steel"}, {"Fairy"},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(7),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.BorderStyle(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("#FFCC00")).BorderBottom(true)
	s.Selected = s.Selected.Foreground(lipgloss.Color("#000")).Background(lipgloss.Color("#FFCC00"))
	t.SetStyles(s)

	m := model{table: t}
	programModel, err := tea.NewProgram(m).Run()
	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	// Access the selected option from the model
	finalModel, ok := programModel.(model)
	if !ok {
		fmt.Println("Error: could not retrieve final model")
		os.Exit(1)
	}

	if finalModel.selectedOption != "quit" {
		typesName := strings.ToLower(finalModel.selectedOption)
		DamageTable(typesName, endpoint) // Call function to display type details
	}

	return t
}

func TypesCommand() {
	flag.Usage = func() {
		helpMessage := styling.HelpBorder.Render(
			"Get details about a specific typing.\n\n",
			styling.StyleBold.Render("USAGE:"),
			fmt.Sprintf("\n\t%s %s %s", "poke-cli", styling.StyleBold.Render("types"), "[flag]"),
			"\n\n",
			styling.StyleBold.Render("FLAGS:"),
			fmt.Sprintf("\n\t%-30s %s", "-h, --help", "Prints out the help menu."),
		)
		fmt.Println(helpMessage)
	}

	flag.Parse()

	if err := cmd.ValidateTypesArgs(os.Args); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	endpoint := strings.ToLower(os.Args[1])[0:4]
	tableGeneration(endpoint)
}
