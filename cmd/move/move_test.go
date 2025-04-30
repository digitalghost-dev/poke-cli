package move

import (
	"github.com/digitalghost-dev/poke-cli/styling"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func loadGolden(t *testing.T, filename string) string {
	t.Helper()
	goldenPath := filepath.Join("../..", "testdata", filename)
	content, err := os.ReadFile(goldenPath)
	if err != nil {
		t.Fatalf("failed to read golden file: %v", err)
	}
	return string(content)
}

func TestMoveCommand(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedOutput string
		expectedError  bool
	}{
		{
			name:           "Select 'Shadow-Ball' as move",
			args:           []string{"move", "shadow-ball"},
			expectedOutput: loadGolden(t, "moves.golden"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalArgs := os.Args
			os.Args = append([]string{"poke-cli"}, tt.args...)
			defer func() { os.Args = originalArgs }()

			output := MoveCommand()
			cleanOutput := styling.StripANSI(output)

			assert.Equal(t, tt.expectedOutput, cleanOutput, "Output should match expected")
		})
	}
}
