package search

import (
	"os"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/digitalghost-dev/poke-cli/styling"
	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

func TestSearchCommand(t *testing.T) {
	tests := []struct {
		name          string
		args          []string
		contains      string
		expectedError bool
	}{
		{
			name:          "Help flag short",
			args:          []string{"search", "-h"},
			contains:      "USAGE:",
			expectedError: false,
		},
		{
			name:          "Help flag long",
			args:          []string{"search", "--help"},
			contains:      "USAGE:",
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalArgs := os.Args
			defer func() { os.Args = originalArgs }()

			os.Args = append([]string{"poke-cli"}, tt.args...)

			output, err := SearchCommand()
			cleanOutput := styling.StripANSI(output)

			if tt.expectedError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			if !strings.Contains(cleanOutput, tt.contains) {
				t.Errorf("expected output to contain %q, got:\n%s", tt.contains, cleanOutput)
			}
		})
	}
}

func TestModelInit(t *testing.T) {
	m := model{}
	cmd := m.Init()
	assert.Nil(t, cmd, "Init() should return nil")
}

func TestModelQuit(t *testing.T) {
	m := model{}

	// Simulate pressing Esc
	msg := tea.KeyMsg{Type: tea.KeyEsc}
	newModel, cmd := m.Update(msg)

	assert.True(t, newModel.(model).Quitting, "Model should be set to quitting")

	if cmd != nil {
		assert.Equal(t, cmd(), tea.Quit(), "Update() should return tea.Quit command")
	} else {
		t.Errorf("Expected tea.Quit command, but got nil")
	}
}

func TestSearchCommandValidationError(t *testing.T) {
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()

	// Set os.Args with extra argument to trigger validation error
	os.Args = []string{"poke-cli", "search", "pokemon", "extra-arg"}

	_, err := SearchCommand()
	assert.Error(t, err, "SearchCommand should return error for invalid args")
}

func TestModelViewQuitting(t *testing.T) {
	m := model{Quitting: true}
	view := m.View()
	assert.Contains(t, view, "Quitting search", "View should show quitting message")
}

func TestModelViewShowResults(t *testing.T) {
	m := model{
		ShowResults:   true,
		SearchResults: "Test Results",
	}
	view := m.View()
	// View calls RenderInput when ShowResults is true
	assert.NotEmpty(t, view, "View should render results")
}

func TestModelViewNotChosen(t *testing.T) {
	m := model{Chosen: false}
	view := m.View()
	// View calls RenderSelection when not chosen
	assert.Contains(t, view, "Search for a resource", "View should show selection prompt")
}
