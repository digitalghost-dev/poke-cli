package utils

import (
	"errors"
	"testing"

	"github.com/digitalghost-dev/poke-cli/styling"
	"github.com/stretchr/testify/assert"
)

func TestFormatResourceErrors(t *testing.T) {
	tests := []struct {
		name     string
		format   func() string
		contains []string
	}{
		{
			name:   "not found",
			format: func() string { return FormatNotFoundError("Pokémon") },
			contains: []string{
				"Pokémon not found.",
				"Perhaps a typo?",
				"Missing a hyphen instead of a space?",
			},
		},
		{
			name:   "network",
			format: func() string { return FormatNetworkError("Pokémon") },
			contains: []string{
				"Could not reach Pokémon data.",
				"Check your connection and try again.",
			},
		},
		{
			name:   "server",
			format: func() string { return FormatServerError("Pokémon") },
			contains: []string{
				"Pokémon data source returned a server error.",
				"Please try again later.",
			},
		},
		{
			name:   "unexpected data",
			format: func() string { return FormatUnexpectedDataError("Pokémon") },
			contains: []string{
				"Pokémon data source returned data in an unexpected format.",
			},
		},
		{
			name:   "fetch with error",
			format: func() string { return FormatFetchError("Pokémon", errors.New("request failed")) },
			contains: []string{
				"Could not fetch Pokémon data.",
				"request failed",
			},
		},
		{
			name:   "fetch with nil error",
			format: func() string { return FormatFetchError("Pokémon", nil) },
			contains: []string{
				"Could not fetch Pokémon data.",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := styling.StripANSI(tt.format())

			assert.Contains(t, output, "✖ Error!")
			for _, expected := range tt.contains {
				assert.Contains(t, output, expected)
			}
		})
	}
}
