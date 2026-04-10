package pokemon

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"math"
	"os"
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/digitalghost-dev/poke-cli/cmd/utils"
	"github.com/digitalghost-dev/poke-cli/connections"
	"github.com/digitalghost-dev/poke-cli/flags"
	"github.com/digitalghost-dev/poke-cli/styling"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// PokemonCommand processes the Pokémon command
func PokemonCommand() (string, error) {
	var output strings.Builder

	flag.Usage = func() {
		output.WriteString(
			utils.GenerateHelpMessage(
				utils.HelpConfig{
					Description:    "Get details about a specific Pokémon.",
					CmdName:        "pokemon",
					SubCmdName:     "<pokemon-name> [flag]",
					ShowHyphenHint: true,
					Flags: []utils.FlagHelp{
						{Short: "-a", Long: "--abilities", Description: "Prints the Pokémon's abilities."},
						{Short: "-i=xx", Long: "--image=xx", Description: "Prints out the Pokémon's default sprite.\n\t     " + styling.StyleItalic.Render("options: [sm, md, lg]")},
						{Short: "-m", Long: "--moves", Description: "Prints the Pokémon's learnable moves."},
						{Short: "-s", Long: "--stats", Description: "Prints the Pokémon's base stats."},
						{Short: "-t", Long: "--types", Description: styling.ErrorColor.Render("Deprecated. Types are included with each Pokémon.")},
					},
				},
			),
		)
	}

	pf := flags.SetupPokemonFlagSet()

	args := os.Args

	if utils.CheckHelpFlag(&output, flag.Usage) {
		return output.String(), nil
	}

	err := utils.ValidatePokemonArgs(args)
	if err != nil {
		output.WriteString(err.Error())
		return output.String(), err
	}

	endpoint := strings.ToLower(args[1])
	pokemonName := strings.ToLower(args[2])

	if err := pf.FlagSet.Parse(args[3:]); err != nil {
		fmt.Printf("error parsing flags: %v\n", err)
		pf.FlagSet.Usage()
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

	entry := func(w io.Writer) {
		for _, entry := range pokemonSpeciesStruct.FlavorTextEntries {
			if entry.Language.Name == "en" && (entry.Version.Name == "x" || entry.Version.Name == "shield" || entry.Version.Name == "scarlet") {
				flavorText := strings.ReplaceAll(entry.FlavorText, "\n", " ")
				flavorText = strings.Join(strings.Fields(flavorText), " ")

				wrapped := utils.WrapText(flavorText, 60)
				fmt.Fprintln(w, wrapped)
				return
			}
		}
	}

	eggInformation := func(w io.Writer) {
		var eggInformationSlice []string

		for _, entry := range pokemonSpeciesStruct.EggGroups {
			modernEggInformationNames := map[string]string{
				"indeterminate": "Amorphous",
				"ground":        "Field",
				"humanshape":    "Human-Like",
				"plant":         "Grass",
				"no-eggs":       "Undiscovered",
			}

			if name, exists := modernEggInformationNames[entry.Name]; exists {
				eggInformationSlice = append(eggInformationSlice, name)
			} else {
				capitalizedEggInformation := cases.Title(language.English).String(entry.Name)
				eggInformationSlice = append(eggInformationSlice, capitalizedEggInformation)
			}
		}

		sort.Strings(eggInformationSlice)

		genderRate := pokemonSpeciesStruct.GenderRate
		m := map[int]string{
			-1: "Genderless",
			0:  "0% F",
			1:  "12.5% F",
			2:  "25% F",
			3:  "37.5% F",
			4:  "50% F",
			5:  "62.5% F",
			6:  "75% F",
			7:  "87.5% F",
			8:  "100% F",
		}

		hatchCounter := pokemonSpeciesStruct.HatchCounter

		fmt.Fprintf(w,
			"\n%s %s %s\n%s %s %s\n%s %s %d",
			styling.ColoredBullet,
			"Egg Group(s):", strings.Join(eggInformationSlice, ", "),
			styling.ColoredBullet,
			"Gender Rate:", m[genderRate],
			styling.ColoredBullet,
			"Egg Cycles:", hatchCounter,
		)
	}

	typing := func(w io.Writer) {
		var typeBoxes []string

		for _, pokeType := range pokemonStruct.Types {
			colorHex, exists := styling.ColorMap[pokeType.Type.Name]
			if exists {
				color := lipgloss.Color(colorHex)
				typeColorStyle := lipgloss.NewStyle().
					Align(lipgloss.Center).
					Foreground(lipgloss.Color("#FAFAFA")).
					Background(color).
					Margin(1, 1, 0, 0).
					Height(1).
					Width(14)

				rendered := typeColorStyle.Render(cases.Title(language.English).String(pokeType.Type.Name))
				typeBoxes = append(typeBoxes, rendered)
			}
		}

		joinedTypes := lipgloss.JoinHorizontal(lipgloss.Top, typeBoxes...)
		fmt.Fprintln(w, joinedTypes)
	}

	metrics := func(w io.Writer) {
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

		fmt.Fprintf(w, "\n%s National Pokédex #: %d\n%s Weight: %.1fkg (%.1f lbs)\n%s Height: %.1fm (%d′%02d″)\n",
			styling.ColoredBullet, pokemonStruct.ID,
			styling.ColoredBullet, weightKilograms, weightPounds,
			styling.ColoredBullet, heightMeters, feet, inches)
	}

	species := func(w io.Writer) {
		if pokemonSpeciesStruct.EvolvesFromSpecies.Name != "" {
			evolvesFrom := pokemonSpeciesStruct.EvolvesFromSpecies.Name

			capitalizedPokemonName := styling.CapitalizeResourceName(evolvesFrom)
			fmt.Fprintf(w, "%s %s %s", styling.ColoredBullet, "Evolves from:", capitalizedPokemonName)
		} else {
			fmt.Fprintf(w, "%s %s", styling.ColoredBullet, "Basic Pokémon")
		}
	}

	var (
		entryOutput    bytes.Buffer
		eggGroupOutput bytes.Buffer
		typeOutput     bytes.Buffer
		metricsOutput  bytes.Buffer
		speciesOutput  bytes.Buffer
	)

	entry(&entryOutput)
	eggInformation(&eggGroupOutput)
	typing(&typeOutput)
	metrics(&metricsOutput)
	species(&speciesOutput)

	fmt.Fprintf(&output,
		"Your selected Pokémon: %s\n%s\n%s%s%s%s\n",
		capitalizedString, entryOutput.String(), typeOutput.String(), metricsOutput.String(), speciesOutput.String(), eggGroupOutput.String(),
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
