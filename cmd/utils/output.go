package utils

import (
	"fmt"
	"os"
)

// HandleCommandOutput takes a function that returns (string, error) and wraps it in a no-argument
// function that writes the returned string to stdout if there's no error, or to stderr if there is.
// It returns an exit code: 0 on success, 1 on error.
func HandleCommandOutput(fn func() (string, error)) func() int {
	return func() int {
		output, err := fn()
		if err != nil {
			fmt.Fprintln(os.Stderr, output)
			return 1
		}
		fmt.Println(output)
		return 0
	}
}
