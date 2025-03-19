package cmd

import (
	"flag"
	"fmt"
	"github.com/digitalghost-dev/poke-cli/styling"
	"os"
)

func handleHelpFlag(args []string) {
	if len(args) == 3 && (args[2] == "-h" || args[2] == "--help") {
		flag.Usage()

		if flag.Lookup("test.v") == nil {
			os.Exit(0)
		}
	}
}

// ValidateAbilityArgs validates the command line arguments
func ValidateAbilityArgs(args []string) error {
	handleHelpFlag(args)

	if len(args) > 4 {
		errMessage := styling.ErrorBorder.Render(styling.ErrorColor.Render("Error!"), "\nToo many arguments")
		return fmt.Errorf("%s", errMessage)
	}

	if len(args) == 2 {
		errMessage := styling.ErrorBorder.Render(styling.ErrorColor.Render("Error!"), "\nPlease specify an ability")
		return fmt.Errorf("%s", errMessage)
	}

	return nil
}

// ValidateNaturesArgs validates the command line arguments
func ValidateNaturesArgs(args []string) error {
	handleHelpFlag(args)

	if len(args) > 3 {
		errMessage := styling.ErrorBorder.Render(styling.ErrorColor.Render("Error!"), "\nToo many arguments")
		return fmt.Errorf("%s", errMessage)
	}

	// Check if there are exactly 3 arguments and the third argument is neither '-h' nor '--help'
	// If true, return an error message since only '-h' and '--help' are allowed after 'types'
	if len(args) == 3 && (args[2] != "-h" && args[2] != "--help") {
		errorTitle := styling.ErrorColor.Render("Error!")
		errorString := "\nThe only currently available options\nafter [natures] command are '-h' or '--help'"
		finalErrorMessage := errorTitle + errorString
		renderedError := styling.ErrorBorder.Render(finalErrorMessage)
		return fmt.Errorf("%s", renderedError)
	}

	return nil
}

// ValidatePokemonArgs validates the command line arguments
func ValidatePokemonArgs(args []string) error {
	handleHelpFlag(args)

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

	// Check if there are too many arguments
	if len(args) > 7 {
		errMessage := styling.ErrorBorder.Render(
			styling.ErrorColor.Render("Error!"),
			"\nToo many arguments",
		)
		return fmt.Errorf("%s", errMessage)
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

// ValidateTypesArgs validates the command line arguments
func ValidateTypesArgs(args []string) error {
	handleHelpFlag(args)

	if len(args) > 3 {
		errMessage := styling.ErrorBorder.Render(styling.ErrorColor.Render("Error!"), "\nToo many arguments")
		return fmt.Errorf("%s", errMessage)
	}

	// Check if there are exactly 3 arguments and the third argument is neither '-h' nor '--help'
	// If true, return an error message since only '-h' and '--help' are allowed after 'types'
	if len(args) == 3 && (args[2] != "-h" && args[2] != "--help") {
		errorTitle := styling.ErrorColor.Render("Error!")
		errorString := "\nThe only currently available options\nafter [types] command are '-h' or '--help'"
		finalErrorMessage := errorTitle + errorString
		renderedError := styling.ErrorBorder.Render(finalErrorMessage)
		return fmt.Errorf("%s", renderedError)
	}

	return nil
}
