package item

import (
	"os"
	"testing"

	"github.com/digitalghost-dev/poke-cli/cmd/utils"
	"github.com/digitalghost-dev/poke-cli/styling"
	"github.com/stretchr/testify/assert"
)

func TestItemCommand(t *testing.T) {
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
			name:           "Item help flag",
			args:           []string{"item", "--help"},
			expectedOutput: utils.LoadGolden(t, "item_help.golden"),
		},
		{
			name:           "Item help flag",
			args:           []string{"item", "-h"},
			expectedOutput: utils.LoadGolden(t, "item_help.golden"),
		},
		{
			name:           "Select 'choice-band' as item",
			args:           []string{"item", "choice-band"},
			expectedOutput: utils.LoadGolden(t, "item.golden"),
		},
		{
			name:           "Select 'clear-amulet' as item with missing data",
			args:           []string{"item", "clear-amulet"},
			expectedOutput: utils.LoadGolden(t, "item_missing_data.golden"),
		},
		{
			name:           "Too many arguments",
			args:           []string{"item", "dubious-disc", "--help"},
			expectedOutput: utils.LoadGolden(t, "item_too_many_args.golden"),
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalArgs := os.Args
			os.Args = append([]string{"poke-cli"}, tt.args...)
			defer func() { os.Args = originalArgs }()

			output, _ := ItemCommand()
			cleanOutput := styling.StripANSI(output)

			assert.Equal(t, tt.expectedOutput, cleanOutput, "Output should match expected")
		})
	}
}
