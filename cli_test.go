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
			expectedOutput: "Please declare a Pokémon's name after the CLI name\nRun 'poke-cli --help' for more details\nerror: insufficient arguments\n",
			expectedExit:   1,
		},
		{
			args:           []string{"bulbasaur"},
			expectedOutput: "Your selected Pokémon: Bulbasaur\nNational Pokédex #: 1\n",
			expectedExit:   0,
		},
		{
			args:           []string{"mew", "--types"},
			expectedOutput: "Your selected Pokémon: Mew\nNational Pokédex #: 151\n──────\nTyping\nType 1: psychic\n",
			expectedExit:   0,
		},
		{
			args:           []string{"cacturne", "--types"},
			expectedOutput: "Your selected Pokémon: Cacturne\nNational Pokédex #: 332\n──────\nTyping\nType 1: grass\nType 2: dark\n",
			expectedExit:   0,
		},
		{
			args:           []string{"chimchar", "types"},
			expectedOutput: "Error: Invalid argument 'types'. Only flags are allowed after declaring a Pokémon's name\n",
			expectedExit:   1,
		},
		{
			args:           []string{"flutter-mane", "types"},
			expectedOutput: "Error: Invalid argument 'types'. Only flags are allowed after declaring a Pokémon's name\n",
			expectedExit:   1,
		},
		{
			args:           []string{"AmPhaROs", "--types", "--abilities"},
			expectedOutput: "Your selected Pokémon: Ampharos\nNational Pokédex #: 181\n──────\nTyping\nType 1: electric\n─────────\nAbilities\nAbility 1: static\nAbility 3: plus\n",
			expectedExit:   0,
		},
	}

	for _, test := range tests {
		cmd := exec.Command("poke-cli", test.args...)
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &out

		err := cmd.Run()

		if err != nil {
			// If there's an error, but we expected a successful exit
			if test.expectedExit == 0 {
				t.Errorf("Unexpected error: %v", err)
			}
		}

		if out.String() != test.expectedOutput {
			t.Errorf("Args: %v, Expected output: %q, Got: %q", test.args, test.expectedOutput, out.String())
		}

		if cmd.ProcessState.ExitCode() != test.expectedExit {
			t.Errorf("Args: %v, Expected exit code: %d, Got: %d", test.args, test.expectedExit, cmd.ProcessState.ExitCode())
		}
	}
}
