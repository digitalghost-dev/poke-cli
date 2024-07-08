package subcommands

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/digitalghost-dev/poke-cli/connections"
	"github.com/digitalghost-dev/poke-cli/flags"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"os"
	"strings"
)

const red = lipgloss.Color("#F2055C")

var errorColor = lipgloss.NewStyle().Foreground(red)

// ValidateArgs validates the command line arguments
func ValidateArgs(args []string) error {

	if len(args) > 2 && !strings.HasPrefix(args[2], "-") {
		return fmt.Errorf("error: only flags are allowed after declaring a Pokémon's name")
	}

	if len(args) > 4 {
		return fmt.Errorf("error: too many arguments")
	}

	return nil
}

// PokemonCommand processes the Pokémon command
func PokemonCommand() {
	pokeFlags, typesFlag, abilitiesFlag := flags.SetupPokemonFlagSet()

	args := os.Args

	PokemonName := strings.ToLower(args[1])

	err := ValidateArgs(args)
	if err != nil {
		fmt.Println(errorColor.Render(err.Error()))
		os.Exit(1)
	}

	if err := pokeFlags.Parse(args[2:]); err != nil {
		fmt.Printf("error parsing flags: %v\n", err)
		os.Exit(1)
	}

	_, pokemonName, pokemonID := connections.PokemonApiCall(PokemonName, "https://pokeapi.co/api/v2/pokemon/")
	capitalizedString := cases.Title(language.English).String(pokemonName)

	fmt.Printf("Your selected Pokémon: %s\nNational Pokédex #: %d\n", capitalizedString, pokemonID)

	if *typesFlag {
		if err := flags.TypesFlag(); err != nil {
			fmt.Printf("Error: %s\n", err)
			os.Exit(1)
		}
	}

	if *abilitiesFlag {
		if err := flags.AbilitiesFlag(); err != nil {
			fmt.Printf("Error: %s\n", err)
			os.Exit(1)
		}
	}
}
