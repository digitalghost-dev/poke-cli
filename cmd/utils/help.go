package utils

import (
	"fmt"
	"strings"

	"github.com/digitalghost-dev/poke-cli/styling"
)

type HelpConfig struct {
	Description    string
	CmdName        string
	SubCmdName     string
	Flags          []FlagHelp
	ShowHyphenHint bool
}

type FlagHelp struct {
	Short, Long string
	Description string
}

func GenerateHelpMessage(cfg HelpConfig) string {
	var flagsBuilder strings.Builder
	for _, f := range cfg.Flags {
		fmt.Fprintf(&flagsBuilder, "\n\t%-30s %s", f.Short+", "+f.Long, f.Description)
	}
	flagsList := flagsBuilder.String()

	hyphenHint := ""
	if cfg.ShowHyphenHint {
		hyphenHint = fmt.Sprintf("\n\t%-30s", styling.StyleItalic.Render(styling.HyphenHint))
	}

	helpMessage := styling.HelpBorder.Render(
		cfg.Description+"\n\n",
		styling.StyleBold.Render("USAGE:"),
		fmt.Sprintf("\n\t%s %s %s", "poke-cli", styling.StyleBold.Render(cfg.CmdName), cfg.SubCmdName),
		hyphenHint,
		"\n\n",
		styling.StyleBold.Render("FLAGS:"),
		fmt.Sprintf("\n\t%-30s %s", "-h, --help", "Prints the help menu."),
		flagsList,
	)

	return helpMessage
}
