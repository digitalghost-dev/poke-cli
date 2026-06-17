package champions

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	"github.com/digitalghost-dev/poke-cli/connections"
)

func Run() (back bool, err error) {
	final, err := tea.NewProgram(newDashboard(connections.CallTCGData)).Run()
	if err != nil {
		return false, fmt.Errorf("error running champions dashboard: %w", err)
	}

	result, ok := final.(dashboardModel)
	if !ok {
		return false, fmt.Errorf("unexpected model type from champions dashboard: got %T, want dashboardModel", final)
	}

	return result.goBack, nil
}
