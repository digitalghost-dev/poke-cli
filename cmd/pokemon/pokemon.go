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

	hintMessage := styling.StyleItalic.Render("options: [sm, md, lg]")

	flag.Usage = func() {
		styledFlag := styling.ErrorColor.Render(fmt.Sprintf("%-30s", "-t, --types"))
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
			fmt.Sprintf("\n\t%s %s", styledFlag, styling.ErrorColor.Render("Deprecated. Types are included with each Pokémon.")),
			fmt.Sprintf("\n\t%-30s %s", "-h, --help", "Prints the help menu."),
		)
		output.WriteString(helpMessage)
	}

	pokeFlags, abilitiesFlag, shortAbilitiesFlag, defenseFlag, shortDefenseFlag, imageFlag, shortImageFlag, moveFlag, shortMoveFlag, statsFlag, shortStatsFlag, typesFlag, shortTypesFlag := flags.SetupPokemonFlagSet()

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

	pokemonSpeciesStruct, err := connections.PokemonSpeciesApiCall("pokemon-species", pokemonStruct.Species.Name, connections.APIURL)
	if err != nil {
		output.WriteString(err.Error())
		return output.String(), err
	}

	capitalizedString := cases.Title(language.English).String(strings.ReplaceAll(pokemonName, "-", " "))

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

	eggGroup := func(w io.Writer) {
		var eggGroupSlice []string

		for _, entry := range pokemonSpeciesStruct.EggGroups {
			modernEggGroupNames := map[string]string{
				"indeterminate": "Amorphous",
			}

			if name, exists := modernEggGroupNames[entry.Name]; exists {
				eggGroupSlice = append(eggGroupSlice, name)
			} else {
				capitalizedEggGroup := cases.Title(language.English).String(entry.Name)
				eggGroupSlice = append(eggGroupSlice, capitalizedEggGroup)
			}

		}

		sort.Strings(eggGroupSlice)
		fmt.Fprintf(w, "\n%s %s %s", styling.ColoredBullet, "Egg Group(s):", strings.Join(eggGroupSlice, ", "))
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

			capitalizedPokemonName := cases.Title(language.English).String(strings.ReplaceAll(evolvesFrom, "-", " "))
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
	eggGroup(&eggGroupOutput)
	typing(&typeOutput)
	metrics(&metricsOutput)
	species(&speciesOutput)

	output.WriteString(fmt.Sprintf(
		"Your selected Pokémon: %s\n%s\n%s%s%s%s\n",
		capitalizedString, entryOutput.String(), typeOutput.String(), metricsOutput.String(), speciesOutput.String(), eggGroupOutput.String(),
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

	flagChecks := []struct {
		condition bool
		flagFunc  func(io.Writer, string, string) error
	}{
		{*abilitiesFlag || *shortAbilitiesFlag, flags.AbilitiesFlag},
		{*defenseFlag || *shortDefenseFlag, flags.DefenseFlag},
		{*moveFlag || *shortMoveFlag, flags.MovesFlag},
		{*typesFlag || *shortTypesFlag, flags.TypesFlag},
		{*statsFlag || *shortStatsFlag, flags.StatsFlag},
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
