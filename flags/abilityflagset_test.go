package flags

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestSetupAbilityFlagSet(t *testing.T) {
	abilityFlags, pokemonFlag, shortPokemonFlag := SetupAbilityFlagSet()

	assert.NotNil(t, abilityFlags, "Flag set should not be nil")
	assert.Equal(t, "AbilityFlagSet", abilityFlags.Name(), "Flag set name should be 'AbilityFlagSet'")

	flagTests := []struct {
		flag     interface{}
		expected interface{}
		name     string
	}{
		{pokemonFlag, false, "Pokemon flag should be 'pokemon'"},
		{shortPokemonFlag, false, "Short pokemon flag should be 'p'"},
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

	err := PokemonFlag("ability", "stench")

	if closeErr := w.Close(); closeErr != nil {
		t.Fatalf("Failed to close pipe writer: %v", closeErr)
	}
	os.Stdout = stdout

	_, readErr := output.ReadFrom(r)
	if readErr != nil {
		t.Fatalf("Failed to read from pipe: %v", readErr)
	}

	assert.NoError(t, err)

	expectedOutput := "\nPokémon with Stench\n\n"
	actualOutput := stripANSI(output.String())

	assert.True(t, strings.Contains(actualOutput, expectedOutput), "Output should contain expected ability description")
}
