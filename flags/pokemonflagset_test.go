package flags

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestSetupPokemonFlagSet(t *testing.T) {
	// Call the function to get the flag set and types flag
	pokeFlags, typesFlag, shortTypesFlag, abilitiesFlag, shortAbilitiesFlag := SetupPokemonFlagSet()

	// Assertions
	assert.NotNil(t, pokeFlags, "Flag set should not be nil")
	assert.Equal(t, "pokeFlags", pokeFlags.Name(), "Flag set name should be 'pokeFlags'")
	//assert.Equal(t, flag.ExitOnError, pokeFlags.NFlag(), "Flag set should have ExitOnError behavior")

	// Check types flag
	assert.NotNil(t, typesFlag, "Types flag should not be nil")
	assert.Equal(t, false, *typesFlag, "Types flag name should be 'types'")

	// Check short types flag
	assert.NotNil(t, shortTypesFlag, "Short types flag should not be nil")
	assert.Equal(t, false, *shortTypesFlag, "Short types flag name should be 't'")

	// Check abilities flag
	assert.NotNil(t, abilitiesFlag, "Abilities flag should not be nil")
	assert.Equal(t, false, *abilitiesFlag, "Abilities flag name should be 'abilities'")

	// Check short abilities flag
	assert.NotNil(t, shortAbilitiesFlag, "Short abilities flag should not be nil")
	assert.Equal(t, false, *shortAbilitiesFlag, "Short abilities flag name should be 'a'")
}

func TestAbilitiesFlag(t *testing.T) {
	// Capture standard output
	var output bytes.Buffer
	stdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Call the function with a known Pokémon (e.g., bulbasaur)
	err := AbilitiesFlag("pokemon", "bulbasaur")

	// Close and restore stdout
	if closeErr := w.Close(); closeErr != nil {
		t.Fatalf("Failed to close pipe writer: %v", closeErr)
	}
	os.Stdout = stdout

	_, readErr := output.ReadFrom(r)
	if readErr != nil {
		t.Fatalf("Failed to read from pipe: %v", readErr)
	}

	// Assert no errors occurred during execution
	assert.NoError(t, err)

	// Define the expected output based on the API response
	expectedOutput := `─────────
Abilities
Ability 1: overgrow
Hidden Ability: chlorophyll
`

	// Assert the actual output matches the expected output
	assert.Contains(t, output.String(), "Abilities", "Output should contain 'Abilities'")
	assert.Contains(t, output.String(), "Ability 1: overgrow", "Output should contain 'Ability 1: overgrow'")
	assert.Contains(t, output.String(), "Hidden Ability: chlorophyll", "Output should contain 'Ability 2: chlorophyll'")
	assert.Equal(t, expectedOutput, output.String(), "Output does not match the expected formatting")
}

func TestTypesFlag(t *testing.T) {
	// Capture standard output
	var output bytes.Buffer
	stdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Call the TypesFlag function with a valid Pokémon
	err := TypesFlag("pokemon", "bulbasaur")

	// Close and restore stdout
	if closeErr := w.Close(); closeErr != nil {
		t.Fatalf("Failed to close pipe writer: %v", closeErr)
	}
	os.Stdout = stdout

	_, readErr := output.ReadFrom(r)
	if readErr != nil {
		t.Fatalf("Failed to read from pipe: %v", readErr)
	}

	// Assert no errors occurred
	assert.NoError(t, err, "TypesFlag should not return an error for a valid Pokémon")

	// Define expected output components
	expectedHeader := "Typing"
	expectedType1 := "Type 1: grass"
	expectedType2 := "Type 2: poison"

	// Assert output contains the expected header and typing information
	assert.Contains(t, output.String(), expectedHeader, "Output should contain the 'Typing' header")
	assert.Contains(t, output.String(), expectedType1, "Output should contain the Pokémon's first type")
	assert.Contains(t, output.String(), expectedType2, "Output should contain the Pokémon's second type")
}
