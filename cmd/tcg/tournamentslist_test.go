package tcg

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/exp/teatest"
	"github.com/digitalghost-dev/poke-cli/styling"
)

func noopConn(_ string) ([]byte, error) { return []byte("[]"), nil }

// helpers

func loadedModel() tournamentsModel {
	tournaments := []tournamentData{
		{Location: "London", TextDate: "January 10-12, 2025"},
		{Location: "Dallas", TextDate: "February 1-2, 2025"},
	}
	var items []list.Item
	for _, td := range tournaments {
		items = append(items, styling.Item(td.Location+" · "+td.TextDate))
	}
	l := list.New(items, styling.ItemDelegate{}, 40, 16)
	l.SetFilteringEnabled(false)
	return tournamentsModel{
		conn:        noopConn,
		tournaments: tournaments,
		list:        l,
		loading:     false,
	}
}

// tournamentsList factory

func TestTournamentsList_InitialState(t *testing.T) {
	m := tournamentsList(noopConn)
	if !m.loading {
		t.Error("expected loading=true on init")
	}
	if m.selected != nil {
		t.Error("expected selected=nil on init")
	}
}

// Init

func TestTournamentsModel_Init_ReturnsCmd(t *testing.T) {
	m := tournamentsList(noopConn)
	cmd := m.Init()
	if cmd == nil {
		t.Error("expected Init() to return a non-nil cmd (spinner tick + fetch batch)")
	}
}

// fetchTournaments

func TestFetchTournaments_ConnectionError(t *testing.T) {
	mock := func(_ string) ([]byte, error) { return nil, errors.New("connection refused") }
	msg := fetchTournaments(mock)()
	result, ok := msg.(tournamentsDataMsg)
	if !ok {
		t.Fatalf("expected tournamentsDataMsg, got %T", msg)
	}
	if result.err == nil {
		t.Error("expected error, got nil")
	}
	if result.tournaments != nil {
		t.Error("expected nil tournaments on error")
	}
}

func TestFetchTournaments_InvalidJSON(t *testing.T) {
	mock := func(_ string) ([]byte, error) { return []byte("not json"), nil }
	msg := fetchTournaments(mock)()
	result, ok := msg.(tournamentsDataMsg)
	if !ok {
		t.Fatalf("expected tournamentsDataMsg, got %T", msg)
	}
	if result.err == nil {
		t.Error("expected unmarshal error, got nil")
	}
}

func TestFetchTournaments_Success(t *testing.T) {
	mock := func(_ string) ([]byte, error) {
		return []byte(`[{"location":"London","text_date":"January 10-12, 2025"},{"location":"Dallas","text_date":"February 1-2, 2025"}]`), nil
	}
	msg := fetchTournaments(mock)()
	result, ok := msg.(tournamentsDataMsg)
	if !ok {
		t.Fatalf("expected tournamentsDataMsg, got %T", msg)
	}
	if result.err != nil {
		t.Errorf("expected no error, got %v", result.err)
	}
	if len(result.tournaments) != 2 {
		t.Errorf("expected 2 tournaments, got %d", len(result.tournaments))
	}
	if result.tournaments[0].Location != "London" {
		t.Errorf("expected first location to be London, got %q", result.tournaments[0].Location)
	}
}

// Update — key messages

func TestTournamentsModel_Update_CtrlC(t *testing.T) {
	tests := []struct {
		name string
		key  tea.KeyType
	}{
		{name: "ctrl+c", key: tea.KeyCtrlC},
		{name: "esc", key: tea.KeyEsc},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := loadedModel()
			tm := teatest.NewTestModel(t, m, teatest.WithInitialTermSize(80, 24))
			tm.Send(tea.KeyMsg{Type: tt.key})
			tm.WaitFinished(t, teatest.WithFinalTimeout(300*time.Millisecond))
			final := tm.FinalModel(t).(tournamentsModel)
			if !final.quitting {
				t.Errorf("expected quitting=true after %s", tt.name)
			}
		})
	}
}

func TestTournamentsModel_Update_Enter_SetsSelected(t *testing.T) {
	m := loadedModel()
	tm := teatest.NewTestModel(t, m, teatest.WithInitialTermSize(80, 24))
	tm.Send(tea.KeyMsg{Type: tea.KeyEnter})
	tm.WaitFinished(t, teatest.WithFinalTimeout(300*time.Millisecond))
	final := tm.FinalModel(t).(tournamentsModel)
	if final.selected == nil {
		t.Fatal("expected selected to be set after enter")
	}
	if final.selected.Location != "London" {
		t.Errorf("expected Location=London, got %q", final.selected.Location)
	}
}

// Update — tournamentsDataMsg

func TestTournamentsModel_Update_DataMsg_Success(t *testing.T) {
	m := tournamentsList(noopConn)
	msg := tournamentsDataMsg{
		tournaments: []tournamentData{
			{Location: "London", TextDate: "January 10-12, 2025"},
			{Location: "Dallas", TextDate: "February 1-2, 2025"},
		},
	}
	newModel, _ := m.Update(msg)
	result := newModel.(tournamentsModel)
	if result.loading {
		t.Error("expected loading=false after data received")
	}
	if len(result.tournaments) != 2 {
		t.Errorf("expected 2 tournaments, got %d", len(result.tournaments))
	}
	if result.list.Items() == nil {
		t.Error("expected list to be populated")
	}
}

func TestTournamentsModel_Update_DataMsg_Error(t *testing.T) {
	m := tournamentsList(noopConn)
	msg := tournamentsDataMsg{err: errors.New("fetch failed")}
	newModel, _ := m.Update(msg)
	result := newModel.(tournamentsModel)
	if result.loading {
		t.Error("expected loading=false after error")
	}
	if result.error == nil {
		t.Error("expected error to be set")
	}
}

// Update — window resize

func TestTournamentsModel_Update_WindowResize_WhenLoaded(t *testing.T) {
	m := loadedModel()
	newModel, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	result := newModel.(tournamentsModel)
	if result.list.Width() != 120 {
		t.Errorf("expected list width=120, got %d", result.list.Width())
	}
}

func TestTournamentsModel_Update_WindowResize_WhenLoading(t *testing.T) {
	m := tournamentsList(noopConn) // loading=true
	newModel, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	result := newModel.(tournamentsModel)
	// list should not be updated while loading (it hasn't been created yet)
	if result.list.Width() == 120 {
		t.Error("expected list width to remain unchanged while loading")
	}
}

// View

func TestTournamentsModel_View_Loading(t *testing.T) {
	m := tournamentsList(noopConn)
	view := m.View()
	if !strings.Contains(view, "Loading tournaments") {
		t.Errorf("expected loading message, got %q", view)
	}
}

func TestTournamentsModel_View_Error(t *testing.T) {
	m := tournamentsList(noopConn)
	m.loading = false
	m.error = errors.New("something went wrong")
	view := m.View()
	if !strings.Contains(view, "something went wrong") {
		t.Errorf("expected error message in view, got %q", view)
	}
}

func TestTournamentsModel_View_Quitting(t *testing.T) {
	m := tournamentsList(noopConn)
	m.quitting = true
	view := m.View()
	if !strings.Contains(view, "Quitting") {
		t.Errorf("expected quitting message, got %q", view)
	}
}

func TestTournamentsModel_View_Selected(t *testing.T) {
	m := loadedModel()
	td := m.tournaments[0]
	m.selected = &td
	view := m.View()
	if !strings.Contains(view, "London") {
		t.Errorf("expected selected tournament in view, got %q", view)
	}
}

func TestTournamentsModel_View_Normal(t *testing.T) {
	m := loadedModel()
	view := m.View()
	if view == "" {
		t.Error("expected non-empty view for loaded model")
	}
}
