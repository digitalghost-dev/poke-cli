package shell

import (
	"errors"
	"strings"
	"testing"
	"time"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/exp/teatest/v2"
	"github.com/digitalghost-dev/poke-cli/styling"
)

func loadedPicker() pickerModel {
	tournaments := []TournamentRef{
		{Location: "London", TextDate: "January 10-12, 2025"},
		{Location: "Dallas", TextDate: "February 1-2, 2025"},
	}
	var items []list.Item
	for _, td := range tournaments {
		items = append(items, styling.Item(td.Location+" · "+td.TextDate))
	}
	l := list.New(items, styling.ItemDelegate{}, 40, 16)
	l.SetFilteringEnabled(false)
	return pickerModel{conn: noopConn, tournaments: tournaments, list: l, loading: false}
}

func TestPicker_InitialState(t *testing.T) {
	m := newPicker(testSpec(), noopConn)
	if !m.loading {
		t.Error("expected loading=true on init")
	}
	if m.selected != nil {
		t.Error("expected selected=nil on init")
	}
	if m.Init() == nil {
		t.Error("expected Init() to return a non-nil cmd")
	}
}

func TestFetchTournaments_ConnectionError(t *testing.T) {
	mock := func(_ string) ([]byte, error) { return nil, errors.New("refused") }
	msg := fetchTournaments("https://x.test", mock)()
	result := msg.(tournamentsDataMsg)
	if result.err == nil {
		t.Error("expected error")
	}
	if result.tournaments != nil {
		t.Error("expected nil tournaments on error")
	}
}

func TestFetchTournaments_InvalidJSON(t *testing.T) {
	mock := func(_ string) ([]byte, error) { return []byte("not json"), nil }
	if fetchTournaments("https://x.test", mock)().(tournamentsDataMsg).err == nil {
		t.Error("expected unmarshal error")
	}
}

func TestFetchTournaments_Success(t *testing.T) {
	var capturedURL string
	mock := func(url string) ([]byte, error) {
		capturedURL = url
		return []byte(`[{"location":"London","text_date":"Jan 10-12"}]`), nil
	}
	result := fetchTournaments("https://x.test/list", mock)().(tournamentsDataMsg)
	if result.err != nil {
		t.Fatalf("unexpected error: %v", result.err)
	}
	if len(result.tournaments) != 1 || result.tournaments[0].Location != "London" {
		t.Errorf("unexpected tournaments: %+v", result.tournaments)
	}
	if capturedURL != "https://x.test/list" {
		t.Errorf("expected the spec's ListURL to be fetched, got %q", capturedURL)
	}
}

func TestPicker_Update_Quit(t *testing.T) {
	tm := teatest.NewTestModel(t, loadedPicker(), teatest.WithInitialTermSize(80, 24))
	tm.Send(tea.KeyPressMsg{Code: tea.KeyEscape})
	tm.WaitFinished(t, teatest.WithFinalTimeout(300*time.Millisecond))
	if !tm.FinalModel(t).(pickerModel).quitting {
		t.Error("expected quitting=true after esc")
	}
}

func TestPicker_Update_Enter_SetsSelected(t *testing.T) {
	tm := teatest.NewTestModel(t, loadedPicker(), teatest.WithInitialTermSize(80, 24))
	tm.Send(tea.KeyPressMsg{Code: tea.KeyEnter})
	tm.WaitFinished(t, teatest.WithFinalTimeout(300*time.Millisecond))
	final := tm.FinalModel(t).(pickerModel)
	if final.selected == nil || final.selected.Location != "London" {
		t.Errorf("expected London selected, got %+v", final.selected)
	}
}

func TestPicker_Update_DataMsg_Success(t *testing.T) {
	m := newPicker(testSpec(), noopConn)
	newModel, _ := m.Update(tournamentsDataMsg{tournaments: []TournamentRef{{Location: "London"}}})
	result := newModel.(pickerModel)
	if result.loading {
		t.Error("expected loading=false after data")
	}
	if len(result.tournaments) != 1 {
		t.Errorf("expected 1 tournament, got %d", len(result.tournaments))
	}
}

func TestPicker_Update_DataMsg_Error(t *testing.T) {
	m := newPicker(testSpec(), noopConn)
	newModel, _ := m.Update(tournamentsDataMsg{err: errors.New("boom")})
	result := newModel.(pickerModel)
	if result.loading || result.error == nil {
		t.Error("expected loading=false and error set")
	}
}

func TestPicker_Update_WindowResize_WhenLoaded(t *testing.T) {
	newModel, _ := loadedPicker().Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	if newModel.(pickerModel).list.Width() != 120 {
		t.Errorf("expected list width 120, got %d", newModel.(pickerModel).list.Width())
	}
}

func TestPicker_View_States(t *testing.T) {
	loading := newPicker(testSpec(), noopConn)
	if !strings.Contains(loading.View().Content, "Loading tournaments") {
		t.Error("expected loading message")
	}

	errM := newPicker(testSpec(), noopConn)
	errM.loading = false
	errM.error = errors.New("something went wrong")
	if !strings.Contains(errM.View().Content, "something went wrong") {
		t.Error("expected error message in view")
	}

	quit := newPicker(testSpec(), noopConn)
	quit.quitting = true
	if !strings.Contains(quit.View().Content, "Quitting") {
		t.Error("expected quitting message")
	}

	sel := loadedPicker()
	td := sel.tournaments[0]
	sel.selected = &td
	if !strings.Contains(sel.View().Content, "London") {
		t.Error("expected selected tournament in view")
	}

	if loadedPicker().View().Content == "" {
		t.Error("expected non-empty normal view")
	}
}
