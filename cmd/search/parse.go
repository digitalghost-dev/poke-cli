package search

import (
	"encoding/json"
	"fmt"
	"github.com/digitalghost-dev/poke-cli/connections"
	"io/ioutil"
	"log"
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

func do(endpoint string, obj interface{}) error {
	req, err := http.NewRequest(http.MethodGet, connections.APIURL+endpoint, nil)
	if err != nil {
		return err
	}
	client := &http.Client{Timeout: 10 * time.Second}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, &obj)
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
			if (result.Name[0:len(substr)]) != (substr) {
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
	err = do(fmt.Sprintf("%s?offset=0&limit=9999", endpoint), &result)
	result.Results = parseSearch(result.Results, search)
	result.Count = len(result.Results)
	return
}

func SearchResults() {
	endpoint := "pokemon"
	searchTerm := "zac"

	result, err := query(endpoint, searchTerm) // Using the Search function from resources
	if err != nil {
		log.Fatalf("Error fetching search results: %v", err)
	}

	fmt.Println("Search Results:")
	for _, r := range result.Results {
		fmt.Println("-", r.Name)
	}
}
