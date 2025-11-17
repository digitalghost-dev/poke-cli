package card

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestImageModel_Init(t *testing.T) {
	model := ImageModel{
		CardName: "001/198 - Pineco",
		ImageURL: "test-sixel-data",
	}

	cmd := model.Init()
	if cmd != nil {
		t.Error("Init() should return nil")
	}
}

func TestImageModel_Update_EscKey(t *testing.T) {
	model := ImageModel{
		CardName: "001/198 - Pineco",
		ImageURL: "test-sixel-data",
	}

	// Test ESC key
	msg := tea.KeyMsg{Type: tea.KeyEsc}
	newModel, cmd := model.Update(msg)

	// Should return quit command
	if cmd == nil {
		t.Error("Update with ESC should return tea.Quit command")
	}

	// Model should be returned (even if quitting)
	if _, ok := newModel.(ImageModel); !ok {
		t.Error("Update should return ImageModel")
	}
}

func TestImageModel_Update_CtrlC(t *testing.T) {
	model := ImageModel{
		CardName: "001/198 - Pineco",
		ImageURL: "test-sixel-data",
	}

	msg := tea.KeyMsg{Type: tea.KeyCtrlC}
	_, cmd := model.Update(msg)

	// Should return quit command
	if cmd == nil {
		t.Error("Update with Ctrl+C should return tea.Quit command")
	}
}

func TestImageModel_Update_DifferentKey(t *testing.T) {
	model := ImageModel{
		CardName: "001/198 - Pineco",
		ImageURL: "test-sixel-data",
	}

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}}
	_, cmd := model.Update(msg)

	if cmd != nil {
		t.Error("Update with non-quit key should not return a command")
	}
}

func TestImageModel_View(t *testing.T) {
	expectedURL := "test-sixel-data-123"
	model := ImageModel{
		CardName: "001/198 - Pineco",
		ImageURL: expectedURL,
	}

	result := model.View()

	if result != expectedURL {
		t.Errorf("View() = %v, want %v", result, expectedURL)
	}
}

func TestImageModel_View_Empty(t *testing.T) {
	model := ImageModel{
		CardName: "001/198 - Pineco",
		ImageURL: "",
	}

	result := model.View()

	if result != "" {
		t.Errorf("View() with empty ImageURL should return empty string, got %v", result)
	}
}
