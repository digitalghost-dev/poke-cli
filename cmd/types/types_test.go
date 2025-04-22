package types

import (
	"bytes"
	"fmt"
	"github.com/digitalghost-dev/poke-cli/styling"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func captureTypesOutput(f func()) string {
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

func TestTypesCommand(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedOutput string
		expectedError  bool
	}{
		{
			name: "Help flag",
			args: []string{"types", "-h"},
			expectedOutput: styling.StripANSI(
				"╭───────────────────────────────────────────────────────────╮\n" +
					"│Get details about a specific typing.                       │\n" +
					"│                                                           │\n" +
					"│ USAGE:                                                    │\n" +
					"│    poke-cli types [flag]                                  │\n" +
					"│                                                           │\n" +
					"│ FLAGS:                                                    │\n" +
					"│    -h, --help                     Prints out the help menu│\n" +
					"╰───────────────────────────────────────────────────────────╯\n"),
			expectedError: false,
		},
		{
			name: "Help flag",
			args: []string{"types", "--help"},
			expectedOutput: styling.StripANSI(
				"╭───────────────────────────────────────────────────────────╮\n" +
					"│Get details about a specific typing.                       │\n" +
					"│                                                           │\n" +
					"│ USAGE:                                                    │\n" +
					"│    poke-cli types [flag]                                  │\n" +
					"│                                                           │\n" +
					"│ FLAGS:                                                    │\n" +
					"│    -h, --help                     Prints out the help menu│\n" +
					"╰───────────────────────────────────────────────────────────╯\n"),
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalArgs := os.Args
			defer func() { os.Args = originalArgs }()

			os.Args = append([]string{"poke-cli"}, tt.args...)

			output := captureTypesOutput(func() {
				defer func() {
					if r := recover(); r != nil && !tt.expectedError {
						t.Fatalf("Unexpected error: %v", r)
					}
				}()
				TypesCommand()
			})

			cleanOutput := styling.StripANSI(output)

			assert.Equal(t, tt.expectedOutput, cleanOutput, "Output should contain the expected string")
		})
	}
}

func TestModelInit(t *testing.T) {
	m := model{}
	cmd := m.Init()
	assert.Nil(t, cmd, "Init() should return nil")
}
