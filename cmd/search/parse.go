package search

import (
	"encoding/json"
	"fmt"
	"github.com/digitalghost-dev/poke-cli/connections"
	"io"
	"net/http"
	"strings"
	"time"
)

type Resource struct {
	Count    int         `json:"count"`
	Next     string      `json:"next"`
	Previous interface{} `json:"previous"`
	Results  []Result    `json:"results"`
}

// Result is a resources list result.
type Result struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

func callAPI(endpoint string, obj interface{}) error {
	url := connections.APIURL + endpoint
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("User-Agent", "poke-cli/1.0")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request to %s failed: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body) // error intentionally ignored here
		return fmt.Errorf("API returned status %d (%s): %s", resp.StatusCode, http.StatusText(resp.StatusCode), string(bodyBytes))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if err := json.Unmarshal(body, &obj); err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %w\nResponse body: %s", err, string(body))
	}

	return nil
}

func parseSearch(results []Result, search string) []Result {
	var x int
	var substr string

	for _, result := range results {
		if string(search[0]) == "^" {
			substr = search[1:]
			if len(substr) > len(result.Name) {
				continue
			}
			if result.Name[0:len(substr)] != substr {
				continue
			}
		} else {
			if !strings.Contains(result.Name, search) {
				continue
			}
		}
		results[x] = result
		x++
	}
	return results[:x]
}

// Search returns resources list, filtered by resources term.
func query(endpoint string, search string) (result Resource,
	err error) {
	err = callAPI(fmt.Sprintf("%s?offset=0&limit=9999", endpoint), &result)
	result.Results = parseSearch(result.Results, search)
	result.Count = len(result.Results)
	return
}
