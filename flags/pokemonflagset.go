// This file holds all the flags used by the <pokemonName> subcommand

package flags

import (
	"flag"
	"github.com/digitalghost-dev/poke-cli/connections"
	"os"
)

func SetupPokemonFlagSet() (*flag.FlagSet, *bool) {
	pokeFlags := flag.NewFlagSet("pokeFlags", flag.ExitOnError)

	typesFlag := pokeFlags.Bool("types", false, "Print the declared Pokémon's typing")

	return pokeFlags, typesFlag
}

func TypesFlag() error {
	pokemonName := os.Args[1]

	connections.TypeApiCall(pokemonName, "https://pokeapi.co/api/v2/pokemon/")

	return nil
}
