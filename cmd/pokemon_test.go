package cmd

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidatePokemonArgs_ValidInput(t *testing.T) {
	validInputs := [][]string{
		{"poke-cli", "pokemon", "pikachu"},
		{"poke-cli", "pokemon", "bulbasaur", "--types"},
		{"poke-cli", "pokemon", "cloyster", "--abilities"},
		{"poke-cli", "pokemon", "mewtwo", "--types", "--abilities"},
		{"poke-cli", "pokemon", "BlaZiKen", "-a", "-t"},
	}

	for _, input := range validInputs {
		err := ValidatePokemonArgs(input)
		assert.NoError(t, err, "Expected no error for valid input")
	}
}

func TestValidatePokemonArgs_InvalidFlag(t *testing.T) {
	invalidInputs := [][]string{
		{"poke-cli", "pokemon", "bulbasaur", "types"},
		{"poke-cli", "pokemon", "mewtwo", "--types", "abilities"},
	}
	expectedErrors := []string{
		"Error: Invalid argument 'types'. Only flags are allowed after declaring a Pokémon's name",
		"Error: Invalid argument 'abilities'. Only flags are allowed after declaring a Pokémon's name",
	}

	for i, input := range invalidInputs {
		err := ValidatePokemonArgs(input)
		assert.Error(t, err, "Expected error for invalid flag")
		assert.NotEmpty(t, expectedErrors[i], err.Error())
	}
}

func TestValidatePokemonArgs_TooManyArgs(t *testing.T) {
	invalidInput := [][]string{
		{"poke-cli", "pikachu", "--types", "all", "normal"},
	}
	expectedError := "error: too many arguments\n"

	for _, input := range invalidInput {
		err := ValidatePokemonArgs(input)
		assert.Error(t, err, "Expected error for too many arguments")
		assert.NotEqual(t, expectedError, err.Error())
	}
}
