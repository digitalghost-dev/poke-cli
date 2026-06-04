package utils

import (
	"testing"

	"github.com/digitalghost-dev/poke-cli/styling"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCheckLength(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		maxLength   int
		wantErr     bool
		expectedErr string
	}{
		{name: "empty slice", args: []string{}, maxLength: 1},
		{name: "within limit", args: []string{"arg1", "arg2"}, maxLength: 3},
		{name: "exactly at limit", args: []string{"arg1", "arg2", "arg3"}, maxLength: 3},
		{name: "exceeds limit", args: []string{"arg1", "arg2", "arg3", "arg4"}, maxLength: 3, wantErr: true, expectedErr: "Too many arguments"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := checkLength(tt.args, tt.maxLength)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, styling.StripANSI(err.Error()), tt.expectedErr)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestValidateArgs(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		validator Validator
		wantErr   bool
		contains  string
	}{
		{
			name:      "ability accepts name and flag",
			args:      []string{"ability", "technician", "--pokemon"},
			validator: Validator{MaxArgs: 3, CmdName: "ability", RequireName: true, HasFlags: true},
		},
		{
			name:      "ability rejects missing name",
			args:      []string{"ability"},
			validator: Validator{MaxArgs: 3, CmdName: "ability", RequireName: true, HasFlags: true},
			wantErr:   true,
			contains:  "Please declare",
		},
		{
			name:      "ability rejects too many args",
			args:      []string{"ability", "strong-jaw", "all", "pokemon"},
			validator: Validator{MaxArgs: 3, CmdName: "ability", RequireName: true, HasFlags: true},
			wantErr:   true,
			contains:  "Too many arguments",
		},
		{
			name:      "berry accepts no name",
			args:      []string{"berry"},
			validator: Validator{MaxArgs: 3, CmdName: "berry"},
		},
		{
			name:      "berry accepts name",
			args:      []string{"berry", "oran"},
			validator: Validator{MaxArgs: 3, CmdName: "berry"},
		},
		{
			name:      "berry accepts help",
			args:      []string{"berry", "--help"},
			validator: Validator{MaxArgs: 3, CmdName: "berry"},
		},
		{
			name:      "berry rejects extra arg",
			args:      []string{"berry", "oran", "sitrus"},
			validator: Validator{MaxArgs: 3, CmdName: "berry"},
			wantErr:   true,
			contains:  "only available options",
		},
		{
			name:      "card accepts no args",
			args:      []string{"card"},
			validator: Validator{MaxArgs: 2, CmdName: "card"},
		},
		{
			name:      "card accepts help",
			args:      []string{"card", "--help"},
			validator: Validator{MaxArgs: 2, CmdName: "card"},
		},
		{
			name:      "card rejects extra arg",
			args:      []string{"card", "scarlet"},
			validator: Validator{MaxArgs: 2, CmdName: "card"},
			wantErr:   true,
			contains:  "only available options",
		},
		{
			name:      "item accepts name",
			args:      []string{"item", "potion"},
			validator: Validator{MaxArgs: 2, CmdName: "item", RequireName: true},
		},
		{
			name:      "item rejects missing name",
			args:      []string{"item"},
			validator: Validator{MaxArgs: 2, CmdName: "item", RequireName: true},
			wantErr:   true,
			contains:  "Please declare",
		},
		{
			name:      "item rejects too many args",
			args:      []string{"item", "potion", "extra"},
			validator: Validator{MaxArgs: 2, CmdName: "item", RequireName: true},
			wantErr:   true,
			contains:  "Too many arguments",
		},
		{
			name:      "move accepts name",
			args:      []string{"move", "thunderbolt"},
			validator: Validator{MaxArgs: 2, CmdName: "move", RequireName: true},
		},
		{
			name:      "move rejects missing name",
			args:      []string{"move"},
			validator: Validator{MaxArgs: 2, CmdName: "move", RequireName: true},
			wantErr:   true,
			contains:  "Please declare",
		},
		{
			name:      "move rejects too many args",
			args:      []string{"move", "tackle", "scratch"},
			validator: Validator{MaxArgs: 2, CmdName: "move", RequireName: true},
			wantErr:   true,
			contains:  "Too many arguments",
		},
		{
			name:      "mechanics accepts no args",
			args:      []string{"mechanics"},
			validator: Validator{MaxArgs: 2, CmdName: "mechanics", HasFlags: true},
		},
		{
			name:      "mechanics accepts natures flag",
			args:      []string{"mechanics", "--natures"},
			validator: Validator{MaxArgs: 2, CmdName: "mechanics", HasFlags: true},
		},
		{
			name:      "mechanics rejects too many args",
			args:      []string{"mechanics", "--natures", "extra"},
			validator: Validator{MaxArgs: 2, CmdName: "mechanics", HasFlags: true},
			wantErr:   true,
			contains:  "Too many arguments",
		},
		{
			name:      "search accepts no args",
			args:      []string{"search"},
			validator: Validator{MaxArgs: 2, CmdName: "search"},
		},
		{
			name:      "search accepts help",
			args:      []string{"search", "--help"},
			validator: Validator{MaxArgs: 2, CmdName: "search"},
		},
		{
			name:      "search rejects extra arg",
			args:      []string{"search", "pokemon"},
			validator: Validator{MaxArgs: 2, CmdName: "search"},
			wantErr:   true,
			contains:  "only available options",
		},
		{
			name:      "speed accepts no args",
			args:      []string{"speed"},
			validator: Validator{MaxArgs: 2, CmdName: "speed"},
		},
		{
			name:      "speed accepts help",
			args:      []string{"speed", "--help"},
			validator: Validator{MaxArgs: 2, CmdName: "speed"},
		},
		{
			name:      "speed rejects extra arg",
			args:      []string{"speed", "100"},
			validator: Validator{MaxArgs: 2, CmdName: "speed"},
			wantErr:   true,
			contains:  "only available options",
		},
		{
			name:      "tcg accepts no args",
			args:      []string{"tcg"},
			validator: Validator{MaxArgs: 2, CmdName: "tcg", HasFlags: true},
		},
		{
			name:      "tcg accepts web flag",
			args:      []string{"tcg", "--web"},
			validator: Validator{MaxArgs: 2, CmdName: "tcg", HasFlags: true},
		},
		{
			name:      "tcg rejects too many args",
			args:      []string{"tcg", "--web", "extra"},
			validator: Validator{MaxArgs: 2, CmdName: "tcg", HasFlags: true},
			wantErr:   true,
			contains:  "Too many arguments",
		},
		{
			name:      "types accepts no args",
			args:      []string{"types"},
			validator: Validator{MaxArgs: 2, CmdName: "types"},
		},
		{
			name:      "types accepts help",
			args:      []string{"types", "--help"},
			validator: Validator{MaxArgs: 2, CmdName: "types"},
		},
		{
			name:      "types rejects extra arg",
			args:      []string{"types", "rock"},
			validator: Validator{MaxArgs: 2, CmdName: "types"},
			wantErr:   true,
			contains:  "only available options",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateArgs(tt.args, tt.validator)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, styling.StripANSI(err.Error()), tt.contains)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestValidatePokemonArgs(t *testing.T) {
	validInputs := [][]string{
		{"pokemon", "--help"},
		{"pokemon", "mankey"},
		{"pokemon", "talonflame", "--stats", "--types"},
		{"pokemon", "passimian", "--abilities", "-t"},
		{"pokemon", "dodrio", "-a", "-s", "-t"},
		{"pokemon", "dragalge", "-a", "-s", "-t", "--image=sm"},
		{"pokemon", "squirtle", "-a", "-s"},
		{"pokemon", "dragapult", "-s", "-a"},
	}

	for _, input := range validInputs {
		err := ValidatePokemonArgs(input)
		require.NoError(t, err, "expected no error for valid input %v", input)
	}

	invalidInputs := []struct {
		args     []string
		contains string
	}{
		{args: []string{"pokemon"}, contains: "Please declare"},
		{args: []string{"pokemons"}, contains: "Please declare"},
		{args: []string{"pokemon", "mewtwo", "--"}, contains: "Empty flag"},
		{args: []string{"pokemon", "baxcalibur", "-"}, contains: "Empty flag"},
		{args: []string{"pokemon", "charizard", "extraArg"}, contains: "Invalid argument"},
		{args: []string{"pokemon", "hypo", "--abilities", "-s", "--types", "--image=sm", "-m", "-p"}, contains: "Too many arguments"},
	}

	for _, input := range invalidInputs {
		err := ValidatePokemonArgs(input.args)
		require.Error(t, err, "expected error for invalid input %v", input.args)
		assert.Contains(t, styling.StripANSI(err.Error()), input.contains)
	}
}
