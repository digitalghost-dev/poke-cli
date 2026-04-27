package berry

import (
	"os"
	"strings"
	"testing"
	"time"

	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/exp/teatest/v2"
	"github.com/digitalghost-dev/poke-cli/cmd/utils"
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
		{
			name:     "invalid berry name",
			args:     []string{"poke-cli", "berry", "fakemon"},
			wantErr:  true,
			contains: "not found",
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
		table.WithWidth(16),
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
	if view.Content == "" {
		t.Errorf("View() should not return empty string for normal state")
	}

	// Test quitting state
	m.quitting = true
	view = m.View()
	if !strings.Contains(view.Content, "Goodbye") {
		t.Errorf("View() should contain 'Goodbye' when quitting, got %q", view.Content)
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
		table.WithWidth(16),
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
		if !strings.Contains(view.Content, element) {
			t.Errorf("View() should contain %q, got %q", element, view.Content)
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
		table.WithWidth(16),
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
	testModel.Send(tea.KeyPressMsg{Code: tea.KeyDown})
	testModel.Send(tea.KeyPressMsg{Code: tea.KeyDown})

	// Navigate back up once
	testModel.Send(tea.KeyPressMsg{Code: tea.KeyUp})

	// Quit the program
	testModel.Send(tea.KeyPressMsg{Code: 'c', Mod: tea.ModCtrl})
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
	testModel.Send(tea.KeyPressMsg{Code: tea.KeyEscape})
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
	testModel.Send(tea.KeyPressMsg{Code: 'c', Mod: tea.ModCtrl})
	testModel.WaitFinished(t, teatest.WithFinalTimeout(300*time.Millisecond))

	final := testModel.FinalModel(t).(model)

	// First row should be selected by default
	if final.selectedOption != "Aguav" {
		t.Errorf("Expected initial selectedOption to be 'Aguav', got %q", final.selectedOption)
	}
}

func TestBerryCommandOutput(t *testing.T) {
	err := os.Setenv("GO_TESTING", "1")
	if err != nil {
		t.Fatalf("Failed to set GO_TESTING env var: %v", err)
	}

	defer func() {
		err := os.Unsetenv("GO_TESTING")
		if err != nil {
			t.Logf("Warning: failed to unset GO_TESTING: %v", err)
		}
	}()

	tests := []struct {
		name           string
		args           []string
		expectedOutput string
	}{
		{
			name:           "Select 'Cheri' berry",
			args:           []string{"berry", "Cheri"},
			expectedOutput: utils.LoadGolden(t, "berry_cheri.golden"),
		},
		{
			name:           "Select 'Oran' berry",
			args:           []string{"berry", "Oran"},
			expectedOutput: utils.LoadGolden(t, "berry_oran.golden"),
		},
		{
			name:           "Select 'Sitrus' berry",
			args:           []string{"berry", "Sitrus"},
			expectedOutput: utils.LoadGolden(t, "berry_sitrus.golden"),
		},
		{
			name:           "Select 'Aguav' berry",
			args:           []string{"berry", "Aguav"},
			expectedOutput: utils.LoadGolden(t, "berry_aguav.golden"),
		},
		{
			name:           "Select 'Chople' berry",
			args:           []string{"berry", "Chople"},
			expectedOutput: utils.LoadGolden(t, "berry_chople.golden"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalArgs := os.Args
			os.Args = append([]string{"poke-cli"}, tt.args...)
			defer func() { os.Args = originalArgs }()

			output, _ := BerryCommand()
			cleanOutput := styling.StripANSI(output)

			assert.Equal(t, tt.expectedOutput, cleanOutput, "Output should match expected")
		})
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
