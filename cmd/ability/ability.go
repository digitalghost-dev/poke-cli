package ability

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/digitalghost-dev/poke-cli/cmd/utils"
	"github.com/digitalghost-dev/poke-cli/connections"
	"github.com/digitalghost-dev/poke-cli/flags"
	"github.com/digitalghost-dev/poke-cli/styling"
)

func AbilityCommand() (string, error) {
	var output strings.Builder

	flag.Usage = func() {
		helpMessage := styling.HelpBorder.Render(
			"Get details about a specific ability.\n\n",
			styling.StyleBold.Render("USAGE:"),
			fmt.Sprintf("\n\t%s %s %s %s", "poke-cli", styling.StyleBold.Render("ability"), "<ability-name>", "[flag]"),
			fmt.Sprintf("\n\t%-30s", styling.StyleItalic.Render(styling.HyphenHint)),
			"\n\n",
			styling.StyleBold.Render("FLAGS:"),
			fmt.Sprintf("\n\t%-30s %s", "-p, --pokemon", "Prints Pokémon that learn this ability."),
			fmt.Sprintf("\n\t%-30s %s", "-h, --help", "Prints the help menu."),
		)
		output.WriteString(helpMessage)
	}

	af := flags.SetupAbilityFlagSet()

	args := os.Args

	if utils.CheckHelpFlag(&output, flag.Usage) {
		return output.String(), nil
	}

	flag.Parse()

	if err := utils.ValidateAbilityArgs(args); err != nil {
		output.WriteString(err.Error())
		return output.String(), err
	}

	endpoint := strings.ToLower(args[1])
	abilityName := strings.ToLower(args[2])

	if err := af.FlagSet.Parse(args[3:]); err != nil {
		output.WriteString(fmt.Sprintf("error parsing flags: %v\n", err))
		af.FlagSet.Usage()

		return output.String(), err
	}

	abilitiesStruct, abilityName, err := connections.AbilityApiCall(endpoint, abilityName, connections.APIURL)
	if err != nil {
		output.WriteString(err.Error())
		return output.String(), err
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

	capitalizedAbility := styling.CapitalizeResourceName(abilityName)
	output.WriteString(styling.StyleBold.Render(capitalizedAbility) + "\n")

	generationParts := strings.Split(abilitiesStruct.Generation.Name, "-")
	if len(generationParts) > 1 {
		generationUpper := strings.ToUpper(generationParts[1])
		output.WriteString(fmt.Sprintf("%s First introduced in generation "+generationUpper+"\n", styling.ColoredBullet))
	} else {
		output.WriteString(fmt.Sprintf("%s Generation: Unknown\n", styling.ColoredBullet))
	}

	// API is missing some data for the short_effect for abilities from Generation 9.
	// If short_effect is empty, fallback to the move's flavor_text_entry.
	if englishShortEffect == "" {
		output.WriteString(fmt.Sprintf("%s Effect: %s", styling.ColoredBullet, englishFlavorEntry))
	} else {
		output.WriteString(fmt.Sprintf("%s Effect: %s", styling.ColoredBullet, englishShortEffect))
	}

	if *af.Pokemon || *af.ShortPokemon {
		if err := flags.PokemonAbilitiesFlag(&output, endpoint, abilityName); err != nil {
			output.WriteString(fmt.Sprintf("error parsing flags: %v\n", err))
			return utils.HandleFlagError(&output, err)
		}
	}

	return output.String(), nil
}
