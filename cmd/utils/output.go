package utils

import (
	"fmt"
	"os"
)

// HandleCommandOutput wraps a function that returns (string, error) into a no-arg function
// that prints the output to stdout or stderr depending on whether an error occurred.
func HandleCommandOutput(fn func() (string, error)) func() {
	return func() {
		output, err := fn()
		if err != nil {
			fmt.Fprintln(os.Stderr, output)
			return
		}
		fmt.Println(output)
	}
}
