package search

import (
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/exp/teatest/v2"
)

func TestSelection(t *testing.T) {
	m := initialModel()
	testModel := teatest.NewTestModel(t, m, teatest.WithInitialTermSize(500, 600))

	testModel.Send(tea.KeyPressMsg{Code: tea.KeyDown})
	testModel.Send(tea.KeyPressMsg{Code: tea.KeyUp})
	testModel.Send(tea.KeyPressMsg{Code: tea.KeyEnter})

	testModel.Send(tea.KeyPressMsg{Code: 'c', Mod: tea.ModCtrl})
	testModel.WaitFinished(t, teatest.WithFinalTimeout(300*time.Millisecond))

	final := testModel.FinalModel(t).(model)

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
	m := initialModel()
	testModel := teatest.NewTestModel(t, m)

	// Move down twice, this should attempt to exceed max Choice
	testModel.Send(tea.KeyPressMsg{Code: tea.KeyDown}) // 0 → 1
	testModel.Send(tea.KeyPressMsg{Code: tea.KeyDown}) // 1 → 2, but should clamp to 1

	// Move up three times, this should attempt to go below 0
	testModel.Send(tea.KeyPressMsg{Code: tea.KeyUp}) // 1 → 0
	testModel.Send(tea.KeyPressMsg{Code: tea.KeyUp}) // 0 → -1, clamp to 0
	testModel.Send(tea.KeyPressMsg{Code: tea.KeyUp}) // stays at 0

	// Simulate enter and quit to finish
	testModel.Send(tea.KeyPressMsg{Code: tea.KeyEnter})
	testModel.Send(tea.KeyPressMsg{Code: 'c', Mod: tea.ModCtrl})
	testModel.WaitFinished(t)

	final := testModel.FinalModel(t).(model)

	if final.Choice != 0 && final.Choice != 1 {
		t.Errorf("Choice should be clamped between 0 and 1, got %d", final.Choice)
	}
}
