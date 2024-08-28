package main

import (
	"flag"
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/digitalghost-dev/poke-cli/cmd"
	"github.com/digitalghost-dev/poke-cli/flags"
	"os"
)

func main() {
	var styleBold = lipgloss.NewStyle().Bold(true)

	latestFlag := flag.Bool("latest", false, "Prints the program's latest Docker Image and Release versions.")
	shortLatestFlag := flag.Bool("l", false, "Prints the program's latest Docker Image and Release versions.")

	flag.Usage = func() {
		fmt.Println("Welcome! This tool displays data about a selected Pokémon in the terminal!")

		// Usage section
		fmt.Println(styleBold.Render("\nUSAGE:"))
		fmt.Println("\t", "poke-cli [flag]")
		fmt.Println("\t", "poke-cli [command] [flag]")
		fmt.Println("\t", "poke-cli [command] [subcommand] [flag]")

		// Flags section
		fmt.Println(styleBold.Render("\nFLAGS:"))
		fmt.Println("\t", "-h, --help", "\t\t", "Shows the help menu")
		fmt.Println("\t", "-l, --latest", "\t\t", "Prints the latest version of the program")
		fmt.Print("\n")

		// Commands section
		fmt.Println(styleBold.Render("\nCOMMANDS"))
		fmt.Println("\t", "pokemon", "\t\t", "Get details of a specific Pokémon")
		fmt.Print("\n")
	}

	flag.Parse()

	commands := map[string]func(){
		"pokemon": cmd.PokemonCommand,
		"types":   cmd.TypesCommand,
	}

	if *latestFlag || *shortLatestFlag {
		flags.LatestFlag()
	} else if cmdFunc, exists := commands[os.Args[1]]; exists {
		cmdFunc()
	} else {
		fmt.Println("Unknown command")
	}
}
