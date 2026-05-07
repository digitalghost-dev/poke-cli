package pokemon

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"strings"

	"github.com/digitalghost-dev/poke-cli/cmd/utils"
	"github.com/digitalghost-dev/poke-cli/connections"
	"github.com/digitalghost-dev/poke-cli/flags"
	"github.com/digitalghost-dev/poke-cli/styling"
)

// PokemonCommand processes the Pokémon command
func PokemonCommand(args []string) (string, error) {
	var output strings.Builder

	usage := func() {
		output.WriteString(
			utils.GenerateHelpMessage(
				utils.HelpConfig{
					Description:    "Get details about a specific Pokémon.",
					CmdName:        "pokemon",
					SubCmdName:     "<pokemon-name> [flag]",
					ShowHyphenHint: true,
					Flags: []utils.FlagHelp{
						{Short: "-a", Long: "--abilities", Description: "Prints the Pokémon's abilities."},
						{Short: "-d", Long: "--defense", Description: "Prints the Pokémon's type defenses."},
						{Short: "-i=xx", Long: "--image=xx", Description: "Prints out the Pokémon's default sprite.\n\t     " + styling.StyleItalic.Render("options: [sm, md, lg]")},
						{Short: "-m", Long: "--moves", Description: "Prints the Pokémon's learnable moves."},
						{Short: "-s", Long: "--stats", Description: "Prints the Pokémon's base stats."},
						{Short: "-t", Long: "--types", Description: styling.ErrorColor.Render("Deprecated. Typing is included by default.")},
					},
				},
			),
		)
	}

	pf := flags.SetupPokemonFlagSet()

	if utils.CheckHelpFlag(args, usage) {
		return output.String(), nil
	}

	validationArgs := append([]string{"poke-cli"}, args...)
	err := utils.ValidatePokemonArgs(validationArgs)
	if err != nil {
		output.WriteString(err.Error())
		return output.String(), err
	}

	endpoint := strings.ToLower(args[0])
	pokemonName := strings.ToLower(args[1])

	if err := pf.FlagSet.Parse(args[2:]); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return output.String(), nil
		}
		fmt.Fprintf(&output, "error parsing flags: %v\n", err)
		return output.String(), err
	}

	pokemonStruct, pokemonName, err := connections.PokemonApiCall(endpoint, pokemonName, connections.APIURL)
	if err != nil {
		output.WriteString(err.Error())
		return output.String(), err
	}

	pokemonSpeciesStruct, _, err := connections.PokemonSpeciesApiCall("pokemon-species", pokemonStruct.Species.Name, connections.APIURL)
	if err != nil {
		output.WriteString(err.Error())
		return output.String(), err
	}

	capitalizedString := styling.CapitalizeResourceName(pokemonName)

	var (
		entryOutput        bytes.Buffer
		eggGroupOutput     bytes.Buffer
		typeOutput         bytes.Buffer
		metricsOutput      bytes.Buffer
		speciesOutput      bytes.Buffer
		effortValuesOutput bytes.Buffer
	)

	renderEntry(&entryOutput, pokemonSpeciesStruct)
	renderEggInformation(&eggGroupOutput, pokemonSpeciesStruct)
	renderTyping(&typeOutput, pokemonStruct)
	renderMetrics(&metricsOutput, pokemonStruct)
	renderSpecies(&speciesOutput, pokemonSpeciesStruct)
	renderEffortValues(&effortValuesOutput, pokemonStruct)

	fmt.Fprintf(&output,
		"Your selected Pokémon: %s\n%s\n%s%s%s%s%s\n",
		capitalizedString, entryOutput.String(), typeOutput.String(), metricsOutput.String(), speciesOutput.String(), eggGroupOutput.String(), effortValuesOutput.String(),
	)

	if *pf.Image != "" || *pf.ShortImage != "" {
		// Determine the size based on the provided flags
		size := *pf.Image
		if *pf.ShortImage != "" {
			size = *pf.ShortImage
		}

		// Call the ImageFlag function with the specified size
		if err := flags.ImageFlag(&output, endpoint, pokemonName, size); err != nil {
			fmt.Fprintf(&output, "%v\n", err)
			return output.String(), fmt.Errorf("%w", err)
		}
	}

	flagChecks := []struct {
		condition bool
		flagFunc  func(io.Writer, string, string) error
	}{
		{*pf.Abilities || *pf.ShortAbilities, flags.AbilitiesFlag},
		{*pf.Defense || *pf.ShortDefense, flags.DefenseFlag},
		{*pf.Move || *pf.ShortMove, flags.MovesFlag},
		{*pf.Stats || *pf.ShortStats, flags.StatsFlag},
		{*pf.Types || *pf.ShortTypes, flags.TypesFlag},
	}

	for _, check := range flagChecks {
		if check.condition {
			if err := check.flagFunc(&output, endpoint, pokemonName); err != nil {
				return utils.HandleFlagError(&output, err)
			}
		}
	}

	return output.String(), nil
}
