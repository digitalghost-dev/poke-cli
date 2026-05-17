package natures

import (
	"github.com/digitalghost-dev/poke-cli/cmd/utils"
	"github.com/digitalghost-dev/poke-cli/styling"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNaturesCommand(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedOutput string
		wantError      bool
	}{
		{
			name:           "Natures help flag",
			args:           []string{"natures", "--help"},
			expectedOutput: utils.LoadGolden(t, "natures_help.golden"),
		},
		{
			name:           "Natures help flag",
			args:           []string{"natures", "-h"},
			expectedOutput: utils.LoadGolden(t, "natures_help.golden"),
		},
		{
			name:           "Invalid extra argument",
			args:           []string{"natures", "brave"},
			expectedOutput: utils.LoadGolden(t, "natures_invalid_extra_arg.golden"),
			wantError:      true,
		},
		{
			name:           "Full Natures output with table",
			args:           []string{"natures"},
			expectedOutput: utils.LoadGolden(t, "natures.golden"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := NaturesCommand(tt.args)
			if tt.wantError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			cleanOutput := styling.StripANSI(output)

			assert.Equal(t, tt.expectedOutput, cleanOutput, "Output should match expected")
		})
	}
}
