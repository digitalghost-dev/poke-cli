package styling

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetTypeColor(t *testing.T) {
	// Test known types
	for typeName, expectedColor := range ColorMap {
		t.Run(typeName, func(t *testing.T) {
			color := GetTypeColor(typeName)
			assert.Equal(t, expectedColor, color, "Expected color for type %s to be %s", typeName, expectedColor)
		})
	}

	// Test unknown type
	t.Run("unknown type", func(t *testing.T) {
		color := GetTypeColor("unknown")
		assert.Equal(t, "", color, "Expected color for unknown type to be an empty string")
	})
}
