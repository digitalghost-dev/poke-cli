package cmd

import (
	"github.com/digitalghost-dev/poke-cli/styling"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func loadGolden(t *testing.T, filename string) string {
	t.Helper()
	goldenPath := filepath.Join("..", "testdata", filename)
	content, err := os.ReadFile(goldenPath)
	if err != nil {
		t.Fatalf("failed to read golden file: %v", err)
	}
	return string(content)
}

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
			expectedOutput: loadGolden(t, "natures_help.golden"),
		},
		{
			name:           "Invalid extra argument",
			args:           []string{"natures", "brave"},
			expectedOutput: loadGolden(t, "natures_invalid_extra_arg.golden"),
			wantError:      true,
		},
		{
			name:           "Full Natures output with table",
			args:           []string{"natures"},
			expectedOutput: loadGolden(t, "natures.golden"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalArgs := os.Args
			os.Args = append([]string{"poke-cli"}, tt.args...)
			defer func() { os.Args = originalArgs }()

			output := NaturesCommand()
			cleanOutput := styling.StripANSI(output)

			assert.Equal(t, tt.expectedOutput, cleanOutput, "Output should match expected")
		})
	}
}
