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

// ValidateArgs validates the command line arguments
func ValidateArgs(args []string, errorColor lipgloss.Style) error {
	if len(args) > 4 {
		return fmt.Errorf("error: too many arguments")
	}

	if len(args) < 2 {
		fmt.Println(errorColor.Render("Please declare a Pokémon's name after the CLI name"))
		fmt.Println(errorColor.Render("Run 'poke-cli --help' for more details"))
		return fmt.Errorf("error: insufficient arguments")
	}

	if len(args) > 2 {
		for _, arg := range args[2:] {
			if arg[0] != '-' {
				errorMsg := fmt.Sprintf("Error: Invalid argument '%s'. Only flags are allowed after declaring a Pokémon's name", arg)
				return fmt.Errorf(errorColor.Render(strings.TrimSpace(errorMsg)))
			}
		}
	}

	return nil
}

// PokemonCommand processes the Pokémon command
func PokemonCommand() {
	const red = lipgloss.Color("#F2055C")
	var errorColor = lipgloss.NewStyle().Foreground(red)

	pokeFlags, typesFlag, shortTypesFlag, abilitiesFlag, shortAbilitiesFlag := flags.SetupPokemonFlagSet()

	args := os.Args

	err := ValidateArgs(args, errorColor)
	if err != nil {
		fmt.Println(errorColor.Render(err.Error()))
		os.Exit(1)
	}

	pokemonName := strings.ToLower(args[1])

	if err := pokeFlags.Parse(args[2:]); err != nil {
		fmt.Printf("error parsing flags: %v\n", err)
		os.Exit(1)
	}

	_, pokemonName, pokemonID := connections.PokemonApiCall(pokemonName, "https://pokeapi.co/api/v2/pokemon/")
	capitalizedString := cases.Title(language.English).String(pokemonName)

	fmt.Printf("Your selected Pokémon: %s\nNational Pokédex #: %d\n", capitalizedString, pokemonID)

	if *typesFlag || *shortTypesFlag {
		if err := flags.TypesFlag(pokemonName); err != nil {
			fmt.Printf("Error: %s\n", err)
			os.Exit(1)
		}
	}

	if *abilitiesFlag || *shortAbilitiesFlag {
		if err := flags.AbilitiesFlag(pokemonName); err != nil {
			fmt.Printf("Error: %s\n", err)
			os.Exit(1)
		}
	}
}
