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

func TestValidateTypesArgs_TooManyArgs(t *testing.T) {
	invalidInputs := [][]string{
		{"poke-cli", "types", "ground"},
	}
	expectedError := "error, too many arguemnts\n"

	for _, input := range invalidInputs {
		err := ValidateTypesArgs(input)
		assert.Error(t, err, "Expected error for too many arguments")
		assert.NotEqual(t, expectedError, err.Error())
	}
}
