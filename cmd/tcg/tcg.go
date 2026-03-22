package tcg

import (
	"flag"
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/digitalghost-dev/poke-cli/cmd/utils"
	"github.com/digitalghost-dev/poke-cli/connections"
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

	if err := utils.ValidateArgs(os.Args, utils.Validator{MaxArgs: 3, CmdName: "tcg", RequireName: false, HasFlags: false}); err != nil {
		output.WriteString(err.Error())
		return output.String(), err
	}

	conn := connections.CallTCGData

	for {
		// Program 1: Tournament selection
		finalModel, err := tea.NewProgram(tournamentsList(conn), tea.WithAltScreen()).Run()
		if err != nil {
			return "", fmt.Errorf("error running tournament selection program: %w", err)
		}

		result, ok := finalModel.(tournamentsModel)
		if !ok {
			return "", fmt.Errorf("unexpected model type from tournament selection: got %T, want TournamentsModel", finalModel)
		}

		if result.selected == nil {
			break
		}

		// Program 2: Dashboard
		tabs := []string{"Overview", "Standings", "Decks", "Countries"}
		dashboardFinal, err := tea.NewProgram(model{
			conn:       conn,
			tabs:       tabs,
			styles:     newStyles(),
			tournament: result.selected.Location,
		}, tea.WithAltScreen()).Run()
		if err != nil {
			return "", fmt.Errorf("error running dashboard program: %w", err)
		}

		dashboard, ok := dashboardFinal.(model)
		if !ok {
			return "", fmt.Errorf("unexpected model type from dashboard: got %T, want model", dashboardFinal)
		}

		if !dashboard.goBack {
			break
		}
	}

	return output.String(), nil
}
