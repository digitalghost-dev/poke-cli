package types

import (
	"github.com/digitalghost-dev/poke-cli/cmd/utils"
	"github.com/digitalghost-dev/poke-cli/styling"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestTypesCommand(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedOutput string
		wantError      bool
	}{
		{
			name:           "Types help flag",
			args:           []string{"types", "--help"},
			expectedOutput: utils.LoadGolden(t, "types_help.golden"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalArgs := os.Args
			os.Args = append([]string{"poke-cli"}, tt.args...)
			defer func() { os.Args = originalArgs }()

			output := TypesCommand()
			cleanOutput := styling.StripANSI(output)

			assert.Equal(t, tt.expectedOutput, cleanOutput, "Output should match expected")
		})
	}
}
func TestModelInit(t *testing.T) {
	m := model{}
	cmd := m.Init()
	assert.Nil(t, cmd, "Init() should return nil")
}
