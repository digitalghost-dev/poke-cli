package utils

import "github.com/digitalghost-dev/poke-cli/styling"

func FormatError(message string) string {
	return styling.ErrorBorder.Render(
		styling.ErrorColor.Render("✖ Error!"),
		"\n"+message,
	)
}

func FormatNotFoundError(resourceType string) string {
	return FormatError(resourceType + " not found.\n• Perhaps a typo?\n• Missing a hyphen instead of a space?")
}

func FormatNetworkError(resourceType string) string {
	return FormatError("Could not reach " + resourceType + " data.\nCheck your connection and try again.")
}

func FormatServerError(resourceType string) string {
	return FormatError(resourceType + " data source returned a server error.\nPlease try again later.")
}

func FormatUnexpectedDataError(resourceType string) string {
	return FormatError(resourceType + " data source returned data in an unexpected format.")
}

func FormatFetchError(resourceType string, err error) string {
	if err == nil {
		return FormatError("Could not fetch " + resourceType + " data.")
	}
	return FormatError("Could not fetch " + resourceType + " data.\n" + err.Error())
}
