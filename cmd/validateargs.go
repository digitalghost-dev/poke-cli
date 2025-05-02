package cmd

import (
	"fmt"
	"github.com/digitalghost-dev/poke-cli/styling"
)

// checkLength checks if the number of arguments is lower than the max value
func checkLength(args []string, max int) error {
	if len(args) > max {
		errMessage := styling.ErrorBorder.Render(
			styling.ErrorColor.Render("Error!") + "\nToo many arguments",
		)
		return fmt.Errorf("%s", errMessage)
	}
	return nil
}

// ValidateAbilityArgs validates the command line arguments
func ValidateAbilityArgs(args []string) error {
	if err := checkLength(args, 4); err != nil {
		return err
	}

	if len(args) == 2 {
		errMessage := styling.ErrorBorder.Render(styling.ErrorColor.Render("Error!"), "\nPlease specify an ability")
		return fmt.Errorf("%s", errMessage)
	}

	return nil
}

// ValidateMoveArgs validates the command line arguments
func ValidateMoveArgs(args []string) error {
	if err := checkLength(args, 3); err != nil {
		return err
	}

	if len(args) == 2 {
		errMessage := styling.ErrorBorder.Render(styling.ErrorColor.Render("Error!"), "\nPlease specify a move")
		return fmt.Errorf("%s", errMessage)
	}

	return nil
}

// ValidateNaturesArgs validates the command line arguments
func ValidateNaturesArgs(args []string) error {
	if err := checkLength(args, 3); err != nil {
		return err
	}

	// Check if there are exactly 3 arguments and the third argument is neither '-h' nor '--help'
	// If true, return an error message since only '-h' and '--help' are allowed after 'types'
	if len(args) == 3 && args[2] != "-h" && args[2] != "--help" {
		errMsg := styling.ErrorColor.Render("Error!") +
			"\nThe only currently available options\nafter <natures> command are '-h' or '--help'"
		return fmt.Errorf("%s %s", styling.ErrorBorder.Render(errMsg), "\n")
	}

	return nil
}

// ValidatePokemonArgs validates the command line arguments
func ValidatePokemonArgs(args []string) error {
	// Check if the number of arguments is less than 3
	if len(args) < 3 {
		errMessage := styling.ErrorBorder.Render(
			styling.ErrorColor.Render("Error!"),
			"\nPlease declare a Pokémon's name after the <pokemon> command",
			"\nRun 'poke-cli pokemon -h' for more details",
			"\nerror: insufficient arguments",
		)
		return fmt.Errorf("%s", errMessage)
	}

	if err := checkLength(args, 7); err != nil {
		return err
	}

	// Validate each argument after the Pokémon's name
	if len(args) > 3 {
		for _, arg := range args[3:] {
			// Check for single `-` or `--` which are invalid
			if arg == "-" || arg == "--" {
				errorTitle := styling.ErrorColor.Render("Error!")
				errorString := fmt.Sprintf(
					"\nInvalid argument '%s'. Single '-' or '--' is not allowed.\nPlease use valid flags.",
					arg,
				)
				finalErrorMessage := errorTitle + errorString
				renderedError := styling.ErrorBorder.Render(finalErrorMessage)
				return fmt.Errorf("%s", renderedError)
			}

			// Check if the argument starts with a flag prefix but is invalid
			if arg[0] != '-' {
				errorTitle := styling.ErrorColor.Render("Error!")
				errorString := fmt.Sprintf(
					"\nInvalid argument '%s'.\nOnly flags are allowed after declaring a Pokémon's name",
					arg,
				)
				finalErrorMessage := errorTitle + errorString
				renderedError := styling.ErrorBorder.Render(finalErrorMessage)
				return fmt.Errorf("%s", renderedError)
			}
		}
	}

	return nil
}

// ValidateSearchArgs validates the command line arguments
func ValidateSearchArgs(args []string) error {
	if err := checkLength(args, 3); err != nil {
		return err
	}

	// Check if there are exactly 3 arguments and the third argument is neither '-h' nor '--help'
	// If true, return an error message since only '-h' and '--help' are allowed after <search>
	if len(args) == 3 && args[2] != "-h" && args[2] != "--help" {
		errMsg := styling.ErrorColor.Render("Error!") +
			"\nThe only currently available options\nafter <search> command are '-h' or '--help'"
		return fmt.Errorf("%s", styling.ErrorBorder.Render(errMsg))
	}

	return nil
}

// ValidateTypesArgs validates the command line arguments
func ValidateTypesArgs(args []string) error {
	if err := checkLength(args, 3); err != nil {
		return err
	}

	// Check if there are exactly 3 arguments and the third argument is neither '-h' nor '--help'
	// If true, return an error message since only '-h' and '--help' are allowed after <types>
	if len(args) == 3 && args[2] != "-h" && args[2] != "--help" {
		errMsg := styling.ErrorColor.Render("Error!") +
			"\nThe only currently available options\nafter <types> command are '-h' or '--help'"
		return fmt.Errorf("%s", styling.ErrorBorder.Render(errMsg))
	}

	return nil
}
