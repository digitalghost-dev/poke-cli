package search

import (
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/exp/teatest"
)

func TestSelection(t *testing.T) {
	model := initialModel()
	testModel := teatest.NewTestModel(t, model, teatest.WithInitialTermSize(500, 600))

	testModel.Send(tea.KeyMsg{Type: tea.KeyDown})
	testModel.Send(tea.KeyMsg{Type: tea.KeyUp})
	testModel.Send(tea.KeyMsg{Type: tea.KeyEnter})

	testModel.Send(tea.KeyMsg{Type: tea.KeyCtrlC})
	testModel.WaitFinished(t, teatest.WithFinalTimeout(300*time.Millisecond)) // optional timeout safety

	final := testModel.FinalModel(t).(Model)

	if !final.Chosen {
		t.Errorf("Expected model to be in Chosen state after pressing enter")
	}
	if final.Choice != 0 {
		t.Errorf("Expected Choice to be 0, got %d", final.Choice)
	}
	if !final.TextInput.Focused() {
		t.Errorf("Expected TextInput to be focused after selection")
	}
	if !final.Quitting {
		t.Errorf("Expected model to be quitting after ctrl+c")
	}
}

func TestChoiceClamping(t *testing.T) {
	model := initialModel()
	testModel := teatest.NewTestModel(t, model)

	// Move down twice: this should attempt to exceed max Choice
	testModel.Send(tea.KeyMsg{Type: tea.KeyDown}) // 0 → 1
	testModel.Send(tea.KeyMsg{Type: tea.KeyDown}) // 1 → 2, but should clamp to 1

	// Move up three times: this should attempt to go below 0
	testModel.Send(tea.KeyMsg{Type: tea.KeyUp}) // 1 → 0
	testModel.Send(tea.KeyMsg{Type: tea.KeyUp}) // 0 → -1, clamp to 0
	testModel.Send(tea.KeyMsg{Type: tea.KeyUp}) // stays at 0

	// Now simulate enter and quit to finish
	testModel.Send(tea.KeyMsg{Type: tea.KeyEnter})
	testModel.Send(tea.KeyMsg{Type: tea.KeyCtrlC})
	testModel.WaitFinished(t)

	final := testModel.FinalModel(t).(Model)

	if final.Choice != 0 && final.Choice != 1 {
		t.Errorf("Choice should be clamped between 0 and 1, got %d", final.Choice)
	}
}
