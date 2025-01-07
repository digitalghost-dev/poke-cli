package cmd

import (
	"flag"
	"fmt"
	"github.com/digitalghost-dev/poke-cli/connections"
	"github.com/digitalghost-dev/poke-cli/flags"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"math"
	"os"
	"strings"
)

// PokemonCommand processes the Pokémon command
func PokemonCommand() {

	hintMessage := styleItalic.Render("options: [sm, md, lg]")

	flag.Usage = func() {
		helpMessage := helpBorder.Render(
			"Get details about a specific Pokémon.\n\n",
			styleBold.Render("USAGE:"),
			fmt.Sprintf("\n\t%s %s %s %s", "poke-cli", styleBold.Render("pokemon"), "<pokemon-name>", "[flag]"),
			fmt.Sprintf("\n\t%-30s", styleItalic.Render("Use a hyphen when typing a name with a space.")),
			"\n\n",
			styleBold.Render("FLAGS:"),
			fmt.Sprintf("\n\t%-30s %s", "-a, --abilities", "Prints the Pokémon's abilities."),
			fmt.Sprintf("\n\t%-30s %s", "-i=xx, --image=xx", "Prints out the Pokémon's default sprite."),
			fmt.Sprintf("\n\t%5s%-15s", "", hintMessage),
			fmt.Sprintf("\n\t%-30s %s", "-s, --stats", "Prints the Pokémon's base stats."),
			fmt.Sprintf("\n\t%-30s %s", "-t, --types", "Prints the Pokémon's typing."),
			fmt.Sprintf("\n\t%-30s %s", "-h, --help", "Prints the help menu."),
		)
		fmt.Println(helpMessage)
	}

	flag.Parse()

	pokeFlags, abilitiesFlag, shortAbilitiesFlag, imageFlag, shortImageFlag, statsFlag, shortStatsFlag, typesFlag, shortTypesFlag := flags.SetupPokemonFlagSet()

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

	_, pokemonName, pokemonID, pokemonWeight, pokemonHeight := connections.PokemonApiCall(endpoint, pokemonName, "https://pokeapi.co/api/v2/")
	capitalizedString := cases.Title(language.English).String(strings.Replace(pokemonName, "-", " ", -1))

	// Weight calculation
	weightKilograms := float64(pokemonWeight) / 10
	weightPounds := float64(weightKilograms) * 2.20462

	// Height calculation
	heightMeters := float64(pokemonHeight) / 10
	heightFeet := heightMeters * 3.28084
	feet := int(heightFeet)
	inches := int(math.Round((heightFeet - float64(feet)) * 12)) // Use math.Round to avoid truncation

	// Adjust for rounding to 12 inches (carry over to the next foot)
	if inches == 12 {
		feet++
		inches = 0
	}

	fmt.Printf(
		"Your selected Pokémon: %s\nNational Pokédex #: %d\nWeight: %.1fkg (%.1f lbs)\nHeight: %.1fm (%d′%02d″)\n",
		capitalizedString, pokemonID, weightKilograms, weightPounds, heightFeet, feet, inches,
	)

	if *imageFlag != "" || *shortImageFlag != "" {
		// Determine the size based on the provided flags
		size := *imageFlag
		if *shortImageFlag != "" {
			size = *shortImageFlag
		}

		// Call the ImageFlag function with the specified size
		if err := flags.ImageFlag(endpoint, pokemonName, size); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	if *abilitiesFlag || *shortAbilitiesFlag {
		if err := flags.AbilitiesFlag(endpoint, pokemonName); err != nil {
			fmt.Printf("Error: %s\n", err)
			os.Exit(1)
		}
	}

	if *typesFlag || *shortTypesFlag {
		if err := flags.TypesFlag(endpoint, pokemonName); err != nil {
			fmt.Printf("Error: %s\n", err)
			os.Exit(1)
		}
	}

	if *statsFlag || *shortStatsFlag {
		if err := flags.StatsFlag(endpoint, pokemonName); err != nil {
			fmt.Printf("Error: %s\n", err)
			os.Exit(1)
		}
	}
}
