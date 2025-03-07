package cmd

import (
	"bytes"
	"fmt"
	"github.com/digitalghost-dev/poke-cli/styling"
	"log"
	"os"
	"strings"
	"testing"
)

func captureAbilityOutput(f func()) string {
	r, w, _ := os.Pipe()
	defer func(r *os.File) {
		err := r.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(r)

	oldStdout := os.Stdout
	os.Stdout = w
	defer func() { os.Stdout = oldStdout }()

	f()

	err := w.Close()
	if err != nil {
		log.Fatal(err)
	}

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	return buf.String()
}

func TestAbilityCommand(t *testing.T) {
	err := os.Setenv("GO_TESTING", "1")
	if err != nil {
		return
	}
	defer func() {
		err := os.Unsetenv("GO_TESTING")
		if err != nil {
			fmt.Println(err)
		}
	}()

	tests := []struct {
		name           string
		args           []string
		expectedOutput string
		expectError    bool
	}{
		{
			name:           "Help flag",
			args:           []string{"ability", "-h"},
			expectedOutput: "Get details about a specific ability.",
			expectError:    false,
		},
		{
			name:           "Valid Execution",
			args:           []string{"ability", "stench"},
			expectedOutput: styling.StripANSI("Stench\nEffect: Has a 10% chance of making target Pokémon flinch with each hit."),
			expectError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalArgs := os.Args
			defer func() { os.Args = originalArgs }()

			os.Args = append([]string{"poke-cli"}, tt.args...)

			output := captureAbilityOutput(func() {
				defer func() {
					if r := recover(); r != nil {
						if !tt.expectError {
							t.Fatalf("Unexpected error: %v", r)
						}
					}
				}()
				AbilityCommand()
			})

			if !strings.Contains(output, tt.expectedOutput) {
				t.Errorf("Output mismatch. Expected to contain:\n%s\nGot:\n%s", tt.expectedOutput, output)
			}
		})
	}
}
