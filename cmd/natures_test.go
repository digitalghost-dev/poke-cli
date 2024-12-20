package cmd

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"
)

func TestNaturesCommand(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedOutput string
		expectError    bool
	}{
		{
			name: "Help flag",
			args: []string{"natures", "-h"},
			expectedOutput: "╭────────────────────────────────────────────────────────────╮\n" +
				"│Get details about Pokémon natures.                          │\n" +
				"│                                                            │\n" +
				"│ USAGE:                                                     │\n" +
				"│    poke-cli natures [flag]                                 │\n" +
				"│                                                            │\n" +
				"│ FLAGS:                                                     │\n" +
				"│    -h, --help                     Prints out the help menu.│\n" +
				"╰────────────────────────────────────────────────────────────╯\n",
			expectError: false,
		},
		{
			name: "Valid Execution",
			args: []string{"natures"},
			expectedOutput: "Natures affect the growth of a Pokémon.\n" +
				"Each nature increases one of its stats by 10% and decreases one by 10%.\n" +
				"Five natures increase and decrease the same stat and therefore have no effect.\n\n" +
				"Nature Chart:\n" +
				"┌──────────┬─────────┬──────────┬──────────┬──────────┬─────────┐\n" +
				"│          │ -Attack │ -Defense │ -Sp. Atk │ -Sp. Def │ Speed   │\n" +
				"├──────────┼─────────┼──────────┼──────────┼──────────┼─────────┤\n" +
				"│ +Attack  │ Hardy   │ Lonely   │ Adamant  │ Naughty  │ Brave   │\n" +
				"├──────────┼─────────┼──────────┼──────────┼──────────┼─────────┤\n" +
				"│ +Defense │ Bold    │ Docile   │ Impish   │ Lax      │ Relaxed │\n" +
				"├──────────┼─────────┼──────────┼──────────┼──────────┼─────────┤\n" +
				"│ +Sp. Atk │ Modest  │ Mild     │ Bashful  │ Rash     │ Quiet   │\n" +
				"├──────────┼─────────┼──────────┼──────────┼──────────┼─────────┤\n" +
				"│ +Sp. Def │ Calm    │ Gentle   │ Careful  │ Quirky   │ Sassy   │\n" +
				"├──────────┼─────────┼──────────┼──────────┼──────────┼─────────┤\n" +
				"│ Speed    │ Timid   │ Hasty    │ Jolly    │ Naive    │ Serious │\n" +
				"└──────────┴─────────┴──────────┴──────────┴──────────┴─────────┘\n",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command("poke-cli", tt.args...)
			var out bytes.Buffer
			cmd.Stdout = &out
			cmd.Stderr = &out

			err := cmd.Run()

			if tt.expectError {
				if err == nil {
					t.Fatalf("Expected an error but got none.\nOutput: %s", out.String())
				}
			} else {
				if err != nil {
					t.Fatalf("Did not expect an error but got: %v\nOutput: %s", err, out.String())
				}
			}

			output := out.String()
			if !strings.Contains(output, tt.expectedOutput) {
				t.Errorf("Output mismatch.\nExpected to contain:\n%s\nGot:\n%s", tt.expectedOutput, output)
			}
		})
	}
}
