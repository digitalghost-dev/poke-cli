package flags

import (
	"bytes"
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/digitalghost-dev/poke-cli/styling"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetupAbilityFlagSet(t *testing.T) {
	af := SetupAbilityFlagSet()

	assert.NotNil(t, af, "Flag set should not be nil")
	assert.Equal(t, "abilityFlags", af.FlagSet.Name(), "Flag set name should be 'abilityFlags'")

	flagTests := []struct {
		flag     interface{}
		expected interface{}
		name     string
	}{
		{af.Pokemon, false, "Pokemon flag should be 'pokemon'"},
		{af.ShortPokemon, false, "Short pokemon flag should be 'p'"},
	}

	for _, tt := range flagTests {
		assert.NotNil(t, tt.flag, tt.name)
		assert.Equal(t, tt.expected, reflect.ValueOf(tt.flag).Elem().Interface(), tt.name)
	}
}

func TestPokemonFlag(t *testing.T) {
	var output bytes.Buffer
	stdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := PokemonAbilitiesFlag(&output, "ability", "stench")

	if closeErr := w.Close(); closeErr != nil {
		t.Fatalf("Failed to close pipe writer: %v", closeErr)
	}
	os.Stdout = stdout

	_, readErr := output.ReadFrom(r)
	if readErr != nil {
		t.Fatalf("Failed to read from pipe: %v", readErr)
	}

	require.NoError(t, err)

	expectedOutput := styling.StripANSI(fmt.Sprintf(
		"Pokemon with Stench\n\n"+
			"%2d. %-30s%2d. %-30s%2d. %-30s\n"+
			"%2d. %-30s%2d. %-30s%2d. %-30s\n"+
			"%2d. %-30s%2d. %-30s%2d. %-30s\n"+
			"%2d. %-30s\n",
		1, "Gloom", 2, "Grimer", 3, "Muk",
		4, "Koffing", 5, "Weezing", 6, "Stunky",
		7, "Skuntank", 8, "Trubbish", 9, "Garbodor",
		10, "Garbodor-Gmax"),
	)
	actualOutput := strings.TrimSpace(styling.StripANSI(output.String()))
	expectedOutput = strings.TrimSpace(expectedOutput)

	if !strings.Contains(actualOutput, expectedOutput) {
		t.Logf("Actual Output:\n%s\n", actualOutput)
		t.Logf("Expected Output:\n%s\n", expectedOutput)
	}
	assert.Contains(t, actualOutput, expectedOutput, "Output should contain Pokémon with the ability")
}
