package shell

import (
	"fmt"
	"image/color"
	"strconv"
	"strings"

	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"
)

type ConnFunc func(string) ([]byte, error)

type Decoded struct {
	TableRows []table.Row
	Countries []BarChartItem
	Overview  func(contentWidth int, highlight color.Color) string
	ExtraTab  func(contentWidth int) string
}

type Spec struct {
	Tabs         []string
	ListURL      string
	DashboardURL func(location string) string
	Columns      func(width int) []table.Column
	Decode       func(body []byte) (Decoded, error)
}

func Run(spec Spec, conn ConnFunc) (back bool, err error) {
	runPicker := func(m pickerModel) (pickerModel, error) {
		final, err := tea.NewProgram(m).Run()
		if err != nil {
			return pickerModel{}, err
		}
		result, ok := final.(pickerModel)
		if !ok {
			return pickerModel{}, fmt.Errorf("unexpected model type from tournament selection: got %T, want pickerModel", final)
		}
		return result, nil
	}

	runDashboard := func(m dashboardModel) (dashboardModel, error) {
		final, err := tea.NewProgram(m).Run()
		if err != nil {
			return dashboardModel{}, err
		}
		result, ok := final.(dashboardModel)
		if !ok {
			return dashboardModel{}, fmt.Errorf("unexpected model type from dashboard: got %T, want dashboardModel", final)
		}
		return result, nil
	}

	if err := loop(spec, conn, runPicker, runDashboard); err != nil {
		return false, err
	}
	return true, nil
}

func loop(
	spec Spec,
	conn ConnFunc,
	runPicker func(pickerModel) (pickerModel, error),
	runDashboard func(dashboardModel) (dashboardModel, error),
) error {
	for {
		result, err := runPicker(newPicker(spec, conn))
		if err != nil {
			return fmt.Errorf("error running tournament selection program: %w", err)
		}
		if result.selected == nil {
			break
		}

		dash, err := runDashboard(newDashboard(spec, conn, result.selected.Location))
		if err != nil {
			return fmt.Errorf("error running dashboard program: %w", err)
		}
		if !dash.goBack {
			break
		}
	}
	return nil
}

func FormatInt(n int) string {
	s := strconv.Itoa(n)
	var result strings.Builder
	for i, c := range s {
		if i > 0 && (len(s)-i)%3 == 0 {
			result.WriteRune(',')
		}
		result.WriteRune(c)
	}
	return result.String()
}
