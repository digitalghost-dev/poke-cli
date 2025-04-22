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

func TypesCommand() {
	flag.Usage = func() {
		helpMessage := styling.HelpBorder.Render(
			"Get details about a specific typing.\n\n",
			styling.StyleBold.Render("USAGE:"),
			fmt.Sprintf("\n\t%s %s %s", "poke-cli", styling.StyleBold.Render("types"), "[flag]"),
			"\n\n",
			styling.StyleBold.Render("FLAGS:"),
			fmt.Sprintf("\n\t%-30s %s", "-h, --help", "Prints out the help menu"),
		)
		fmt.Println(helpMessage)
	}

	flag.Parse()

	if len(os.Args) == 3 && (os.Args[2] == "-h" || os.Args[2] == "--help") {
		flag.Usage()
		return
	}

	if err := cmd.ValidateTypesArgs(os.Args); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	endpoint := strings.ToLower(os.Args[1])[0:4]
	tableGeneration(endpoint)
}

type model struct {
	table          table.Model
	selectedOption string
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var bubbleCmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "ctrl+c":
			return m, tea.Quit
		case "enter":
			m.selectedOption = m.table.SelectedRow()[0]
			return m, tea.Quit
		}
	}

	m.table, bubbleCmd = m.table.Update(msg)
	return m, bubbleCmd
}

func (m model) View() string {
	if m.selectedOption != "" {
		return ""
	}

	return fmt.Sprintf("Select a type!\n%s\n%s",
		styling.TypesTableBorder.Render(m.table.View()),
		styling.KeyMenu.Render("↑ (move up) • ↓ (move down)\nenter (select) • ctrl+c | esc (quit)"))
}

// Function that generates and handles the type selection table
func tableGeneration(endpoint string) table.Model {
	types := []string{"Normal", "Fire", "Water", "Electric", "Grass", "Ice",
		"Fighting", "Poison", "Ground", "Flying", "Psychic", "Bug",
		"Rock", "Ghost", "Dragon", "Steel", "Fairy"}

	rows := make([]table.Row, len(types))
	for i, t := range types {
		rows[i] = []string{t}
	}

	t := table.New(
		table.WithColumns([]table.Column{{Title: "Type", Width: 16}}),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(7),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("#FFCC00")).
		BorderBottom(true)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("#000")).
		Background(lipgloss.Color("#FFCC00"))
	t.SetStyles(s)

	m := model{table: t}
	if programModel, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	} else if finalModel, ok := programModel.(model); ok && finalModel.selectedOption != "quit" {
		DamageTable(strings.ToLower(finalModel.selectedOption), endpoint)
	}

	return t
}
