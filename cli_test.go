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
			args: []string{"pokemons"},
			expectedOutput: "╭──────────────────────────────────────────────────────╮\n" +
				"│Error!                                                │\n" +
				"│Available Commands:                                   │\n" +
				"│    pokemon         Get details of a specific Pokémon │\n" +
				"│    types           Get details of a specific typing  │\n" +
				"│                                                      │\n" +
				"│Also run [poke-cli -h] for more info!                 │\n" +
				"╰──────────────────────────────────────────────────────╯\n",
			expectedExit: 0,
		},
		{
			args: []string{"pokemon"},
			expectedOutput: "╭────────────────────────────────────────────────────────────╮\n" +
				"│Error!                                                      │\n" +
				"│Please declare a Pokémon's name after the [pokemon] command │\n" +
				"│Run 'poke-cli pokemon -h' for more details                  │\n" +
				"│error: insufficient arguments                               │\n" +
				"╰────────────────────────────────────────────────────────────╯\n",
			expectedExit: 1,
		},
		{
			args:           []string{"pokemon", "bulbasaur"},
			expectedOutput: "Your selected Pokémon: Bulbasaur\nNational Pokédex #: 1\n",
			expectedExit:   0,
		},
		{
			args:           []string{"pokemon", "mew", "--types"},
			expectedOutput: "Your selected Pokémon: Mew\nNational Pokédex #: 151\n──────\nTyping\nType 1: psychic\n",
			expectedExit:   0,
		},
		{
			args:           []string{"pokemon", "cacturne", "--types"},
			expectedOutput: "Your selected Pokémon: Cacturne\nNational Pokédex #: 332\n──────\nTyping\nType 1: grass\nType 2: dark\n",
			expectedExit:   0,
		},
		{
			args: []string{"pokemon", "chimchar", "types"},
			expectedOutput: "╭─────────────────────────────────────────────────────────────────────────────────╮\n" +
				"│Error!                                                                           │\n" +
				"│Invalid argument 'types'. Only flags are allowed after declaring a Pokémon's name│\n" +
				"╰─────────────────────────────────────────────────────────────────────────────────╯\n",
			expectedExit: 1,
		},
		{
			args: []string{"pokemon", "flutter-mane", "types"},
			expectedOutput: "╭─────────────────────────────────────────────────────────────────────────────────╮\n" +
				"│Error!                                                                           │\n" +
				"│Invalid argument 'types'. Only flags are allowed after declaring a Pokémon's name│\n" +
				"╰─────────────────────────────────────────────────────────────────────────────────╯\n",
			expectedExit: 1,
		},
		{
			args: []string{
				"pokemon", "AmPhaROs", "--types", "--abilities",
			},
			expectedOutput: "Your selected Pokémon: Ampharos\n" +
				"National Pokédex #: 181\n" +
				"──────\n" +
				"Typing\n" +
				"Type 1: electric\n" +
				"─────────\n" +
				"Abilities\n" +
				"Ability 1: static\n" +
				"Hidden Ability: plus\n",
			expectedExit: 0,
		},
		{
			args: []string{
				"pokemon", "CLOysTeR", "-t", "-a",
			},
			expectedOutput: "Your selected Pokémon: Cloyster\n" +
				"National Pokédex #: 91\n" +
				"──────\n" +
				"Typing\n" +
				"Type 1: water\n" +
				"Type 2: ice\n" +
				"─────────\n" +
				"Abilities\n" +
				"Ability 1: shell-armor\n" +
				"Ability 2: skill-link\n" +
				"Hidden Ability: overcoat\n",
			expectedExit: 0,
		},
		{
			args: []string{"pokemon", "gyarados", "--help"},
			expectedOutput: "╭──────────────────────────────────────────────────────────────╮\n" +
				"│poke-cli pokemon <pokemon-name> [flags]                       │\n" +
				"│                                                              │\n" +
				"│FLAGS:                                                        │\n" +
				"│     -a, --abilities      Prints out the Pokémon's abilities. │\n" +
				"│     -t, --types          Prints out the Pokémon's typing.    │\n" +
				"╰──────────────────────────────────────────────────────────────╯\n",
			expectedExit: 0,
		},
		{
			args: []string{"--help"},
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
			args: []string{"types", "ground", "all"},
			expectedOutput: "╭──────────────────╮\n" +
				"│Error!            │\n" +
				"│Too many arguments│\n" +
				"╰──────────────────╯\n",
			expectedExit: 1,
		},
		{
			args: []string{"types", "--help"},
			expectedOutput: "╭───────────────────────────────────────────────────────────────╮\n" +
				"│USAGE:                                                         │\n" +
				"│    poke-cli types [flag]                                      │\n" +
				"│    Get details about a specific typing                        │\n" +
				"│    ----------                                                 │\n" +
				"│    Examples:                                                  │\n" +
				"│    poke-cli types                                             │\n" +
				"│    A table will then display with the option to select a type.│\n" +
				"╰───────────────────────────────────────────────────────────────╯\n",
			expectedExit: 0,
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
