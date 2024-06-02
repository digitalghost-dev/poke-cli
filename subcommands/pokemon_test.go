package subcommands

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidateArgs_ValidInput(t *testing.T) {
	validInputs := [][]string{
		{"pikachu"},
		{"bulbasaur", "--types"},
	}

	for _, input := range validInputs {
		err := ValidateArgs(input)
		assert.NoError(t, err, "Expected no error for valid input")
	}
}

func TestValidateArgs_InvalidFlag(t *testing.T) {
	invalidInputs := [][]string{
		{"poke-cli", "pikachu", "wartortle"},
		{"poke-cli", "bulbasaur", "types"},
	}

	for _, input := range invalidInputs {
		err := ValidateArgs(input)
		assert.Error(t, err, "Expected error for invalid flag")
	}
}

func TestValidateArgs_TooManyArgs(t *testing.T) {
	invalidInput := [][]string{
		{"mewtwo", "--types", "all"},
		{"pikachu", "--types", "all", "normal"},
	}
	expectedError := "error: too many arguments\n"

	for _, input := range invalidInput {
		err := ValidateArgs(input)
		assert.Error(t, err, "Expected error for too many arguments")
		assert.NotEqual(t, expectedError, err.Error())
	}
}
