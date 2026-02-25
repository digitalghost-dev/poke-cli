package utils

import (
	"strings"
	"testing"

	"github.com/digitalghost-dev/poke-cli/styling"
)

func TestGenerateHelpMessage_ContainsDescription(t *testing.T) {
	cfg := HelpConfig{
		Description: "Get details about a specific item.",
		CmdName:     "item",
		SubCmdName:  "<item-name>",
	}

	output := styling.StripANSI(GenerateHelpMessage(cfg))

	if !strings.Contains(output, cfg.Description) {
		t.Errorf("expected output to contain description %q, got:\n%s", cfg.Description, output)
	}
}

func TestGenerateHelpMessage_ContainsUsage(t *testing.T) {
	cfg := HelpConfig{
		Description: "Get details about a specific item.",
		CmdName:     "item",
		SubCmdName:  "<item-name>",
	}

	output := styling.StripANSI(GenerateHelpMessage(cfg))

	for _, want := range []string{"USAGE:", "poke-cli", cfg.CmdName, cfg.SubCmdName} {
		if !strings.Contains(output, want) {
			t.Errorf("expected output to contain %q, got:\n%s", want, output)
		}
	}
}

func TestGenerateHelpMessage_AlwaysContainsHelpFlag(t *testing.T) {
	cfg := HelpConfig{
		Description: "Some command.",
		CmdName:     "berry",
		SubCmdName:  "<berry-name>",
	}

	output := styling.StripANSI(GenerateHelpMessage(cfg))

	for _, want := range []string{"-h, --help", "Prints the help menu."} {
		if !strings.Contains(output, want) {
			t.Errorf("expected output to always contain %q, got:\n%s", want, output)
		}
	}
}

func TestGenerateHelpMessage_HyphenHint(t *testing.T) {
	tests := []struct {
		name            string
		showHyphenHint  bool
		wantHyphenHint  bool
	}{
		{
			name:           "hyphen hint shown when enabled",
			showHyphenHint: true,
			wantHyphenHint: true,
		},
		{
			name:           "hyphen hint hidden when disabled",
			showHyphenHint: false,
			wantHyphenHint: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := HelpConfig{
				Description:    "Get details about a specific move.",
				CmdName:        "move",
				SubCmdName:     "<move-name>",
				ShowHyphenHint: tt.showHyphenHint,
			}

			output := styling.StripANSI(GenerateHelpMessage(cfg))
			contains := strings.Contains(output, styling.HyphenHint)

			if contains != tt.wantHyphenHint {
				t.Errorf("HyphenHint presence = %v, want %v\noutput:\n%s", contains, tt.wantHyphenHint, output)
			}
		})
	}
}

func TestGenerateHelpMessage_CustomFlags(t *testing.T) {
	cfg := HelpConfig{
		Description: "Get details about a Pokémon.",
		CmdName:     "pokemon",
		SubCmdName:  "<pokemon-name>",
		Flags: []FlagHelp{
			{Short: "-a", Long: "--abilities", Description: "Prints the Pokémon's abilities."},
			{Short: "-s", Long: "--stats", Description: "Prints the Pokémon's base stats."},
		},
	}

	output := styling.StripANSI(GenerateHelpMessage(cfg))

	for _, f := range cfg.Flags {
		for _, want := range []string{f.Short, f.Long, f.Description} {
			if !strings.Contains(output, want) {
				t.Errorf("expected output to contain flag field %q, got:\n%s", want, output)
			}
		}
	}
}

func TestGenerateHelpMessage_NoCustomFlags(t *testing.T) {
	cfg := HelpConfig{
		Description: "Get details about a berry.",
		CmdName:     "berry",
		SubCmdName:  "<berry-name>",
	}

	output := styling.StripANSI(GenerateHelpMessage(cfg))

	// Should still have the hardcoded help flag
	if !strings.Contains(output, "-h, --help") {
		t.Errorf("expected output to contain '-h, --help', got:\n%s", output)
	}
}

func TestGenerateHelpMessage_MultipleFlags(t *testing.T) {
	flags := []FlagHelp{
		{Short: "-f", Long: "--flag-one", Description: "First flag."},
		{Short: "-g", Long: "--flag-two", Description: "Second flag."},
		{Short: "-x", Long: "--flag-three", Description: "Third flag."},
	}

	cfg := HelpConfig{
		Description: "A command with many flags.",
		CmdName:     "pokemon",
		SubCmdName:  "<pokemon-name>",
		Flags:       flags,
	}

	output := styling.StripANSI(GenerateHelpMessage(cfg))

	for _, f := range flags {
		if !strings.Contains(output, f.Short) {
			t.Errorf("expected output to contain short flag %q", f.Short)
		}
		if !strings.Contains(output, f.Long) {
			t.Errorf("expected output to contain long flag %q", f.Long)
		}
		if !strings.Contains(output, f.Description) {
			t.Errorf("expected output to contain flag description %q", f.Description)
		}
	}
}