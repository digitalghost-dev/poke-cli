package utils

import (
	"fmt"
	"strings"
)

type Validator struct {
	MaxArgs     int
	CmdName     string
	RequireName bool
	HasFlags    bool
}

// checkLength checks if the number of arguments is lower than the max value.  Helper Function.
func checkLength(args []string, max int) error {
	if len(args) > max {
		return fmt.Errorf("%s", FormatError("Too many arguments"))
	}
	return nil
}

// checkNoOtherOptions checks if there are exactly 3 arguments and the third argument is neither '-h' nor '--help'
func checkNoOtherOptions(args []string, max int, commandName string) error {
	if len(args) == max && args[2] != "-h" && args[2] != "--help" {
		return fmt.Errorf("%s", FormatError(fmt.Sprintf("The only available options after the\n<%s> command are '-h' or '--help'", commandName)))
	}
	return nil
}

func ValidateArgs(args []string, v Validator) error {
	if err := checkLength(args, v.MaxArgs); err != nil {
		return err
	}
	if v.RequireName && len(args) == 2 {
		return fmt.Errorf("%s", FormatError(fmt.Sprintf(
			"Please declare a(n) %s's name after the <%s> command\nRun 'poke-cli %s -h' for more details\nerror: insufficient arguments",
			v.CmdName, v.CmdName, v.CmdName,
		)))
	}
	if !v.HasFlags && !v.RequireName {
		if err := checkNoOtherOptions(args, v.MaxArgs, v.CmdName); err != nil {
			return err
		}
	}
	return nil
}

// ValidatePokemonArgs validates the command line arguments
func ValidatePokemonArgs(args []string) error {
	// Check if the number of arguments is less than 3
	if len(args) < 3 {
		return fmt.Errorf("%s", FormatError(
			"Please declare a Pokémon's name after the <pokemon> command\nRun 'poke-cli pokemon -h' for more details\nerror: insufficient arguments",
		))
	}

	if err := checkLength(args, 8); err != nil {
		return err
	}

	printImageFlagError := func() error {
		return fmt.Errorf("%s", FormatError("The image flag (-i or --image) requires a non-empty value.\nValid sizes are: lg, md, sm."))
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
				return fmt.Errorf("%s", FormatError(fmt.Sprintf("Empty flag '%s'.\nPlease specify valid flag(s).", arg)))
			}

			// Check if the argument after Pokémon's name is an attempted flag
			if arg[0] != '-' {
				return fmt.Errorf("%s", FormatError(fmt.Sprintf("Invalid argument '%s'.\nOnly flags are allowed after declaring a Pokémon's name", arg)))
			}
		}
	}

	return nil
}
