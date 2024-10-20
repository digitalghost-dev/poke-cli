package cmd

import (
	"flag"
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
			// Fully quiting the program while
			// the table is in selection mode
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
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m model) View() string {
	// When an option is selected, no longer display the table.
	if m.selectedOption != "" {
		return ""
	}
	// Otherwise, display the table
	return "Select a type! Hit 'Q' or 'CTRL-C' to quit.\n" + typesTableBorder.Render(m.table.View()) + "\n"
}

// Creating a separate function that handles all the logic for building the table.
func tableGeneration(endpoint string) table.Model {
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
	programModel, err := tea.NewProgram(m).Run()
	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	// Type assert to model and access the selected option
	finalModel, ok := programModel.(model)
	if !ok {
		fmt.Println("Error: could not retrieve final model")
		os.Exit(1)
	}

	// Check if the user quit the program by pressing 'Q' or 'CTRL-C'
	if finalModel.selectedOption == "quit" {
		// Return early to prevent further execution
		return t
	}

	// Proceed with the API call only if a type was selected
	typesName := strings.ToLower(finalModel.selectedOption)

	baseURL := "https://pokeapi.co/api/v2/"
	typeResponse, typeName, _ := connections.TypesApiCall(endpoint, typesName, baseURL)

	capitalizedString := cases.Title(language.English).String(typeName)

	pokemonCount := len(typeResponse.Pokemon)

	fmt.Printf("You selected Type: %s\nNumber of Pokémon with type: %d\n", capitalizedString, pokemonCount)
	return t
}

func TypesCommand() {
	styleBold = lipgloss.NewStyle().Bold(true)
	styleItalic = lipgloss.NewStyle().Italic(true)

	flag.Usage = func() {
		helpMessage := helpBorder.Render(
			styleBold.Render("USAGE:"),
			fmt.Sprintf("\n\t%s %s %s", "poke-cli", styleBold.Render("types"), "[flag]"),
			fmt.Sprintf("\n\t%-30s", "Get details about a specific typing"),
			fmt.Sprintf("\n\t%-30s", "----------"),
			fmt.Sprintf("\n\t%-30s", styleItalic.Render("Examples:")),
			fmt.Sprintf("\n\t%-30s", "poke-cli types"),
			fmt.Sprintf("\n\t%-30s", "A table will then display with the option to select a type."),
		)
		fmt.Println(helpMessage)
	}

	flag.Parse()

	err := ValidateTypesArgs(os.Args)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	endpoint := strings.ToLower(os.Args[1])[0:4]

	tableGeneration(endpoint)
}
