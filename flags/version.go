package flags

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/digitalghost-dev/poke-cli/styling"
)

func LatestFlag() (string, error) {
	var output strings.Builder

	err := latestRelease(&output)

	result := output.String()
	fmt.Print(result)

	return result, err
}

func latestRelease(output *strings.Builder) error {
	type Release struct {
		TagName string `json:"tag_name"`
	}

	parsedURL, err := url.Parse("https://api.github.com/repos/digitalghost-dev/poke-cli/releases/latest")
	if err != nil {
		err = fmt.Errorf("invalid URL: %w", err)
		fmt.Fprintln(output, err)
		return err
	}

	if flag.Lookup("test.v") == nil {
		if parsedURL.Scheme != "https" {
			err := fmt.Errorf("only HTTPS URLs are allowed for security reasons")
			fmt.Fprintln(output, err)
			return err
		}
		if parsedURL.Host != "api.github.com" {
			err := fmt.Errorf("url host is not allowed")
			fmt.Fprintln(output, err)
			return err
		}
	}

	response, err := http.Get(parsedURL.String())
	if err != nil {
		err = fmt.Errorf("error fetching data: %w", err)
		fmt.Fprintln(output, err)
		return err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		err = fmt.Errorf("error reading response body: %w", err)
		fmt.Fprintln(output, err)
		return err
	}

	var release Release
	if err := json.Unmarshal(body, &release); err != nil {
		err = fmt.Errorf("error unmarshalling JSON: %w", err)
		fmt.Fprintln(output, err)
		return err
	}

	releaseString := "Latest available release on GitHub:"
	releaseTag := styling.ColoredBullet.Render("") + release.TagName

	docStyle := lipgloss.NewStyle().
		Padding(1, 2).
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(lipgloss.AdaptiveColor{Light: "#444", Dark: "#EEE"}).
		Width(30)

	fullDoc := lipgloss.JoinVertical(lipgloss.Top, releaseString, releaseTag)

	output.WriteString(docStyle.Render(fullDoc))
	output.WriteString("\n")

	return nil
}
