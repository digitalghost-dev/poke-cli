package move

import (
	"github.com/digitalghost-dev/poke-cli/cmd/utils"
	"github.com/digitalghost-dev/poke-cli/styling"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestMoveCommand(t *testing.T) {
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
		expectedError  bool
	}{
		{
			name:           "Move help flag",
			args:           []string{"move", "--help"},
			expectedOutput: utils.LoadGolden(t, "move_help.golden"),
		},
		{
			name:           "Move help flag",
			args:           []string{"move", "-h"},
			expectedOutput: utils.LoadGolden(t, "move_help.golden"),
		},
		{
			name:           "Select 'shadow-ball' as move",
			args:           []string{"move", "shadow-ball"},
			expectedOutput: utils.LoadGolden(t, "move.golden"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := MoveCommand(tt.args)
			if tt.expectedError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			cleanOutput := styling.StripANSI(output)

			assert.Equal(t, tt.expectedOutput, cleanOutput, "Output should match expected")
		})
	}
}
