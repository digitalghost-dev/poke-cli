package main

import (
	"flag"
	"fmt"
	"github.com/digitalghost-dev/poke-cli/cmd/ability"
	"github.com/digitalghost-dev/poke-cli/cmd/item"
	"github.com/digitalghost-dev/poke-cli/cmd/move"
	"github.com/digitalghost-dev/poke-cli/cmd/natures"
	"github.com/digitalghost-dev/poke-cli/cmd/pokemon"
	"github.com/digitalghost-dev/poke-cli/cmd/search"
	"github.com/digitalghost-dev/poke-cli/cmd/types"
	"github.com/digitalghost-dev/poke-cli/cmd/utils"
	"github.com/digitalghost-dev/poke-cli/flags"
	"github.com/digitalghost-dev/poke-cli/styling"
	"os"
	"runtime/debug"
	"strings"
)

var version = "(devel)"

func currentVersion() {
	if version != "(devel)" {
		// Use version injected by -ldflags
		fmt.Printf("Version: %s\n", version)
		return
	}

	// Fallback to build info when the version is not set
	buildInfo, ok := debug.ReadBuildInfo()
	if !ok {
		fmt.Println("Version: unknown (unable to read build info)")
		return
	}

	if buildInfo.Main.Version != "" {
		fmt.Printf("Version: %s\n", buildInfo.Main.Version)
	} else {
		fmt.Println("Version: (devel)")
	}
}

func runCLI(args []string) int {
	var output strings.Builder

	mainFlagSet := flag.NewFlagSet("poke-cli", flag.ContinueOnError)

	// -l, --latest flag retrieves the latest Docker image and GitHub release versions available
	latestFlag := mainFlagSet.Bool("latest", false, "Prints the program's latest Docker image and release versions.")
	shortLatestFlag := mainFlagSet.Bool("l", false, "Prints the program's latest Docker image and release versions.")

	// -v, --version flag retrieves the currently installed version
	currentVersionFlag := mainFlagSet.Bool("version", false, "Prints the current version")
	shortCurrentVersionFlag := mainFlagSet.Bool("v", false, "Prints the current version")

	mainFlagSet.Usage = func() {
		helpMessage := styling.HelpBorder.Render(
			"Welcome! This tool displays data related to Pokémon!",
			"\n\n", styling.StyleBold.Render("USAGE:"),
			fmt.Sprintf("\n\t%-15s %s", "poke-cli [flag]", ""),
			fmt.Sprintf("\n\t%-15s %s", "poke-cli <command> [flag]", ""),
			fmt.Sprintf("\n\t%-15s %s", "poke-cli <command> <subcommand> [flag]", ""),
			"\n\n", styling.StyleBold.Render("FLAGS:"),
			fmt.Sprintf("\n\t%-15s %s", "-h, --help", "Shows the help menu"),
			fmt.Sprintf("\n\t%-15s %s", "-l, --latest", "Prints the latest version available"),
			fmt.Sprintf("\n\t%-15s %s", "-v, --version", "Prints the current version"),
			"\n\n", styling.StyleBold.Render("COMMANDS:"),
			fmt.Sprintf("\n\t%-15s %s", "ability", "Get details about an ability"),
			fmt.Sprintf("\n\t%-15s %s", "item", "Get details about an item"),
			fmt.Sprintf("\n\t%-15s %s", "move", "Get details about a move"),
			fmt.Sprintf("\n\t%-15s %s", "natures", "Get details about all natures"),
			fmt.Sprintf("\n\t%-15s %s", "pokemon", "Get details about a Pokémon"),
			fmt.Sprintf("\n\t%-15s %s", "search", "Search for a resource"),
			fmt.Sprintf("\n\t%-15s %s", "types", "Get details about a typing"),
			"\n\n", styling.StyleItalic.Render("hint: when calling a resource with a space, use a hyphen"),
			"\n", styling.StyleItalic.Render("example: poke-cli ability strong-jaw"),
			"\n", styling.StyleItalic.Render("example: poke-cli pokemon flutter-mane"),
			"\n\n", fmt.Sprintf("%s %s", "↓ ctrl/cmd + click for docs/guides\n", styling.DocsLink),
		)
		fmt.Println(helpMessage)
	}

	switch {
	case len(args) == 0:
		mainFlagSet.Usage()
		return 0
	case len(args) > 0:
		if args[0] == "-h" || args[0] == "--help" {
			mainFlagSet.Usage()
			return 0
		}
	}

	err := mainFlagSet.Parse(args)
	if err != nil {
		return 2
	}

	remainingArgs := mainFlagSet.Args()

	commands := map[string]func() int{
		"ability": utils.HandleCommandOutput(ability.AbilityCommand),
		"item":    utils.HandleCommandOutput(item.ItemCommand),
		"move":    utils.HandleCommandOutput(move.MoveCommand),
		"natures": utils.HandleCommandOutput(natures.NaturesCommand),
		"pokemon": utils.HandleCommandOutput(pokemon.PokemonCommand),
		"types":   utils.HandleCommandOutput(types.TypesCommand),
		"search": func() int {
			search.SearchCommand()
			return 0
		},
	}

	cmdArg := ""
	if len(remainingArgs) >= 1 {
		cmdArg = remainingArgs[0]
	}
	cmdFunc, exists := commands[cmdArg]

	switch {
	case len(remainingArgs) == 0 && !*latestFlag && !*shortLatestFlag && !*currentVersionFlag && !*shortCurrentVersionFlag:
		mainFlagSet.Usage()
		return 1
	case *latestFlag || *shortLatestFlag:
		flags.LatestFlag()
		return 0
	case *currentVersionFlag || *shortCurrentVersionFlag:
		currentVersion()
		return 0
	case exists:
		return cmdFunc()
	default:
		errMessage := styling.ErrorBorder.Render(
			styling.ErrorColor.Render("Error!"),
			fmt.Sprintf("\n\t%-15s", fmt.Sprintf("'%s' is not a valid command.\n", cmdArg)),
			styling.StyleBold.Render("\nCommands:"),
			fmt.Sprintf("\n\t%-15s %s", "ability", "Get details about an ability"),
			fmt.Sprintf("\n\t%-15s %s", "item", "Get details about an item"),
			fmt.Sprintf("\n\t%-15s %s", "move", "Get details about a move"),
			fmt.Sprintf("\n\t%-15s %s", "natures", "Get details about all natures"),
			fmt.Sprintf("\n\t%-15s %s", "pokemon", "Get details about a Pokémon"),
			fmt.Sprintf("\n\t%-15s %s", "search", "Search for a resource"),
			fmt.Sprintf("\n\t%-15s %s", "types", "Get details about a typing"),
			fmt.Sprintf("\n\nAlso run %s for more info!", styling.StyleBold.Render("poke-cli -h")),
		)
		output.WriteString(errMessage)

		fmt.Println(output.String())

		return 1
	}
}

var exit = os.Exit

func main() {
	exit(runCLI(os.Args[1:]))
}
