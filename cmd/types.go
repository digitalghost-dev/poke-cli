package cmd

import (
	"fmt"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"os"
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
			if len(selectedRow) > 0 {
				m.selectedOption = selectedRow[0]
				return m, tea.Batch(
					tea.Quit,
				)
			}
		}
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.selectedOption != "" {
		return fmt.Sprintf("Selected type: %s\n%s", m.selectedOption)
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
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
