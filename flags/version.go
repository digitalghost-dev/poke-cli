package flags

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"charm.land/lipgloss/v2"
	"github.com/digitalghost-dev/poke-cli/cmd/utils"
	"github.com/digitalghost-dev/poke-cli/connections"
	"github.com/digitalghost-dev/poke-cli/styling"
)

const maxLatestReleaseBytes = 1 * 1024 * 1024 // 1 MiB

var latestReleaseHTTPClient = connections.NewDefaultHTTPClient()

func LatestFlag() (string, error) {
	var output strings.Builder

	err := latestRelease(&output)

	result := output.String()
	fmt.Print(result)

	return result, err
}

func latestRelease(output *strings.Builder) error {
	return latestReleaseFromURL(output, "https://api.github.com/repos/digitalghost-dev/poke-cli/releases/latest", latestReleaseHTTPClient)
}

func latestReleaseFromURL(output *strings.Builder, releaseURL string, client *http.Client) error {
	type Release struct {
		TagName string `json:"tag_name"`
	}

	parsedURL, err := url.Parse(releaseURL)
	if err != nil {
		err = fmt.Errorf("invalid URL: %w", err)
		fmt.Fprintln(output, err)
		return err
	}

	if flag.Lookup("test.v") == nil {
		if parsedURL.Scheme != "https" {
			err := errors.New("only HTTPS URLs are allowed for security reasons")
			fmt.Fprintln(output, err)
			return err
		}
		if parsedURL.Host != "api.github.com" {
			err := errors.New("url host is not allowed")
			fmt.Fprintln(output, err)
			return err
		}
	}

	req, err := http.NewRequest(http.MethodGet, parsedURL.String(), nil)
	if err != nil {
		err = fmt.Errorf("error creating request: %w", err)
		fmt.Fprintln(output, err)
		return err
	}
	req.Header.Set("User-Agent", "poke-cli")

	response, err := client.Do(req)
	if err != nil {
		err = fmt.Errorf("error fetching data: %w", err)
		fmt.Fprintln(output, err)
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		err := gitHubStatusError(response)
		fmt.Fprintln(output, err)
		return err
	}

	body, err := io.ReadAll(io.LimitReader(response.Body, maxLatestReleaseBytes+1))
	if err != nil {
		err = fmt.Errorf("error reading response body: %w", err)
		fmt.Fprintln(output, err)
		return err
	}
	if len(body) > maxLatestReleaseBytes {
		err := fmt.Errorf("response body exceeds %d bytes", maxLatestReleaseBytes)
		fmt.Fprintln(output, err)
		return err
	}

	var release Release
	if err := json.Unmarshal(body, &release); err != nil {
		err = fmt.Errorf("error unmarshalling JSON: %w", err)
		fmt.Fprintln(output, err)
		return err
	}
	if release.TagName == "" {
		err := errors.New("latest release response did not include a tag name")
		fmt.Fprintln(output, err)
		return err
	}

	releaseString := "Latest available release on GitHub:"
	releaseTag := styling.ColoredBullet.Render("") + release.TagName

	isDark := styling.HasDarkBackground()
	ld := lipgloss.LightDark(isDark)
	docStyle := lipgloss.NewStyle().
		Padding(1, 2).
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(ld(lipgloss.Color("#444"), lipgloss.Color("#EEE"))).
		Width(32)

	fullDoc := lipgloss.JoinVertical(lipgloss.Top, releaseString, releaseTag)

	output.WriteString(docStyle.Render(fullDoc))
	output.WriteString("\n")

	return nil
}

func gitHubStatusError(response *http.Response) error {
	rateLimited := (response.StatusCode == http.StatusForbidden || response.StatusCode == http.StatusTooManyRequests) &&
		response.Header.Get("X-RateLimit-Remaining") == "0"
	if !rateLimited {
		return fmt.Errorf("unexpected GitHub response status: %d", response.StatusCode)
	}

	msg := "GitHub API rate limit reached (60 requests/hour for unauthenticated requests)."
	if reset := response.Header.Get("X-RateLimit-Reset"); reset != "" {
		if secs, err := strconv.ParseInt(reset, 10, 64); err == nil {
			msg += "\nTry again after " + time.Unix(secs, 0).Format("3:04 PM") + "."
		}
	}
	return errors.New(utils.FormatError(msg))
}
