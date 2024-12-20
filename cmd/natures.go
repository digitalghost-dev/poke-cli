package cmd

import (
	"flag"
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"os"
)

func NaturesCommand() {

	flag.Usage = func() {
		helpMessage := helpBorder.Render(
			"Get details about Pokémon natures.\n\n",
			styleBold.Render("USAGE:"),
			fmt.Sprintf("\n\t%s %s %s", "poke-cli", styleBold.Render("natures"), "[flag]"),
			"\n\n",
			styleBold.Render("FLAGS:"),
			fmt.Sprintf("\n\t%-30s %s", "-h, --help", "Prints out the help menu."),
		)
		fmt.Println(helpMessage)
	}

	flag.Parse()

	if err := ValidateNaturesArgs(os.Args); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Println("Natures affect the growth of a Pokémon.\n" +
		"Each nature increases one of its stats by 10% and decreases one by 10%.\n" +
		"Five natures increase and decrease the same stat and therefore have no effect.\n\n" +
		styleBold.Render("Nature Chart:"))

	chart := [][]string{
		{" ", red.Render("-Attack"), red.Render("-Defense"), red.Render("-Sp. Atk"), red.Render("-Sp. Def"), red.Render("Speed")},
		{green.Render("+Attack"), "Hardy", "Lonely", "Adamant", "Naughty", "Brave"},
		{green.Render("+Defense"), "Bold", "Docile", "Impish", "Lax", "Relaxed"},
		{green.Render("+Sp. Atk"), "Modest", "Mild", "Bashful", "Rash", "Quiet"},
		{green.Render("+Sp. Def"), "Calm", "Gentle", "Careful", "Quirky", "Sassy"},
		{green.Render("Speed"), "Timid", "Hasty", "Jolly", "Naive", "Serious"},
	}

	t := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(gray)).
		BorderRow(true).
		BorderColumn(true).
		Rows(chart...).
		StyleFunc(func(row, col int) lipgloss.Style {
			return lipgloss.NewStyle().
				Padding(0, 1) // This styles the border color
		})

	fmt.Println(t.Render())
}
