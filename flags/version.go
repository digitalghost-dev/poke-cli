package flags

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os/exec"
	"runtime"
)

type commandRunner func(name string, args ...string) *exec.Cmd

func latestDockerImage(run commandRunner) {
	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		cmd = run("powershell", "-Command", `
            $tags = Invoke-RestMethod -Uri "https://hub.docker.com/v2/repositories/digitalghostdev/poke-cli/tags/?page_size=1";
            $tags.results[0].name
        `)
	} else {
		_, err := exec.LookPath("curl")
		if err == nil {
			cmd = run("sh", "-c", `curl -s https://hub.docker.com/v2/repositories/digitalghostdev/poke-cli/tags/?page_size=1 | grep -o '"name":"[^"]*"' | cut -d '"' -f 4`)
		} else {
			cmd = run("sh", "-c", `wget -qO- https://hub.docker.com/v2/repositories/digitalghostdev/poke-cli/tags/?page_size=1 | grep -o '"name":"[^"]*"' | cut -d '"' -f 4`)
		}
	}

	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("Error running command: %v\n", err)
		return
	}

	fmt.Printf("Latest Docker image version: %s", string(output))
}

func latestRelease(githubAPIURL string) {
	type Release struct {
		TagName string `json:"tag_name"`
	}

	// Parse and validate the URL
	parsedURL, err := url.Parse(githubAPIURL)
	if err != nil {
		fmt.Println("Invalid URL:", err)
		return
	}

	// Enforce HTTPS and specific host unless in test mode
	if flag.Lookup("test.v") == nil { // Check if not in test mode
		if parsedURL.Scheme != "https" {
			fmt.Println("Only HTTPS URLs are allowed for security reasons")
			return
		}
		if parsedURL.Host != "api.github.com" {
			fmt.Println("URL host is not allowed")
			return
		}
	}

	// Make the HTTP GET request
	response, err := http.Get(parsedURL.String())
	if err != nil {
		fmt.Println("Error fetching data:", err)
		return
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	var release Release
	if err := json.Unmarshal(body, &release); err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return
	}

	fmt.Println("Latest release tag:", release.TagName)
}

func LatestFlag() {
	// cmd := exec.Command("git", "describe", "--tags", "--abbrev=0")
	latestDockerImage(exec.Command)
	latestRelease("https://api.github.com/repos/digitalghost-dev/poke-cli/releases/latest")
}
