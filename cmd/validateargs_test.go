package cmd

import (
	"regexp"
	"strings"
	"testing"
)

func stripANSI(input string) string {
	ansiEscape := regexp.MustCompile(`\x1b\[[0-9;]*m`)
	return ansiEscape.ReplaceAllString(input, "")
}

// TestValidatePokemonArgs tests the ValidatePokemonArgs function
func TestValidatePokemonArgs(t *testing.T) {
	args := []string{"poke-cli", "pokemon"}
	expectedError := "╭────────────────────────────────────────────────────────────╮\n" +
		"│Error!                                                      │\n" +
		"│Please declare a Pokémon's name after the [pokemon] command │\n" +
		"│Run 'poke-cli pokemon -h' for more details                  │\n" +
		"│error: insufficient arguments                               │\n" +
		"╰────────────────────────────────────────────────────────────╯"
	err := ValidatePokemonArgs(args)
	if err == nil {
		t.Errorf("Expected error for too few arguments, got nil")
		return
	}

	// Strip ANSI codes for comparison
	actualError := stripANSI(err.Error())
	expectedError = strings.TrimSpace(expectedError)

	if actualError != expectedError {
		t.Errorf("Expected error:\n%s\nGot:\n%s", expectedError, actualError)
	}
}

// TestValidateTypesArgs tests the ValidateTypesArgs function
func TestValidateTypesArgs(t *testing.T) {
	// Test case: Help flag (-h)
	args := []string{"poke-cli", "types", "-h"}
	err := ValidateTypesArgs(args)
	if err != nil {
		t.Errorf("Expected no error for help flag, but got: %v", err)
	}

	// Test case: Valid args (e.g., no subcommands or flags after 'types')
	args = []string{"poke-cli", "types"}
	err = ValidateTypesArgs(args)
	// Ensure no error is returned for valid arguments
	if err != nil {
		t.Errorf("Expected no error for valid args, got: %v", err)
	}
}
