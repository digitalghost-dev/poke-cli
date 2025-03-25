package cmd

import (
	"flag"
	"github.com/digitalghost-dev/poke-cli/styling"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHandleHelpFlag(t *testing.T) {
	// Mock flag.Usage to avoid actual printing
	flag.Usage = func() {}

	// Test cases
	tests := []struct {
		name string
		args []string
	}{
		{"Valid short help flag", []string{"cmd", "subcmd", "-h"}},
		{"Valid long help flag", []string{"cmd", "subcmd", "--help"}},
		{"Invalid case (no flag)", []string{"cmd", "subcmd"}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			handleHelpFlag(tc.args)
		})
	}
}

func TestValidateAbilityArgs(t *testing.T) {
	// Testing valid arguments
	validInputs := [][]string{
		{"poke-cli", "ability", "--help"},
		{"poke-cli", "ability", "inner-focus"},
		{"poke-cli", "ability", "unaware", "-h"},
		{"poke-cli", "ability", "technician", "--pokemon"},
	}

	for _, input := range validInputs {
		err := ValidateAbilityArgs(input)
		assert.NoError(t, err, "Expected no error for valid input")
	}

	// Testing invalid arguments
	invalidInputs := [][]string{
		{"poke-cli", "abilities"},
	}

	for _, input := range invalidInputs {
		err := ValidateAbilityArgs(input)
		assert.Error(t, err, "Expected error for invalid input")
	}

	// Testing too many arguments
	tooManyArgs := [][]string{
		{"poke-cli", "ability", "strong-jaw", "all", "pokemon"},
	}

	expectedError := styling.StripANSI("╭──────────────────╮\n│Error!            │\n│Too many arguments│\n╰──────────────────╯")

	for _, input := range tooManyArgs {
		err := ValidateAbilityArgs(input)

		if err == nil {
			t.Fatalf("Expected an error for input %v, but got nil", input)
		}

		strippedErr := styling.StripANSI(err.Error())
		assert.Equal(t, expectedError, strippedErr, "Unexpected error message for invalid input")
	}
}

func TestValidateNaturesArgs(t *testing.T) {
	// Testing valid arguments
	validInputs := [][]string{
		{"poke-cli", "natures"},
		{"poke-cli", "natures", "--help"},
	}

	for _, input := range validInputs {
		err := ValidateNaturesArgs(input)
		assert.NoError(t, err, "Expected no error for valid input")
	}

	// Testing invalid arguments
	invalidInputs := [][]string{
		{"poke-cli", "natures", "docile"},
		{"poke-cli", "natures", "brave", "--help"},
	}

	for _, input := range invalidInputs {
		err := ValidateNaturesArgs(input)
		assert.Error(t, err, "Expected error for invalid input")
	}
}

// TestValidatePokemonArgs tests the ValidatePokemonArgs function
func TestValidatePokemonArgs(t *testing.T) {
	// Testing valid arguments
	validInputs := [][]string{
		{"poke-cli", "pokemon", "--help"},
		{"poke-cli", "pokemon", "mankey"},
		{"poke-cli", "pokemon", "talonflame", "--stats", "--types"},
		{"poke-cli", "pokemon", "passimian", "--abilities", "-t"},
		{"poke-cli", "pokemon", "dodrio", "-a", "-s", "-t"},
		{"poke-cli", "pokemon", "dragalge", "-a", "-s", "-t", "--image=sm"},
		{"poke-cli", "pokemon", "squirtle", "-a", "-s"},
		{"poke-cli", "pokemon", "squirtle", "-s", "-a"},
	}

	for _, input := range validInputs {
		err := ValidatePokemonArgs(input)
		assert.NoError(t, err, "Expected no error for valid input")
	}

	// Testing invalid arguments
	invalidInputs := [][]string{
		{"poke-cli"},
		{"poke-cli", "pokemon"},
		{"poke-cli", "pokemons"},
		{"poke-cli", "pokemon", "mewtwo", "--"},
		{"poke-cli", "pokemon", "baxcalibur", "-"},
		{"poke-cli", "pokemon", "charizard", "extraArg"},
	}

	for _, input := range invalidInputs {
		err := ValidatePokemonArgs(input)
		assert.Error(t, err, "Expected error for invalid input")
	}

	// Testing too many arguments
	tooManyArgs := [][]string{
		{"poke-cli", "pokemon", "hypno", "--abilities", "-s", "--types", "--image=sm", "-m"},
	}

	expectedError := styling.StripANSI("╭──────────────────╮\n│Error!            │\n│Too many arguments│\n╰──────────────────╯")

	for _, input := range tooManyArgs {
		err := ValidatePokemonArgs(input)

		if err == nil {
			t.Fatalf("Expected an error for input %v, but got nil", input)
		}

		strippedErr := styling.StripANSI(err.Error())
		assert.Equal(t, expectedError, strippedErr, "Unexpected error message for invalid input")
	}
}

// TestValidateSearchArgs tests the ValidateSearchArgs function
func TestValidateSearchArgs(t *testing.T) {
	validInputs := [][]string{
		{"poke-cli", "search"},
		{"poke-cli", "search", "--help"},
	}

	for _, input := range validInputs {
		err := ValidateSearchArgs(input)
		assert.NoError(t, err, "Expected no error for valid input")
	}

	invalidInputs := [][]string{
		{"poke-cli", "search", "pokemon"},
	}

	for _, input := range invalidInputs {
		err := ValidateSearchArgs(input)
		assert.Error(t, err, "Expected error for invalid input")
	}

	tooManyArgs := [][]string{
		{"poke-cli", "search", "pokemon", "meowscarada"},
	}

	expectedError := styling.StripANSI("╭──────────────────╮\n│Error!            │\n│Too many arguments│\n╰──────────────────╯")

	for _, input := range tooManyArgs {
		err := ValidateSearchArgs(input)

		if err == nil {
			t.Fatalf("Expected an error for input %v, but got nil", input)
		}

		strippedErr := styling.StripANSI(err.Error())
		assert.Equal(t, expectedError, strippedErr, "Unexpected error message for invalid input")
	}
}

// TestValidateTypesArgs tests the ValidateTypesArgs function
func TestValidateTypesArgs(t *testing.T) {
	// Testing valid arguments
	validInputs := [][]string{
		{"poke-cli", "types"},
		{"poke-cli", "types", "--help"},
	}

	for _, input := range validInputs {
		err := ValidateTypesArgs(input)
		assert.NoError(t, err, "Expected no error for valid input")
	}

	// Testing invalid arguments
	invalidInputs := [][]string{
		{"poke-cli", "types", "rock"},
	}

	for _, input := range invalidInputs {
		err := ValidateTypesArgs(input)
		assert.Error(t, err, "Expected error for invalid input")
	}

	// Testing too many arguments
	tooManyArgs := [][]string{
		{"poke-cli", "types", "rock", "pokemon"},
	}

	expectedError := styling.StripANSI("╭──────────────────╮\n│Error!            │\n│Too many arguments│\n╰──────────────────╯")

	for _, input := range tooManyArgs {
		err := ValidateTypesArgs(input)

		if err == nil {
			t.Fatalf("Expected an error for input %v, but got nil", input)
		}

		strippedErr := styling.StripANSI(err.Error())
		assert.Equal(t, expectedError, strippedErr, "Unexpected error message for invalid input")
	}
}
