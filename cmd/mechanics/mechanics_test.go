package mechanics

import (
	"testing"

	"github.com/digitalghost-dev/poke-cli/cmd/utils"
	"github.com/digitalghost-dev/poke-cli/styling"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMechanicsCommand(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedOutput string
		wantError      bool
	}{
		{
			name:           "Mechanics help flag --help",
			args:           []string{"mechanics", "--help"},
			expectedOutput: utils.LoadGolden(t, "mechanics_help.golden"),
		},
		{
			name:           "Mechanics help flag -h",
			args:           []string{"mechanics", "-h"},
			expectedOutput: utils.LoadGolden(t, "mechanics_help.golden"),
		},
		{
			name:           "No flag shows help",
			args:           []string{"mechanics"},
			expectedOutput: utils.LoadGolden(t, "mechanics_help.golden"),
		},
		{
			name:           "Natures flag --natures",
			args:           []string{"mechanics", "--natures"},
			expectedOutput: utils.LoadGolden(t, "mechanics_natures.golden"),
		},
		{
			name:           "Natures flag -n",
			args:           []string{"mechanics", "-n"},
			expectedOutput: utils.LoadGolden(t, "mechanics_natures.golden"),
		},
		{
			name:           "Too many arguments",
			args:           []string{"mechanics", "--natures", "extra"},
			expectedOutput: utils.LoadGolden(t, "mechanics_too_many_args.golden"),
			wantError:      true,
		},
		{
			name:           "Invalid flag",
			args:           []string{"mechanics", "--bogus"},
			expectedOutput: utils.LoadGolden(t, "mechanics_invalid_flag.golden"),
			wantError:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := MechanicsCommand(tt.args)
			if tt.wantError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			cleanOutput := styling.StripANSI(output)

			assert.Equal(t, tt.expectedOutput, cleanOutput, "Output should match expected")
		})
	}
}
