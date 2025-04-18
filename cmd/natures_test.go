package cmd

import (
	"bytes"
	"fmt"
	"github.com/digitalghost-dev/poke-cli/styling"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func captureNaturesOutput(f func()) string {
	r, w, _ := os.Pipe()
	defer func(r *os.File) {
		err := r.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(r)

	oldStdout := os.Stdout
	os.Stdout = w
	defer func() { os.Stdout = oldStdout }()

	f()

	err := w.Close()
	if err != nil {
		return ""
	}

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

			assert.Equal(t, tt.expectedOutput, cleanOutput, "Output should equal the expected string")
		})
	}
}
