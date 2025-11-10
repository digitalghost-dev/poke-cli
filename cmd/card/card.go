package card

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/digitalghost-dev/poke-cli/cmd/utils"
	"github.com/digitalghost-dev/poke-cli/styling"
)

func CardCommand() (string, error) {
	var output strings.Builder

	flag.Usage = func() {
		helpMessage := styling.HelpBorder.Render(
			"View data about cards from the TCG!\n\n",
			styling.StyleBold.Render("USAGE:"),
			fmt.Sprintf("\n\t%s %s %s", "poke-cli", styling.StyleBold.Render("card"), "[flag]"),
			"\n\n",
			styling.StyleBold.Render("FLAGS:"),
			fmt.Sprintf("\n\t%-30s %s", "-h, --help", "Prints out the help menu"),
		)
		output.WriteString(helpMessage)
	}

	flag.Parse()

	// Handle help flag
	if len(os.Args) == 3 && (os.Args[2] == "-h" || os.Args[2] == "--help") {
		flag.Usage()
		return output.String(), nil
	}

	// Validate arguments
	if err := utils.ValidateCardArgs(os.Args); err != nil {
		output.WriteString(err.Error())
		return output.String(), err
	}

	items := []list.Item{
		item("Mega Evolution"),
		item("Scarlet & Violet"),
		item("Sword & Shield"),
	}

	const listWidth = 20
	const listHeight = 12

	l := list.New(items, itemDelegate{}, listWidth, listHeight)
	l.Title = "First, pick a series"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	m := SeriesModel{List: l}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	return output.String(), nil
}
