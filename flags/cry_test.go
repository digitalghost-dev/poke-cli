//go:build !nocry

package flags

import (
	"testing"

	"github.com/digitalghost-dev/poke-cli/styling"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCryFlag_PokemonNotFound(t *testing.T) {
	err := CryFlag("pokemon", "phantumpf")
	require.Error(t, err)

	actual := styling.StripANSI(err.Error())
	assert.Contains(t, actual, "Pokémon not found")
	assert.Contains(t, actual, "Perhaps a typo?")
}
