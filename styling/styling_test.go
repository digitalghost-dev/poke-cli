package styling

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetTypeColor(t *testing.T) {
	// Test known types
	for typeName, expectedColor := range ColorMap {
		t.Run(typeName, func(t *testing.T) {
			color := GetTypeColor(typeName)
			assert.Equal(t, expectedColor, color, "Expected color for type %s to be %s", typeName, expectedColor)
		})
	}

	// Test unknown type
	t.Run("unknown type", func(t *testing.T) {
		color := GetTypeColor("unknown")
		assert.Equal(t, "", color, "Expected color for unknown type to be an empty string")
	})
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
