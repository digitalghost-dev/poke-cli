package cmd

import (
	"flag"
	"fmt"
	"github.com/digitalghost-dev/poke-cli/connections"
	"github.com/digitalghost-dev/poke-cli/flags"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"os"
	"strings"
)

// PokemonCommand processes the Pokémon command
func PokemonCommand() {

	flag.Usage = func() {
		helpMessage := helpBorder.Render(
			"Get details about a specific Pokémon.\n\n",
			styleBold.Render("USAGE:"),
			fmt.Sprintf("\n\t%s %s %s %s", "poke-cli", styleBold.Render("pokemon"), "<pokemon-name>", "[flag]"),
			fmt.Sprintf("\n\t%-30s", styleItalic.Render("Use a hyphen when typing a name with a space.")),
			"\n\n",
			styleBold.Render("FLAGS:"),
			fmt.Sprintf("\n\t%-30s %s", "-a, --abilities", "Prints out the Pokémon's abilities."),
			fmt.Sprintf("\n\t%-30s %s", "-t, --types", "Prints out the Pokémon's typing."),
			fmt.Sprintf("\n\t%-30s %s", "-h, --help", "Prints out the help menu."),
		)
		fmt.Println(helpMessage)
	}

	flag.Parse()

	pokeFlags, typesFlag, shortTypesFlag, abilitiesFlag, shortAbilitiesFlag := flags.SetupPokemonFlagSet()

	args := os.Args

	err := ValidatePokemonArgs(args)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	endpoint := strings.ToLower(args[1])
	pokemonName := strings.ToLower(args[2])

	if err := pokeFlags.Parse(args[3:]); err != nil {
		fmt.Printf("error parsing flags: %v\n", err)
		pokeFlags.Usage()
		os.Exit(1)
	}

	_, pokemonName, pokemonID := connections.PokemonApiCall(endpoint, pokemonName, "https://pokeapi.co/api/v2/")
	capitalizedString := cases.Title(language.English).String(pokemonName)

	fmt.Printf("Your selected Pokémon: %s\nNational Pokédex #: %d\n", capitalizedString, pokemonID)

	if *typesFlag || *shortTypesFlag {
		if err := flags.TypesFlag(endpoint, pokemonName); err != nil {
			fmt.Printf("Error: %s\n", err)
			os.Exit(1)
		}
	}

	if *abilitiesFlag || *shortAbilitiesFlag {
		if err := flags.AbilitiesFlag(endpoint, pokemonName); err != nil {
			fmt.Printf("Error: %s\n", err)
			os.Exit(1)
		}
	}
}
