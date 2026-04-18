package search

import (
	"testing"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
)

func TestUpdateInput(t *testing.T) {
	ti := textinput.New()
	ti.SetValue("mewtwo")

	m := model{
		ShowResults: true,
		TextInput:   ti,
	}

	msg := tea.KeyPressMsg{Code: 'b', Text: "b"}
	mUpdated, _ := UpdateInput(msg, m)

	updated := mUpdated.(model)

	if updated.ShowResults {
		t.Errorf("expected ShowResults to be false after pressing 'b'")
	}

	if updated.TextInput.Value() != "" {
		t.Errorf("expected TextInput to be reset")
	}
}
