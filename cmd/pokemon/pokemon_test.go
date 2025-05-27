package pokemon

import (
	"github.com/digitalghost-dev/poke-cli/cmd/utils"
	"github.com/digitalghost-dev/poke-cli/styling"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestPokemonCommand(t *testing.T) {
	err := os.Setenv("GO_TESTING", "1")
	if err != nil {
		t.Fatalf("Failed to set GO_TESTING env var: %v", err)
	}

	defer func() {
		err := os.Unsetenv("GO_TESTING")
		if err != nil {
			t.Logf("Warning: failed to unset GO_TESTING: %v", err)
		}
	}()

	tests := []struct {
		name           string
		args           []string
		expectedOutput string
		wantError      bool
	}{
		{
			name:           "Pokemon help flag",
			args:           []string{"pokemon", "--help"},
			expectedOutput: utils.LoadGolden(t, "pokemon_help.golden"),
		},
		{
			name:           "Pokemon abilities flag",
			args:           []string{"pokemon", "metagross", "--abilities"},
			expectedOutput: utils.LoadGolden(t, "pokemon_abilities.golden"),
		},
		{
			name:           "Pokemon image flag",
			args:           []string{"pokemon", "skeledirge", "--image=md"},
			expectedOutput: utils.LoadGolden(t, "pokemon_image.golden"),
		},
		{
			name:           "Pokemon invalid image flag",
			args:           []string{"pokemon", "tryanitar", "--image="},
			expectedOutput: utils.LoadGolden(t, "pokemon_invalid_image_flag.golden"),
			wantError:      true,
		},
		{
			name:           "Pokemon stats flag",
			args:           []string{"pokemon", "toxicroak", "--stats"},
			expectedOutput: utils.LoadGolden(t, "pokemon_stats.golden"),
		},
		{
			name:           "Pokemon typed flags",
			args:           []string{"pokemon", "armarouge", "--types"},
			expectedOutput: utils.LoadGolden(t, "pokemon_types.golden"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalArgs := os.Args
			os.Args = append([]string{"poke-cli"}, tt.args...)
			defer func() { os.Args = originalArgs }()

			output, _ := PokemonCommand()
			cleanOutput := styling.StripANSI(output)

			assert.Equal(t, tt.expectedOutput, cleanOutput, "Output should match expected")
		})
	}
}
