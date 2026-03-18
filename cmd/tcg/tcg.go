package tcg

import (
	"flag"
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/digitalghost-dev/poke-cli/cmd/utils"
)

func TcgCommand() (string, error) {
	var output strings.Builder

	flag.Usage = func() {
		output.WriteString(
			utils.GenerateHelpMessage(
				utils.HelpConfig{
					Description: "Get details about TCG tournaments.",
					CmdName:     "tcg",
				},
			),
		)
	}

	if utils.CheckHelpFlag(&output, flag.Usage) {
		return output.String(), nil
	}

	flag.Parse()

	if err := utils.ValidateArgs(os.Args, utils.Validator{MaxArgs: 3, CmdName: "search", RequireName: false, HasFlags: false}); err != nil {
		output.WriteString(err.Error())
		return output.String(), err
	}

	tournamentsModel := TournamentsList()

	// Program 1: Tournament selection
	finalModel, err := tea.NewProgram(tournamentsModel, tea.WithAltScreen()).Run()
	if err != nil {
		return "", fmt.Errorf("error running tournamen selection program: %w", err)
	}

	result, ok := finalModel.(TournamentsModel)
	if !ok {
		return "", fmt.Errorf("unexpected model type from tournament selection: got %T, want TournamentsModel", finalModel)
	}

	if result.Choice != "" {
		// Program 2: Dashboard
		tabs := []string{"Overview", "Standings", "Decks", "Countries"}
		tabContent := []string{"Overview Tab", "Standings Tab", "Decks Tab", "Countries Tab"}
		dashboardModel := model{
			Tabs:       tabs,
			TabContent: tabContent,
			styles:     newStyles(),
		}

		_, err = tea.NewProgram(dashboardModel, tea.WithAltScreen()).Run()
		if err != nil {
			return "", fmt.Errorf("error running dashboard program: %w", err)
		}
	}

	return output.String(), nil
}
