// This file holds all the flags used by the <pokemonName> subcommand

package flags

import (
	"flag"
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/digitalghost-dev/poke-cli/connections"
	"strings"
)

var (
	helpBorder = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#FFCC00"))
	styleBold = lipgloss.NewStyle().Bold(true)
)

func header(header string) {
	HeaderBold := lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("#FFCC00")).
		BorderTop(true).
		Bold(true).
		Render(header)

	fmt.Println(HeaderBold)
}

func SetupPokemonFlagSet() (*flag.FlagSet, *bool, *bool, *bool, *bool, *bool, *bool) {
	pokeFlags := flag.NewFlagSet("pokeFlags", flag.ExitOnError)

	abilitiesFlag := pokeFlags.Bool("abilities", false, "Print the Pokémon's abilities")
	shortAbilitiesFlag := pokeFlags.Bool("a", false, "Print the Pokémon's abilities")

	statsFlag := pokeFlags.Bool("stats", false, "Print the Pokémon's base stats")
	shortStatsFlag := pokeFlags.Bool("s", false, "Print the Pokémon's base stats")

	typesFlag := pokeFlags.Bool("types", false, "Print the Pokémon's typing")
	shortTypesFlag := pokeFlags.Bool("t", false, "Prints the Pokémon's typing")

	pokeFlags.Usage = func() {
		fmt.Println(
			helpBorder.Render("poke-cli pokemon <pokemon-name> [flags]",
				styleBold.Render("\n\nFLAGS:"), "\n\t", "-a, --abilities", "\t", "Prints out the Pokémon's abilities.",
				"\n\t", "-t, --types", "\t\t", "Prints out the Pokémon's typing.", "\n\t", "-s, --stats", "\t\t",
				"Prints out the Pokémon's base stats."),
		)
	}

	return pokeFlags, abilitiesFlag, shortAbilitiesFlag, statsFlag, shortStatsFlag, typesFlag, shortTypesFlag
}

func AbilitiesFlag(endpoint string, pokemonName string) error {
	baseURL := "https://pokeapi.co/api/v2/"
	pokemonStruct, _, _, _, _ := connections.PokemonApiCall(endpoint, pokemonName, baseURL)

	header("Abilities")

	for _, pokeAbility := range pokemonStruct.Abilities {
		if pokeAbility.Slot == 1 {
			fmt.Printf("Ability %d: %s\n", pokeAbility.Slot, pokeAbility.Ability.Name)
		} else if pokeAbility.Slot == 2 {
			fmt.Printf("Ability %d: %s\n", pokeAbility.Slot, pokeAbility.Ability.Name)
		} else {
			fmt.Printf("Hidden Ability: %s\n", pokeAbility.Ability.Name)
		}
	}

	return nil
}

func StatsFlag(endpoint string, pokemonName string) error {
	baseURL := "https://pokeapi.co/api/v2/"
	pokemonStruct, _, _, _, _ := connections.PokemonApiCall(endpoint, pokemonName, baseURL)

	header("Base Stats")

	// Helper function to map stat values to categories
	getStatCategory := func(value int) string {
		switch {
		case value < 20:
			return "lowest"
		case value < 60:
			return "lower"
		case value < 90:
			return "low"
		case value < 120:
			return "high"
		case value < 150:
			return "higher"
		default:
			return "highest"
		}
	}

	// Helper function to print the bar for a stat
	printBar := func(label string, value, maxWidth, maxValue int, style lipgloss.Style) {
		scaledValue := (value * maxWidth) / maxValue
		bar := strings.Repeat("▇", scaledValue)
		coloredBar := style.Render(bar)
		fmt.Printf("%-10s %s %d\n", label, coloredBar, value)
	}

	// Mapping from API stat names to custom display names
	nameMapping := map[string]string{
		"hp":              "HP",
		"attack":          "Atk",
		"defense":         "Def",
		"special-attack":  "Sp. Atk",
		"special-defense": "Sp. Def",
		"speed":           "Speed",
	}

	statColorMap := map[string]string{
		"lowest":  "#F34444",
		"lower":   "#FF7F0F",
		"low":     "#FFDD57",
		"high":    "#A0E515",
		"higher":  "#22C65A",
		"highest": "#00C2B8",
	}

	// Find the maxium stat value
	maxValue := 0
	for _, stat := range pokemonStruct.Stats {
		if stat.BaseStat > maxValue {
			maxValue = stat.BaseStat
		}
	}

	maxWidth := 45

	// Print bars for each stat
	for _, stat := range pokemonStruct.Stats {
		apiName := stat.Stat.Name
		customName, exists := nameMapping[apiName]
		if !exists {
			continue
		}

		category := getStatCategory(stat.BaseStat)
		color := statColorMap[category]
		style := lipgloss.NewStyle().Foreground(lipgloss.Color(color))

		printBar(customName, stat.BaseStat, maxWidth, maxValue, style)
	}

	return nil
}

func TypesFlag(endpoint string, pokemonName string) error {
	baseURL := "https://pokeapi.co/api/v2/"
	pokemonStruct, _, _, _, _ := connections.PokemonApiCall(endpoint, pokemonName, baseURL)

	colorMap := map[string]string{
		"normal":   "#B7B7A9",
		"fire":     "#FF4422",
		"water":    "#3499FF",
		"electric": "#FFCC33",
		"grass":    "#77CC55",
		"ice":      "#66CCFF",
		"fighting": "#BB5544",
		"poison":   "#AA5699",
		"ground":   "#DEBB55",
		"flying":   "#889AFF",
		"psychic":  "#FF5599",
		"bug":      "#AABC22",
		"rock":     "#BBAA66",
		"ghost":    "#6666BB",
		"dragon":   "#7766EE",
		"dark":     "#775544",
		"steel":    "#AAAABB",
		"fairy":    "#EE99EE",
	}

	header("Typing")

	for _, pokeType := range pokemonStruct.Types {
		colorHex, exists := colorMap[pokeType.Type.Name]
		if exists {
			color := lipgloss.Color(colorHex)
			fmt.Printf("Type %d: %s\n", pokeType.Slot, lipgloss.NewStyle().Bold(true).Foreground(color).Render(pokeType.Type.Name))
		} else {
			fmt.Printf("Type %d: %s\n", pokeType.Slot, pokeType.Type.Name)
		}
	}

	return nil
}
