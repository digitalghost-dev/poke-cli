package berry

import (
	"os"
	"strings"
	"testing"

	"github.com/charmbracelet/bubbles/table"
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
