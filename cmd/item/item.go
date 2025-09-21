package item

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/digitalghost-dev/poke-cli/cmd/utils"
	"github.com/digitalghost-dev/poke-cli/connections"
	"github.com/digitalghost-dev/poke-cli/structs"
	"github.com/digitalghost-dev/poke-cli/styling"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func ItemCommand() (string, error) {
	var output strings.Builder

	flag.Usage = func() {
		helpMessage := styling.HelpBorder.Render(
			"Get details about a specific item.\n\n",
			styling.StyleBold.Render("USAGE:"),
			fmt.Sprintf("\n\t%s %s %s %s", "poke-cli", styling.StyleBold.Render("item"), "<item-name>", "[flag]"),
			fmt.Sprintf("\n\t%-30s", styling.StyleItalic.Render("Use a hyphen when typing a name with a space.")),
			"\n\n",
			styling.StyleBold.Render("FLAGS:"),
			fmt.Sprintf("\n\t%-30s %s", "-h, --help", "Prints the help menu."),
		)
		output.WriteString(helpMessage)
	}

	args := os.Args

	flag.Parse()

	if len(os.Args) == 3 && (os.Args[2] == "-h" || os.Args[2] == "--help") {
		flag.Usage()
		return output.String(), nil
	}

	if err := utils.ValidateItemArgs(os.Args); err != nil {
		output.WriteString(err.Error())
		return output.String(), err
	}

	endpoint := strings.ToLower(args[1])
	itemName := strings.ToLower(args[2])

	itemStruct, itemName, err := connections.ItemApiCall(endpoint, itemName, connections.APIURL)
	if err != nil {
		if os.Getenv("GO_TESTING") != "1" {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		return err.Error(), nil
	}

	itemInfoContainer(&output, itemStruct, itemName)

	return output.String(), nil
}

func itemInfoContainer(output *strings.Builder, itemStruct structs.ItemJSONStruct, itemName string) {
	capitalizedItem := styling.StyleBold.Render(cases.Title(language.English).String(strings.ReplaceAll(itemName, "-", " ")))
	itemCost := fmt.Sprintf("Cost: %d", itemStruct.Cost)
	itemCategory := "Category: " + cases.Title(language.English).String(strings.ReplaceAll(itemStruct.Category.Name, "-", " "))

	docStyle := lipgloss.NewStyle().
		Padding(1, 2).
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(lipgloss.AdaptiveColor{Light: "#444", Dark: "#EEE"}).
		Width(32)

	var flavorTextEntry string
	var missingData string
	var fullDoc string

	if len(itemStruct.FlavorTextEntries) == 0 {
		missingData = styling.StyleItalic.Render("Missing data from API")
		fullDoc = lipgloss.JoinVertical(lipgloss.Top, capitalizedItem, itemCost, itemCategory, "---", "Description:", missingData)
	} else {
		for _, entry := range itemStruct.FlavorTextEntries {
			if entry.Language.Name == "en" && entry.VersionGroup.Name == "sword-shield" {
				if entry.Text != "" {
					flavorTextEntry = entry.Text
					fullDoc = lipgloss.JoinVertical(lipgloss.Top, capitalizedItem, itemCost, itemCategory, "---", "Description:", flavorTextEntry)
					break
				}
			}
		}
	}

	output.WriteString(docStyle.Render(fullDoc))
	output.WriteString("\n")
}
