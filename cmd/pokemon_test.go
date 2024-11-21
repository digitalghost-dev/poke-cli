package cmd

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"os"
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

func TestPokemonCommand(t *testing.T) {
	// Capture standard output
	var output bytes.Buffer
	stdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Set up test arguments (focus only on Pokémon name and command)
	os.Args = []string{"poke-cli", "pokemon", "bulbasaur"}

	// Call the function
	PokemonCommand()

	// Close and restore stdout
	if closeErr := w.Close(); closeErr != nil {
		t.Fatalf("Failed to close pipe writer: %v", closeErr)
	}
	os.Stdout = stdout

	_, readErr := output.ReadFrom(r)
	if readErr != nil {
		t.Fatalf("Failed to read from pipe: %v", readErr)
	}

	// Define expected output based on actual API response
	expectedName := "Bulbasaur"
	expectedID := "1"

	// Assert output contains expected Pokémon details
	assert.Contains(t, output.String(), "Your selected Pokémon:", "Output should contain Pokémon details header")
	assert.Contains(t, output.String(), expectedName, "Output should contain the Pokémon name")
	assert.Contains(t, output.String(), "National Pokédex #:", "Output should contain Pokémon ID header")
	assert.Contains(t, output.String(), expectedID, "Output should contain the Pokémon ID")
}
