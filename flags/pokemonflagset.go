// This file holds all the flags used by the <pokemonName> subcommand

package flags

import (
	"flag"
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/digitalghost-dev/poke-cli/connections"
)

func SetupPokemonFlagSet() (*flag.FlagSet, *bool, *bool) {
	pokeFlags := flag.NewFlagSet("pokeFlags", flag.ExitOnError)

	typesFlag := pokeFlags.Bool("types", false, "Print the declared Pokémon's typing")
	abilitiesFlag := pokeFlags.Bool("abilities", false, "Print the declared Pokémon's abilities")

	return pokeFlags, typesFlag, abilitiesFlag
}

func AbilitiesFlag(pokemonName string) error {
	pokemonStruct, _, _ := connections.PokemonApiCall(pokemonName, "https://pokeapi.co/api/v2/pokemon/")

	abilitiesHeaderBold := lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("#FFCC00")).
		BorderTop(true).
		Bold(true).
		Render("Abilities")

	fmt.Println(abilitiesHeaderBold)
	for _, pokeAbility := range pokemonStruct.Abilities {
		fmt.Printf("Ability %d: %s\n", pokeAbility.Slot, pokeAbility.Ability.Name)
	}

	return nil
}

func TypesFlag(pokemonName string) error {
	pokemonStruct, _, _ := connections.PokemonApiCall(pokemonName, "https://pokeapi.co/api/v2/pokemon/")

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
