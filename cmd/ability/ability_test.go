package ability

import (
	"github.com/digitalghost-dev/poke-cli/cmd/utils"
	"github.com/digitalghost-dev/poke-cli/styling"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestAbilityCommand(t *testing.T) {
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
			name:           "Ability help flag",
			args:           []string{"ability", "--help"},
			expectedOutput: utils.LoadGolden(t, "ability_help.golden"),
		},
		{
			name:           "Ability command: clear-body",
			args:           []string{"ability", "clear-body"},
			expectedOutput: utils.LoadGolden(t, "ability.golden"),
		},
		{
			name:           "Misspelled ability name",
			args:           []string{"ability", "bulletproff"},
			expectedOutput: utils.LoadGolden(t, "ability_misspelled.golden"),
			wantError:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalArgs := os.Args
			os.Args = append([]string{"poke-cli"}, tt.args...)
			defer func() { os.Args = originalArgs }()

			output := AbilityCommand()
			cleanOutput := styling.StripANSI(output)

			assert.Equal(t, tt.expectedOutput, cleanOutput, "Output should match expected")
		})
	}
}
