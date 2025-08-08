package flags

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/digitalghost-dev/poke-cli/styling"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func LatestFlag() (string, error) {
	var output strings.Builder

	latestRelease(&output)

	result := output.String()
	fmt.Print(result)

	return result, nil
}

func latestRelease(output *strings.Builder) {
	type Release struct {
		TagName string `json:"tag_name"`
	}

	// Parse and validate the URL
	parsedURL, err := url.Parse("https://api.github.com/repos/digitalghost-dev/poke-cli/releases/latest")
	if err != nil {
		fmt.Fprintf(output, "invalid URL: %v\n", err)
		return
	}

	// Implementing gosec error
	if flag.Lookup("test.v") == nil {
		if parsedURL.Scheme != "https" {
			fmt.Fprint(output, "only HTTPS URLs are allowed for security reasons\n")
			return
		}
		if parsedURL.Host != "api.github.com" {
			fmt.Fprint(output, "url host is not allowed\n")
			return
		}
	}

	response, err := http.Get(parsedURL.String())
	if err != nil {
		fmt.Fprintf(output, "error fetching data: %v\n", err)
		return
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Fprintf(output, "error reading response body: %v\n", err)
		return
	}

	var release Release
	if err := json.Unmarshal(body, &release); err != nil {
		fmt.Fprintf(output, "error unmarshalling JSON: %v\n", err)
		return
	}

	releaseString := "Latest available version:"
	releaseTag := styling.ColoredBullet.Render("") + release.TagName

	docStyle := lipgloss.NewStyle().
		Padding(1, 2).
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(lipgloss.AdaptiveColor{Light: "#444", Dark: "#EEE"}).
		Width(30)

	fullDoc := lipgloss.JoinVertical(lipgloss.Top, releaseString, releaseTag)

	output.WriteString(docStyle.Render(fullDoc))
	output.WriteString("\n")
}
