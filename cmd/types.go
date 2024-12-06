package cmd

import (
	"flag"
	"fmt"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/term"
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
			// Quit the program when in selection mode
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

// Function to display type details after a type is selected
func displayTypeDetails(typesName string, endpoint string) {

	// Setting up variables to style the list
	var columnWidth = 13
	var subtle = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
	var list = lipgloss.NewStyle().Border(lipgloss.NormalBorder(), false, true, false, false).BorderForeground(subtle).MarginRight(2).Height(8).Width(columnWidth + 1)
	var listHeader = lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).BorderBottom(true).BorderForeground(subtle).MarginRight(2).Render
	var listItem = lipgloss.NewStyle().Render
	var docStyle = lipgloss.NewStyle().Padding(1, 1, 1, 1)

	baseURL := "https://pokeapi.co/api/v2/"
	typesStruct, typeName, _ := connections.TypesApiCall(endpoint, typesName, baseURL)

	// Format selected type
	selectedType := cases.Title(language.English).String(typeName)
	coloredType := lipgloss.NewStyle().Foreground(lipgloss.Color(getTypeColor(typeName))).Render(selectedType)

	fmt.Printf("You selected the %s type.\nNumber of Pokémon with type: %d\nNumber of moves with type: %d\n", coloredType, len(typesStruct.Pokemon), len(typesStruct.Moves))
	fmt.Println("----------")
	fmt.Println(styleBold.Render("Damage Chart:"))

	physicalWidth, _, _ := term.GetSize(uintptr(int(os.Stdout.Fd())))
	doc := strings.Builder{}

	// Helper function to build list items
	buildListItems := func(items []struct{ Name, URL string }) string {
		var itemList []string
		for _, item := range items {
			color := getTypeColor(item.Name)
			coloredStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(color))
			coloredItem := coloredStyle.Render(cases.Title(language.English).String(item.Name))
			itemList = append(itemList, listItem(coloredItem))
		}
		return lipgloss.JoinVertical(lipgloss.Left, itemList...)
	}

	// Render lists based on Damage Relations
	lists := lipgloss.JoinHorizontal(lipgloss.Top,
		list.Render(
			lipgloss.JoinVertical(lipgloss.Left,
				listHeader("Weakness"),
				buildListItems([]struct{ Name, URL string }(typesStruct.DamageRelations.DoubleDamageFrom)),
			),
		),
		list.Width(columnWidth).Render(
			lipgloss.JoinVertical(lipgloss.Left,
				listHeader("x2 Damage"),
				buildListItems([]struct{ Name, URL string }(typesStruct.DamageRelations.DoubleDamageTo)),
			),
		),
		list.Width(columnWidth).Render(
			lipgloss.JoinVertical(lipgloss.Left,
				listHeader("Resists"),
				buildListItems([]struct{ Name, URL string }(typesStruct.DamageRelations.HalfDamageFrom)),
			),
		),
		list.Width(columnWidth).Render(
			lipgloss.JoinVertical(lipgloss.Left,
				listHeader("x0.5 Damage"),
				buildListItems([]struct{ Name, URL string }(typesStruct.DamageRelations.HalfDamageTo)),
			),
		),
		list.Width(columnWidth).Render(
			lipgloss.JoinVertical(lipgloss.Left,
				listHeader("Immune"),
				buildListItems([]struct{ Name, URL string }(typesStruct.DamageRelations.NoDamageFrom)),
			),
		),
		list.Width(columnWidth).Render(
			lipgloss.JoinVertical(lipgloss.Left,
				listHeader("x0 Damage"),
				buildListItems([]struct{ Name, URL string }(typesStruct.DamageRelations.NoDamageTo)),
			),
		),
	)

	// Append lists to document
	doc.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, lists))

	if physicalWidth > 0 {
		docStyle = docStyle.MaxWidth(physicalWidth)
	}

	// Print the rendered document
	fmt.Println(docStyle.Render(doc.String()))
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
		displayTypeDetails(typesName, endpoint) // Call function to display type details
	}

	return t
}

func TypesCommand() {

	flag.Usage = func() {
		helpMessage := helpBorder.Render(
			"Get details about a specific typing.\n\n",
			styleBold.Render("USAGE:"),
			fmt.Sprintf("\n\t%s %s %s", "poke-cli", styleBold.Render("types"), "[flag]"),
			"\n\n",
			styleBold.Render("FLAGS:"),
			fmt.Sprintf("\n\t%-30s %s", "-h, --help", "Prints out the help menu."),
		)
		fmt.Println(helpMessage)
	}

	flag.Parse()

	if err := ValidateTypesArgs(os.Args); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	endpoint := strings.ToLower(os.Args[1])[0:4]
	tableGeneration(endpoint)
}
