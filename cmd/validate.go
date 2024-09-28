package cmd

import (
	"flag"
	"fmt"
)

// ValidatePokemonArgs validates the command line arguments
func ValidatePokemonArgs(args []string) error {

	if len(args) > 5 {
		return fmt.Errorf(errorBorder.Render(errorColor.Render("Error!"), "\nToo many arguments"))
	}

	if len(args) < 3 {
		return fmt.Errorf(errorBorder.Render(errorColor.Render("Error!"), "\nPlease declare a Pokémon's name after the [pokemon] command", "\nRun 'poke-cli pokemon -h' for more details", "\nerror: insufficient arguments"))
	}

	if len(args) > 3 {
		for _, arg := range args[3:] {
			if arg[0] != '-' {
				errorTitle := errorColor.Render("Error!")
				errorString := fmt.Sprintf("\nInvalid argument '%s'. Only flags are allowed after declaring a Pokémon's name", arg)
				formattedString := errorTitle + errorString
				return fmt.Errorf(errorBorder.Render(formattedString))
			}
		}
	}

	if args[2] == "-h" || args[2] == "--help" {
		flag.Usage()
		return fmt.Errorf("")
	}

	return nil
}

// ValidateTypesArgs validates the command line arguments
func ValidateTypesArgs(args []string) error {
	if len(args) > 3 {
		return fmt.Errorf(errorBorder.Render(errorColor.Render("Error!"), "\nToo many arguments"))
	}

	if len(args) == 3 && (args[2] == "-h" || args[2] == "--help") {
		fmt.Println(errorBorder.Render("Error! The only currently available options\nafter [types] command is '-h' or '--help'"))
		return nil
	}

	return nil
}
