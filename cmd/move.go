package cmd

import (
	"flag"
	"fmt"
	"github.com/digitalghost-dev/poke-cli/connections"
	"github.com/digitalghost-dev/poke-cli/flags"
	"github.com/digitalghost-dev/poke-cli/styling"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"os"
	"strings"
)

func MoveCommand() {
	flag.Usage = func() {
		helpMessage := styling.HelpBorder.Render(
			"Get details about a specific move.\n\n",
			styling.StyleBold.Render("USAGE:"),
			fmt.Sprintf("\n\t%s %s %s %s", "poke-cli", styling.StyleBold.Render("move"), "<move-name>", "[flag"),
			fmt.Sprintf("\n\t%-30s", styling.StyleItalic.Render("Use a hyphen when typing a name with a space.")),
			"\n\n",
			styling.StyleBold.Render("FLAGS:"),
			fmt.Sprintf("\n\t%-30s %s", "-p, --pokemon", "Prints Pokémon that learn this move."),
			fmt.Sprintf("\n\t%-30s %s", "-h, --help", "Prints the help menu."),
		)
		fmt.Println(helpMessage)
	}

	moveFlags, _, _ := flags.SetupMoveFlagSet()

	args := os.Args

	flag.Parse()

	if err := ValidateMoveArgs(args); err != nil {
		fmt.Println(err.Error())
		if os.Getenv("GO_TESTING") != "1" {
			os.Exit(1)
		}
	}

	endpoint := strings.ToLower(args[1])
	moveName := strings.ToLower(args[2])

	if err := moveFlags.Parse(args[3:]); err != nil {
		fmt.Printf("error parsing flags: %v\n", err)
		moveFlags.Usage()
		if os.Getenv("GO_TESTING") != "1" {
			os.Exit(1)
		}
	}

	moveStruct, moveName, err := connections.MoveApiCall(endpoint, moveName, connections.APIURL)
	if err != nil {
		fmt.Println(err)
		if os.Getenv("GO_TESTING") != "1" {
			os.Exit(1)
		}
	}

	// Extract English effect_entries
	var englishEffectEntry string
	for _, entry := range moveStruct.EffectEntries {
		if entry.Language.Name == "en" {
			englishEffectEntry = entry.Effect
			break
		}
	}

	capitalizedMove := cases.Title(language.English).String(strings.ReplaceAll(moveName, "_", " "))
	fmt.Println(styling.StyleBold.Render(capitalizedMove))

	fmt.Println(englishEffectEntry)
	fmt.Println(moveStruct.Type.Name)
	fmt.Println(moveStruct.Power)
}
