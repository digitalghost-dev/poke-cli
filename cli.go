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

func main() {
	latestFlag := flag.Bool("latest", false, "Prints the program's latest Docker Image and Release versions.")
	shortLatestFlag := flag.Bool("l", false, "Prints the program's latest Docker Image and Release versions.")

	flag.Usage = func() {
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

	flag.Parse()

	commands := map[string]func(){
		"pokemon": cmd.PokemonCommand,
		"types":   cmd.TypesCommand,
	}

	if len(os.Args) < 2 {
		flag.Usage()
	} else if *latestFlag || *shortLatestFlag {
		flags.LatestFlag()
	} else if cmdFunc, exists := commands[os.Args[1]]; exists {
		cmdFunc()
	} else {
		errMessage := errorBorder.Render(
			errorColor.Render("Error!"),
			styleBold.Render("\nAvailable Commands:"),
			fmt.Sprintf("\n\t%-15s %s", "pokemon", "Get details of a specific Pokémon"),
			fmt.Sprintf("\n\t%-15s %s", "types", "Get details of a specific typing\n"),
			fmt.Sprintf("\nAlso run %s for more info!", styleBold.Render("[poke-cli -h]")),
		)
		fmt.Printf("%s\n", errMessage)
	}
}
