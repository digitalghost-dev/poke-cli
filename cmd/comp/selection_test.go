package comp

import (
	"strings"
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/exp/teatest/v2"
)

func TestCompList_BuildsModel(t *testing.T) {
	m := CompList()
	if len(m.list.Items()) != 3 {
		t.Errorf("expected 3 items, got %d", len(m.list.Items()))
	}
	if m.list.Title != "Pick a competition type" {
		t.Errorf("unexpected title: %q", m.list.Title)
	}
	if m.compID != "" || m.choice != "" || m.quitting {
		t.Error("expected a clean initial model")
	}
}

func TestCompIDMap(t *testing.T) {
	tests := map[string]string{
		"TCG Competition Data":   "tcg",
		"VGC Competition Data":   "vgc",
		"Pokémon Champions Data": "champions",
	}
	for label, want := range tests {
		if got := compIDMap[label]; got != want {
			t.Errorf("compIDMap[%q] = %q, want %q", label, got, want)
		}
	}
}

func TestPicker_Init(t *testing.T) {
	if CompList().Init() != nil {
		t.Error("expected Init() to return nil")
	}
}

func TestPicker_Update_Quit(t *testing.T) {
	for _, key := range []tea.KeyPressMsg{
		{Code: tea.KeyEscape},
		{Code: 'c', Mod: tea.ModCtrl},
	} {
		newModel, cmd := CompList().Update(key)
		result := newModel.(pickerModel)
		if !result.quitting {
			t.Errorf("expected quitting=true after %v", key)
		}
		if cmd == nil {
			t.Error("expected a quit command")
		}
	}
}

func TestPicker_Update_Enter_SetsChoice(t *testing.T) {
	tm := teatest.NewTestModel(t, CompList(), teatest.WithInitialTermSize(80, 24))
	tm.Send(tea.KeyPressMsg{Code: tea.KeyEnter})
	tm.WaitFinished(t, teatest.WithFinalTimeout(300*time.Millisecond))
	final := tm.FinalModel(t).(pickerModel)
	if final.choice != "TCG Competition Data" {
		t.Errorf("expected first item chosen, got %q", final.choice)
	}
	if final.compID != "tcg" {
		t.Errorf("expected compID=tcg, got %q", final.compID)
	}
}

func TestPicker_Update_WindowResize(t *testing.T) {
	newModel, _ := CompList().Update(tea.WindowSizeMsg{Width: 100, Height: 40})
	if newModel.(pickerModel).list.Width() != 100 {
		t.Errorf("expected list width 100, got %d", newModel.(pickerModel).list.Width())
	}
}

func TestPicker_View_States(t *testing.T) {
	quit := CompList()
	quit.quitting = true
	if !strings.Contains(quit.View().Content, "Quitting") {
		t.Error("expected quitting message")
	}

	chosen := CompList()
	chosen.choice = "TCG Competition Data"
	if !strings.Contains(chosen.View().Content, "TCG Competition Data") {
		t.Error("expected chosen competition in view")
	}

	normal := CompList()
	v := normal.View()
	if v.Content == "" {
		t.Error("expected non-empty normal view")
	}
	if !v.AltScreen {
		t.Error("expected AltScreen enabled")
	}
	if !strings.Contains(v.Content, "TCG Competition Data") {
		t.Error("expected list items in normal view")
	}
}
