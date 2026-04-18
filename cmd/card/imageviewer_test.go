package card

import (
	"strings"
	"testing"

	spinnerpkg "charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
)

func TestImageModel_Init(t *testing.T) {
	model := ImageRenderer("001/198 - Pineco", "http://example.com/image.png")

	cmd := model.Init()
	if cmd == nil {
		t.Error("Init() should return a command (batch of spinner tick + fetch)")
	}
}

func TestImageModel_Update_EscKey(t *testing.T) {
	model := imageModel{
		CardName: "001/198 - Pineco",
		ImageURL: "test-sixel-data",
	}

	// Test ESC key
	msg := tea.KeyPressMsg{Code: tea.KeyEscape}
	newModel, cmd := model.Update(msg)

	// Should return quit command
	if cmd == nil {
		t.Error("Update with ESC should return tea.Quit command")
	}

	// Model should be returned (even if quitting)
	if _, ok := newModel.(imageModel); !ok {
		t.Error("Update should return ImageModel")
	}
}

func TestImageModel_Update_CtrlC(t *testing.T) {
	model := imageModel{
		CardName: "001/198 - Pineco",
		ImageURL: "test-sixel-data",
	}

	msg := tea.KeyPressMsg{Code: 'c', Mod: tea.ModCtrl}
	_, cmd := model.Update(msg)

	if cmd == nil {
		t.Error("Update with Ctrl+C should return tea.Quit command")
	}
}

func TestImageModel_Update_DifferentKey(t *testing.T) {
	model := imageModel{
		CardName: "001/198 - Pineco",
		ImageURL: "test-sixel-data",
	}

	msg := tea.KeyPressMsg{Code: 'a', Text: "a"}
	_, cmd := model.Update(msg)

	if cmd != nil {
		t.Error("Update with non-quit key should not return a command")
	}
}

func TestImageModel_View_Loading(t *testing.T) {
	model := ImageRenderer("001/198 - Pineco", "http://example.com/image.png")

	result := model.View()

	// When loading, should show spinner and card name
	if result.Content == "" {
		t.Error("View() should not be empty when loading")
	}
	// Can't check exact spinner output as it's dynamic, but should contain card name
	if !strings.Contains(result.Content, "001/198 - Pineco") {
		t.Error("View() should contain card name when loading")
	}
}

func TestImageModel_View_Loaded(t *testing.T) {
	expectedData := "test-sixel-data-123"
	model := imageModel{
		CardName:  "001/198 - Pineco",
		ImageURL:  "http://example.com/image.png",
		Loading:   false,
		ImageData: expectedData,
	}

	result := model.View()

	if result.Content != expectedData {
		t.Errorf("View() = %v, want %v", result.Content, expectedData)
	}
}

func TestImageModel_View_Empty(t *testing.T) {
	model := imageModel{
		CardName:  "001/198 - Pineco",
		ImageURL:  "",
		Loading:   false,
		ImageData: "",
	}

	result := model.View()

	if result.Content != "" {
		t.Errorf("View() with empty ImageData should return empty string, got %v", result.Content)
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

	msg := imageReadyMsg{imageData: "test-image-data-456", protocol: "Kitty"}
	newModel, cmd := model.Update(msg)

	if cmd != nil {
		t.Error("Update with imageReadyMsg should return nil command")
	}

	updatedModel := newModel.(imageModel)
	if updatedModel.Loading {
		t.Error(`Update with imageReadyMsg should set Loading to false`)
	}

	if updatedModel.ImageData != "test-image-data-456" {
		t.Errorf("Update with imageReadyMsg should set ImageData, got %v", updatedModel.ImageData)
	}

	if updatedModel.Protocol != "Kitty" {
		t.Errorf("Update with imageReadyMsg should set Protocol, got %v", updatedModel.Protocol)
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
	if _, ok := newModel.(imageModel); !ok {
		t.Error("Update should return ImageModel")
	}
}
