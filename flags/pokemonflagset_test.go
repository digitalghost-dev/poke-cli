package flags

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"os"
	"regexp"
	"testing"
)

func stripANSI(input string) string {
	ansiRegex := regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)
	return ansiRegex.ReplaceAllString(input, "")
}

func TestSetupPokemonFlagSet(t *testing.T) {
	// Call the function to get the flag set and types flag
	pokeFlags, abilitiesFlag, shortAbilitiesFlag, statsFlag, shortStatsFlag, typesFlag, shortTypesFlag := SetupPokemonFlagSet()

	// Assertions
	assert.NotNil(t, pokeFlags, "Flag set should not be nil")
	assert.Equal(t, "pokeFlags", pokeFlags.Name(), "Flag set name should be 'pokeFlags'")
	//assert.Equal(t, flag.ExitOnError, pokeFlags.NFlag(), "Flag set should have ExitOnError behavior")

	// Check abilities flag
	assert.NotNil(t, abilitiesFlag, "Abilities flag should not be nil")
	assert.Equal(t, false, *abilitiesFlag, "Abilities flag name should be 'abilities'")

	// Check short abilities flag
	assert.NotNil(t, shortAbilitiesFlag, "Short abilities flag should not be nil")
	assert.Equal(t, false, *shortAbilitiesFlag, "Short abilities flag name should be 'a'")

	// Check types flag
	assert.NotNil(t, typesFlag, "Types flag should not be nil")
	assert.Equal(t, false, *typesFlag, "Types flag name should be 'types'")

	// Check short types flag
	assert.NotNil(t, shortTypesFlag, "Short types flag should not be nil")
	assert.Equal(t, false, *shortTypesFlag, "Short types flag name should be 't'")

	// Check abilities flag
	assert.NotNil(t, statsFlag, "Stats flag should not be nil")
	assert.Equal(t, false, *statsFlag, "Stats flag name should be 'abilities'")

	// Check short abilities flag
	assert.NotNil(t, shortStatsFlag, "Short stats flag should not be nil")
	assert.Equal(t, false, *shortStatsFlag, "Short stats flag name should be 'a'")
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
Ability 1: Overgrow
Hidden Ability: Chlorophyll
`

	// Assert the actual output matches the expected output
	actualOutput := stripANSI(output.String())

	assert.Equal(t, expectedOutput, actualOutput, "Output should contain data for the abilities flag")
}

func TestStatsFlag(t *testing.T) {
	// Capture standard output
	var output bytes.Buffer
	stdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Call the StatsFlag function with a valid Pokémon
	err := StatsFlag("pokemon", "bulbasaur")

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
	assert.NoError(t, err, "StatsFlag should not return an error for a valid Pokémon")

	// Define expected output components
	expectedOutput := `──────────
Base Stats
HP         ▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇ 45
Atk        ▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇ 49
Def        ▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇ 49
Sp. Atk    ▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇ 65
Sp. Def    ▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇ 65
Speed      ▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇ 45
`

	// Assert output contains the expected header and typing information
	actualOutput := stripANSI(output.String())

	assert.Equal(t, expectedOutput, actualOutput, "Output should contain data for the stats flag")

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
	expectedOutput := `──────
Typing
Type 1: Grass
Type 2: Poison
`

	// Assert output contains the expected header and typing information
	actualOutput := stripANSI(output.String())

	assert.Equal(t, expectedOutput, actualOutput, "Output should contain data for the types flag")
}
