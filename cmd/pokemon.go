package cmd

import (
	"flag"
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

	if len(args) > 5 {
		return fmt.Errorf("error: too many arguments")
	}

	if len(args) < 3 {
		fmt.Println(errorColor.Render("Please declare a Pokémon's name after [pokemon] command"))
		fmt.Println(errorColor.Render("Run 'poke-cli --help' for more details"))
		return fmt.Errorf("error: insufficient arguments")
	}

	if len(args) > 3 {
		for _, arg := range args[3:] {
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

	var styleBold = lipgloss.NewStyle().Bold(true)
	var styleItalic = lipgloss.NewStyle().Italic(true)

	flag.Usage = func() {
		// Usage section
		fmt.Println(styleBold.Render("\nUSAGE:"))
		fmt.Println("\t", "poke-cli", styleBold.Render("pokemon"), "[flag]")
		fmt.Println("\t", "Get details about a specific Pokémon")
		fmt.Println("\t", "----------")
		fmt.Println("\t", styleItalic.Render("Examples:"), "\t", "poke-cli pokemon bulbasaur")
		fmt.Println("\t\t\t", "poke-cli pokemon flutter-mane --types")
		fmt.Println("\t\t\t", "poke-cli pokemon excadrill -t -a")

		// Flags section
		fmt.Println(styleBold.Render("\nFLAGS:"))
		fmt.Println("\t", "-a, --abilities", "\t", "Prints out the Pokémon's abilities.")
		fmt.Println("\t", "-t, --types", "\t\t", "Prints out the Pokémon's typing.")
		fmt.Print("\n")
	}

	flag.Parse()

	pokeFlags, typesFlag, shortTypesFlag, abilitiesFlag, shortAbilitiesFlag := flags.SetupPokemonFlagSet()

	args := os.Args

	err := ValidateArgs(args, errorColor)
	if err != nil {
		fmt.Println(errorColor.Render(err.Error()))
		os.Exit(1)
	}

	if args[2] == "-h" || args[2] == "--help" {
		flag.Usage()
		os.Exit(0)
	}

	endpoint := os.Args[1]
	pokemonName := strings.ToLower(args[2])

	if err := pokeFlags.Parse(args[3:]); err != nil {
		fmt.Printf("error parsing flags: %v\n", err)
		os.Exit(1)
	}

	_, pokemonName, pokemonID := connections.PokemonApiCall(endpoint, pokemonName, "https://pokeapi.co/api/v2/")
	capitalizedString := cases.Title(language.English).String(pokemonName)

	fmt.Printf("Your selected Pokémon: %s\nNational Pokédex #: %d\n", capitalizedString, pokemonID)

	if *typesFlag || *shortTypesFlag {
		if err := flags.TypesFlag(endpoint, pokemonName); err != nil {
			fmt.Printf("Error: %s\n", err)
			os.Exit(1)
		}
	}

	if *abilitiesFlag || *shortAbilitiesFlag {
		if err := flags.AbilitiesFlag(endpoint, pokemonName); err != nil {
			fmt.Printf("Error: %s\n", err)
			os.Exit(1)
		}
	}
}
