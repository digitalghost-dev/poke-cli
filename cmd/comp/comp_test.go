package comp

import (
	"strings"
	"testing"

	"github.com/digitalghost-dev/poke-cli/styling"
)

func TestCompCommand_Help(t *testing.T) {
	for _, flag := range []string{"-h", "--help"} {
		t.Run(flag, func(t *testing.T) {
			output, err := CompCommand([]string{"comp", flag})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			clean := styling.StripANSI(output)
			for _, want := range []string{"USAGE:", "comp", "competitive"} {
				if !strings.Contains(clean, want) {
					t.Errorf("expected help output to contain %q, got:\n%s", want, clean)
				}
			}
		})
	}
}

func TestCompCommand_TooManyArgs(t *testing.T) {
	output, err := CompCommand([]string{"comp", "one", "two"})
	if err == nil {
		t.Fatal("expected error for too many arguments")
	}
	if !strings.Contains(styling.StripANSI(output), "Too many arguments") {
		t.Errorf("expected 'Too many arguments' in output, got:\n%s", styling.StripANSI(output))
	}
}

func TestCompCommand_InvalidOption(t *testing.T) {
	output, err := CompCommand([]string{"comp", "bogus"})
	if err == nil {
		t.Fatal("expected error for invalid option")
	}
	if !strings.Contains(styling.StripANSI(output), "only available options") {
		t.Errorf("expected invalid-option error in output, got:\n%s", styling.StripANSI(output))
	}
}
