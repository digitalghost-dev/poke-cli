package main

import (
	"bytes"
	"os/exec"
	"testing"
)

func TestCLI(t *testing.T) {

	tests := []struct {
		args           []string
		expectedOutput string
		expectedExit   int
	}{
		{
			args:           []string{},
			expectedOutput: "Please declare a Pokémon's name after the CLI name\nRun 'poke-cli --help' for more details\n",
			expectedExit:   1,
		},
		{
			args:           []string{"bulbasaur"},
			expectedOutput: "Selected Pokémon: Bulbasaur\n",
			expectedExit:   0,
		},
		{
			args:           []string{"mew", "--types"},
			expectedOutput: "Selected Pokémon: Mew\nType 1: psychic\n",
			expectedExit:   0,
		},
		{
			args:           []string{"cacturne", "--types"},
			expectedOutput: "Selected Pokémon: Cacturne\nType 1: grass\nType 2: dark\n",
			expectedExit:   0,
		},
		{
			args:           []string{"chimchar", "types"},
			expectedOutput: "error: only flags are allowed after declaring a Pokémon's name\n",
			expectedExit:   1,
		},
		{
			args:           []string{"flutter-mane", "types"},
			expectedOutput: "Selected Pokémon: Flutter-Mane\nType 1: ghost\nType 2: fairy\n",
			expectedExit:   0,
		},
	}

	for _, test := range tests {
		cmd := exec.Command("poke-cli", test.args...)
		var out bytes.Buffer
		cmd.Stdout = &out

		err := cmd.Run()
		if err != nil {
			return
		}

		if out.String() != test.expectedOutput {
			t.Errorf("Expected output: %s, Got: %s", test.expectedOutput, out.String())
		}

		if cmd.ProcessState.ExitCode() != test.expectedExit {
			t.Errorf("Expected exit code: %d, Got: %d", test.expectedExit, cmd.ProcessState.ExitCode())
		}
	}
}
