package tcg

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// runTcgLoop

func TestRunTcgLoop_NoTournamentSelected(t *testing.T) {
	// User quits tournament selection without picking → loop exits immediately.
	runTournaments := func(m tournamentsModel) (tournamentsModel, error) {
		return tournamentsModel{selected: nil}, nil
	}
	dashboardCalled := false
	runDashboard := func(m model) (model, error) {
		dashboardCalled = true
		return model{}, nil
	}
	err := runTcgLoop(noopConn, runTournaments, runDashboard)
	require.NoError(t, err)
	require.False(t, dashboardCalled, "dashboard should not be launched when no tournament is selected")
}

func TestRunTcgLoop_TournamentSelected_DashboardExits(t *testing.T) {
	// User picks a tournament, views the dashboard, then quits (goBack=false).
	td := tournamentData{Location: "London"}
	runTournaments := func(m tournamentsModel) (tournamentsModel, error) {
		return tournamentsModel{selected: &td}, nil
	}
	runDashboard := func(m model) (model, error) {
		assert.Equal(t, "London", m.tournament)
		return model{goBack: false}, nil
	}
	err := runTcgLoop(noopConn, runTournaments, runDashboard)
	assert.NoError(t, err)
}

func TestRunTcgLoop_GoBack_LoopsToTournamentSelection(t *testing.T) {
	// User presses b in the dashboard → goes back to the tournament list.
	// On the second visit, they quit without selecting.
	td := tournamentData{Location: "London"}
	calls := 0
	runTournaments := func(m tournamentsModel) (tournamentsModel, error) {
		calls++
		if calls == 1 {
			return tournamentsModel{selected: &td}, nil
		}
		return tournamentsModel{selected: nil}, nil
	}
	runDashboard := func(m model) (model, error) {
		return model{goBack: true}, nil
	}
	err := runTcgLoop(noopConn, runTournaments, runDashboard)
	require.NoError(t, err)
	require.Equal(t, 2, calls, "expected tournament selection to run twice")
}

func TestRunTcgLoop_TournamentRunnerError(t *testing.T) {
	runTournaments := func(m tournamentsModel) (tournamentsModel, error) {
		return tournamentsModel{}, errors.New("program crashed")
	}
	err := runTcgLoop(noopConn, runTournaments, nil)
	assert.ErrorContains(t, err, "tournament selection")
}

func TestRunTcgLoop_DashboardRunnerError(t *testing.T) {
	td := tournamentData{Location: "London"}
	runTournaments := func(m tournamentsModel) (tournamentsModel, error) {
		return tournamentsModel{selected: &td}, nil
	}
	runDashboard := func(m model) (model, error) {
		return model{}, errors.New("dashboard crashed")
	}
	err := runTcgLoop(noopConn, runTournaments, runDashboard)
	assert.ErrorContains(t, err, "dashboard")
}
