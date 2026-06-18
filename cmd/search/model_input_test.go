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
		showResults: true,
		textInput:   ti,
	}

	msg := tea.KeyPressMsg{Code: 'b', Text: "b"}
	mUpdated, _ := UpdateInput(msg, m)

	updated := mUpdated.(model)

	if updated.showResults {
		t.Errorf("expected showResults to be false after pressing 'b'")
	}

	if updated.textInput.Value() != "" {
		t.Errorf("expected textInput to be reset")
	}
}
