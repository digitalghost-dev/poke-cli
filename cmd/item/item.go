package item

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/digitalghost-dev/poke-cli/cmd/utils"
	"github.com/digitalghost-dev/poke-cli/connections"
	"github.com/digitalghost-dev/poke-cli/constants"
	"github.com/digitalghost-dev/poke-cli/structs"
	"github.com/digitalghost-dev/poke-cli/styling"
)

func ItemCommand() (string, error) {
	var output strings.Builder

	flag.Usage = func() {
		output.WriteString(
			utils.GenerateHelpMessage(
				utils.HelpConfig{
					Description:    "Get details about a specific item.",
					CmdName:        "item",
					SubCmdName:     "<item-name>",
					ShowHyphenHint: true,
				},
			),
		)
	}

	args := os.Args

	if utils.CheckHelpFlag(&output, flag.Usage) {
		return output.String(), nil
	}

	flag.Parse()

	if err := utils.ValidateArgs(os.Args, utils.Validator{MaxArgs: 3, CmdName: "item", RequireName: true, HasFlags: false}); err != nil {
		output.WriteString(err.Error())
		return output.String(), err
	}

	endpoint := strings.ToLower(args[1])
	itemName := strings.ToLower(args[2])

	itemStruct, itemName, err := connections.ItemApiCall(endpoint, itemName, connections.APIURL)
	if err != nil {
		output.WriteString(err.Error())
		return output.String(), err
	}

	itemInfoContainer(&output, itemStruct, itemName)

	return output.String(), nil
}

func itemInfoContainer(output *strings.Builder, itemStruct structs.ItemJSONStruct, itemName string) {
	capitalizedItem := styling.StyleBold.Render(styling.CapitalizeResourceName(itemName))
	itemCost := fmt.Sprintf("Cost: %d", itemStruct.Cost)
	itemCategory := "Category: " + styling.CapitalizeResourceName(itemStruct.Category.Name)

	isDark := lipgloss.HasDarkBackground(os.Stdin, os.Stdout)
	ld := lipgloss.LightDark(isDark)
	docStyle := lipgloss.NewStyle().
		Padding(1, 2).
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(ld(lipgloss.Color("#444"), lipgloss.Color("#EEE"))).
		Width(34)

	var flavorTextEntry string
	var missingData string
	var fullDoc string

	if len(itemStruct.FlavorTextEntries) == 0 {
		missingData = styling.StyleItalic.Render("Missing data from API")
		fullDoc = lipgloss.JoinVertical(lipgloss.Top, capitalizedItem, itemCost, itemCategory, "---", "Description:", missingData)
	} else {
		for _, entry := range itemStruct.FlavorTextEntries {
			if entry.Language.Name == "en" && entry.VersionGroup.Name == constants.VersionSwordShield {
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
