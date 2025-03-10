package main

import (
	"bytes"
	"os"
	"regexp"
	"testing"
)

func TestCurrentVersion(t *testing.T) {
	tests := []struct {
		name           string
		version        string
		expectedOutput string
	}{
		{
			name:           "Version set by ldflags",
			version:        "v1.0.0",
			expectedOutput: "Version: v1.0.0\n",
		},
		{
			name:           "Version set to (devel)",
			version:        "(devel)",
			expectedOutput: "Version: (devel)\n", // Simplified assumption
		},
	}

	// Save the original version and restore it after tests
	originalVersion := version
	defer func() { version = originalVersion }()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set the version for this test case
			version = tt.version

			// Redirect stdout to capture the output
			r, w, _ := os.Pipe()
			oldStdout := os.Stdout
			os.Stdout = w

			// Call the function
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

			// Compare the output with the expected result
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
			name: "No Arguments",
			args: []string{},
			expectedOutput: "╭──────────────────────────────────────────────────────────╮\n" +
				"│Welcome! This tool displays data related to Pokémon!      │\n" +
				"│                                                          │\n" +
				"│ USAGE:                                                   │\n" +
				"│    poke-cli [flag]                                       │\n" +
				"│    poke-cli <command> [flag]                             │\n" +
				"│    poke-cli <command> <subcommand> [flag]                │\n" +
				"│                                                          │\n" +
				"│ FLAGS:                                                   │\n" +
				"│    -h, --help      Shows the help menu                   │\n" +
				"│    -l, --latest    Prints the latest version available   │\n" +
				"│    -v, --version   Prints the current version            │\n" +
				"│                                                          │\n" +
				"│ COMMANDS:                                                │\n" +
				"│    ability         Get details about an ability          │\n" +
				"│    natures         Get details about all natures         │\n" +
				"│    pokemon         Get details about a Pokémon           │\n" +
				"│    types           Get details about a typing            │\n" +
				"│                                                          │\n" +
				"│ hint: when calling a resource with a space, use a hyphen │\n" +
				"│ example: poke-cli ability strong-jaw                     │\n" +
				"│ example: poke-cli pokemon flutter-mane -s                │\n" +
				"╰──────────────────────────────────────────────────────────╯",
			expectedCode: 0,
		},
		{
			name: "Help Flag Short",
			args: []string{"-h"},
			expectedOutput: "╭──────────────────────────────────────────────────────────╮\n" +
				"│Welcome! This tool displays data related to Pokémon!      │\n" +
				"│                                                          │\n" +
				"│ USAGE:                                                   │\n" +
				"│    poke-cli [flag]                                       │\n" +
				"│    poke-cli <command> [flag]                             │\n" +
				"│    poke-cli <command> <subcommand> [flag]                │\n" +
				"│                                                          │\n" +
				"│ FLAGS:                                                   │\n" +
				"│    -h, --help      Shows the help menu                   │\n" +
				"│    -l, --latest    Prints the latest version available   │\n" +
				"│    -v, --version   Prints the current version            │\n" +
				"│                                                          │\n" +
				"│ COMMANDS:                                                │\n" +
				"│    ability         Get details about an ability          │\n" +
				"│    natures         Get details about all natures         │\n" +
				"│    pokemon         Get details about a Pokémon           │\n" +
				"│    types           Get details about a typing            │\n" +
				"│                                                          │\n" +
				"│ hint: when calling a resource with a space, use a hyphen │\n" +
				"│ example: poke-cli ability strong-jaw                     │\n" +
				"│ example: poke-cli pokemon flutter-mane -s                │\n" +
				"╰──────────────────────────────────────────────────────────╯",
			expectedCode: 0,
		},
		{
			name: "Help Flag Long",
			args: []string{"--help"},
			expectedOutput: "╭──────────────────────────────────────────────────────────╮\n" +
				"│Welcome! This tool displays data related to Pokémon!      │\n" +
				"│                                                          │\n" +
				"│ USAGE:                                                   │\n" +
				"│    poke-cli [flag]                                       │\n" +
				"│    poke-cli <command> [flag]                             │\n" +
				"│    poke-cli <command> <subcommand> [flag]                │\n" +
				"│                                                          │\n" +
				"│ FLAGS:                                                   │\n" +
				"│    -h, --help      Shows the help menu                   │\n" +
				"│    -l, --latest    Prints the latest version available   │\n" +
				"│    -v, --version   Prints the current version            │\n" +
				"│                                                          │\n" +
				"│ COMMANDS:                                                │\n" +
				"│    ability         Get details about an ability          │\n" +
				"│    natures         Get details about all natures         │\n" +
				"│    pokemon         Get details about a Pokémon           │\n" +
				"│    types           Get details about a typing            │\n" +
				"│                                                          │\n" +
				"│ hint: when calling a resource with a space, use a hyphen │\n" +
				"│ example: poke-cli ability strong-jaw                     │\n" +
				"│ example: poke-cli pokemon flutter-mane -s                │\n" +
				"╰──────────────────────────────────────────────────────────╯",

			expectedCode: 0,
		},
		{
			name:           "Invalid Command",
			args:           []string{"invalid"},
			expectedOutput: "Error!",
			expectedCode:   1,
		},
		{
			name:           "Latest Flag",
			args:           []string{"-l"},
			expectedOutput: "Latest Docker image version: v1.0.0\nLatest release tag: v1.0.0\n",
			expectedCode:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exit = func(code int) {}
			output := captureOutput(func() {
				code := runCLI(tt.args)
				if code != tt.expectedCode {
					t.Errorf("Expected exit code %d, got %d", tt.expectedCode, code)
				}
			})

			output = stripANSI(output)

			if !bytes.Contains([]byte(output), []byte(tt.expectedOutput)) {
				t.Errorf("Expected output to contain %q, got %q", tt.expectedOutput, output)
			}
		})
	}
}

func TestMainFunction(t *testing.T) {
	originalExit := exit
	defer func() { exit = originalExit }() // Restore original exit after test

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
