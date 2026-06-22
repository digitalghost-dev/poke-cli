package ability

import (
	"errors"
	"fmt"
	"strings"

	"github.com/digitalghost-dev/poke-cli/cmd/utils"
	"github.com/digitalghost-dev/poke-cli/connections"
	"github.com/digitalghost-dev/poke-cli/flags"
	"github.com/digitalghost-dev/poke-cli/styling"
	flag "github.com/spf13/pflag"
)

func AbilityCommand(args []string) (string, error) {
	var output strings.Builder

	usage := func() {
		output.WriteString(
			utils.GenerateHelpMessage(
				utils.HelpConfig{
					Description:    "Get details about a specific ability.",
					CmdName:        "ability",
					SubCmdName:     "<ability-name>",
					ShowHyphenHint: true,
					Flags: []utils.FlagHelp{
						{Short: "-p", Long: "--pokemon", Description: "Prints Pokémon that learn this ability."},
					},
				},
			),
		)
	}

	af := flags.SetupAbilityFlagSet()

	if utils.CheckHelpFlag(args, usage) {
		return output.String(), nil
	}

	if err := utils.ValidateArgs(
		args,
		utils.Validator{MaxArgs: 3, CmdName: "ability", RequireName: true, HasFlags: true},
	); err != nil {
		output.WriteString(err.Error())
		return output.String(), err
	}

	endpoint := strings.ToLower(args[0])
	abilityName := strings.ToLower(args[1])

	if err := af.FlagSet.Parse(args[2:]); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return output.String(), nil
		}
		fmt.Fprintf(&output, "error parsing flags: %v\n", err)
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
	output.WriteString(styling.StyleBold.Render(capitalizedAbility))
	output.WriteByte('\n')

	generationParts := strings.Split(abilitiesStruct.Generation.Name, "-")
	if len(generationParts) > 1 {
		generationUpper := strings.ToUpper(generationParts[1])
		fmt.Fprintf(&output, "%s First introduced in generation %s\n", styling.ColoredBullet, generationUpper)
	} else {
		fmt.Fprintf(&output, "%s Generation: Unknown\n", styling.ColoredBullet)
	}

	// API is missing some data for the short_effect for abilities from Generation 9.
	// If short_effect is empty, fallback to the move's flavor_text_entry.
	if englishShortEffect == "" {
		fmt.Fprintf(&output, "%s Effect: %s", styling.ColoredBullet, englishFlavorEntry)
	} else {
		fmt.Fprintf(&output, "%s Effect: %s", styling.ColoredBullet, englishShortEffect)
	}

	if *af.Pokemon {
		if err := flags.PokemonAbilitiesFlag(&output, endpoint, abilityName); err != nil {
			return utils.HandleFlagError(&output, err)
		}
	}

	return output.String(), nil
}
