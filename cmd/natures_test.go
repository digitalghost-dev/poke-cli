package cmd

import (
	"bytes"
	"github.com/digitalghost-dev/poke-cli/styling"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var exitCode int

func fakeExit(code int) {
	exitCode = code
	panic("exit")
}

func captureNaturesOutput(f func()) string {
	r, w, _ := os.Pipe()
	defer func() {
		_ = r.Close()
	}()

	oldStdout := os.Stdout
	os.Stdout = w
	defer func() { os.Stdout = oldStdout }()

	f()

	_ = w.Close()

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
			name:           "Invalid extra argument",
			args:           []string{"natures", "extra"},
			expectedOutput: styling.StripANSI(styling.ErrorBorder.Render(styling.ErrorColor.Render("Error!")+"\nThe only currently available options\nafter <natures> command are '-h' or '--help'")) + "\n",
			expectedError:  true,
		},
		{
			name: "Full Natures output with table",
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
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Override osExit
			oldExit := osExit
			osExit = fakeExit
			defer func() { osExit = oldExit }()

			// Reset captured exit code
			exitCode = 0

			// Save original os.Args
			originalArgs := os.Args
			defer func() { os.Args = originalArgs }()
			os.Args = append([]string{"poke-cli"}, tt.args...)

			// Capture output
			output := captureNaturesOutput(func() {
				defer func() {
					if r := recover(); r != nil {
						if r != "exit" {
							t.Fatalf("Unexpected panic: %v", r)
						}
					}
				}()
				NaturesCommand()
			})

			cleanOutput := styling.StripANSI(output)

			// Logging expected vs actual
			t.Logf("Expected Output:\n%s", tt.expectedOutput)
			t.Logf("Actual Output:\n%s", cleanOutput)

			// Assertions
			assert.Equal(t, tt.expectedOutput, cleanOutput, "Output should match expected")
			if tt.expectedError {
				assert.Equal(t, 1, exitCode, "Expected exit code 1 on error")
			} else {
				assert.Equal(t, 0, exitCode, "Expected no exit (code 0) on success")
			}
		})
	}
}
