package utils

import "github.com/digitalghost-dev/poke-cli/styling"

func FormatError(message string) string {
	return styling.ErrorBorder.Render(
		styling.ErrorColor.Render("✖ Error!"),
		"\n"+message,
	)
}
