package main

import (
	"bytes"
	"os"
	"regexp"
	"testing"

	"github.com/digitalghost-dev/poke-cli/cmd/utils"
	"github.com/digitalghost-dev/poke-cli/styling"
	"github.com/stretchr/testify/assert"
)

func TestCurrentVersion(t *testing.T) {
	tests := []struct {
		name           string
		version        string
		expectedOutput string
	}{
		{
			name:           "Version set by ldflags",
			version:        "v1.0.2",
			expectedOutput: "Version: v1.0.2\n",
		},
		{
			name:           "Version set to (devel)",
			version:        "(devel)",
			expectedOutput: "Version: (devel)\n",
		},
	}

	// Save the original version and restore it after tests
	originalVersion := version
	defer func() { version = originalVersion }()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			version = tt.version

			r, w, _ := os.Pipe()
			oldStdout := os.Stdout
			os.Stdout = w

			currentVersion()

			// Close the writer and restore stdout
			err := w.Close()
			if err != nil {
				t.Fatalf("Failed to close pipe: %v", err)
			}
			os.Stdout = oldStdout

			// Read the output from the pipe
			var buf bytes.Buffer
			if _, err := buf.ReadFrom(r); err != nil {
				t.Fatalf("Failed to read from pipe: %v", err)
			}

			got := buf.String()
			if got != tt.expectedOutput {
				t.Errorf("Expected %q, got %q", tt.expectedOutput, got)
			}
		})
	}
}

func captureOutput(f func()) string {
	var buf bytes.Buffer
	stdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	_ = w.Close()
	os.Stdout = stdout
	_, _ = buf.ReadFrom(r)

	return buf.String()
}

func stripANSI(input string) string {
	// Regular expression to match ANSI escape sequences
	ansiRegex := regexp.MustCompile(`\x1b\[[0-9;]*m`)
	return ansiRegex.ReplaceAllString(input, "")
}

func TestRunCLI(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedOutput string
		expectedCode   int
	}{
		{
			name:           "Test CLI help flag",
			args:           []string{"--help"},
			expectedOutput: utils.LoadGolden(t, "cli_help.golden"),
			expectedCode:   0,
		},
		{
			name:           "Non-valid command",
			args:           []string{"movesets"},
			expectedOutput: utils.LoadGolden(t, "cli_incorrect_command.golden"),
			expectedCode:   1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalArgs := os.Args
			os.Args = append([]string{"poke-cli"}, tt.args...)
			defer func() { os.Args = originalArgs }()

			var exitCode int
			output := captureOutput(func() {
				exitCode = runCLI(tt.args)
			})

			cleanOutput := styling.StripANSI(output)

			assert.Equal(t, tt.expectedOutput, cleanOutput, "Output should match expected")
			assert.Equal(t, tt.expectedCode, exitCode, "Exit code should match expected")
		})
	}
}

// TODO: finish testing different commands?
func TestRunCLI_VariousCommands(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected int
	}{
		//{"Invalid command", []string{"foobar"}, 1},
		{"Latest flag long", []string{"--latest"}, 0},
		{"Latest flag short", []string{"-l"}, 0},
		{"Version flag long", []string{"--version"}, 0},
		{"Version flag short", []string{"-v"}, 0},
		{"Search command with invalid args", []string{"search", "pokemon", "extra-arg"}, 1},
		//{"Missing Pokémon name", []string{"pokemon"}, 1},
		//{"Another invalid command", []string{"invalid"}, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exitCode := runCLI(tt.args)
			if exitCode != tt.expected {
				t.Errorf("expected %d, got %d for args %v", tt.expected, exitCode, tt.args)
			}
		})
	}
}

func TestMainFunction(t *testing.T) {
	originalExit := exit
	defer func() { exit = originalExit }()

	tests := []struct {
		name           string
		args           []string
		expectedOutput string
		expectedCode   int
	}{
		{
			name:           "Run main command",
			args:           []string{"poke-cli"},
			expectedOutput: "Welcome! This tool displays data related to Pokémon!",
			expectedCode:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exitCode := 0
			exit = func(code int) { exitCode = code }

			output := captureOutput(func() {
				os.Args = tt.args
				main()
			})

			output = stripANSI(output)

			if exitCode != tt.expectedCode {
				t.Errorf("Expected exit code %d, got %d", tt.expectedCode, exitCode)
			}

			if !bytes.Contains([]byte(output), []byte(tt.expectedOutput)) {
				t.Errorf("Expected output to contain %q, got %q", tt.expectedOutput, output)
			}
		})
	}
}
