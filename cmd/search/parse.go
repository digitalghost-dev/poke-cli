package search

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/digitalghost-dev/poke-cli/connections"
	"github.com/schollz/closestmatch"
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

func containsRegexChars(s string) bool {
	return strings.ContainsAny(s, "^$.*+?[]{}()|\\")
}

func parseRegex(results []Result, pattern string) ([]Result, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("invalid regex pattern: %w", err)
	}

	filtered := make([]Result, 0)
	for _, r := range results {
		if re.MatchString(r.Name) {
			filtered = append(filtered, r)
		}
	}
	return filtered, nil
}

func parseFuzzy(results []Result, search string) []Result {
	name := make([]string, len(results))
	for i, r := range results {
		name[i] = r.Name
	}

	bagSizes := []int{2, 3, 4}
	cm := closestmatch.New(name, bagSizes)

	matches := cm.ClosestN(search, 20)

	matchSet := make(map[string]struct{}, len(matches))
	for _, m := range matches {
		matchSet[m] = struct{}{}
	}

	filtered := results[:0]
	for _, r := range results {
		if _, ok := matchSet[r.Name]; ok {
			filtered = append(filtered, r)
		}
	}
	return filtered
}

func parseSearch(results []Result, search string) ([]Result, error) {
	if containsRegexChars(search) {
		return parseRegex(results, search)
	}

	return parseFuzzy(results, search), nil
}

var apiCall = connections.ApiCallSetup // set as a var for testability

// Search returns resources list, filtered by resources term.
func query(endpoint string, search string) (result Resource, err error) {
	url := connections.APIURL + endpoint + "/?offset=0&limit=9999"
	err = apiCall(url, &result, false)
	if err != nil {
		return
	}
	result.Results, err = parseSearch(result.Results, search)
	if err != nil {
		return
	}
	result.Count = len(result.Results)
	return
}
