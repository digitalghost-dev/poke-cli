package flags

import (
	"flag"
	"fmt"
	"github.com/digitalghost-dev/poke-cli/connections"
	"github.com/digitalghost-dev/poke-cli/styling"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"io"
	"strings"
)

func SetupAbilityFlagSet() (*flag.FlagSet, *bool, *bool) {
	abilityFlags := flag.NewFlagSet("AbilityFlagSet", flag.ExitOnError)

	pokemonFlag := abilityFlags.Bool("pokemon", false, "List all Pokémon with chosen ability")
	shortPokemonFlag := abilityFlags.Bool("p", false, "List all Pokémon with chosen ability")

	abilityFlags.Usage = func() {
		helpMessage := styling.HelpBorder.Render("poke-cli ability <ability-name> [flags]\n\n",
			styling.StyleBold.Render("FLAGS:"),
			fmt.Sprintf("\n\t%-30s %s", "-p, --pokemon", "List all Pokémon with chosen ability."),
		)
		fmt.Println(helpMessage)
	}

	return abilityFlags, pokemonFlag, shortPokemonFlag
}

func PokemonAbilitiesFlag(w io.Writer, endpoint string, abilityName string) error {
	abilitiesStruct, _, _ := connections.AbilityApiCall(endpoint, abilityName, connections.APIURL)

	capitalizedEffect := cases.Title(language.English).String(strings.ReplaceAll(abilityName, "-", " "))

	if _, err := fmt.Fprintf(w, "\n%s\n\n", styling.StyleUnderline.Render("Pokemon with "+capitalizedEffect)); err != nil {
		return err
	}

	// Extract Pokémon names and capitalize them
	var pokemonNames []string
	for _, pokemon := range abilitiesStruct.Pokemon {
		pokemonNames = append(pokemonNames, cases.Title(language.English).String(pokemon.PokemonName.Name))
	}

	// Print names in a grid format
	const cols = 3
	for i, name := range pokemonNames {
		entry := fmt.Sprintf("%2d. %-30s", i+1, name)
		_, err := fmt.Fprint(w, entry)
		if err != nil {
			return err
		}
		if (i+1)%cols == 0 {
			_, err := fmt.Fprintln(w)
			if err != nil {
				return err
			}
		}
	}
	if _, err := fmt.Fprint(w); err != nil {
		return err
	}

	return nil
}
