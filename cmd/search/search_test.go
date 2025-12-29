package search

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/digitalghost-dev/poke-cli/styling"

	"github.com/stretchr/testify/assert"
)

func captureSearchOutput(f func()) string {
	// Create a pipe to capture standard output
	r, w, _ := os.Pipe()
	defer func(r *os.File) {
		err := r.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(r)

	// Redirect os.Stdout to the write end of the pipe
	oldStdout := os.Stdout
	os.Stdout = w
	defer func() { os.Stdout = oldStdout }()

	// Run the function
	f()

	// Close the write end of the pipe
	err := w.Close()
	if err != nil {
		return ""
	}

	// Read the captured output
	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	return buf.String()
}

func TestSearchCommand(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedOutput string
		expectedError  bool
	}{
		{
			name: "Help flag",
			args: []string{"search", "-h"},
			expectedOutput: "╭──────────────────────────────────────────────────────────────╮\n" +
				"│Search for a resource by name or partial match.               │\n" +
				"│                                                              │\n" +
				"│ USAGE:                                                       │\n" +
				"│    poke-cli search [flag]                                    │\n" +
				"│                                                              │\n" +
				"│ FLAGS:                                                       │\n" +
				"│    -h, --help      Prints out the help menu.                 │\n" +
				"│                                                              │\n" +
				"│ Supports prefix matching using ^ (example: ^char → charizard)│\n" +
				"╰──────────────────────────────────────────────────────────────╯\n",
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original os.Args
			originalArgs := os.Args
			defer func() { os.Args = originalArgs }()

			// Set os.Args for the test
			os.Args = append([]string{"poke-cli"}, tt.args...)

			// Capture the output
			output := captureSearchOutput(func() {
				defer func() {
					// Recover from os.Exit calls
					if r := recover(); r != nil {
						if !tt.expectedError {
							t.Fatalf("Unexpected error: %v", r)
						}
					}
				}()
				err := SearchCommand()
				if tt.expectedError {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
			})

			strippedOutput := styling.StripANSI(output)
			if !strings.Contains(strippedOutput, tt.expectedOutput) {
				t.Errorf("Output mismatch.\nExpected to contain:\n%s\nGot:\n%s", tt.expectedOutput, strippedOutput)
			}
		})
	}
}

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

func TestSearchCommandValidationError(t *testing.T) {
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()

	// Set os.Args with extra argument to trigger validation error
	os.Args = []string{"poke-cli", "search", "pokemon", "extra-arg"}

	err := SearchCommand()
	assert.Error(t, err, "SearchCommand should return error for invalid args")
}

func TestModelViewQuitting(t *testing.T) {
	m := Model{Quitting: true}
	view := m.View()
	assert.Contains(t, view, "Quitting search", "View should show quitting message")
}

func TestModelViewShowResults(t *testing.T) {
	m := Model{
		ShowResults:   true,
		SearchResults: "Test Results",
	}
	view := m.View()
	// View calls RenderInput when ShowResults is true
	assert.NotEmpty(t, view, "View should render results")
}

func TestModelViewNotChosen(t *testing.T) {
	m := Model{Chosen: false}
	view := m.View()
	// View calls RenderSelection when not chosen
	assert.Contains(t, view, "Search for a resource", "View should show selection prompt")
}
