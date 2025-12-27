package berry

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/digitalghost-dev/poke-cli/cmd/utils"
	"github.com/digitalghost-dev/poke-cli/connections"
	"github.com/digitalghost-dev/poke-cli/styling"
)

func BerryCommand() (string, error) {
	var output strings.Builder

	flag.Usage = func() {
		helpMessage := styling.HelpBorder.Render(
			"Get details about a specific berry.\n\n",
			styling.StyleBold.Render("USAGE:"),
			fmt.Sprintf("\n\t%s %s %s", "poke-cli", styling.StyleBold.Render("berry"), "[flag]"),
			"\n\n",
			styling.StyleBold.Render("FLAGS:"),
			fmt.Sprintf("\n\t%-30s %s", "-h, --help", "Prints out the help menu"),
		)
		output.WriteString(helpMessage)
	}

	flag.Parse()

	// Handle help flag
	if len(os.Args) == 3 && (os.Args[2] == "-h" || os.Args[2] == "--help") {
		flag.Usage()
		return output.String(), nil
	}

	// Validate arguments
	if err := utils.ValidateBerryArgs(os.Args); err != nil {
		output.WriteString(err.Error())
		return output.String(), err
	}

	if err := tableGeneration(); err != nil {
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
		}
	}

	m.table, bubbleCmd = m.table.Update(msg)

	// Keep the selected option in sync on every update
	if row := m.table.SelectedRow(); len(row) > 0 {
		name := row[0]
		if name != m.selectedOption {
			m.selectedOption = name
		}
	}

	return m, bubbleCmd
}

// View renders the current UI
func (m model) View() string {
	if m.quitting {
		return "\n Goodbye! \n"
	}

	selectedBerry := ""
	if row := m.table.SelectedRow(); len(row) > 0 {
		selectedBerry = BerryName(row[0]) + "\n---\n" + BerryEffect(row[0]) + "\n---\n" + BerryInfo(row[0]) + "\n---\nImage\n" + BerryImage(row[0])
	}

	leftPanel := styling.TypesTableBorder.Render(m.table.View())

	rightPanel := lipgloss.NewStyle().
		Width(50).
		Height(29).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#FFCC00")).
		Padding(1).
		Render(selectedBerry)

	screen := lipgloss.JoinHorizontal(lipgloss.Top, leftPanel, rightPanel)

	return fmt.Sprintf("Highlight a berry!\n%s\n%s",
		screen,
		styling.KeyMenu.Render("↑ (move up) • ↓ (move down)\nctrl+c | esc (quit)"))
}

func tableGeneration() error {
	namesList, err := connections.QueryBerryData(`
		SELECT 
		    UPPER(SUBSTR(name, 1, 1)) || SUBSTR(name, 2)
		FROM 
		    berries 
		ORDER BY 
		    name`)
	if err != nil {
		log.Fatalf("Failed to get berry names: %v", err)
	}

	rows := make([]table.Row, len(namesList))
	for i, n := range namesList {
		rows[i] = []string{n}
	}

	t := table.New(
		table.WithColumns([]table.Column{{Title: "Berry", Width: 16}}),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(28),
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
	_, err = tea.NewProgram(m).Run()

	if err != nil {
		return fmt.Errorf("error running program: %w", err)
	}

	return nil
}
