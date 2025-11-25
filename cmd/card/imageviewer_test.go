package card

import (
	"image"
	"image/color"
	"image/png"
	"net/http"
	"net/http/httptest"
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

func TestImageRenderer_Success(t *testing.T) {
	// Create a test HTTP server that serves a valid PNG image
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		img := image.NewRGBA(image.Rect(0, 0, 10, 10))
		blue := color.RGBA{R: 0, G: 0, B: 255, A: 255}
		for y := 0; y < 10; y++ {
			for x := 0; x < 10; x++ {
				img.Set(x, y, blue)
			}
		}

		w.Header().Set("Content-Type", "image/png")
		w.WriteHeader(http.StatusOK)
		_ = png.Encode(w, img)
	}))
	defer server.Close()

	model := ImageRenderer("Pikachu", server.URL)

	if model.CardName != "Pikachu" {
		t.Errorf("ImageRenderer() CardName = %v, want %v", model.CardName, "Pikachu")
	}

	if model.Error != nil {
		t.Errorf("ImageRenderer() Error should be nil on success, got %v", model.Error)
	}

	if model.ImageURL == "" {
		t.Error("ImageRenderer() ImageURL should not be empty on success")
	}
}

func TestImageRenderer_Error(t *testing.T) {
	// Create a test HTTP server that returns an error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	model := ImageRenderer("Charizard", server.URL)

	if model.CardName != "Charizard" {
		t.Errorf("ImageRenderer() CardName = %v, want %v", model.CardName, "Charizard")
	}

	if model.Error == nil {
		t.Error("ImageRenderer() Error should not be nil when image fetch fails")
	}

	if model.ImageURL != "" {
		t.Errorf("ImageRenderer() ImageURL should be empty on error, got %v", model.ImageURL)
	}
}

func TestImageRenderer_InvalidImage(t *testing.T) {
	// Create a test HTTP server that returns invalid image data
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/png")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("not a valid image"))
	}))
	defer server.Close()

	model := ImageRenderer("Mewtwo", server.URL)

	if model.CardName != "Mewtwo" {
		t.Errorf("ImageRenderer() CardName = %v, want %v", model.CardName, "Mewtwo")
	}

	if model.Error == nil {
		t.Error("ImageRenderer() Error should not be nil when image decoding fails")
	}

	if model.ImageURL != "" {
		t.Errorf("ImageRenderer() ImageURL should be empty on error, got %v", model.ImageURL)
	}
}
