package main

import (
	"flag"
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/digitalghost-dev/poke-cli/subcommands"
	"os"
)

var styleBold = lipgloss.NewStyle().Bold(true)
var styleItalic = lipgloss.NewStyle().Italic(true)

func main() {

	flag.Usage = func() {
		fmt.Println("Welcome! This tool displays data about a selected Pokémon in the terminal!")

		fmt.Println(styleBold.Render("\nUSAGE:"))
		fmt.Println("\t", "poke-cli [flag]")
		fmt.Println("\t", "poke-cli [pokemon name] [flag]")
		fmt.Println("\t", "----------")
		fmt.Println("\t", styleItalic.Render("Example:"), "poke-cli bulbasaur", styleItalic.Render("or"), "poke-cli flutter-mane --types")

		fmt.Println(styleBold.Render("\nGLOBAL FLAGS:"))
		fmt.Println("\t", "-h, --help", "\t", "Shows the help menu")
		fmt.Print("\n")

		fmt.Println(styleBold.Render("POKEMON NAME FLAGS:"))
		fmt.Println("\t", "Add a flag after declaring a Pokémon's name for more details!")
		fmt.Print("\t", "--types", "\t\t", "Prints out the Pokémon's typing.\n\n")
	}

	flag.Parse()

	if len(os.Args) < 2 {
		fmt.Println("Please declare a Pokémon's name after the CLI name")
		fmt.Println("Run 'poke-cli --help' for more details")
		os.Exit(1)
	}

	subcommands.PokemonCommand()
}
