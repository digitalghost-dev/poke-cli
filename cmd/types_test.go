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
	expectedError := "error, too many arguments\n"

	for _, input := range invalidInputs {
		err := ValidateTypesArgs(input)
		assert.Error(t, err, "Expected error for too many arguments")
		assert.NotEqual(t, expectedError, err.Error())
	}
}

func TestModelInit(t *testing.T) {
	m := model{}
	result := m.Init()

	assert.Nil(t, result, "Expected Init() to return nil")
}

func TestModelView_SelectedOption(t *testing.T) {
	m := model{selectedOption: "someOption"}

	output := m.View()

	assert.Equal(t, "", output, "Expected output to be an empty string when selectedOption is set")
}

func TestModelView_DisplayTable(t *testing.T) {
	m := model{selectedOption: ""}

	// Construct the expected output exactly as `View()` should render it
	expectedOutput := "Select a type!\n" +
		typesTableBorder.Render(m.table.View()) +
		"\n" +
		keyMenu.Render("↑ (move up) • ↓ (move down)\nctrl+c | esc (quit) • enter (select)")

	output := m.View()

	assert.Equal(t, expectedOutput, output, "Expected View output to include table view")
}
