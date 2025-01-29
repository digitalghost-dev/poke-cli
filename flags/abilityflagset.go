package flags

import (
	"flag"
	"fmt"
	"github.com/digitalghost-dev/poke-cli/connections"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func SetupAbilityFlagSet() (*flag.FlagSet, *bool, *bool) {
	abilityFlags := flag.NewFlagSet("AbilityFlagSet", flag.ExitOnError)

	pokemonFlag := abilityFlags.Bool("pokemon", false, "List all Pokémon with chosen ability")
	shortPokemonFlag := abilityFlags.Bool("p", false, "List all Pokémon with chosen ability")

	abilityFlags.Usage = func() {
		helpMessage := helpBorder.Render("poke-cli pokemon <pokemon-name> [flags]\n\n",
			styleBold.Render("FLAGS:"),
			fmt.Sprintf("\n\t%-30s %s", "-p, --pokemon", "List all Pokémon with chosen ability."),
		)
		fmt.Println(helpMessage)
	}

	return abilityFlags, pokemonFlag, shortPokemonFlag
}

func PokemonFlag(endpoint string, abilityName string) error {
	abilitiesStruct, _, _ := connections.AbilityApiCall(endpoint, abilityName, "https://pokeapi.co/api/v2/")

	capitalizedEffect := cases.Title(language.English).String(abilityName)

	fmt.Printf("\nPokémon with %s\n\n", capitalizedEffect)

	// Extract Pokémon names and capitalize them
	var pokemonNames []string
	for _, pokemon := range abilitiesStruct.Pokemon {
		pokemonNames = append(pokemonNames, cases.Title(language.English).String(pokemon.PokemonName.Name))
	}

	// Print names in a grid format
	const cols = 4
	maxWidth := 26

	for i, name := range pokemonNames {
		entry := fmt.Sprintf("%2d. %-*s", i+1, maxWidth-5, name) // Numbered entry with padding
		fmt.Print(entry)
		if (i+1)%cols == 0 {
			fmt.Println() // New line after every `cols` entries
		}
	}
	fmt.Println()

	return nil
}
