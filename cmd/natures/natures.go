package natures

import (
	"flag"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/digitalghost-dev/poke-cli/cmd/utils"
	"github.com/digitalghost-dev/poke-cli/styling"
)

func NaturesCommand() (string, error) {
	var output strings.Builder

	flag.Usage = func() {
		output.WriteString(
			utils.GenerateHelpMessage(
				utils.HelpConfig{
					Description: "Get details about all natures.",
					CmdName:     "natures",
				},
			),
		)
	}

	if utils.CheckHelpFlag(&output, flag.Usage) {
		return output.String(), nil
	}

	flag.Parse()

	if err := utils.ValidateArgs(os.Args, utils.Validator{MaxArgs: 3, CmdName: "natures", RequireName: false, HasFlags: false}); err != nil {
		output.WriteString(err.Error())
		return output.String(), err
	}

	output.WriteString("Natures affect the growth of a Pokémon.\n" +
		"Each nature increases one of its stats by 10% and decreases one by 10%.\n" +
		"Five natures increase and decrease the same stat and therefore have no effect.\n\n" +
		styling.StyleBold.Render("Nature Chart:") + "\n")

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

	output.WriteString(t.Render() + "\n")

	return output.String(), nil
}
