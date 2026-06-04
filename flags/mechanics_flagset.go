package flags

import (
	"flag"
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/table"
	"github.com/digitalghost-dev/poke-cli/styling"
)

type MechanicsFlags struct {
	FlagSet      *flag.FlagSet
	Natures      *bool
	ShortNatures *bool
}

func SetupMechanicsFlagSet() *MechanicsFlags {
	mf := &MechanicsFlags{}
	mf.FlagSet = flag.NewFlagSet("mechanicsFlags", flag.ContinueOnError)

	mf.Natures = mf.FlagSet.Bool("natures", false, "Show a table with natures.")
	mf.ShortNatures = mf.FlagSet.Bool("n", false, "Show a table with natures.")

	mf.FlagSet.Usage = func() {
		helpMessage := styling.HelpBorder.Render("poke-cli mechanics [flags]\n\n",
			styling.StyleBold.Render("FLAGS:"),
			fmt.Sprintf("\n\t%-30s %s", "-n, --natures", "Show a table with natures."),
		)
		fmt.Println(helpMessage)
	}

	return mf
}

func NaturesFlag() string {
	var output strings.Builder

	output.WriteString("Natures affect the growth of a Pokémon.\n" +
		"Each nature increases one of its stats by 10% and decreases one by 10%.\n" +
		"Five natures increase and decrease the same stat and therefore have no effect.\n\n")
	output.WriteString(styling.StyleBold.Render("Nature Chart:"))
	output.WriteString("\n")

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
				Padding(0, 1)
		})

	output.WriteString(t.Render())
	output.WriteString("\n")

	return output.String()
}
