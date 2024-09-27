package cmd

import (
	"testing"
)

// TestValidatePokemonArgs tests the ValidatePokemonArgs function
func TestValidatePokemonArgs(t *testing.T) {

	// Test case: Too few arguments
	args := []string{"poke-cli", "pokemon"}
	expectedError := "╭────────────────────────────────────────────────────────────╮\n" +
		"│Error!                                                      │\n" +
		"│Please declare a Pokémon's name after the [pokemon] command │\n" +
		"│Run 'poke-cli pokemon -h' for more details                  │\n" +
		"│error: insufficient arguments                               │\n" +
		"╰────────────────────────────────────────────────────────────╯"
	err := ValidatePokemonArgs(args)
	if err == nil || err.Error() != expectedError {
		t.Errorf("Expected error for too few arguments, got: %v", err)
	}
}

// TestValidateTypesArgs tests the ValidateTypesArgs function
func TestValidateTypesArgs(t *testing.T) {
	// Test case: Help flag (-h)
	args := []string{"poke-cli", "types", "-h"}
	err := ValidateTypesArgs(args)
	if err == nil || err.Error() != "" {
		t.Errorf("Expected no error for help flag, got: %v", err)
	}

	// Test case: Valid args
	args = []string{"poke-cli", "types"}
	err = ValidateTypesArgs(args)
	if err != nil {
		t.Errorf("Expected no error for valid args, got: %v", err)
	}
}
