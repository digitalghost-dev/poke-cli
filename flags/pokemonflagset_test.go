package flags

import (
	"bytes"
	"github.com/digitalghost-dev/poke-cli/styling"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestSetupPokemonFlagSet(t *testing.T) {
	// Call the function to get the flag set and flags
	pokeFlags, abilitiesFlag, shortAbilitiesFlag, imageFlag, shortImageFlag, moveFlag, shortMoveFlag, statsFlag, shortStatsFlag, typesFlag, shortTypesFlag := SetupPokemonFlagSet()

	// Check flag set properties
	assert.NotNil(t, pokeFlags, "Flag set should not be nil")
	assert.Equal(t, "pokeFlags", pokeFlags.Name(), "Flag set name should be 'pokeFlags'")

	// Define test cases for flag assertions
	flagTests := []struct {
		flag     interface{}
		expected interface{}
		name     string
	}{
		{abilitiesFlag, false, "Abilities flag should be 'abilities'"},
		{shortAbilitiesFlag, false, "Short abilities flag should be 'a'"},
		{imageFlag, "", "Image flag default value should be 'md'"},
		{shortImageFlag, "", "Short image flag default value should be 'md'"},
		{moveFlag, false, "Move flag default value should be 'moves'"},
		{shortMoveFlag, false, "Short move flag default value should be 'm'"},
		{typesFlag, false, "Types flag should be 'types'"},
		{shortTypesFlag, false, "Short types flag should be 't'"},
		{statsFlag, false, "Stats flag should be 'stats'"},
		{shortStatsFlag, false, "Short stats flag should be 's'"},
	}

	// Run assertions for all flags
	for _, tt := range flagTests {
		assert.NotNil(t, tt.flag, tt.name)
		assert.Equal(t, tt.expected, reflect.ValueOf(tt.flag).Elem().Interface(), tt.name)
	}
}

func TestAbilitiesFlag(t *testing.T) {
	// Capture standard output
	var output bytes.Buffer
	stdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Call the function with a known Pokémon (e.g., bulbasaur)
	err := AbilitiesFlag(&output, "pokemon", "bulbasaur")

	// Close and restore stdout
	if closeErr := w.Close(); closeErr != nil {
		t.Fatalf("Failed to close pipe writer: %v", closeErr)
	}
	os.Stdout = stdout

	_, readErr := output.ReadFrom(r)
	if readErr != nil {
		t.Fatalf("Failed to read from pipe: %v", readErr)
	}

	require.NoError(t, err)

	// Define the expected output based on the API response
	expectedOutput := `─────────
Abilities
Ability 1: Overgrow
Hidden Ability: Chlorophyll
`

	// Assert the actual output matches the expected output
	actualOutput := styling.StripANSI(output.String())

	assert.Equal(t, expectedOutput, actualOutput, "Output should contain data for the abilities flag")
}

func TestImageFlag(t *testing.T) {
	// Capture standard output
	var output bytes.Buffer
	stdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Call the function with a known Pokémon (e.g., bulbasaur)
	err := ImageFlag(&output, "pokemon", "bulbasaur", "sm")

	// Close and restore stdout
	if closeErr := w.Close(); closeErr != nil {
		t.Fatalf("Failed to close pipe writer: %v", closeErr)
	}
	os.Stdout = stdout

	_, readErr := output.ReadFrom(r)
	if readErr != nil {
		t.Fatalf("Failed to read from pipe: %v", readErr)
	}

	require.NoError(t, err)

	// Validate that the output contains some expected patterns
	actualOutput := styling.StripANSI(output.String())

	// Since the output is an ASCII image, we can't hardcode the expected output,
	// but we can check that it contains some general expected structure
	if !strings.Contains(actualOutput, "▀") {
		t.Errorf("Output does not contain the expected ASCII art characters.")
	}

	if len(actualOutput) == 0 {
		t.Errorf("Output is empty; expected ASCII art representation of the Pokémon image.")
	}
}

func TestImageFlagOptions(t *testing.T) {
	// Define valid options as a slice
	validOptions := []string{"lg", "md", "sm"}

	// Test valid options
	for _, option := range validOptions {
		t.Run("ValidOption_"+option, func(t *testing.T) {
			var buf bytes.Buffer
			err := ImageFlag(&buf, "pokemon", "bulbasaur", option)
			assert.NoError(t, err, "ImageFlag should not return an error for valid option '%s'", option)
		})
	}

	// Define invalid options as a slice
	invalidOptions := []string{"s", "med", "large"}

	// Test invalid options
	for _, option := range invalidOptions {
		t.Run("InvalidOption_"+option, func(t *testing.T) {
			var buf bytes.Buffer
			err := ImageFlag(&buf, "pokemon", "bulbasaur", option)
			assert.Error(t, err, "ImageFlag should return an error for invalid option '%s'", option)
		})
	}
}

func TestStatsFlag(t *testing.T) {
	// Capture standard output
	var output bytes.Buffer
	stdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Call the StatsFlag function with a valid Pokémon
	err := StatsFlag(&output, "pokemon", "bulbasaur")

	// Close and restore stdout
	if closeErr := w.Close(); closeErr != nil {
		t.Fatalf("Failed to close pipe writer: %v", closeErr)
	}
	os.Stdout = stdout

	_, readErr := output.ReadFrom(r)
	if readErr != nil {
		t.Fatalf("Failed to read from pipe: %v", readErr)
	}

	require.NoError(t, err, "StatsFlag should not return an error for a valid Pokémon")

	// Define expected output components
	expectedOutput := `──────────
Base Stats
HP         ▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇ 45
Atk        ▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇ 49
Def        ▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇ 49
Sp. Atk    ▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇ 65
Sp. Def    ▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇ 65
Speed      ▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇ 45
Total      318
`

	// Assert output contains the expected header and typing information
	actualOutput := styling.StripANSI(output.String())

	assert.Equal(t, expectedOutput, actualOutput, "Output should contain data for the stats flag")
}

func TestTypesFlag(t *testing.T) {
	// Capture standard output
	var output bytes.Buffer
	stdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Call the TypesFlag function with a valid Pokémon
	err := TypesFlag(&output, "pokemon", "bulbasaur")

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
	require.NoError(t, err, "TypesFlag should not return an error for a valid Pokémon")

	// Define expected output components
	expectedOutput := `──────
Typing
Type 1: Grass
Type 2: Poison
`

	// Assert output contains the expected header and typing information
	actualOutput := styling.StripANSI(output.String())

	assert.Equal(t, expectedOutput, actualOutput, "Output should contain data for the types flag")
}
