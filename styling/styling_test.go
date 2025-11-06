package styling

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"image/color"
	"testing"
)

func TestGetTypeColor(t *testing.T) {
	// Test known types
	for typeName, expectedColor := range ColorMap {
		t.Run(typeName, func(t *testing.T) {
			typeColor := GetTypeColor(typeName)
			assert.Equal(t, expectedColor, typeColor, "Expected color for type %s to be %s", typeName, expectedColor)
		})
	}
}

func TestStripANSI(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "No ANSI codes",
			input:    "Hello, World!",
			expected: "Hello, World!",
		},
		{
			name:     "Simple ANSI color code",
			input:    "\x1b[31mHello\x1b[0m",
			expected: "Hello",
		},
		{
			name:     "Multiple ANSI codes",
			input:    "\x1b[1;34mBold Blue\x1b[0m Text",
			expected: "Bold Blue Text",
		},
		{
			name:     "Nested ANSI codes",
			input:    "\x1b[1mBold \x1b[31mRed\x1b[0m",
			expected: "Bold Red",
		},
		{
			name:     "Only ANSI codes",
			input:    "\x1b[1;32m\x1b[0m",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := StripANSI(tt.input)
			if output != tt.expected {
				t.Errorf("StripANSI(%q) = %q; want %q", tt.input, output, tt.expected)
			}
		})
	}
}

func TestColor_RGBA(t *testing.T) {
	col := Color{R: 1.0, G: 0.5, B: 0.0}
	r, g, b, a := col.RGBA()
	if r != 65535 || g != 32768 || b != 0 || a != 65535 {
		t.Errorf("Unexpected RGBA values: got (%d, %d, %d, %d)", r, g, b, a)
	}
}

func TestMakeColor(t *testing.T) {
	// color.RGBA uses 8-bit values, which are multiplied by 0x101 in RGBA() to get 16-bit range.
	c := color.RGBA{R: 255, G: 128, B: 0, A: 255}
	col, ok := MakeColor(c)
	if !ok {
		t.Fatal("Expected color to be valid (alpha != 0)")
	}

	// Allowing small float tolerance due to conversion
	if diff := func(a, b float64) bool { return fmt.Sprintf("%.2f", a) != fmt.Sprintf("%.2f", b) }; diff(col.R, 1.0) || diff(col.G, 0.5) || diff(col.B, 0.0) {
		t.Errorf("Unexpected color values: got %+v", col)
	}

	// Test alpha = 0 case
	cTransparent := color.RGBA{0, 0, 0, 0}
	_, ok = MakeColor(cTransparent)
	if ok {
		t.Error("Expected MakeColor to return false for fully transparent color")
	}
}

func TestColor_Hex(t *testing.T) {
	col := Color{R: 1.0, G: 0.0, B: 0.5}
	hex := col.Hex()
	expected := "#ff0080"
	if hex != expected {
		t.Errorf("Expected %s, got %s", expected, hex)
	}
}

func TestFormTheme(t *testing.T) {
	theme := FormTheme()

	assert.NotNil(t, theme, "FormTheme should return a non-nil theme")
	assert.NotNil(t, theme.Focused, "Focused state should be configured")
	assert.NotNil(t, theme.Blurred, "Blurred state should be configured")
	assert.NotNil(t, theme.Group, "Group state should be configured")

	focusedButtonStyle := theme.Focused.FocusedButton
	assert.NotNil(t, focusedButtonStyle, "Focused button style should be set")

	assert.Equal(t, theme.Focused.FocusedButton, theme.Focused.Next, "Next button should use focused button style")
	assert.NotNil(t, theme.Blurred.Base, "Blurred base should be configured")
	assert.Equal(t, theme.Focused.Title, theme.Group.Title, "Group title should match focused title")
	assert.Equal(t, theme.Focused.Description, theme.Group.Description, "Group description should match focused description")
}
