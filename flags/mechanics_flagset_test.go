package flags

import (
	"reflect"
	"testing"

	"github.com/digitalghost-dev/poke-cli/styling"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetupMechanicsFlagSet(t *testing.T) {
	mf := SetupMechanicsFlagSet()

	assert.NotNil(t, mf, "Flag set should not be nil")
	assert.Equal(t, "mechanicsFlags", mf.FlagSet.Name(), "Flag set name should be 'mechanicsFlags'")

	flagTests := []struct {
		flag     interface{}
		expected interface{}
		name     string
	}{
		{mf.Natures, false, "Natures flag should default to false"},
		{mf.ShortNatures, false, "Short natures flag should default to false"},
	}

	for _, tt := range flagTests {
		assert.NotNil(t, tt.flag, tt.name)
		assert.Equal(t, tt.expected, reflect.ValueOf(tt.flag).Elem().Interface(), tt.name)
	}
}

func TestMechanicsFlagSetParse(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		wantLong  bool
		wantShort bool
	}{
		{"long natures flag", []string{"--natures"}, true, false},
		{"short natures flag", []string{"-n"}, false, true},
		{"no flags", []string{}, false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mf := SetupMechanicsFlagSet()
			err := mf.FlagSet.Parse(tt.args)
			require.NoError(t, err)
			assert.Equal(t, tt.wantLong, *mf.Natures, "long flag value")
			assert.Equal(t, tt.wantShort, *mf.ShortNatures, "short flag value")
		})
	}
}

func TestNaturesFlag(t *testing.T) {
	output := styling.StripANSI(NaturesFlag())

	assert.Contains(t, output, "Natures affect the growth of a Pokémon.")
	assert.Contains(t, output, "Nature Chart:")

	// Spot-check natures from each row/column of the chart.
	for _, nature := range []string{"Hardy", "Adamant", "Brave", "Timid", "Serious"} {
		assert.Contains(t, output, nature, "chart should contain nature %q", nature)
	}
}
