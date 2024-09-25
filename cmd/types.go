package cmd

import (
	"fmt"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/digitalghost-dev/poke-cli/connections"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"os"
	"strings"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("#FFCC00"))

type model struct {
	table          table.Model
	selectedOption string // Track the selected option
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			selectedRow := m.table.SelectedRow()
			m.selectedOption = selectedRow[0]
			return m, tea.Batch(
				tea.Quit,
			)
		}
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.selectedOption != "" {
		return ""
	}
	// Otherwise, display the table
	return "Select a type! Hit 'Q' or 'CTRL-C' to quit.\n" + baseStyle.Render(m.table.View()) + "\n"
}

func TypesCommand() {
	columns := []table.Column{
		{Title: "Type", Width: 16},
	}

	rows := []table.Row{
		{"Normal"},
		{"Fire"},
		{"Water"},
		{"Electric"},
		{"Grass"},
		{"Ice"},
		{"Fighting"},
		{"Poison"},
		{"Ground"},
		{"Flying"},
		{"Psychic"},
		{"Bug"},
		{"Rock"},
		{"Ghost"},
		{"Dragon"},
		{"Steel"},
		{"Fairy"},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(7),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("#FFCC00")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("#000")).
		Background(lipgloss.Color("#FFCC00")).
		Bold(false)
	t.SetStyles(s)

	m := model{table: t}

	// Run the program and capture the final state
	programModel, err := tea.NewProgram(m).Run()
	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	// Type assert to model and access the selected option
	finalModel := programModel.(model)

	// Convert the selected option to lowercase
	typesName := strings.ToLower(finalModel.selectedOption)

	// Use the TypesApiCall to fetch the details of the selected type
	_, typeName, typeID := connections.TypesApiCall("type", typesName, "https://pokeapi.co/api/v2/")
	capitalizedString := cases.Title(language.English).String(typeName)

	// Display the result
	fmt.Printf("You selected Type: %s\nType ID: %d\n", capitalizedString, typeID)
}
