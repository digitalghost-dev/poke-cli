package flags

import (
	"github.com/stretchr/testify/assert"
	"reflect"
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
