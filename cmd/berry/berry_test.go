package berry

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/exp/teatest"
	"github.com/digitalghost-dev/poke-cli/styling"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBerryCommand(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		wantErr  bool
		contains string
	}{
		{
			name:     "help flag short",
			args:     []string{"poke-cli", "berry", "-h"},
			wantErr:  false,
			contains: "USAGE:",
		},
		{
			name:     "help flag long",
			args:     []string{"poke-cli", "berry", "--help"},
			wantErr:  false,
			contains: "FLAGS:",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up os.Args for the test
			oldArgs := os.Args
			os.Args = tt.args
			defer func() { os.Args = oldArgs }()

			output, err := BerryCommand()

			if (err != nil) != tt.wantErr {
				t.Errorf("BerryCommand() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.contains != "" && !strings.Contains(output, tt.contains) {
				t.Errorf("BerryCommand() output should contain %q, got %q", tt.contains, output)
			}
		})
	}
}

func TestModelInit(t *testing.T) {
	m := model{}
	cmd := m.Init()
	if cmd != nil {
		t.Errorf("Init() should return nil, got %v", cmd)
	}
}

func TestModelUpdate(t *testing.T) {
	// Create a simple table for testing
	columns := []table.Column{{Title: "Berry", Width: 16}}
	rows := []table.Row{{"TestBerry"}}
	testTable := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(5),
	)

	m := model{
		table: testTable,
	}

	tests := []struct {
		name        string
		keyMsg      string
		shouldQuit  bool
		expectError bool
	}{
		{
			name:       "escape key",
			keyMsg:     "esc",
			shouldQuit: true,
		},
		{
			name:       "ctrl+c key",
			keyMsg:     "ctrl+c",
			shouldQuit: true,
		},
		{
			name:       "other key",
			keyMsg:     "j",
			shouldQuit: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updated, _ := m.Update(nil)

			if tt.shouldQuit {
				if updated == nil {
					t.Errorf("Update() returned nil model")
				}
			}
		})
	}
}

func TestModelView(t *testing.T) {
	// Test with empty table
	m := model{
		quitting: false,
		table:    table.New(),
	}

	view := m.View()
	if view == "" {
		t.Errorf("View() should not return empty string for normal state")
	}

	// Test quitting state
	m.quitting = true
	view = m.View()
	if !strings.Contains(view, "Goodbye") {
		t.Errorf("View() should contain 'Goodbye' when quitting, got %q", view)
	}
}

func TestModelViewWithSelectedBerry(t *testing.T) {
	// Create a table with test data
	columns := []table.Column{{Title: "Berry", Width: 16}}
	rows := []table.Row{{"Aguav"}}
	testTable := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(5),
	)

	m := model{
		table: testTable,
	}

	view := m.View()

	// Should contain the main UI elements
	expectedElements := []string{
		"Highlight a berry!",
		"↑ (move up) • ↓ (move down)",
		"ctrl+c | esc (quit)",
	}

	for _, element := range expectedElements {
		if !strings.Contains(view, element) {
			t.Errorf("View() should contain %q, got %q", element, view)
		}
	}
}

// createTestModel creates a model with test data for testing without database calls
func createTestModel() model {
	rows := []table.Row{
		{"Aguav"},
		{"Aspear"},
		{"Cheri"},
		{"Chesto"},
		{"Oran"},
	}

	t := table.New(
		table.WithColumns([]table.Column{{Title: "Berry", Width: 16}}),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(5),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(styling.YellowColor).
		BorderBottom(true)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("#000")).
		Background(styling.YellowColor)
	t.SetStyles(s)

	return model{table: t}
}

func TestTableNavigation(t *testing.T) {
	m := createTestModel()
	testModel := teatest.NewTestModel(t, m, teatest.WithInitialTermSize(100, 50))

	// Navigate down twice
	testModel.Send(tea.KeyMsg{Type: tea.KeyDown})
	testModel.Send(tea.KeyMsg{Type: tea.KeyDown})

	// Navigate back up once
	testModel.Send(tea.KeyMsg{Type: tea.KeyUp})

	// Quit the program
	testModel.Send(tea.KeyMsg{Type: tea.KeyCtrlC})
	testModel.WaitFinished(t, teatest.WithFinalTimeout(300*time.Millisecond))

	final := testModel.FinalModel(t).(model)

	if !final.quitting {
		t.Errorf("Expected model to be quitting after ctrl+c")
	}

	// After down, down, up from first row, we should be on second row (index 1 = "Aspear")
	if final.selectedOption != "Aspear" {
		t.Errorf("Expected selectedOption to be 'Aspear', got %q", final.selectedOption)
	}
}

func TestTableQuitWithEscape(t *testing.T) {
	m := createTestModel()
	testModel := teatest.NewTestModel(t, m, teatest.WithInitialTermSize(100, 50))

	// Quit with escape
	testModel.Send(tea.KeyMsg{Type: tea.KeyEsc})
	testModel.WaitFinished(t, teatest.WithFinalTimeout(300*time.Millisecond))

	final := testModel.FinalModel(t).(model)

	if !final.quitting {
		t.Errorf("Expected model to be quitting after escape")
	}
}

func TestTableInitialSelection(t *testing.T) {
	m := createTestModel()
	testModel := teatest.NewTestModel(t, m, teatest.WithInitialTermSize(100, 50))

	// Don't navigate, just quit immediately
	testModel.Send(tea.KeyMsg{Type: tea.KeyCtrlC})
	testModel.WaitFinished(t, teatest.WithFinalTimeout(300*time.Millisecond))

	final := testModel.FinalModel(t).(model)

	// First row should be selected by default
	if final.selectedOption != "Aguav" {
		t.Errorf("Expected initial selectedOption to be 'Aguav', got %q", final.selectedOption)
	}
}

func TestBerryCommandValidationError(t *testing.T) {
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()

	// Set os.Args with extra argument to trigger validation error
	os.Args = []string{"poke-cli", "berry", "cheri", "extra-arg"}

	output, err := BerryCommand()
	require.Error(t, err, "BerryCommand should return error for invalid args")
	assert.Contains(t, output, "Error", "Output should contain error message")
}
