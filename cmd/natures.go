package cmd

import (
	"flag"
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/digitalghost-dev/poke-cli/styling"
	"os"
)

func NaturesCommand() {
	flag.Usage = func() {
		helpMessage := styling.HelpBorder.Render(
			"Get details about all natures.\n\n",
			styling.StyleBold.Render("USAGE:"),
			fmt.Sprintf("\n\t%s %s", "poke-cli", styling.StyleBold.Render("natures")),
		)
		fmt.Println(helpMessage)
	}

	flag.Parse()

	// Check for help flag
	for _, arg := range os.Args[1:] {
		if arg == "-h" || arg == "--help" {
			flag.Usage()
			
			return
		}
	}

	if err := ValidateNaturesArgs(os.Args); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Println("Natures affect the growth of a Pokémon.\n" +
		"Each nature increases one of its stats by 10% and decreases one by 10%.\n" +
		"Five natures increase and decrease the same stat and therefore have no effect.\n\n" +
		styling.StyleBold.Render("Nature Chart:"))

	chart := [][]string{
		{" ", styling.Red.Render("-Attack"), styling.Red.Render("-Defense"), styling.Red.Render("-Sp. Atk"), styling.Red.Render("-Sp. Def"), styling.Red.Render("Speed")},
		{styling.Green.Render("+Attack"), "Hardy", "Lonely", "Adamant", "Naughty", "Brave"},
		{styling.Green.Render("+Defense"), "Bold", "Docile", "Impish", "Lax", "Relaxed"},
		{styling.Green.Render("+Sp. Atk"), "Modest", "Mild", "Bashful", "Rash", "Quiet"},
		{styling.Green.Render("+Sp. Def"), "Calm", "Gentle", "Careful", "Quirky", "Sassy"},
		{styling.Green.Render("Speed"), "Timid", "Hasty", "Jolly", "Naive", "Serious"},
	}

	t := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(styling.Gray)).
		BorderRow(true).
		BorderColumn(true).
		Rows(chart...).
		StyleFunc(func(row, col int) lipgloss.Style {
			return lipgloss.NewStyle().
				Padding(0, 1) // This styles the border color
		})

	fmt.Println(t.Render())
}
