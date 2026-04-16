package tcg

import (
	"flag"
	"fmt"
	"os"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/digitalghost-dev/poke-cli/cmd/utils"
	"github.com/digitalghost-dev/poke-cli/connections"
	"github.com/digitalghost-dev/poke-cli/flags"
)

func TcgCommand() (string, error) {
	var output strings.Builder

	flag.Usage = func() {
		output.WriteString(
			utils.GenerateHelpMessage(
				utils.HelpConfig{
					Description: "Get details about TCG tournaments.",
					CmdName:     "tcg",
					Flags: []utils.FlagHelp{
						{Short: "-w", Long: "--web", Description: "Opens the Streamlit dashboard in your default browser."},
					},
				},
			),
		)
	}

	if utils.CheckHelpFlag(&output, flag.Usage) {
		return output.String(), nil
	}

	flag.Parse()

	if err := utils.ValidateArgs(os.Args, utils.Validator{MaxArgs: 3, CmdName: "tcg", RequireName: false, HasFlags: true}); err != nil {
		output.WriteString(err.Error())
		return output.String(), err
	}

	tf := flags.SetupTcgFlagSet()
	if err := tf.FlagSet.Parse(os.Args[2:]); err != nil {
		fmt.Fprintf(&output, "error parsing flags: %v\n", err)
		return output.String(), err
	}

	if *tf.Web || *tf.ShortWeb {
		msg, err := flags.WebFlag("https://web.poke-cli.com/")
		if err != nil {
			return "", err
		}
		output.WriteString(msg)
		return output.String(), nil
	}

	conn := connections.CallTCGData

	runTournaments := func(m tournamentsModel) (tournamentsModel, error) {
		final, err := tea.NewProgram(m).Run()
		if err != nil {
			return tournamentsModel{}, err
		}
		result, ok := final.(tournamentsModel)
		if !ok {
			return tournamentsModel{}, fmt.Errorf("unexpected model type from tournament selection: got %T, want tournamentsModel", final)
		}
		return result, nil
	}

	runDashboard := func(m model) (model, error) {
		final, err := tea.NewProgram(m).Run()
		if err != nil {
			return model{}, err
		}
		result, ok := final.(model)
		if !ok {
			return model{}, fmt.Errorf("unexpected model type from dashboard: got %T, want model", final)
		}
		return result, nil
	}

	if err := runTcgLoop(conn, runTournaments, runDashboard); err != nil {
		return "", err
	}

	return output.String(), nil
}

func runTcgLoop(
	conn func(string) ([]byte, error),
	runTournaments func(tournamentsModel) (tournamentsModel, error),
	runDashboard func(model) (model, error),
) error {
	for {
		result, err := runTournaments(tournamentsList(conn))
		if err != nil {
			return fmt.Errorf("error running tournament selection program: %w", err)
		}
		if result.selected == nil {
			break
		}

		tabs := []string{"Overview", "Standings", "Decks", "Countries"}
		dashboard, err := runDashboard(model{
			conn:       conn,
			tabs:       tabs,
			styles:     newStyles(),
			tournament: result.selected.Location,
		})
		if err != nil {
			return fmt.Errorf("error running dashboard program: %w", err)
		}
		if !dashboard.goBack {
			break
		}
	}
	return nil
}
