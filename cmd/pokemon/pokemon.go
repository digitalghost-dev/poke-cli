package pokemon

import (
	"flag"
	"fmt"
	"github.com/digitalghost-dev/poke-cli/cmd/utils"
	"github.com/digitalghost-dev/poke-cli/connections"
	"github.com/digitalghost-dev/poke-cli/flags"
	"github.com/digitalghost-dev/poke-cli/styling"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"math"
	"os"
	"strings"
)

// PokemonCommand processes the Pokémon command
func PokemonCommand() (string, error) {
	var output strings.Builder

	hintMessage := styling.StyleItalic.Render("options: [sm, md, lg]")

	flag.Usage = func() {
		helpMessage := styling.HelpBorder.Render(
			"Get details about a specific Pokémon.\n\n",
			styling.StyleBold.Render("USAGE:"),
			fmt.Sprintf("\n\t%s %s %s %s", "poke-cli", styling.StyleBold.Render("pokemon"), "<pokemon-name>", "[flag]"),
			fmt.Sprintf("\n\t%-30s", styling.StyleItalic.Render("Use a hyphen when typing a name with a space.")),
			"\n\n",
			styling.StyleBold.Render("FLAGS:"),
			fmt.Sprintf("\n\t%-30s %s", "-a, --abilities", "Prints the Pokémon's abilities."),
			fmt.Sprintf("\n\t%-30s %s", "-i=xx, --image=xx", "Prints out the Pokémon's default sprite."),
			fmt.Sprintf("\n\t%5s%-15s", "", hintMessage),
			fmt.Sprintf("\n\t%-30s %s", "-m, --moves", "Prints the Pokemon's learnable moves."),
			fmt.Sprintf("\n\t%-30s %s", "-s, --stats", "Prints the Pokémon's base stats."),
			fmt.Sprintf("\n\t%-30s %s", "-t, --types", "Prints the Pokémon's typing."),
			fmt.Sprintf("\n\t%-30s %s", "-h, --help", "Prints the help menu."),
		)
		output.WriteString(helpMessage)
	}

	pokeFlags, abilitiesFlag, shortAbilitiesFlag, imageFlag, shortImageFlag, moveFlag, shortMoveFlag, statsFlag, shortStatsFlag, typesFlag, shortTypesFlag := flags.SetupPokemonFlagSet()

	args := os.Args

	flag.Parse()

	if len(os.Args) == 3 && (os.Args[2] == "-h" || os.Args[2] == "--help") {
		flag.Usage()
		return output.String(), nil
	}

	err := utils.ValidatePokemonArgs(args)
	if err != nil {
		output.WriteString(err.Error()) // This is the styled error
		return output.String(), err
	}

	endpoint := strings.ToLower(args[1])
	pokemonName := strings.ToLower(args[2])

	if err := pokeFlags.Parse(args[3:]); err != nil {
		fmt.Printf("error parsing flags: %v\n", err)
		pokeFlags.Usage()
		os.Exit(1)
	}

	pokemonStruct, pokemonName, err := connections.PokemonApiCall(endpoint, pokemonName, connections.APIURL)
	if err != nil {
		output.WriteString(err.Error())
		return output.String(), err
	}
	capitalizedString := cases.Title(language.English).String(strings.ReplaceAll(pokemonName, "-", " "))

	// Weight calculation
	weightKilograms := float64(pokemonStruct.Weight) / 10
	weightPounds := float64(weightKilograms) * 2.20462

	// Height calculation
	heightMeters := float64(pokemonStruct.Height) / 10
	heightFeet := heightMeters * 3.28084
	feet := int(heightFeet)
	inches := int(math.Round((heightFeet - float64(feet)) * 12)) // Use math.Round to avoid truncation

	// Adjust for rounding to 12 inches (carry over to the next foot)
	if inches == 12 {
		feet++
		inches = 0
	}

	output.WriteString(fmt.Sprintf(
		"Your selected Pokémon: %s\n%s National Pokédex #: %d\n%s Weight: %.1fkg (%.1f lbs)\n%s Height: %.1fm (%d′%02d″)\n",
		capitalizedString, styling.ColoredBullet, pokemonStruct.ID,
		styling.ColoredBullet, weightKilograms, weightPounds,
		styling.ColoredBullet, heightFeet, feet, inches,
	))

	if *imageFlag != "" || *shortImageFlag != "" {
		// Determine the size based on the provided flags
		size := *imageFlag
		if *shortImageFlag != "" {
			size = *shortImageFlag
		}

		// Call the ImageFlag function with the specified size
		if err := flags.ImageFlag(&output, endpoint, pokemonName, size); err != nil {
			output.WriteString(fmt.Sprintf("%v\n", err))
			return output.String(), fmt.Errorf("%w", err)
		}
	}

	if *abilitiesFlag || *shortAbilitiesFlag {
		if err := flags.AbilitiesFlag(&output, endpoint, pokemonName); err != nil {
			output.WriteString(fmt.Sprintf("error parsing flags: %v\n", err))
			return "", fmt.Errorf("error parsing flags: %w", err)
		}
	}

	if *moveFlag || *shortMoveFlag {
		if err := flags.MovesFlag(&output, endpoint, pokemonName); err != nil {
			output.WriteString(fmt.Sprintf("error parsing flags: %v\n", err))
			return "", fmt.Errorf("error parsing flags: %w", err)
		}
	}

	if *typesFlag || *shortTypesFlag {
		if err := flags.TypesFlag(&output, endpoint, pokemonName); err != nil {
			output.WriteString(fmt.Sprintf("error parsing flags: %v\n", err))
			return "", fmt.Errorf("error parsing flags: %w", err)
		}
	}

	if *statsFlag || *shortStatsFlag {
		if err := flags.StatsFlag(&output, endpoint, pokemonName); err != nil {
			output.WriteString(fmt.Sprintf("error parsing flags: %v\n", err))
			return "", fmt.Errorf("error parsing flags: %w", err)
		}
	}

	return output.String(), nil
}
