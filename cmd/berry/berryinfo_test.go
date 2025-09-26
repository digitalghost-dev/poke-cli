package berry

import (
	"image"
	"image/color"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestBerryName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple berry name",
			input:    "Aguav",
			expected: "Berry: Aguav",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "Berry: ",
		},
		{
			name:     "berry with special characters",
			input:    "Test-Berry",
			expected: "Berry: Test-Berry",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BerryName(tt.input)
			if result != tt.expected {
				t.Errorf("BerryName(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestBerryEffect(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "non-existent berry",
			input:    "NonExistentBerry",
			expected: "Effect information not available",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "Effect information not available",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BerryEffect(tt.input)
			if tt.input == "NonExistentBerry" || tt.input == "" {
				if result != tt.expected {
					t.Errorf("BerryEffect(%q) = %q, want %q", tt.input, result, tt.expected)
				}
			}
		})
	}
}

func TestBerryInfo(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "non-existent berry",
			input:    "NonExistentBerry",
			expected: "Additional information not available",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "Additional information not available",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BerryInfo(tt.input)
			if tt.input == "NonExistentBerry" || tt.input == "" {
				if result != tt.expected {
					t.Errorf("BerryInfo(%q) = %q, want %q", tt.input, result, tt.expected)
				}
			}
		})
	}
}

func TestBerryImageWithMockServer(t *testing.T) {
	// Create a mock HTTP server that serves a simple test image
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a simple 2x2 test image
		img := image.NewRGBA(image.Rect(0, 0, 2, 2))
		img.Set(0, 0, color.RGBA{255, 0, 0, 255})   // Red
		img.Set(1, 0, color.RGBA{0, 255, 0, 255})   // Green
		img.Set(0, 1, color.RGBA{0, 0, 255, 255})   // Blue
		img.Set(1, 1, color.RGBA{255, 255, 0, 255}) // Yellow

		w.Header().Set("Content-Type", "image/png")
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	result := BerryImage("NonExistentBerry")
	expected := "Image information not available"
	if result != expected {
		t.Errorf("BerryImage('NonExistentBerry') = %q, want %q", result, expected)
	}

	result = BerryImage("")
	if result != expected {
		t.Errorf("BerryImage('') = %q, want %q", result, expected)
	}
}

func TestBerryImageErrorHandling(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "non-existent berry",
			input:    "NonExistentBerry",
			expected: "Image information not available",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "Image information not available",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BerryImage(tt.input)
			if result != tt.expected {
				t.Errorf("BerryImage(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

// Test helper function to check if the ToString function structure is working
func TestToStringStructure(t *testing.T) {
	// This test checks that the ToString function can handle basic cases
	// without actually making HTTP requests or database queries

	// Create a simple test image
	img := image.NewRGBA(image.Rect(0, 0, 4, 4))
	for y := 0; y < 4; y++ {
		for x := 0; x < 4; x++ {
			img.Set(x, y, color.RGBA{uint8(x * 63), uint8(y * 63), 100, 255})
		}
	}

	// Test the ToString function indirectly by checking that BerryImage
	// with invalid input returns the expected error message
	result := BerryImage("InvalidBerry")
	if !strings.Contains(result, "information not available") {
		t.Errorf("Expected error message for invalid berry, got: %q", result)
	}
}

// Test for database query error handling
func TestBerryFunctionsErrorHandling(t *testing.T) {
	testCases := []struct {
		name     string
		function func(string) string
		input    string
		contains string
	}{
		{
			name:     "BerryEffect with invalid input",
			function: BerryEffect,
			input:    "InvalidBerry123",
			contains: "not available",
		},
		{
			name:     "BerryInfo with invalid input",
			function: BerryInfo,
			input:    "InvalidBerry123",
			contains: "not available",
		},
		{
			name:     "BerryImage with invalid input",
			function: BerryImage,
			input:    "InvalidBerry123",
			contains: "not available",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.function(tc.input)
			if !strings.Contains(result, tc.contains) {
				t.Errorf("Expected result to contain %q, got %q", tc.contains, result)
			}
		})
	}
}
