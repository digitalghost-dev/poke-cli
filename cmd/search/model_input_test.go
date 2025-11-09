package search

import (
	"testing"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func TestUpdateInput(t *testing.T) {
	ti := textinput.New()
	ti.SetValue("mewtwo")

	m := Model{
		ShowResults: true,
		TextInput:   ti,
	}

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'b'}}
	mUpdated, _ := UpdateInput(msg, m)

	updated := mUpdated.(Model)

	if updated.ShowResults {
		t.Errorf("expected ShowResults to be false after pressing 'b'")
	}

	if updated.TextInput.Value() != "" {
		t.Errorf("expected TextInput to be reset")
	}
}
