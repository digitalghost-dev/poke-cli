package flags

import (
	"os"
	"testing"

	"github.com/digitalghost-dev/poke-cli/cmd/utils"
	"github.com/digitalghost-dev/poke-cli/styling"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLatestVersionFlag(t *testing.T) {
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
			name:           "Get latest version with short flag",
			args:           []string{"-l"},
			expectedOutput: utils.LoadGolden(t, "main_latest_flag.golden"),
			expectedError:  false,
		},
		{
			name:           "Get latest version with long flag",
			args:           []string{"--latest"},
			expectedOutput: utils.LoadGolden(t, "main_latest_flag.golden"),
			expectedError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalArgs := os.Args
			os.Args = append([]string{"poke-cli"}, tt.args...)
			defer func() { os.Args = originalArgs }()

			output, err := LatestFlag()
			cleanOutput := styling.StripANSI(output)

			if tt.expectedError {
				require.Error(t, err, "Expected an error")
			} else {
				require.NoError(t, err, "Expected no error")
			}

			assert.Equal(t, tt.expectedOutput, cleanOutput, "Output should match expected")
		})
	}
}
