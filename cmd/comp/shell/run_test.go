package shell

import (
	"errors"
	"image/color"
	"testing"

	"charm.land/bubbles/v2/table"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func noopConn(_ string) ([]byte, error) { return []byte("[]"), nil }

func testSpec() Spec {
	return Spec{
		Tabs:         []string{"Overview", "Standings", "Extra", "Countries"},
		ListURL:      "https://example.test/list",
		DashboardURL: func(loc string) string { return "https://example.test/dash?location=" + loc },
		Columns: func(width int) []table.Column {
			return []table.Column{{Title: "Rank", Width: 4}, {Title: "Name", Width: 20}}
		},
		Decode: func(_ []byte) (Decoded, error) { return testDecoded(), nil },
	}
}

func testDecoded() Decoded {
	return Decoded{
		TableRows: []table.Row{{"1", "Ash"}, {"2", "Misty"}},
		Countries: []Tally{{Label: "USA", Count: 5}},
		Overview:  func(_ int, _ color.Color) string { return "OVERVIEW-BODY" },
		Extra: Frequency{
			NameHeader:  "Deck",
			CountHeader: "Players",
			Caption:     "EXTRA-CAPTION",
			Items:       []Tally{{Label: "EXTRA-DECK", Count: 2}},
		},
	}
}

func TestFormatInt(t *testing.T) {
	tests := []struct {
		n    int
		want string
	}{
		{0, "0"},
		{999, "999"},
		{1000, "1,000"},
		{4010, "4,010"},
		{10000, "10,000"},
		{1000000, "1,000,000"},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := FormatInt(tt.n); got != tt.want {
				t.Errorf("FormatInt(%d) = %q, want %q", tt.n, got, tt.want)
			}
		})
	}
}

func TestLoop_NoTournamentSelected(t *testing.T) {
	runPicker := func(_ pickerModel) (pickerModel, error) { return pickerModel{selected: nil}, nil }
	dashCalled := false
	runDashboard := func(_ dashboardModel) (dashboardModel, error) {
		dashCalled = true
		return dashboardModel{}, nil
	}
	err := loop(testSpec(), noopConn, runPicker, runDashboard)
	require.NoError(t, err)
	require.False(t, dashCalled, "dashboard should not launch when no tournament is selected")
}

func TestLoop_TournamentSelected_DashboardExits(t *testing.T) {
	td := TournamentRef{Location: "London"}
	runPicker := func(_ pickerModel) (pickerModel, error) { return pickerModel{selected: &td}, nil }
	runDashboard := func(m dashboardModel) (dashboardModel, error) {
		assert.Equal(t, "London", m.tournament)
		return dashboardModel{goBack: false}, nil
	}
	assert.NoError(t, loop(testSpec(), noopConn, runPicker, runDashboard))
}

func TestLoop_GoBack_LoopsToPicker(t *testing.T) {
	td := TournamentRef{Location: "London"}
	calls := 0
	runPicker := func(_ pickerModel) (pickerModel, error) {
		calls++
		if calls == 1 {
			return pickerModel{selected: &td}, nil
		}
		return pickerModel{selected: nil}, nil
	}
	runDashboard := func(_ dashboardModel) (dashboardModel, error) { return dashboardModel{goBack: true}, nil }
	require.NoError(t, loop(testSpec(), noopConn, runPicker, runDashboard))
	require.Equal(t, 2, calls, "expected the picker to run twice")
}

func TestLoop_PickerError(t *testing.T) {
	runPicker := func(_ pickerModel) (pickerModel, error) { return pickerModel{}, errors.New("boom") }
	err := loop(testSpec(), noopConn, runPicker, nil)
	assert.ErrorContains(t, err, "tournament selection")
}

func TestLoop_DashboardError(t *testing.T) {
	td := TournamentRef{Location: "London"}
	runPicker := func(_ pickerModel) (pickerModel, error) { return pickerModel{selected: &td}, nil }
	runDashboard := func(_ dashboardModel) (dashboardModel, error) {
		return dashboardModel{}, errors.New("boom")
	}
	assert.ErrorContains(t, loop(testSpec(), noopConn, runPicker, runDashboard), "dashboard")
}
