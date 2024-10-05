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
	return "Select a type! Hit 'Q' or 'CTRL-C' to quit.\n" + baseStyle.Render(m.table.View()) + "\n"
}

func selectionResult(endpoint, typesName string) error {
	baseURL := "https://pokeapi.co/api/v2/"

	typeResponse, typeName, _ := connections.TypesApiCall(endpoint, typesName, baseURL)

	capitalizedString := cases.Title(language.English).String(typeName)

	pokemonCount := len(typeResponse.Pokemon)

	fmt.Printf("You selected Type: %s\nNumber of Pokémon with type: %d\n", capitalizedString, pokemonCount)
	return nil
}

func TypesCommand() {
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

	args := os.Args

	err := ValidateTypesArgs(args)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	if len(args) == 3 && (args[2] == "-h" || args[2] == "--help") {
		flag.Usage()
		return
	}

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
	finalModel := programModel.(model)
	typesName := strings.ToLower(finalModel.selectedOption)

	// Extract the first 4 characters of the endpoint from the argument
	endpoint := strings.ToLower(args[1])[0:4]

	// Call the selectionTable function with the selected type and endpoint
	if err := selectionResult(endpoint, typesName); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
