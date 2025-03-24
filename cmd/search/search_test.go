package search

import (
	tea "github.com/charmbracelet/bubbletea"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestModelInit(t *testing.T) {
	m := Model{}
	cmd := m.Init()
	assert.Nil(t, cmd, "Init() should return nil")
}

func TestModelQuit(t *testing.T) {
	m := Model{}

	// Simulate pressing Esc
	msg := tea.KeyMsg{Type: tea.KeyEsc}
	newModel, cmd := m.Update(msg)

	assert.True(t, newModel.(Model).Quitting, "Model should be set to quitting")

	if cmd != nil {
		assert.Equal(t, cmd(), tea.Quit(), "Update() should return tea.Quit command")
	} else {
		t.Errorf("Expected tea.Quit command, but got nil")
	}
}
