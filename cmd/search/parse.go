package search

import (
	"github.com/digitalghost-dev/poke-cli/connections"
	"strings"
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

var apiCall = connections.ApiCallSetup // set as a var for testability

// Search returns resources list, filtered by resources term.
func query(endpoint string, search string) (result Resource, err error) {

	url := connections.APIURL + endpoint + "/?offset=0&limit=9999"
	err = apiCall(url, &result, false)
	if err != nil {
		return
	}
	result.Results = parseSearch(result.Results, search)
	result.Count = len(result.Results)
	return
}
