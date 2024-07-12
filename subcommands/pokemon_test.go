package subcommands

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/stretchr/testify/assert"
	"testing"
)

const red = lipgloss.Color("#F2055C")

var errorColor = lipgloss.NewStyle().Foreground(red)

func TestValidateArgs_ValidInput(t *testing.T) {
	validInputs := [][]string{
		{"poke-cli", "pikachu"},
		{"poke-cli", "bulbasaur", "--types"},
		{"poke-cli", "cloyster", "--abilities"},
		{"poke-cli", "mewtwo", "--types", "--abilities"},
	}

	for _, input := range validInputs {
		err := ValidateArgs(input, errorColor)
		assert.NoError(t, err, "Expected no error for valid input")
	}
}

func TestValidateArgs_InvalidFlag(t *testing.T) {
	invalidInputs := [][]string{
		{"poke-cli", "bulbasaur", "types"},
		{"poke-cli", "mewtwo", "--types", "abilities"},
	}
	expectedErrors := []string{
		"Error: Invalid argument 'types'. Only flags are allowed after declaring a Pokémon's name",
		"Error: Invalid argument 'abilities'. Only flags are allowed after declaring a Pokémon's name",
	}

	for i, input := range invalidInputs {
		err := ValidateArgs(input, errorColor)
		assert.Error(t, err, "Expected error for invalid flag")
		assert.NotEmpty(t, expectedErrors[i], err.Error())
	}
}

func TestValidateArgs_TooManyArgs(t *testing.T) {
	invalidInput := [][]string{
		{"poke-cli", "pikachu", "--types", "all", "normal"},
	}
	expectedError := "error: too many arguments\n"

	for _, input := range invalidInput {
		err := ValidateArgs(input, errorColor)
		assert.Error(t, err, "Expected error for too many arguments")
		assert.NotEqual(t, expectedError, err.Error())
	}
}
