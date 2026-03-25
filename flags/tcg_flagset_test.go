package flags

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetupTcgFlagSet(t *testing.T) {
	tf := SetupTcgFlagSet()

	assert.NotNil(t, tf, "Flag set should not be nil")
	assert.Equal(t, "tcgFlags", tf.FlagSet.Name(), "Flag set name should be 'tcgFlags'")

	flagTests := []struct {
		flag     interface{}
		expected interface{}
		name     string
	}{
		{tf.Web, false, "Web flag default should be false"},
		{tf.ShortWeb, false, "ShortWeb flag default should be false"},
	}

	for _, tt := range flagTests {
		assert.NotNil(t, tt.flag, tt.name)
		assert.Equal(t, tt.expected, reflect.ValueOf(tt.flag).Elem().Interface(), tt.name)
	}
}

func TestSetupTcgFlagSet_ParseWebFlag(t *testing.T) {
	tests := []struct {
		name string
		args []string
		web  bool
	}{
		{"long flag", []string{"--web"}, true},
		{"short flag", []string{"-w"}, true},
		{"no flag", []string{}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tf := SetupTcgFlagSet()
			err := tf.FlagSet.Parse(tt.args)
			require.NoError(t, err)
			assert.Equal(t, tt.web, *tf.Web || *tf.ShortWeb)
		})
	}
}
