package utils

import (
	"fmt"
	"strings"

	"github.com/digitalghost-dev/poke-cli/styling"
)

// checkLength checks if the number of arguments is lower than the max value.  Helper Function.
func checkLength(args []string, max int) error {
	if len(args) > max {
		errMessage := styling.ErrorBorder.Render(
			styling.ErrorColor.Render("✖ Error!") + "\nToo many arguments",
		)
		return fmt.Errorf("%s", errMessage)
	}
	return nil
}

// checkNoOtherOptions checks if there are exactly 3 arguments and the third argument is neither '-h' nor '--help'
func checkNoOtherOptions(args []string, max int, commandName string) error {
	if len(args) == max && args[2] != "-h" && args[2] != "--help" {
		errMsg := styling.ErrorColor.Render("✖ Error!") +
			"\nThe only available options after the\n" + commandName + " command are '-h' or '--help'"
		return fmt.Errorf("%s", styling.ErrorBorder.Render(errMsg))
	}
	return nil
}

// ValidateAbilityArgs validates the command line arguments
func ValidateAbilityArgs(args []string) error {
	if err := checkLength(args, 4); err != nil {
		return err
	}

	if len(args) == 2 {
		errMessage := styling.ErrorBorder.Render(styling.ErrorColor.Render("✖ Error!"), "\nPlease specify an ability")
		return fmt.Errorf("%s", errMessage)
	}

	return nil
}

func ValidateBerryArgs(args []string) error {
	if err := checkLength(args, 3); err != nil {
		return err
	}

	if err := checkNoOtherOptions(args, 3, "<berry>"); err != nil {
		return err
	}

	return nil
}

// ValidateItemArgs validates the command line arguments
func ValidateItemArgs(args []string) error {
	if err := checkLength(args, 3); err != nil {
		return err
	}

	if len(args) == 2 {
		errMessage := styling.ErrorBorder.Render(styling.ErrorColor.Render("✖ Error!"), "\nPlease specify an item ")
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
		errMessage := styling.ErrorBorder.Render(styling.ErrorColor.Render("✖ Error!"), "\nPlease specify a move ")
		return fmt.Errorf("%s", errMessage)
	}

	return nil
}

// ValidateNaturesArgs validates the command line arguments
func ValidateNaturesArgs(args []string) error {
	if err := checkLength(args, 3); err != nil {
		return err
	}

	if err := checkNoOtherOptions(args, 3, "<natures>"); err != nil {
		return err
	}

	return nil
}

// ValidatePokemonArgs validates the command line arguments
func ValidatePokemonArgs(args []string) error {
	// Check if the number of arguments is less than 3
	if len(args) < 3 {
		errMessage := styling.ErrorBorder.Render(
			styling.ErrorColor.Render("✖ Error!"),
			"\nPlease declare a Pokémon's name after the <pokemon> command",
			"\nRun 'poke-cli pokemon -h' for more details",
			"\nerror: insufficient arguments",
		)
		return fmt.Errorf("%s", errMessage)
	}

	if err := checkLength(args, 8); err != nil {
		return err
	}

	printImageFlagError := func() error {
		msg := styling.ErrorBorder.Render(
			styling.ErrorColor.Render("✖ Error!") +
				"\nThe image flag (-i or --image) requires a non-empty value.\nValid sizes are: lg, md, sm.",
		)
		return fmt.Errorf("%s", msg)
	}

	for _, arg := range args {
		switch {
		case strings.HasPrefix(arg, "-i=") && strings.TrimPrefix(arg, "-i=") == "":
			return printImageFlagError()
		case strings.HasPrefix(arg, "--image=") && strings.TrimPrefix(arg, "--image=") == "":
			return printImageFlagError()
		case strings.HasPrefix(arg, "-image=") && strings.TrimPrefix(arg, "-image=") == "":
			return printImageFlagError()
		}
	}

	// Validate each argument after the Pokémon's name
	if len(args) > 3 {
		for _, arg := range args[3:] {
			// Check for an empty flag after Pokémon's name
			if arg == "-" || arg == "--" {
				errorTitle := styling.ErrorColor.Render("✖ Error!")
				errorString := fmt.Sprintf(
					"\nEmpty flag '%s'.\nPlease specify valid flag(s).",
					arg,
				)
				finalErrorMessage := errorTitle + errorString
				renderedError := styling.ErrorBorder.Render(finalErrorMessage)
				return fmt.Errorf("%s", renderedError)
			}

			// Check if the argument after Pokémon's name is an attempted flag
			if arg[0] != '-' {
				errorTitle := styling.ErrorColor.Render("✖ Error!")
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

	if err := checkNoOtherOptions(args, 3, "<search>"); err != nil {
		return err
	}

	return nil
}

// ValidateSpeedArgs validates the command line arguments
func ValidateSpeedArgs(args []string) error {
	if err := checkLength(args, 3); err != nil {
		return err
	}

	if err := checkNoOtherOptions(args, 3, "<speed>"); err != nil {
		return err
	}

	return nil
}

// ValidateTypesArgs validates the command line arguments
func ValidateTypesArgs(args []string) error {
	if err := checkLength(args, 3); err != nil {
		return err
	}

	if err := checkNoOtherOptions(args, 3, "<types>"); err != nil {
		return err
	}

	return nil
}
