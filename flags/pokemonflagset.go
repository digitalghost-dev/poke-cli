// This file holds all the flags used by the <pokemonName> subcommand

package flags

import (
	"flag"
	"fmt"
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

	pokemonStruct := connections.PokemonTypeApiCall(pokemonName, "https://pokeapi.co/api/v2/pokemon/")

	for _, pokeType := range pokemonStruct.Types {
		fmt.Printf("Type %d: %s\n", pokeType.Slot, pokeType.Type.Name)
	}

	return nil
}
