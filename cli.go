package main

import (
	"fmt"
	"os"
	"runtime/debug"
	"slices"
	"strings"

	"github.com/digitalghost-dev/poke-cli/cmd/ability"
	"github.com/digitalghost-dev/poke-cli/cmd/berry"
	"github.com/digitalghost-dev/poke-cli/cmd/card"
	"github.com/digitalghost-dev/poke-cli/cmd/comp"
	"github.com/digitalghost-dev/poke-cli/cmd/item"
	"github.com/digitalghost-dev/poke-cli/cmd/mechanics"
	"github.com/digitalghost-dev/poke-cli/cmd/move"
	"github.com/digitalghost-dev/poke-cli/cmd/pokemon"
	"github.com/digitalghost-dev/poke-cli/cmd/search"
	"github.com/digitalghost-dev/poke-cli/cmd/speed"
	"github.com/digitalghost-dev/poke-cli/cmd/types"
	"github.com/digitalghost-dev/poke-cli/cmd/utils"
	"github.com/digitalghost-dev/poke-cli/connections"
	"github.com/digitalghost-dev/poke-cli/flags"
	"github.com/digitalghost-dev/poke-cli/setup"
	"github.com/digitalghost-dev/poke-cli/styling"
	flag "github.com/spf13/pflag"
	"golang.org/x/term"
)

var version = "(devel)"

var commandDescriptions = []struct {
	name string
	desc string
}{
	{"ability", "Get details about an ability"},
	{"berry", "Get details about a berry"},
	{"card", "Get details about a TCG card"},
	{"comp", "Get details about competitive Pokémon"},
	{"item", "Get details about an item"},
	{"mechanics", "Get details about video game mechanics"},
	{"move", "Get details about a move"},
	{"pokemon", "Get details about a Pokémon"},
	{"search", "Search for a resource"},
	{"speed", "Calculate the speed of a Pokémon in battle"},
	{"types", "Get details about a typing"},
}

type mainFlags struct {
	FlagSet *flag.FlagSet
	config  *bool
	latest  *bool
	version *bool
}

func renderCommandList() string {
	var sb strings.Builder
	for _, cmd := range commandDescriptions {
		fmt.Fprintf(&sb, "\n\t%-15s %s", cmd.name, cmd.desc)
	}
	return sb.String()
}

func currentVersion() string {
	if version != "(devel)" {
		// Use version injected by -ldflags
		return "Version: " + version
	}

	// Fallback to build info when the version is not set
	buildInfo, ok := debug.ReadBuildInfo()
	if !ok {
		return "Version: unknown (unable to read build info)"
	}

	if buildInfo.Main.Version != "" {
		return "Version: " + buildInfo.Main.Version
	}
	return "Version: (devel)"
}

func setupMainFlagSet() *mainFlags {
	f := &mainFlags{}
	f.FlagSet = flag.NewFlagSet("mainFlags", flag.ContinueOnError)
	f.FlagSet.SetInterspersed(false)

	f.config = f.FlagSet.BoolP("config", "c", false, "Launch the config settings screen")
	f.latest = f.FlagSet.BoolP("latest", "l", false, "Prints the latest version available")
	f.version = f.FlagSet.BoolP("version", "v", false, "Prints the current version")

	f.FlagSet.Usage = func() {
		helpMessage := styling.HelpBorder.Render(
			"Welcome! This tool displays data related to Pokémon!",
			"\n\n", styling.StyleBold.Render("USAGE:"),
			fmt.Sprintf("\n\t%-15s %s", "poke-cli [flag]", ""),
			fmt.Sprintf("\n\t%-15s %s", "poke-cli <command> [flag]", ""),
			fmt.Sprintf("\n\t%-15s %s", "poke-cli <command> <subcommand> [flag]", ""),
			"\n\n", styling.StyleBold.Render("FLAGS:"),
			fmt.Sprintf("\n\t%-15s %s", "-h, --help", "Shows the help menu"),
			fmt.Sprintf("\n\t%-15s %s", "-c, --config", "Launch the config settings screen"),
			fmt.Sprintf("\n\t%-15s %s", "-l, --latest", "Prints the latest version available"),
			fmt.Sprintf("\n\t%-15s %s", "-v, --version", "Prints the current version"),
			"\n\n", styling.StyleBold.Render("COMMANDS:"),
			renderCommandList(),
			"\n\n", styling.StyleItalic.Render(styling.HyphenHint),
			"\n", styling.StyleItalic.Render("example: poke-cli ability strong-jaw"),
			"\n", styling.StyleItalic.Render("example: poke-cli pokemon flutter-mane"),
			"\n\n", fmt.Sprintf("%s %s", "↓ ctrl/cmd + click for docs/guides\n", styling.DocsLink),
		)
		fmt.Println(helpMessage)
	}

	return f
}

func runCLI(args []string) int {
	var output strings.Builder

	cfg, firstRun, err := flags.Load()
	if err != nil {
		cfg = flags.Defaults()
	}

	styling.ApplyTheme(cfg.Display.Theme)

	wantsConfig := slices.Contains(args, "--config") || slices.Contains(args, "-c")
	if firstRun && !wantsConfig && isInteractive() {
		if updated, saved, runErr := setup.Run(cfg); runErr == nil && saved {
			cfg = updated
			saveConfig(cfg)
		}
	}

	connections.ConfigureCache(cfg.Cache.ShowWarning, cfg.Cache.Path)

	f := setupMainFlagSet()

	switch {
	case len(args) == 0:
		f.FlagSet.Usage()
		return 0
	case len(args) > 0:
		if args[0] == "-h" || args[0] == "--help" {
			f.FlagSet.Usage()
			return 0
		}
	}

	err = f.FlagSet.Parse(args)
	if err != nil {
		return 2
	}

	remainingArgs := f.FlagSet.Args()

	type commandFunc func([]string) (string, error)
	commands := map[string]commandFunc{
		"ability":   ability.AbilityCommand,
		"berry":     berry.BerryCommand,
		"card":      card.CardCommand,
		"comp":      comp.CompCommand,
		"item":      item.ItemCommand,
		"mechanics": mechanics.MechanicsCommand,
		"move":      move.MoveCommand,
		"pokemon":   pokemon.PokemonCommand,
		"search":    search.SearchCommand,
		"speed":     speed.SpeedCommand,
		"types":     types.TypesCommand,
	}

	cmdArg := ""
	if len(remainingArgs) >= 1 {
		cmdArg = remainingArgs[0]
	}
	cmdFunc, exists := commands[cmdArg]

	switch {
	case len(remainingArgs) == 0 && !*f.latest && !*f.version && !*f.config:
		f.FlagSet.Usage()
		return 1
	case *f.latest:
		_, err := flags.LatestFlag()
		if err != nil {
			return 1
		}
		return 0
	case *f.version:
		fmt.Println(currentVersion())
		return 0
	case *f.config:
		updated, saved, runErr := setup.Run(cfg)
		if runErr != nil {
			return 1
		}
		if saved {
			saveConfig(updated)
		}
		return 0
	case exists:
		return utils.HandleCommandOutput(cmdFunc, remainingArgs)()
	default:
		msg := fmt.Sprintf("\t%-15s", fmt.Sprintf("'%s' is not a valid command.\n", cmdArg)) +
			styling.StyleBold.Render("\nCommands:") +
			renderCommandList() +
			fmt.Sprintf("\n\nAlso run %s for more info!", styling.StyleBold.Render("poke-cli -h"))
		output.WriteString(utils.FormatError(msg))

		fmt.Println(output.String())

		return 1
	}
}

var exit = os.Exit

func isInteractive() bool {
	return term.IsTerminal(int(os.Stdin.Fd())) && term.IsTerminal(int(os.Stdout.Fd())) // #nosec G115
}

func saveConfig(cfg flags.Config) {
	if err := flags.Save(cfg); err != nil {
		fmt.Fprintln(os.Stderr, styling.WarningColor.Render(
			"Couldn't save settings to the config file (it may be open elsewhere); changes won't persist."))
	}
}

func main() {
	exit(runCLI(os.Args[1:]))
}
