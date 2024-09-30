package cmd

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidateTypesArgs_ValidInput(t *testing.T) {
	validInputs := [][]string{
		{"poke-cli", "types"},
		{"poke-cli", "types", "-h"},
	}

	for _, input := range validInputs {
		err := ValidateTypesArgs(input)
		assert.NoError(t, err, "Expected no error for valid input")
	}
}
