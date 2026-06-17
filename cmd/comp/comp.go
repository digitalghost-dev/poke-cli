package comp

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/digitalghost-dev/poke-cli/cmd/comp/tcg"
	"github.com/digitalghost-dev/poke-cli/cmd/comp/vgc"
	"github.com/digitalghost-dev/poke-cli/cmd/comp/champions"
	"github.com/digitalghost-dev/poke-cli/cmd/utils"
)

func CompCommand(args []string) (string, error) {
	var output strings.Builder

	usage := func() {
		output.WriteString(
			utils.GenerateHelpMessage(
				utils.HelpConfig{
					Description: "Get details about competitive Pokémon.",
					CmdName:     "comp",
				},
			),
		)
	}

	if utils.CheckHelpFlag(args, usage) {
		return output.String(), nil
	}

	// Validate arguments
	if err := utils.ValidateArgs(
		args,
		utils.Validator{MaxArgs: 2, CmdName: "comp", RequireName: false, HasFlags: false},
	); err != nil {
		output.WriteString(err.Error())
		return output.String(), err
	}

	// Program 1: Competition type selection
	for {
		finalModel, err := tea.NewProgram(CompList()).Run()
		if err != nil {
			return "", fmt.Errorf("error running comp selection program: %w", err)
		}

		result, ok := finalModel.(pickerModel)
		if !ok {
			return "", fmt.Errorf("unexpected model type from competition selection: got %T, want compModel", finalModel)
		}

		if result.CompID == "" {
			break
		}

		var back bool

		switch result.CompID {
		case "tcg":
			back, err = tcg.Run()
		case "vgc":
			back, err = vgc.Run()
		case "champions":
			back, err = champions.Run()
		}
		if err != nil {
			return "", err
		}
		if !back {
			break
		}
	}

	return output.String(), nil
}
