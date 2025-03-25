package cmd

import (
	"bytes"
	"fmt"
	"github.com/digitalghost-dev/poke-cli/styling"
	"os"
	"strings"
	"testing"
)

func captureNaturesOutput(f func()) string {
	// Create a pipe to capture standard output
	r, w, _ := os.Pipe()
	defer func(r *os.File) {
		err := r.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(r)

	// Redirect os.Stdout to the write end of the pipe
	oldStdout := os.Stdout
	os.Stdout = w
	defer func() { os.Stdout = oldStdout }()

	// Run the function
	f()

	// Close the write end of the pipe
	err := w.Close()
	if err != nil {
		return ""
	}

	// Read the captured output
	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	return buf.String()
}

func TestNaturesCommand(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedOutput string
		expectedError  bool
	}{
		{
			name: "Help flag",
			args: []string{"natures", "-h"},
			expectedOutput: styling.StripANSI(
				"╭──────────────────────────────╮\n" +
					"│Get details about all natures.│\n" +
					"│                              │\n" +
					"│ USAGE:                       │\n" +
					"│    poke-cli natures          │\n" +
					"╰──────────────────────────────╯\n"),
			expectedError: false,
		},
		{
			name: "Valid Execution",
			args: []string{"natures"},
			expectedOutput: styling.StripANSI(
				"Natures affect the growth of a Pokémon.\n" +
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
					"└──────────┴─────────┴──────────┴──────────┴──────────┴─────────┘\n"),
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original os.Args
			originalArgs := os.Args
			defer func() { os.Args = originalArgs }()

			// Set os.Args for the test
			os.Args = append([]string{"poke-cli"}, tt.args...)

			// Capture the output
			output := captureNaturesOutput(func() {
				defer func() {
					// Recover from os.Exit calls
					if r := recover(); r != nil {
						if !tt.expectedError {
							t.Fatalf("Unexpected error: %v", r)
						}
					}
				}()
				NaturesCommand()
			})

			cleanOutput := styling.StripANSI(output)

			// Check output
			if !strings.Contains(cleanOutput, tt.expectedOutput) {
				t.Errorf("Output mismatch.\nExpected to contain:\n%s\nGot:\n%s", tt.expectedOutput, output)
			}
		})
	}
}
