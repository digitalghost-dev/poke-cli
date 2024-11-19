package flags

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"
)

func latestDockerImage(fullCommand string) {
	cmd := exec.Command("bash", "-c", fullCommand)
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("error running command: %v\n", err)
		return
	}

	fmt.Print("Latest Docker image version: ", string(output))
}

func latestRelease(githubAPIURL string) {
	type Release struct {
		TagName string `json:"tag_name"`
	}

	response, err := http.Get(githubAPIURL)
	if err != nil {
		fmt.Println("Error fetching data:", err)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("Error closing body:", err)
		}
	}(response.Body)

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
	latestDockerImage(`curl -s https://hub.docker.com/v2/repositories/digitalghostdev/poke-cli/tags/?page_size=1 | grep -o '"name":"[^"]*"' | cut -d '"' -f 4`)
	latestRelease("https://api.github.com/repos/digitalghost-dev/poke-cli/releases/latest")
}
