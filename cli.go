package main

import (
	"flag"
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/digitalghost-dev/poke-cli/cmd"
	"github.com/digitalghost-dev/poke-cli/flags"
	"os"
)

var (
	styleBold  = lipgloss.NewStyle().Bold(true)
	helpBorder = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#FFCC00"))
	errorColor  = lipgloss.NewStyle().Foreground(lipgloss.Color("#F2055C"))
	errorBorder = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#F2055C"))
)

func runCLI(args []string) int {
	mainFlagSet := flag.NewFlagSet("poke-cli", flag.ContinueOnError)
	latestFlag := mainFlagSet.Bool("latest", false, "Prints the program's latest Docker Image and Release versions.")
	shortLatestFlag := mainFlagSet.Bool("l", false, "Prints the program's latest Docker Image and Release versions.")

	mainFlagSet.Usage = func() {
		helpMessage := helpBorder.Render(
			"Welcome! This tool displays data related to Pokémon!",
			"\n\n", styleBold.Render("USAGE:"),
			fmt.Sprintf("\n\t%-15s %s", "poke-cli [flag]", ""),
			fmt.Sprintf("\n\t%-15s %s", "poke-cli [command] [flag]", ""),
			fmt.Sprintf("\n\t%-15s %s", "poke-cli [command] [subcommand] [flag]", ""),
			"\n\n", styleBold.Render("FLAGS:"),
			fmt.Sprintf("\n\t%-15s %s", "-h, --help", "Shows the help menu"),
			fmt.Sprintf("\n\t%-15s %s", "-l, --latest", "Prints the latest available"),
			fmt.Sprintf("\n\t%-15s %s", "", "version of the program"),
			"\n\n", styleBold.Render("AVAILABLE COMMANDS:"),
			fmt.Sprintf("\n\t%-15s %s", "pokemon", "Get details of a specific Pokémon"),
			fmt.Sprintf("\n\t%-15s %s", "types", "Get details of a specific typing"),
		)
		fmt.Println(helpMessage)
	}

	// Check for help flag manually
	for _, arg := range args {
		if arg == "-h" || arg == "--help" {
			mainFlagSet.Usage()
			return 0
		}
	}

	err := mainFlagSet.Parse(args)
	if err != nil {
		return 2
	}

	commands := map[string]func(){
		"pokemon": cmd.PokemonCommand,
		"types":   cmd.TypesCommand,
	}

	if len(os.Args) < 2 {
		mainFlagSet.Usage()
		return 1
	} else if *latestFlag || *shortLatestFlag {
		flags.LatestFlag()
		return 0
	} else if cmdFunc, exists := commands[os.Args[1]]; exists {
		cmdFunc()
		return 0
	} else {
		errMessage := errorBorder.Render(
			errorColor.Render("Error!"),
			styleBold.Render("\nAvailable Commands:"),
			fmt.Sprintf("\n\t%-15s %s", "pokemon", "Get details of a specific Pokémon"),
			fmt.Sprintf("\n\t%-15s %s", "types", "Get details of a specific typing\n"),
			fmt.Sprintf("\nAlso run %s for more info!", styleBold.Render("[poke-cli -h]")),
		)
		fmt.Printf("%s\n", errMessage)
		return 1
	}
}

var exit = os.Exit // Default to os.Exit, but you can override this in tests

func main() {
	exit(runCLI(os.Args[1:]))
}
