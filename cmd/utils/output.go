package utils

import (
	"fmt"
	"os"
	"strings"
)

// HandleCommandOutput takes a function that returns (string, error) and wraps it in a no-argument
// function that writes the returned string to stdout if there's no error, or to stderr if there is.
// It returns an exit code: 0 on success, 1 on error.
func HandleCommandOutput(fn func([]string) (string, error), args []string) func() int {
	return func() int {
		output, err := fn(args)
		if err != nil {
			fmt.Fprintln(os.Stderr, output)
			return 1
		}
		fmt.Println(output)
		return 0
	}
}

func HandleFlagError(output *strings.Builder, err error) (string, error) {
	fmt.Fprintf(output, "error parsing flags: %v\n", err)
	return "", fmt.Errorf("error parsing flags: %w", err)
}

func CheckHelpFlag(args []string, usageFunc func()) bool {
	if len(args) == 2 && (args[1] == "-h" || args[1] == "--help") {
		usageFunc()
		return true
	}
	return false
}

func WrapText(text string, width int) string {
	words := strings.Fields(text)
	if len(words) == 0 {
		return text
	}

	var result strings.Builder
	lineLength := 0

	for _, word := range words {
		wordLen := len(word)

		if lineLength > 0 && lineLength+1+wordLen > width {
			result.WriteString("\n")
			lineLength = 0
		}

		if lineLength > 0 {
			result.WriteString(" ")
			lineLength++
		}

		result.WriteString(word)
		lineLength += wordLen
	}

	return result.String()
}
