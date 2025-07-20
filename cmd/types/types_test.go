package types

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/exp/teatest"
	"github.com/digitalghost-dev/poke-cli/cmd/utils"
	"github.com/digitalghost-dev/poke-cli/styling"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func TestTypesCommand(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedOutput string
		wantError      bool
	}{
		{
			name:           "Types help flag",
			args:           []string{"types", "--help"},
			expectedOutput: utils.LoadGolden(t, "types_help.golden"),
		},
		{
			name:           "Types help flag",
			args:           []string{"types", "-h"},
			expectedOutput: utils.LoadGolden(t, "types_help.golden"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalArgs := os.Args
			os.Args = append([]string{"poke-cli"}, tt.args...)
			defer func() { os.Args = originalArgs }()

			output, _ := TypesCommand()
			cleanOutput := styling.StripANSI(output)

			assert.Equal(t, tt.expectedOutput, cleanOutput, "Output should match expected")
		})
	}
}

func TestModelInit(t *testing.T) {
	m := model{}
	cmd := m.Init()
	assert.Nil(t, cmd, "Init() should return nil")
}

// createTestModel creates a model with a table for testing
func createTestModel() model {
	// Create a simple table with a few types
	types := []string{"Normal", "Fire", "Water"}
	rows := make([]table.Row, len(types))
	for i, t := range types {
		rows[i] = []string{t}
	}

	t := table.New(
		table.WithColumns([]table.Column{{Title: "Type", Width: 16}}),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(7),
	)

	// Set table styles
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("#FFCC00")).
		BorderBottom(true)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("#000")).
		Background(lipgloss.Color("#FFCC00"))
	t.SetStyles(s)

	return model{table: t}
}

func TestUpdate(t *testing.T) {
	t.Run("Escape key should set quitting to true", func(t *testing.T) {
		m := createTestModel()
		testModel := teatest.NewTestModel(t, m)

		// Send escape key
		testModel.Send(tea.KeyMsg{Type: tea.KeyEsc})
		testModel.WaitFinished(t, teatest.WithFinalTimeout(300*time.Millisecond))

		final := testModel.FinalModel(t).(model)
		assert.True(t, final.quitting, "quitting should be true after pressing escape")
	})

	t.Run("Ctrl+C key should set quitting to true", func(t *testing.T) {
		m := createTestModel()
		testModel := teatest.NewTestModel(t, m)

		// Send ctrl+c key
		testModel.Send(tea.KeyMsg{Type: tea.KeyCtrlC})
		testModel.WaitFinished(t, teatest.WithFinalTimeout(300*time.Millisecond))

		final := testModel.FinalModel(t).(model)
		assert.True(t, final.quitting, "quitting should be true after pressing ctrl+c")
	})

	t.Run("Enter key should set selectedOption", func(t *testing.T) {
		m := createTestModel()
		testModel := teatest.NewTestModel(t, m)

		// Send enter key
		testModel.Send(tea.KeyMsg{Type: tea.KeyEnter})
		testModel.WaitFinished(t, teatest.WithFinalTimeout(300*time.Millisecond))

		final := testModel.FinalModel(t).(model)
		assert.Equal(t, "Normal", final.selectedOption, "selectedOption should be set to the selected row")
	})

	t.Run("Arrow keys should update table selection", func(t *testing.T) {
		m := createTestModel()
		testModel := teatest.NewTestModel(t, m)

		// Send down arrow key to select the second row
		testModel.Send(tea.KeyMsg{Type: tea.KeyDown})

		// Then send enter to select it
		testModel.Send(tea.KeyMsg{Type: tea.KeyEnter})
		testModel.WaitFinished(t, teatest.WithFinalTimeout(300*time.Millisecond))

		final := testModel.FinalModel(t).(model)
		assert.Equal(t, "Fire", final.selectedOption, "selectedOption should be updated after arrow navigation")
	})
}

func TestView(t *testing.T) {
	t.Run("View should return goodbye message when quitting", func(t *testing.T) {
		m := createTestModel()
		m.quitting = true

		view := m.View()
		assert.Equal(t, "\n  Goodbye! \n", view, "View should return goodbye message when quitting")
	})

	t.Run("View should return empty string when selectedOption is set", func(t *testing.T) {
		m := createTestModel()
		m.selectedOption = "Fire"
	})

	t.Run("View should render table in normal state", func(t *testing.T) {
		m := createTestModel()

		view := m.View()
		assert.Contains(t, view, "Select a type!", "View should contain the title")
		assert.Contains(t, view, "Type", "View should contain the table header")
		assert.Contains(t, view, "move up", "View should contain the key menu")
	})
}
