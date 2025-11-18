package card

import (
	"image"
	"image/color"
	"image/png"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCardName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple card name",
			input:    "Pikachu",
			expected: "Pikachu",
		},
		{
			name:     "card with number",
			input:    "001/198 - Pineco",
			expected: "001/198 - Pineco",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CardName(tt.input)
			if result != tt.expected {
				t.Errorf("CardName() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestResizeImage(t *testing.T) {
	// Create a simple test image (100x100 red square)
	testImg := image.NewRGBA(image.Rect(0, 0, 100, 100))
	red := color.RGBA{R: 255, G: 0, B: 0, A: 255}
	for y := 0; y < 100; y++ {
		for x := 0; x < 100; x++ {
			testImg.Set(x, y, red)
		}
	}

	tests := []struct {
		name       string
		img        image.Image
		width      int
		height     int
		wantWidth  int
		wantHeight int
	}{
		{
			name:       "resize to smaller dimensions",
			img:        testImg,
			width:      50,
			height:     50,
			wantWidth:  50,
			wantHeight: 50,
		},
		{
			name:       "resize to larger dimensions",
			img:        testImg,
			width:      200,
			height:     200,
			wantWidth:  200,
			wantHeight: 200,
		},
		{
			name:       "resize to card dimensions",
			img:        testImg,
			width:      500,
			height:     675,
			wantWidth:  500,
			wantHeight: 675,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := resizeImage(tt.img, tt.width, tt.height)
			bounds := result.Bounds()

			if bounds.Dx() != tt.wantWidth {
				t.Errorf("resizeImage() width = %v, want %v", bounds.Dx(), tt.wantWidth)
			}
			if bounds.Dy() != tt.wantHeight {
				t.Errorf("resizeImage() height = %v, want %v", bounds.Dy(), tt.wantHeight)
			}
		})
	}
}

func TestCardImage_Success(t *testing.T) {
	// Create a test HTTP server that serves a small PNG image
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a minimal 10x10 PNG image
		img := image.NewRGBA(image.Rect(0, 0, 10, 10))
		blue := color.RGBA{R: 0, G: 0, B: 255, A: 255}
		for y := 0; y < 10; y++ {
			for x := 0; x < 10; x++ {
				img.Set(x, y, blue)
			}
		}

		w.Header().Set("Content-Type", "image/png")
		w.WriteHeader(http.StatusOK)
		png.Encode(w, img)
	}))
	defer server.Close()

	result, err := CardImage(server.URL)

	if err != nil {
		t.Errorf("CardImage() error = %v, want nil", err)
		return
	}

	// Check that result is a valid Sixel string
	if !strings.HasPrefix(result, "\x1bPq") {
		t.Error("CardImage() should return string starting with Sixel header")
	}

	if !strings.HasSuffix(result, "\x1b\\") {
		t.Error("CardImage() should return string ending with Sixel terminator")
	}

	if len(result) == 0 {
		t.Error("CardImage() should return non-empty string")
	}
}

func TestCardImage_EncodingError(t *testing.T) {
	// Create a test HTTP server that serves invalid image data
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/png")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("not a valid PNG"))
	}))
	defer server.Close()

	result, err := CardImage(server.URL)

	if err == nil {
		t.Error("CardImage() should return error for invalid image data")
	}

	if result != "" {
		t.Errorf("CardImage() on error should return empty string, got %v", result)
	}

	if !strings.Contains(err.Error(), "failed to decode image") {
		t.Errorf("Error message should mention 'failed to decode image', got: %v", err)
	}
}

func TestCardImage_Non200Response(t *testing.T) {
	// Create a test HTTP server that returns 404
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	result, err := CardImage(server.URL)

	if err == nil {
		t.Error("CardImage() should return error for non-200 response")
	}

	if result != "" {
		t.Errorf("CardImage() on error should return empty string, got %v", result)
	}

	if !strings.Contains(err.Error(), "non-200 response") {
		t.Errorf("Error message should mention 'non-200 response', got: %v", err)
	}
}
