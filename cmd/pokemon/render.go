package pokemon

import (
	"fmt"
	"io"
	"math"
	"sort"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/digitalghost-dev/poke-cli/cmd/utils"
	"github.com/digitalghost-dev/poke-cli/structs"
	"github.com/digitalghost-dev/poke-cli/styling"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func renderEggInformation(w io.Writer, s structs.PokemonSpeciesJSONStruct) {
	var eggInformationSlice []string

	modernEggInformationNames := map[string]string{
		"indeterminate": "Amorphous",
		"ground":        "Field",
		"humanshape":    "Human-Like",
		"plant":         "Grass",
		"no-eggs":       "Undiscovered",
	}

	for _, entry := range s.EggGroups {
		if name, exists := modernEggInformationNames[entry.Name]; exists {
			eggInformationSlice = append(eggInformationSlice, name)
		} else {
			eggInformationSlice = append(eggInformationSlice, cases.Title(language.English).String(entry.Name))
		}
	}

	sort.Strings(eggInformationSlice)

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

	fmt.Fprintf(w,
		"\n%s %s %s\n%s %s %s\n%s %s %d",
		styling.ColoredBullet,
		"Egg Group(s):", strings.Join(eggInformationSlice, ", "),
		styling.ColoredBullet,
		"Gender Rate:", m[s.GenderRate],
		styling.ColoredBullet,
		"Egg Cycles:", s.HatchCounter,
	)
}

func renderEffortValues(w io.Writer, s structs.PokemonJSONStruct) {
	nameMapping := map[string]string{
		"hp":              "HP",
		"attack":          "Atk",
		"defense":         "Def",
		"special-attack":  "SpA",
		"special-defense": "SpD",
		"speed":           "Spd",
	}

	var evs []string

	for _, effortValue := range s.Stats {
		if effortValue.Effort > 0 {
			name, ok := nameMapping[effortValue.Stat.Name]
			if !ok {
				name = "Missing from API"
			}
			evs = append(evs, fmt.Sprintf("%d %s", effortValue.Effort, name))
		}
	}

	fmt.Fprintf(w, "\n%s Effort Values: %s", styling.ColoredBullet, strings.Join(evs, ", "))
}

func renderEntry(w io.Writer, s structs.PokemonSpeciesJSONStruct) {
	for _, entry := range s.FlavorTextEntries {
		if entry.Language.Name == "en" && (entry.Version.Name == "x" || entry.Version.Name == "shield" || entry.Version.Name == "scarlet") {
			flavorText := strings.ReplaceAll(entry.FlavorText, "\n", " ")
			flavorText = strings.Join(strings.Fields(flavorText), " ")
			fmt.Fprintln(w, utils.WrapText(flavorText, 60))
			return
		}
	}
}

func renderMetrics(w io.Writer, s structs.PokemonJSONStruct) {
	weightKilograms := float64(s.Weight) / 10
	weightPounds := float64(weightKilograms) * 2.20462

	heightMeters := float64(s.Height) / 10
	heightFeet := heightMeters * 3.28084
	feet := int(heightFeet)
	inches := int(math.Round((heightFeet - float64(feet)) * 12))

	if inches == 12 {
		feet++
		inches = 0
	}

	fmt.Fprintf(w, "\n%s National Pokédex #: %d\n%s Weight: %.1fkg (%.1f lbs)\n%s Height: %.1fm (%d′%02d″)\n",
		styling.ColoredBullet, s.ID,
		styling.ColoredBullet, weightKilograms, weightPounds,
		styling.ColoredBullet, heightMeters, feet, inches)
}

func renderSpecies(w io.Writer, s structs.PokemonSpeciesJSONStruct) {
	if s.EvolvesFromSpecies.Name != "" {
		capitalizedPokemonName := styling.CapitalizeResourceName(s.EvolvesFromSpecies.Name)
		fmt.Fprintf(w, "%s %s %s", styling.ColoredBullet, "Evolves from:", capitalizedPokemonName)
	} else {
		fmt.Fprintf(w, "%s %s", styling.ColoredBullet, "Basic Pokémon")
	}
}

func renderTyping(w io.Writer, s structs.PokemonJSONStruct) {
	var typeBoxes []string

	for _, pokeType := range s.Types {
		colorHex, exists := styling.ColorMap[pokeType.Type.Name]
		if exists {
			typeColorStyle := lipgloss.NewStyle().
				Align(lipgloss.Center).
				Foreground(lipgloss.Color("#FAFAFA")).
				Background(lipgloss.Color(colorHex)).
				Margin(1, 1, 0, 0).
				Height(1).
				Width(14)

			rendered := typeColorStyle.Render(cases.Title(language.English).String(pokeType.Type.Name))
			typeBoxes = append(typeBoxes, rendered)
		}
	}

	fmt.Fprintln(w, lipgloss.JoinHorizontal(lipgloss.Top, typeBoxes...))
}
