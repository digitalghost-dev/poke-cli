package card

import (
	"strings"
	"testing"

	spinnerpkg "github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

func TestImageModel_Init(t *testing.T) {
	model := ImageRenderer("001/198 - Pineco", "http://example.com/image.png")

	cmd := model.Init()
	if cmd == nil {
		t.Error("Init() should return a command (batch of spinner tick + fetch)")
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

func TestImageModel_View_Loading(t *testing.T) {
	model := ImageRenderer("001/198 - Pineco", "http://example.com/image.png")

	result := model.View()

	// When loading, should show spinner and card name
	if result == "" {
		t.Error("View() should not be empty when loading")
	}
	// Can't check exact spinner output as it's dynamic, but should contain card name
	if !strings.Contains(result, "001/198 - Pineco") {
		t.Error("View() should contain card name when loading")
	}
}

func TestImageModel_View_Loaded(t *testing.T) {
	expectedData := "test-sixel-data-123"
	model := ImageModel{
		CardName:  "001/198 - Pineco",
		ImageURL:  "http://example.com/image.png",
		Loading:   false,
		ImageData: expectedData,
	}

	result := model.View()

	if result != expectedData {
		t.Errorf("View() = %v, want %v", result, expectedData)
	}
}

func TestImageModel_View_Empty(t *testing.T) {
	model := ImageModel{
		CardName:  "001/198 - Pineco",
		ImageURL:  "",
		Loading:   false,
		ImageData: "",
	}

	result := model.View()

	if result != "" {
		t.Errorf("View() with empty ImageData should return empty string, got %v", result)
	}
}

func TestImageRenderer_InitializesCorrectly(t *testing.T) {
	testURL := "http://example.com/pikachu.png"
	model := ImageRenderer("Pikachu", testURL)

	if model.CardName != "Pikachu" {
		t.Errorf("ImageRenderer() CardName = %v, want %v", model.CardName, "Pikachu")
	}

	if model.ImageURL != testURL {
		t.Errorf("ImageRenderer() ImageURL = %v, want %v", model.ImageURL, testURL)
	}

	if !model.Loading {
		t.Error("ImageRenderer() should initialize with Loading = true")
	}

	if model.ImageData != "" {
		t.Error("ImageRenderer() should initialize with empty ImageData")
	}
}

func TestImageModel_Update_ImageReady(t *testing.T) {
	model := ImageRenderer("Charizard", "http://example.com/charizard.png")

	msg := imageReadyMsg{sixelData: "test-sixel-data-456"}
	newModel, cmd := model.Update(msg)

	if cmd != nil {
		t.Error("Update with imageReadyMsg should return nil command")
	}

	updatedModel := newModel.(ImageModel)
	if updatedModel.Loading {
		t.Error(`Update with imageReadyMsg should set Loading to false`)
	}

	if updatedModel.ImageData != "test-sixel-data-456" {
		t.Errorf("Update with imageReadyMsg should set ImageData, got %v", updatedModel.ImageData)
	}
}

func TestImageModel_Update_SpinnerTick(t *testing.T) {
	model := ImageRenderer("Mewtwo", "http://example.com/mewtwo.png")

	msg := model.Spinner.Tick()

	if _, ok := msg.(spinnerpkg.TickMsg); !ok {
		t.Fatalf("expected spinner.TickMsg, got %T", msg)
	}

	newModel, returnedCmd := model.Update(msg)

	if returnedCmd == nil {
		t.Error("Update with spinner.TickMsg should return a command")
	}

	// Model should still be ImageModel
	if _, ok := newModel.(ImageModel); !ok {
		t.Error("Update should return ImageModel")
	}
}
