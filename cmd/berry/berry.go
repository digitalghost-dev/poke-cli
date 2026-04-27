package berry

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/digitalghost-dev/poke-cli/cmd/utils"
	"github.com/digitalghost-dev/poke-cli/connections"
	"github.com/digitalghost-dev/poke-cli/styling"
)

func BerryCommand() (string, error) {
	var output strings.Builder

	flag.Usage = func() {
		output.WriteString(
			utils.GenerateHelpMessage(
				utils.HelpConfig{
					Description: "Get details about a specific berry.",
					CmdName:     "berry",
				},
			),
		)
	}

	if utils.CheckHelpFlag(&output, flag.Usage) {
		return output.String(), nil
	}

	flag.Parse()

	// Validate arguments
	if err := utils.ValidateArgs(os.Args, utils.Validator{MaxArgs: 4, CmdName: "berry", RequireName: false, HasFlags: false}); err != nil {
		output.WriteString(err.Error())
		return output.String(), err
	}

	if len(os.Args) > 2 {
		berryName := styling.CapitalizeResourceName(os.Args[2])
		if !berryExists(berryName) {
			err := fmt.Errorf("berry %q not found", os.Args[2])
			output.WriteString(utils.FormatError(err.Error()))
			return output.String(), err
		}
		containers := berryContainers(berryName)
		output.WriteString(containers)
		output.WriteString("\n")
		return output.String(), nil
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
	case tea.KeyPressMsg:
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
func (m model) View() tea.View {
	if m.quitting {
		return tea.NewView("\n Goodbye! \n")
	}

	selectedBerry := ""
	if row := m.table.SelectedRow(); len(row) > 0 {
		selectedBerry = berryName(row[0]) + "\n---\n" + berryEffect(row[0]) + "\n---\n" + berryInfo(row[0]) + "\n---\nImage\n" + berryImage(row[0])
	}

	leftPanel := styling.TypesTableBorder.Render(m.table.View())

	rightPanel := lipgloss.NewStyle().
		Width(52).
		Height(29).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styling.YellowColor).
		Padding(1).
		Render(selectedBerry)

	screen := lipgloss.JoinHorizontal(lipgloss.Top, leftPanel, rightPanel)

	return tea.NewView(fmt.Sprintf("Highlight a berry!\n%s\n%s",
		screen,
		styling.KeyMenu.Render("↑ (move up) • ↓ (move down)\nctrl+c | esc (quit)")))
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
		table.WithWidth(16),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(styling.YellowColor).
		BorderBottom(true)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("#000")).
		Background(styling.YellowColor)
	t.SetStyles(s)

	m := model{table: t}
	_, err = tea.NewProgram(m).Run()

	if err != nil {
		return fmt.Errorf("error running program: %w", err)
	}

	return nil
}

func berryContainers(name string) string {
	header := lipgloss.NewStyle().Bold(true).PaddingBottom(1).Render(
		styling.StyleBold.Render(styling.CapitalizeResourceName(name)),
	)
	infoContent := lipgloss.JoinVertical(lipgloss.Top, header, berryInfo(name), "\n"+berryEffect(name))
	imageContent := berryImage(name)

	boxStyle := lipgloss.NewStyle().
		Padding(1, 2).
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(styling.YellowColor).
		Width(34)

	// Render both without height constraints to measure natural heights.
	infoH := lipgloss.Height(boxStyle.Render(infoContent))
	imageH := lipgloss.Height(boxStyle.Render(imageContent))

	// Pad the shorter content with blank lines before final render so both boxes match.
	if infoH < imageH {
		infoContent += strings.Repeat("\n", imageH-infoH)
	} else if imageH < infoH {
		imageContent += strings.Repeat("\n", infoH-imageH)
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, boxStyle.Render(infoContent), boxStyle.Render(imageContent))
}
