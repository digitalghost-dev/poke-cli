package card

import (
	"image"
	"image/color"
	"image/png"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

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
		err := png.Encode(w, img)
		if err != nil {
			return
		}
	}))
	defer server.Close()

	// Set up a supported terminal environment (Sixel)
	os.Setenv("TERM_PROGRAM", "iTerm.app")
	defer os.Unsetenv("TERM_PROGRAM")

	result, protocol, err := CardImage(server.URL)

	if err != nil {
		t.Errorf("CardImage() error = %v, want nil", err)
		return
	}

	if protocol == "" {
		t.Error("CardImage() should return a protocol name")
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
		_, err := w.Write([]byte("not a valid PNG"))
		if err != nil {
			return
		}
	}))
	defer server.Close()

	result, protocol, err := CardImage(server.URL)

	if err == nil {
		t.Error("CardImage() should return error for invalid image data")
	}

	if result != "" {
		t.Errorf("CardImage() on error should return empty string, got %v", result)
	}

	if protocol != "" {
		t.Errorf("CardImage() on error should return empty protocol, got %v", protocol)
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

	result, protocol, err := CardImage(server.URL)

	if err == nil {
		t.Error("CardImage() should return error for non-200 response")
	}

	if result != "" {
		t.Errorf("CardImage() on error should return empty string, got %v", result)
	}

	if protocol != "" {
		t.Errorf("CardImage() on error should return empty protocol, got %v", protocol)
	}

	if !strings.Contains(err.Error(), "non-200 response") {
		t.Errorf("Error message should mention 'non-200 response', got: %v", err)
	}
}

func TestSupportsKittyGraphics(t *testing.T) {
	tests := []struct {
		name        string
		envVars     map[string]string
		wantSupport bool
	}{
		{
			name: "kitty terminal via KITTY_WINDOW_ID",
			envVars: map[string]string{
				"KITTY_WINDOW_ID": "1",
			},
			wantSupport: true,
		},
		{
			name: "kitty via TERM_PROGRAM",
			envVars: map[string]string{
				"TERM_PROGRAM": "kitty",
			},
			wantSupport: true,
		},
		{
			name: "kitty via TERM_PROGRAM uppercase",
			envVars: map[string]string{
				"TERM_PROGRAM": "KITTY",
			},
			wantSupport: true,
		},
		{
			name: "ghostty via TERM_PROGRAM",
			envVars: map[string]string{
				"TERM_PROGRAM": "ghostty",
			},
			wantSupport: true,
		},
		{
			name: "ghostty via TERM_PROGRAM uppercase",
			envVars: map[string]string{
				"TERM_PROGRAM": "Ghostty",
			},
			wantSupport: true,
		},
		{
			name: "wezterm via TERM_PROGRAM",
			envVars: map[string]string{
				"TERM_PROGRAM": "WezTerm",
			},
			wantSupport: true,
		},
		{
			name: "ghostty via TERM variable",
			envVars: map[string]string{
				"TERM": "xterm-ghostty",
			},
			wantSupport: true,
		},
		{
			name: "kitty via TERM variable",
			envVars: map[string]string{
				"TERM": "xterm-kitty",
			},
			wantSupport: true,
		},
		{
			name: "unsupported terminal - Apple Terminal",
			envVars: map[string]string{
				"TERM_PROGRAM": "Apple_Terminal",
				"TERM":         "xterm-256color",
			},
			wantSupport: false,
		},
		{
			name: "unsupported terminal - iTerm2",
			envVars: map[string]string{
				"TERM_PROGRAM": "iTerm.app",
				"TERM":         "xterm-256color",
			},
			wantSupport: false,
		},
		{
			name: "unsupported terminal - GNOME Terminal",
			envVars: map[string]string{
				"TERM": "xterm",
			},
			wantSupport: false,
		},
		{
			name:        "no environment variables set",
			envVars:     map[string]string{},
			wantSupport: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original env vars
			origVars := map[string]string{
				"KITTY_WINDOW_ID": os.Getenv("KITTY_WINDOW_ID"),
				"TERM_PROGRAM":    os.Getenv("TERM_PROGRAM"),
				"TERM":            os.Getenv("TERM"),
			}

			// Clear all relevant env vars first
			os.Unsetenv("KITTY_WINDOW_ID")
			os.Unsetenv("TERM_PROGRAM")
			os.Unsetenv("TERM")

			// Set test env vars
			for key, val := range tt.envVars {
				os.Setenv(key, val)
			}

			// Cleanup after test
			defer func() {
				for key, val := range origVars {
					if val == "" {
						os.Unsetenv(key)
					} else {
						os.Setenv(key, val)
					}
				}
			}()

			got := supportsKittyGraphics()
			if got != tt.wantSupport {
				t.Errorf("supportsKittyGraphics() = %v, want %v", got, tt.wantSupport)
			}
		})
	}
}

func TestSupportsSixelGraphics(t *testing.T) {
	tests := []struct {
		name        string
		envVars     map[string]string
		wantSupport bool
	}{
		{
			name: "iterm2 via TERM_PROGRAM",
			envVars: map[string]string{
				"TERM_PROGRAM": "iTerm.app",
			},
			wantSupport: true,
		},
		{
			name: "wezterm via TERM_PROGRAM",
			envVars: map[string]string{
				"TERM_PROGRAM": "WezTerm",
			},
			wantSupport: true,
		},
		{
			name: "wezterm lowercase",
			envVars: map[string]string{
				"TERM_PROGRAM": "wezterm",
			},
			wantSupport: true,
		},
		{
			name: "rio via TERM_PROGRAM",
			envVars: map[string]string{
				"TERM_PROGRAM": "rio",
			},
			wantSupport: true,
		},
		{
			name: "konsole via TERM_PROGRAM",
			envVars: map[string]string{
				"TERM_PROGRAM": "Konsole",
			},
			wantSupport: true,
		},
		{
			name: "foot via TERM",
			envVars: map[string]string{
				"TERM": "foot",
			},
			wantSupport: true,
		},
		{
			name: "foot with suffix",
			envVars: map[string]string{
				"TERM": "foot-extra",
			},
			wantSupport: true,
		},
		{
			name: "xterm-sixel via TERM",
			envVars: map[string]string{
				"TERM": "xterm-sixel",
			},
			wantSupport: true,
		},
		{
			name: "unsupported terminal - Apple Terminal",
			envVars: map[string]string{
				"TERM_PROGRAM": "Apple_Terminal",
				"TERM":         "xterm-256color",
			},
			wantSupport: false,
		},
		{
			name: "unsupported terminal - Alacritty",
			envVars: map[string]string{
				"TERM": "alacritty",
			},
			wantSupport: false,
		},
		{
			name: "unsupported terminal - standard xterm",
			envVars: map[string]string{
				"TERM": "xterm",
			},
			wantSupport: false,
		},
		{
			name: "unsupported terminal - xterm-256color",
			envVars: map[string]string{
				"TERM": "xterm-256color",
			},
			wantSupport: false,
		},
		{
			name:        "no environment variables set",
			envVars:     map[string]string{},
			wantSupport: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original env vars
			origVars := map[string]string{
				"TERM_PROGRAM": os.Getenv("TERM_PROGRAM"),
				"TERM":         os.Getenv("TERM"),
			}

			// Clear all relevant env vars first
			os.Unsetenv("TERM_PROGRAM")
			os.Unsetenv("TERM")

			// Set test env vars
			for key, val := range tt.envVars {
				os.Setenv(key, val)
			}

			// Cleanup
			defer func() {
				for key, val := range origVars {
					if val == "" {
						os.Unsetenv(key)
					} else {
						os.Setenv(key, val)
					}
				}
			}()

			got := supportsSixelGraphics()
			if got != tt.wantSupport {
				t.Errorf("supportsSixelGraphics() = %v, want %v", got, tt.wantSupport)
			}
		})
	}
}
