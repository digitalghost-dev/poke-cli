package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"regexp"
	"testing"
)

// Strip ANSI color codes
var ansiRegex = regexp.MustCompile(`\x1b\[[0-9;]*m`)

func stripANSI(input string) string {
	return ansiRegex.ReplaceAllString(input, "")
}

func TestMainFunction(t *testing.T) {
	version := "v0.6.5"

	// Backup the original exit function and stdout/stderr
	originalExit := exit
	originalStdout := os.Stdout
	originalStderr := os.Stderr
	defer func() {
		exit = originalExit        // Restore exit
		os.Stdout = originalStdout // Restore stdout
		os.Stderr = originalStderr // Restore stderr
	}()

	// Replace exit with a function that captures the exit code
	exitCode := 0
	exit = func(code int) { exitCode = code }

	tests := []struct {
		args           []string
		expectedOutput string
		expectedExit   int
	}{
		{
			args: []string{"pokemons"},
			expectedOutput: "╭──────────────────────────────────────────────────────╮\n" +
				"│Error!                                                │\n" +
				"│Available Commands:                                   │\n" +
				"│    pokemon         Get details of a specific Pokémon │\n" +
				"│    types           Get details of a specific typing  │\n" +
				"│                                                      │\n" +
				"│Also run [poke-cli -h] for more info!                 │\n" +
				"╰──────────────────────────────────────────────────────╯\n",
			expectedExit: 1,
		},
		{
			args:           []string{"-l"},
			expectedOutput: fmt.Sprintf("Latest Docker image version: %s\nLatest release tag: %s\n", version, version),
			expectedExit:   0,
		},
		{
			args:           []string{"--latest"},
			expectedOutput: fmt.Sprintf("Latest Docker image version: %s\nLatest release tag: %s\n", version, version),
			expectedExit:   0,
		},
		{
			args: []string{"-h"},
			expectedOutput: "╭──────────────────────────────────────────────────────╮\n" +
				"│Welcome! This tool displays data related to Pokémon!  │\n" +
				"│                                                      │\n" +
				"│ USAGE:                                               │\n" +
				"│    poke-cli [flag]                                   │\n" +
				"│    poke-cli [command] [flag]                         │\n" +
				"│    poke-cli [command] [subcommand] [flag]            │\n" +
				"│                                                      │\n" +
				"│ FLAGS:                                               │\n" +
				"│    -h, --help      Shows the help menu               │\n" +
				"│    -l, --latest    Prints the latest available       │\n" +
				"│                    version of the program            │\n" +
				"│                                                      │\n" +
				"│ AVAILABLE COMMANDS:                                  │\n" +
				"│    pokemon         Get details of a specific Pokémon │\n" +
				"│    types           Get details of a specific typing  │\n" +
				"╰──────────────────────────────────────────────────────╯\n",
			expectedExit: 0,
		},
		{
			args:           []string{"pokemon", "kingambit"},
			expectedOutput: "Your selected Pokémon: Kingambit\nNational Pokédex #: 983\n",
			expectedExit:   0,
		},
		{
			args:           []string{"pokemon", "cradily", "--types"},
			expectedOutput: "Your selected Pokémon: Cradily\nNational Pokédex #: 346\n──────\nTyping\nType 1: rock\nType 2: grass\n",
			expectedExit:   0,
		},
		{
			args:           []string{"pokemon", "giratina-altered", "--abilities"},
			expectedOutput: "Your selected Pokémon: Giratina-Altered\nNational Pokédex #: 487\n─────────\nAbilities\nAbility 1: pressure\nHidden Ability: telepathy\n",
			expectedExit:   0,
		},
		{
			args:           []string{"pokemon", "coPPeraJAH", "-t", "-a"},
			expectedOutput: "Your selected Pokémon: Copperajah\nNational Pokédex #: 879\n──────\nTyping\nType 1: steel\n─────────\nAbilities\nAbility 1: sheer-force\nHidden Ability: heavy-metal\n",
			expectedExit:   0,
		},
	}

	for _, test := range tests {
		// Create a pipe to capture output
		r, w, _ := os.Pipe()
		os.Stdout = w
		os.Stderr = w

		// Set os.Args for the test case
		os.Args = append([]string{"poke-cli"}, test.args...)

		// Run the main function
		main()

		// Close the writer and restore stdout and stderr
		err := w.Close()
		if err != nil {
			t.Fatalf("Error closing pipe writer: %v", err)
		}
		os.Stdout = originalStdout
		os.Stderr = originalStderr

		// Read from the pipe
		var buf bytes.Buffer
		if _, err := io.Copy(&buf, r); err != nil {
			t.Errorf("Error copying output: %v", err)
		}

		// Strip ANSI color codes from the actual output
		actualOutput := stripANSI(buf.String())
		if actualOutput != test.expectedOutput {
			t.Errorf("Args: %v\nExpected output: %q\nGot: %q\n", test.args, test.expectedOutput, actualOutput)
		}

		if exitCode != test.expectedExit {
			t.Errorf("Args: %v\nExpected exit code: %d\nGot: %d\n", test.args, test.expectedExit, exitCode)
		}
	}
}
