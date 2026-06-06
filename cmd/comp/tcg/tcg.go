package tcg

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	"github.com/digitalghost-dev/poke-cli/connections"
)

func Run() (back bool, err error) {
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
		return false, err
	}

	return true, nil
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
