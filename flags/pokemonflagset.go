// This file holds all the flags used by the <pokemonName> subcommand

package flags

import (
	"flag"
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/digitalghost-dev/poke-cli/connections"
)

var (
	helpBorder = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#FFCC00"))
	styleBold = lipgloss.NewStyle().Bold(true)
)

func SetupPokemonFlagSet() (*flag.FlagSet, *bool, *bool, *bool, *bool) {
	pokeFlags := flag.NewFlagSet("pokeFlags", flag.ExitOnError)

	typesFlag := pokeFlags.Bool("types", false, "Print the declared Pokémon's typing")
	shortTypesFlag := pokeFlags.Bool("t", false, "Prints the declared Pokémon's typing")

	abilitiesFlag := pokeFlags.Bool("abilities", false, "Print the declared Pokémon's abilities")
	shortAbilitiesFlag := pokeFlags.Bool("a", false, "Print the declared Pokémon's abilities")

	pokeFlags.Usage = func() {
		fmt.Println(
			helpBorder.Render("poke-cli pokemon <pokemon-name> [flags]",
				styleBold.Render("\n\nFLAGS:"), "\n\t", "-a, --abilities", "\t", "Prints out the Pokémon's abilities.",
				"\n\t", "-t, --types", "\t\t", "Prints out the Pokémon's typing."),
		)
	}

	return pokeFlags, typesFlag, shortTypesFlag, abilitiesFlag, shortAbilitiesFlag
}

func AbilitiesFlag(endpoint string, pokemonName string) error {
	baseURL := "https://pokeapi.co/api/v2/"
	pokemonStruct, _, _ := connections.PokemonApiCall(endpoint, pokemonName, baseURL)

	abilitiesHeaderBold := lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("#FFCC00")).
		BorderTop(true).
		Bold(true).
		Render("Abilities")

	fmt.Println(abilitiesHeaderBold)
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

func TypesFlag(endpoint string, pokemonName string) error {
	baseURL := "https://pokeapi.co/api/v2/"
	pokemonStruct, _, _ := connections.PokemonApiCall(endpoint, pokemonName, baseURL)

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

	typingHeaderBold := lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("#FFCC00")).
		BorderTop(true).
		Bold(true).
		Render("Typing")

	fmt.Println(typingHeaderBold)
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
