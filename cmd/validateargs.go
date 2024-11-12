package cmd

import (
	"flag"
	"fmt"
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

// ValidatePokemonArgs validates the command line arguments
func ValidatePokemonArgs(args []string) error {

	handleHelpFlag(args)

	if len(args) < 3 {
		errMessage := errorBorder.Render(
			errorColor.Render("Error!"),
			"\nPlease declare a Pokémon's name after the [pokemon] command",
			"\nRun 'poke-cli pokemon -h' for more details",
			"\nerror: insufficient arguments",
		)
		return fmt.Errorf("%s", errMessage)
	}

	if len(args) > 5 {
		errMessage := errorBorder.Render(errorColor.Render("Error!"), "\nToo many arguments")
		return fmt.Errorf("%s", errMessage)
	}

	if len(args) > 3 {
		for _, arg := range args[3:] {
			if arg[0] != '-' {
				errorTitle := errorColor.Render("Error!")
				errorString := fmt.Sprintf("\nInvalid argument '%s'. Only flags are allowed after declaring a Pokémon's name", arg)
				finalErrorMessage := errorTitle + errorString
				renderedError := errorBorder.Render(finalErrorMessage)
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
		errMessage := errorBorder.Render(errorColor.Render("Error!"), "\nToo many arguments")
		return fmt.Errorf("%s", errMessage)
	}

	// Check if there are exactly 3 arguments and the third argument is neither '-h' nor '--help'
	// If true, return an error message since only '-h' and '--help' are allowed after 'types'
	if len(args) == 3 && (args[2] != "-h" && args[2] != "--help") {
		errorTitle := errorColor.Render("Error!")
		errorString := "\nThe only currently available options\nafter [types] command are '-h' or '--help'"
		finalErrorMessage := errorTitle + errorString
		renderedError := errorBorder.Render(finalErrorMessage)
		return fmt.Errorf("%s", renderedError)
	}

	return nil
}
