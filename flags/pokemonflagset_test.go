package flags

import (
	"bytes"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/digitalghost-dev/poke-cli/styling"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetupPokemonFlagSet(t *testing.T) {
	pokeFlags, abilitiesFlag, shortAbilitiesFlag, defenseFlag, shortDefenseFlag, imageFlag, shortImageFlag, moveFlag, shortMoveFlag, statsFlag, shortStatsFlag, typesFlag, shortTypesFlag := SetupPokemonFlagSet()

	assert.NotNil(t, pokeFlags, "Flag set should not be nil")
	assert.Equal(t, "pokeFlags", pokeFlags.Name(), "Flag set name should be 'pokeFlags'")

	flagTests := []struct {
		flag     interface{}
		expected interface{}
		name     string
	}{
		{abilitiesFlag, false, "Abilities flag should be 'abilities'"},
		{shortAbilitiesFlag, false, "Short abilities flag should be 'a'"},
		{defenseFlag, false, "Defense flag should be 'defense'"},
		{shortDefenseFlag, false, "Short Defense flag should be 'd'"},
		{imageFlag, "", "Image flag default value should be 'md'"},
		{shortImageFlag, "", "Short image flag default value should be 'md'"},
		{moveFlag, false, "Move flag default value should be 'moves'"},
		{shortMoveFlag, false, "Short move flag default value should be 'm'"},
		{typesFlag, false, "Types flag should be 'types'"},
		{shortTypesFlag, false, "Short types flag should be 't'"},
		{statsFlag, false, "Stats flag should be 'stats'"},
		{shortStatsFlag, false, "Short stats flag should be 's'"},
	}

	for _, tt := range flagTests {
		assert.NotNil(t, tt.flag, tt.name)
		assert.Equal(t, tt.expected, reflect.ValueOf(tt.flag).Elem().Interface(), tt.name)
	}
}

func TestAbilitiesFlag(t *testing.T) {
	var output bytes.Buffer
	stdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := AbilitiesFlag(&output, "pokemon", "bulbasaur")

	if closeErr := w.Close(); closeErr != nil {
		t.Fatalf("Failed to close pipe writer: %v", closeErr)
	}
	os.Stdout = stdout

	_, readErr := output.ReadFrom(r)
	if readErr != nil {
		t.Fatalf("Failed to read from pipe: %v", readErr)
	}

	require.NoError(t, err)

	expectedOutput := `─────────
Abilities
Ability 1: Overgrow
Hidden Ability: Chlorophyll
`

	actualOutput := styling.StripANSI(output.String())

	assert.Equal(t, expectedOutput, actualOutput, "Output should contain data for the abilities flag")
}

func TestDefenseFlag(t *testing.T) {
	var output bytes.Buffer
	stdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := DefenseFlag(&output, "pokemon", "bulbasaur")

	if closeErr := w.Close(); closeErr != nil {
		t.Fatalf("Failed to close pipe writer: %v", closeErr)
	}
	os.Stdout = stdout

	_, readErr := output.ReadFrom(r)
	if readErr != nil {
		t.Fatalf("Failed to read from pipe: %v", readErr)
	}

	require.NoError(t, err, "DefenseFlag should not return an error for a valid Pokémon")

	actualOutput := styling.StripANSI(output.String())

	assert.Contains(t, actualOutput, "Type Defenses", "Output should contain the header 'Type Defenses'")

	assert.Contains(t, actualOutput, "0.25×   Damage", "Should include quarter damage category")
	assert.Contains(t, actualOutput, "0.5×    Damage", "Should include half damage category")
	assert.Contains(t, actualOutput, "2.0×    Damage", "Should include double damage category")

	for _, typ := range []string{"Grass"} {
		assert.Contains(t, actualOutput, typ, "Quarter damage should list %s", typ)
	}
	for _, typ := range []string{"Water", "Electric", "Fighting", "Fairy"} {
		assert.Contains(t, actualOutput, typ, "Half damage should list %s", typ)
	}
	for _, typ := range []string{"Fire", "Ice", "Flying", "Psychic"} {
		assert.Contains(t, actualOutput, typ, "Double damage should list %s", typ)
	}

	assert.NotContains(t, actualOutput, "Immune:", "Bulbasaur should not have immunities")
	assert.NotContains(t, actualOutput, "4.0×", "Bulbasaur should not have 4x weaknesses")
}

func TestImageFlag(t *testing.T) {
	var output bytes.Buffer
	stdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := ImageFlag(&output, "pokemon", "bulbasaur", "sm")

	if closeErr := w.Close(); closeErr != nil {
		t.Fatalf("Failed to close pipe writer: %v", closeErr)
	}
	os.Stdout = stdout

	_, readErr := output.ReadFrom(r)
	if readErr != nil {
		t.Fatalf("Failed to read from pipe: %v", readErr)
	}

	require.NoError(t, err)

	actualOutput := styling.StripANSI(output.String())

	if !strings.Contains(actualOutput, "▀") {
		t.Errorf("Output does not contain the expected ASCII art characters.")
	}

	if len(actualOutput) == 0 {
		t.Errorf("Output is empty; expected ASCII art representation of the Pokémon image.")
	}
}

func TestImageFlagOptions(t *testing.T) {
	validOptions := []string{"lg", "md", "sm"}

	for _, option := range validOptions {
		t.Run("ValidOption_"+option, func(t *testing.T) {
			var buf bytes.Buffer
			err := ImageFlag(&buf, "pokemon", "bulbasaur", option)
			assert.NoError(t, err, "ImageFlag should not return an error for valid option '%s'", option)
		})
	}

	invalidOptions := []string{"s", "med", "large"}

	for _, option := range invalidOptions {
		t.Run("InvalidOption_"+option, func(t *testing.T) {
			var buf bytes.Buffer
			err := ImageFlag(&buf, "pokemon", "bulbasaur", option)
			assert.Error(t, err, "ImageFlag should return an error for invalid option '%s'", option)
		})
	}
}

func TestMovesFlag(t *testing.T) {
	var output bytes.Buffer
	stdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := MovesFlag(&output, "pokemon", "bulbasaur")

	if closeErr := w.Close(); closeErr != nil {
		t.Fatalf("Failed to close pipe writer: %v", closeErr)
	}
	os.Stdout = stdout

	_, readErr := output.ReadFrom(r)
	if readErr != nil {
		t.Fatalf("Failed to read from pipe: %v", readErr)
	}

	require.NoError(t, err, "MovesFlag should not return an error for a valid Pokémon")

	actualOutput := styling.StripANSI(output.String())

	assert.Contains(t, actualOutput, "Learnable Moves", "Output should contain the header 'Learnable Moves'")

	assert.Contains(t, actualOutput, "Name", "Output should contain the 'Name' column header")
	assert.Contains(t, actualOutput, "Level", "Output should contain the 'Level' column header")
	assert.Contains(t, actualOutput, "Type", "Output should contain the 'Type' column header")
	assert.Contains(t, actualOutput, "Accuracy", "Output should contain the 'Accuracy' column header")
	assert.Contains(t, actualOutput, "Power", "Output should contain the 'Power' column header")

	assert.NotEmpty(t, actualOutput, "Output should not be empty")

	assert.True(t,
		strings.Contains(actualOutput, "Tackle") ||
			strings.Contains(actualOutput, "Vine Whip") ||
			strings.Contains(actualOutput, "Growl"),
		"Output should contain at least one of Bulbasaur's common moves")
}

func TestStatsFlag(t *testing.T) {
	var output bytes.Buffer
	stdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := StatsFlag(&output, "pokemon", "bulbasaur")

	if closeErr := w.Close(); closeErr != nil {
		t.Fatalf("Failed to close pipe writer: %v", closeErr)
	}
	os.Stdout = stdout

	_, readErr := output.ReadFrom(r)
	if readErr != nil {
		t.Fatalf("Failed to read from pipe: %v", readErr)
	}

	require.NoError(t, err, "StatsFlag should not return an error for a valid Pokémon")

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

	actualOutput := styling.StripANSI(output.String())

	assert.Equal(t, expectedOutput, actualOutput, "Output should contain data for the stats flag")
}

func TestTypesFlag(t *testing.T) {
	var output bytes.Buffer
	stdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := TypesFlag(&output, "pokemon", "bulbasaur")

	if closeErr := w.Close(); closeErr != nil {
		t.Fatalf("Failed to close pipe writer: %v", closeErr)
	}
	os.Stdout = stdout

	_, readErr := output.ReadFrom(r)
	if readErr != nil {
		t.Fatalf("Failed to read from pipe: %v", readErr)
	}

	require.NoError(t, err, "TypesFlag should not return an error for a valid Pokémon")

	expectedOutput := `──────
Typing
Type 1: Grass
Type 2: Poison
╭─────────────────────────────────────╮
│⚠ Warning!                           │
│The '-t | --types' flag is deprecated│
│and will be removed in v2.           │
│                                     │
│Typing is now included by default.   │
│You no longer need this flag.        │
╰─────────────────────────────────────╯
`
	actualOutput := styling.StripANSI(output.String())

	assert.Equal(t, expectedOutput, actualOutput, "Output should contain data for the types flag")
}
