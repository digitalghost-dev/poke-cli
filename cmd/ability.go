package cmd

import (
	"flag"
	"fmt"
	"github.com/digitalghost-dev/poke-cli/connections"
	"github.com/digitalghost-dev/poke-cli/flags"
	"github.com/digitalghost-dev/poke-cli/styling"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"os"
	"strings"
)

func AbilityCommand() {

	flag.Usage = func() {
		helpMessage := styling.HelpBorder.Render(
			"Get details about a specific ability.\n\n",
			styling.StyleBold.Render("USAGE:"),
			fmt.Sprintf("\n\t%s %s %s %s", "poke-cli", styling.StyleBold.Render("ability"), "<ability-name>", "[flag]"),
			fmt.Sprintf("\n\t%-30s", styling.StyleItalic.Render("Use a hyphen when typing a name with a space.")),
			"\n\n",
			styling.StyleBold.Render("FLAGS:"),
			fmt.Sprintf("\n\t%-30s %s", "-p, --pokemon", "Prints Pokémon that learn this ability."),
			fmt.Sprintf("\n\t%-30s %s", "-h, --help", "Prints the help menu."),
		)
		fmt.Println(helpMessage)
	}

	abilityFlags, pokemonFlag, shortPokemonFlag := flags.SetupAbilityFlagSet()

	args := os.Args

	flag.Parse()

	if err := ValidateAbilityArgs(args); err != nil {
		fmt.Println(err.Error())
		if os.Getenv("GO_TESTING") != "1" {
			os.Exit(1)
		}
	}

	endpoint := strings.ToLower(args[1])
	abilityName := strings.ToLower(args[2])

	if err := abilityFlags.Parse(args[3:]); err != nil {
		fmt.Printf("error parsing flags: %v\n", err)
		abilityFlags.Usage()
		if os.Getenv("GO_TESTING") != "1" {
			os.Exit(1)
		}
	}

	abilitiesStruct, abilityName, err := connections.AbilityApiCall(endpoint, abilityName, "https://pokeapi.co/api/v2/")
	if err != nil {
		fmt.Println(err)
		if os.Getenv("GO_TESTING") != "1" {
			os.Exit(1)
		}
	}

	// Extract English short_effect
	var englishShortEffect string
	for _, entry := range abilitiesStruct.EffectEntries {
		if entry.Language.Name == "en" {
			englishShortEffect = entry.ShortEffect
			break
		}
	}

	// Extract English flavor_text_entries
	var englishFlavorEntry string
	for _, entry := range abilitiesStruct.FlavorEntries {
		if entry.Language.Name == "en" {
			englishFlavorEntry = entry.FlavorText
			break
		}
	}

	capitalizedAbility := cases.Title(language.English).String(strings.Replace(abilityName, "-", " ", -1))
	fmt.Println(styling.StyleBold.Render(capitalizedAbility))

	// API is missing some data for the short_effect for abilities from Generation 9.
	// If short_effect is empty, fallback to the move's flavor_text_entry.
	if englishShortEffect == "" {
		fmt.Println("Effect:", englishFlavorEntry)
	} else {
		fmt.Println("Effect:", englishShortEffect)
	}

	// Print the generation where the move was first introduced.
	generationParts := strings.Split(abilitiesStruct.Generation.Name, "-")
	if len(generationParts) > 1 {
		generationUpper := strings.ToUpper(generationParts[1])
		fmt.Println("Generation:", generationUpper)
	} else {
		fmt.Println("Generation: Unknown")
	}

	if *pokemonFlag || *shortPokemonFlag {
		if err := flags.PokemonFlag(endpoint, abilityName); err != nil {
			fmt.Printf("error parsing flags: %v\n", err)
			os.Exit(1)
		}
	}
}
