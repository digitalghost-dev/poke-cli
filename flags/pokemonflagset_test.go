package flags

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSetupPokemonFlagSet(t *testing.T) {
	// Call the function to get the flag set and types flag
	pokeFlags, typesFlag, abilitiesFlag := SetupPokemonFlagSet()

	// Assertions
	assert.NotNil(t, pokeFlags, "Flag set should not be nil")
	assert.Equal(t, "pokeFlags", pokeFlags.Name(), "Flag set name should be 'pokeFlags'")
	//assert.Equal(t, flag.ExitOnError, pokeFlags.NFlag(), "Flag set should have ExitOnError behavior")

	// Check types flag
	assert.NotNil(t, typesFlag, "Types flag should not be nil")
	assert.Equal(t, bool(false), *typesFlag, "Types flag name should be 'types'")

	// Check abilities flag
	assert.NotNil(t, abilitiesFlag, "Abilities flag should not be nil")
	assert.Equal(t, bool(false), *abilitiesFlag, "Abilities flag name should be 'abilities'")
}
