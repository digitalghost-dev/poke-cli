package speed

import (
	"github.com/digitalghost-dev/poke-cli/cmd/utils"
	"github.com/digitalghost-dev/poke-cli/styling"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestSpeedCommand(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedOutput string
		wantError      bool
	}{
		{
			name:           "Speed help flag",
			args:           []string{"speed", "--help"},
			expectedOutput: utils.LoadGolden(t, "speed_help.golden"),
		},
		{
			name:           "Speed help flag",
			args:           []string{"speed", "-h"},
			expectedOutput: utils.LoadGolden(t, "speed_help.golden"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalArgs := os.Args
			os.Args = append([]string{"poke-cli"}, tt.args...)
			defer func() { os.Args = originalArgs }()

			output, _ := SpeedCommand()
			cleanOutput := styling.StripANSI(output)

			assert.Equal(t, tt.expectedOutput, cleanOutput, "Output should match expected")
		})
	}
}

func TestFormula(t *testing.T) {
	// Save the original DefaultSpeedStat function and restore it after the test
	originalSpeedStat := DefaultSpeedStat
	defer func() { DefaultSpeedStat = originalSpeedStat }()

	// Create a mock SpeedStatFunc that always returns 90 (Pikachu's base speed)
	DefaultSpeedStat = func(name string) (string, error) {
		return "90", nil
	}

	tests := []struct {
		name           string
		pokemonDetails PokemonDetails
		expectedSpeed  string
		wantError      bool
	}{
		{
			name: "Basic calculation with default values",
			pokemonDetails: PokemonDetails{
				Name:       "pikachu",
				SpeedStage: "0",
				Nature:     "0%",
				Level:      "50",
				Modifier:   []string{},
				Ability:    "None",
				SpeedEV:    "0",
				SpeedIV:    "0",
			},
			expectedSpeed: "95",
			wantError:     false,
		},
		{
			name: "With positive nature",
			pokemonDetails: PokemonDetails{
				Name:       "pikachu",
				SpeedStage: "0",
				Nature:     "+10%",
				Level:      "50",
				Modifier:   []string{},
				Ability:    "None",
				SpeedEV:    "0",
				SpeedIV:    "0",
			},
			expectedSpeed: "104",
			wantError:     false,
		},
		{
			name: "With negative nature",
			pokemonDetails: PokemonDetails{
				Name:       "pikachu",
				SpeedStage: "0",
				Nature:     "-10%",
				Level:      "50",
				Modifier:   []string{},
				Ability:    "None",
				SpeedEV:    "0",
				SpeedIV:    "0",
			},
			expectedSpeed: "85",
			wantError:     false,
		},
		{
			name: "With ability multiplier",
			pokemonDetails: PokemonDetails{
				Name:       "pikachu",
				SpeedStage: "0",
				Nature:     "0%",
				Level:      "50",
				Modifier:   []string{},
				Ability:    "Swift Swim",
				SpeedEV:    "0",
				SpeedIV:    "0",
			},
			expectedSpeed: "190",
			wantError:     false,
		},
		{
			name: "With positive speed stage",
			pokemonDetails: PokemonDetails{
				Name:       "pikachu",
				SpeedStage: "2",
				Nature:     "0%",
				Level:      "50",
				Modifier:   []string{},
				Ability:    "None",
				SpeedEV:    "0",
				SpeedIV:    "0",
			},
			expectedSpeed: "190",
			wantError:     false,
		},
		{
			name: "With negative speed stage",
			pokemonDetails: PokemonDetails{
				Name:       "pikachu",
				SpeedStage: "-2",
				Nature:     "0%",
				Level:      "50",
				Modifier:   []string{},
				Ability:    "None",
				SpeedEV:    "0",
				SpeedIV:    "0",
			},
			expectedSpeed: "47",
			wantError:     false,
		},
		{
			name: "With modifier",
			pokemonDetails: PokemonDetails{
				Name:       "pikachu",
				SpeedStage: "0",
				Nature:     "0%",
				Level:      "50",
				Modifier:   []string{"Choice Scarf"},
				Ability:    "None",
				SpeedEV:    "0",
				SpeedIV:    "0",
			},
			expectedSpeed: "142",
			wantError:     false,
		},
		{
			name: "With multiple modifiers",
			pokemonDetails: PokemonDetails{
				Name:       "pikachu",
				SpeedStage: "0",
				Nature:     "0%",
				Level:      "50",
				Modifier:   []string{"Choice Scarf", "Tailwind"},
				Ability:    "None",
				SpeedEV:    "0",
				SpeedIV:    "0",
			},
			expectedSpeed: "285",
			wantError:     false,
		},
		{
			name: "With EVs and IVs",
			pokemonDetails: PokemonDetails{
				Name:       "pikachu",
				SpeedStage: "0",
				Nature:     "0%",
				Level:      "50",
				Modifier:   []string{},
				Ability:    "None",
				SpeedEV:    "252",
				SpeedIV:    "31",
			},
			expectedSpeed: "142",
			wantError:     false,
		},
		{
			name: "Complex scenario",
			pokemonDetails: PokemonDetails{
				Name:       "pikachu",
				SpeedStage: "1",
				Nature:     "+10%",
				Level:      "100",
				Modifier:   []string{"Choice Scarf"},
				Ability:    "Quick Feet",
				SpeedEV:    "252",
				SpeedIV:    "31",
			},
			expectedSpeed: "1035",
			wantError:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pokemon = tt.pokemonDetails

			result, err := formula()

			if tt.wantError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)

				cleanOutput := styling.StripANSI(result)

				t.Logf("Expected speed: %s", tt.expectedSpeed)

				// Check if the output contains the expected speed
				assert.Contains(t, cleanOutput, "current speed of "+tt.expectedSpeed,
					"Expected speed "+tt.expectedSpeed+" not found in output: "+cleanOutput)
			}
		})
	}
}
