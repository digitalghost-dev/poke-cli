package connections

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/digitalghost-dev/poke-cli/structs"
	"github.com/digitalghost-dev/poke-cli/styling"
)

const APIURL = "https://pokeapi.co/api/v2/"

var httpClient = &http.Client{Timeout: 30 * time.Second}

type EndpointResource interface {
	GetResourceName() string
}

func FetchEndpoint[T EndpointResource](endpoint, resourceName, baseURL, resourceType string) (T, string, error) {
	var zero T
	fullURL := baseURL + endpoint + "/" + resourceName

	var result T
	err := ApiCallSetup(fullURL, &result, false)

	if err != nil {
		errMessage := styling.ErrorBorder.Render(
			styling.ErrorColor.Render("✖ Error!"),
			fmt.Sprintf("\n%s not found.\n• Perhaps a typo?\n• Missing a hyphen instead of a space?", resourceType),
		)
		return zero, "", fmt.Errorf("%s", errMessage)
	}

	return result, result.GetResourceName(), nil
}

// ApiCallSetup Helper function to handle API calls and JSON unmarshalling
func ApiCallSetup(rawURL string, target interface{}, skipHTTPSCheck bool) error {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("invalid URL provided: %w", err)
	}

	// Check if running in a test environment
	if flag.Lookup("test.v") != nil {
		skipHTTPSCheck = true
	}

	if !skipHTTPSCheck && parsedURL.Scheme != "https" {
		return errors.New("only HTTPS URLs are allowed for security reasons")
	}

	resp, err := httpClient.Get(parsedURL.String())
	if err != nil {
		return fmt.Errorf("error making GET request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("non-200 response: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %w", err)
	}

	err = json.Unmarshal(body, target)
	if err != nil {
		return fmt.Errorf("error unmarshalling JSON: %w", err)
	}

	return nil
}

func AbilityApiCall(endpoint, abilityName, baseURL string) (structs.AbilityJSONStruct, string, error) {
	return FetchEndpoint[structs.AbilityJSONStruct](endpoint, abilityName, baseURL, "Ability")
}

func ItemApiCall(endpoint string, itemName string, baseURL string) (structs.ItemJSONStruct, string, error) {
	return FetchEndpoint[structs.ItemJSONStruct](endpoint, itemName, baseURL, "Item")
}

func MoveApiCall(endpoint string, moveName string, baseURL string) (structs.MoveJSONStruct, string, error) {
	return FetchEndpoint[structs.MoveJSONStruct](endpoint, moveName, baseURL, "Move")
}

func PokemonApiCall(endpoint string, pokemonName string, baseURL string) (structs.PokemonJSONStruct, string, error) {
	return FetchEndpoint[structs.PokemonJSONStruct](endpoint, pokemonName, baseURL, "Pokémon")
}

func PokemonSpeciesApiCall(endpoint string, pokemonSpeciesName string, baseURL string) (structs.PokemonSpeciesJSONStruct, string, error) {
	return FetchEndpoint[structs.PokemonSpeciesJSONStruct](endpoint, pokemonSpeciesName, baseURL, "PokémonSpecies")
}

func TypesApiCall(endpoint string, typesName string, baseURL string) (structs.TypesJSONStruct, string, error) {
	return FetchEndpoint[structs.TypesJSONStruct](endpoint, typesName, baseURL, "Type")
}

func CallTCGData(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Add("apikey", "sb_publishable_oondaaAIQC-wafhEiNgpSQ_reRiEp7j")
	req.Header.Add("Authorization", "Bearer sb_publishable_oondaaAIQC-wafhEiNgpSQ_reRiEp7j")
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making GET request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	return body, nil
}