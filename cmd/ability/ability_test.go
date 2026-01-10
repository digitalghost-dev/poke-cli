package ability

import (
	"os"
	"testing"

	"github.com/digitalghost-dev/poke-cli/cmd/utils"
	"github.com/digitalghost-dev/poke-cli/styling"
	"github.com/stretchr/testify/assert"
)

func TestAbilityCommand(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedOutput string
		wantError      bool
	}{
		{
			name:           "Ability help flag",
			args:           []string{"ability", "--help"},
			expectedOutput: utils.LoadGolden(t, "ability_help.golden"),
		},
		{
			name:           "Ability help flag",
			args:           []string{"ability", "-h"},
			expectedOutput: utils.LoadGolden(t, "ability_help.golden"),
		},
		{
			name:           "Ability command: clear-body",
			args:           []string{"ability", "clear-body"},
			expectedOutput: utils.LoadGolden(t, "ability.golden"),
		},
		{
			name:           "Ability command: beads-of-ruin",
			args:           []string{"ability", "beads-of-ruin"},
			expectedOutput: utils.LoadGolden(t, "ability-ii.golden"),
		},
		{
			name:           "Misspelled ability name",
			args:           []string{"ability", "bulletproff"},
			expectedOutput: utils.LoadGolden(t, "ability_misspelled.golden"),
			wantError:      true,
		},
		{
			name:           "Ability command: --pokemon flag",
			args:           []string{"ability", "anger-point", "--pokemon"},
			expectedOutput: utils.LoadGolden(t, "ability_flag_pokemon.golden"),
		},
		{
			name:           "Ability command: special character in API call",
			args:           []string{"ability", "poison-point"},
			expectedOutput: utils.LoadGolden(t, "ability_poison_point.golden"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalArgs := os.Args
			os.Args = append([]string{"poke-cli"}, tt.args...)
			defer func() { os.Args = originalArgs }()

			output, _ := AbilityCommand()
			cleanOutput := styling.StripANSI(output)

			assert.Equal(t, tt.expectedOutput, cleanOutput, "Output should match expected")
		})
	}
}
