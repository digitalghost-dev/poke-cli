package types

import (
	"github.com/digitalghost-dev/poke-cli/cmd"
	"github.com/digitalghost-dev/poke-cli/styling"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestValidateTypesArgs_ValidInput(t *testing.T) {
	validInputs := [][]string{
		{"poke-cli", "types"},
		{"poke-cli", "types", "-h"},
	}

	for _, input := range validInputs {
		err := cmd.ValidateTypesArgs(input)
		assert.NoError(t, err, "Expected no error for valid input")
	}
}

func TestValidateTypesArgs_TooManyArgs(t *testing.T) {
	invalidInputs := [][]string{
		{"poke-cli", "types", "ground"},
	}
	expectedError := "error, too many arguments\n"

	for _, input := range invalidInputs {
		err := cmd.ValidateTypesArgs(input)
		require.Error(t, err, "Expected error for too many arguments")
		assert.NotEqual(t, expectedError, err.Error())
	}
}

func TestModelInit(t *testing.T) {
	m := model{}
	result := m.Init()

	assert.Nil(t, result, "Expected Init() to return nil")
}

func TestModelView_DisplayTable(t *testing.T) {
	m := model{selectedOption: ""}

	// Construct the expected output exactly as `View()` should render it
	expectedOutput := "Select a type!\n" +
		styling.TypesTableBorder.Render(m.table.View()) +
		"\n" +
		styling.KeyMenu.Render("↑ (move up) • ↓ (move down)\nenter (select) • ctrl+c | esc (quit)")

	output := m.View()

	assert.Equal(t, expectedOutput, output, "Expected View output to include table view")
}
