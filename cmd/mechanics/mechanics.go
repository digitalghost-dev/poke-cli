package mechanics

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/digitalghost-dev/poke-cli/cmd/utils"
	"github.com/digitalghost-dev/poke-cli/flags"
)

func MechanicsCommand(args []string) (string, error) {
	var output strings.Builder

	usage := func() {
		output.WriteString(
			utils.GenerateHelpMessage(
				utils.HelpConfig{
					Description: "Get details about game mechanics.",
					CmdName:     "mechanics",
					Flags: []utils.FlagHelp{
						{Short: "-n", Long: "--natures", Description: "Prints a table with all natures and their respective buffs and debuffs."},
					},
				},
			),
		)
	}

	mf := flags.SetupMechanicsFlagSet()

	if utils.CheckHelpFlag(args, usage) {
		return output.String(), nil
	}

	if err := utils.ValidateArgs(
		args,
		utils.Validator{MaxArgs: 2, CmdName: "mechanics", RequireName: false, HasFlags: true},
	); err != nil {
		output.WriteString(err.Error())
		return output.String(), err
	}

	if err := mf.FlagSet.Parse(args[1:]); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return output.String(), nil
		}
		fmt.Fprintf(&output, "error parsing flags: %v\n", err)
		return output.String(), err
	}

	switch {
	case *mf.Natures || *mf.ShortNatures:
		output.WriteString(flags.NaturesFlag())
	default:
		usage()
	}

	return output.String(), nil
}
