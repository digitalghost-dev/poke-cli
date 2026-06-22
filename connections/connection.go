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

	"github.com/digitalghost-dev/poke-cli/cmd/utils"
	"github.com/digitalghost-dev/poke-cli/structs"
)

const APIURL = "https://pokeapi.co/api/v2/"
const maxAPIResponseBytes = 10 * 1024 * 1024 // 10 MiB

const DefaultHTTPTimeout = 60 * time.Second

var httpClient = NewDefaultHTTPClient()

func NewDefaultHTTPClient() *http.Client {
	return &http.Client{Timeout: DefaultHTTPTimeout}
}

type EndpointResource interface {
	GetResourceName() string
}

type HTTPStatusError struct {
	StatusCode int
	URL        string
}

func (e HTTPStatusError) Error() string {
	return fmt.Sprintf("non-200 response: %d", e.StatusCode)
}

func fetchEndpoint[T EndpointResource](endpoint, resourceName, baseURL, resourceType string) (T, string, error) {
	var zero T
	fullURL := baseURL + endpoint + "/" + resourceName

	var result T
	err := ApiCallSetup(fullURL, &result, false)
	if err != nil {
		return zero, "", formatEndpointError(resourceType, err)
	}

	return result, result.GetResourceName(), nil
}

func formatEndpointError(resourceType string, err error) error {
	var statusErr HTTPStatusError
	if errors.As(err, &statusErr) {
		switch {
		case statusErr.StatusCode == http.StatusNotFound:
			return fmt.Errorf("%s", utils.FormatNotFoundError(resourceType))
		case statusErr.StatusCode >= http.StatusInternalServerError:
			return fmt.Errorf("%s", utils.FormatServerError(resourceType))
		default:
			return fmt.Errorf("%s", utils.FormatFetchError(resourceType, err))
		}
	}

	var urlErr *url.Error
	if errors.As(err, &urlErr) {
		return fmt.Errorf("%s", utils.FormatNetworkError(resourceType))
	}

	var syntaxErr *json.SyntaxError
	var typeErr *json.UnmarshalTypeError
	if errors.As(err, &syntaxErr) || errors.As(err, &typeErr) {
		return fmt.Errorf("%s", utils.FormatUnexpectedDataError(resourceType))
	}

	return fmt.Errorf("%s", utils.FormatFetchError(resourceType, err))
}

// ApiCallSetup Helper function to handle API calls and JSON unmarshalling
func ApiCallSetup(rawURL string, target interface{}, skipHTTPSCheck bool) error {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("invalid URL provided: %w", err)
	}

	if flag.Lookup("test.v") != nil {
		skipHTTPSCheck = true
	}

	if !skipHTTPSCheck && parsedURL.Scheme != "https" {
		return errors.New("only HTTPS URLs are allowed for security reasons")
	}

	body, err := cachedFetch(parsedURL.String())
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, target); err != nil {
		return fmt.Errorf("error unmarshalling JSON: %w", err)
	}

	return nil
}

func directFetch(url string) ([]byte, error) {
	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error making GET request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, HTTPStatusError{StatusCode: resp.StatusCode, URL: url}
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, maxAPIResponseBytes+1))
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}
	if len(body) > maxAPIResponseBytes {
		return nil, fmt.Errorf("response body exceeds %d bytes", maxAPIResponseBytes)
	}
	return body, nil
}

func AbilityApiCall(endpoint, abilityName, baseURL string) (structs.AbilityJSONStruct, string, error) {
	return fetchEndpoint[structs.AbilityJSONStruct](endpoint, abilityName, baseURL, "Ability")
}

func ItemApiCall(endpoint string, itemName string, baseURL string) (structs.ItemJSONStruct, string, error) {
	return fetchEndpoint[structs.ItemJSONStruct](endpoint, itemName, baseURL, "Item")
}

func MoveApiCall(endpoint string, moveName string, baseURL string) (structs.MoveJSONStruct, string, error) {
	return fetchEndpoint[structs.MoveJSONStruct](endpoint, moveName, baseURL, "Move")
}

func PokemonApiCall(endpoint string, pokemonName string, baseURL string) (structs.PokemonJSONStruct, string, error) {
	return fetchEndpoint[structs.PokemonJSONStruct](endpoint, pokemonName, baseURL, "Pokémon")
}

func PokemonSpeciesApiCall(endpoint string, pokemonSpeciesName string, baseURL string) (structs.PokemonSpeciesJSONStruct, string, error) {
	return fetchEndpoint[structs.PokemonSpeciesJSONStruct](endpoint, pokemonSpeciesName, baseURL, "PokémonSpecies")
}

func TypesApiCall(endpoint string, typesName string, baseURL string) (structs.TypesJSONStruct, string, error) {
	return fetchEndpoint[structs.TypesJSONStruct](endpoint, typesName, baseURL, "Type")
}

func CallTCGData(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Add("apikey", "sb_publishable_oondaaAIQC-wafhEiNgpSQ_reRiEp7j")
	req.Header.Add("Authorization", "Bearer sb_publishable_oondaaAIQC-wafhEiNgpSQ_reRiEp7j")
	req.Header.Add("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
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
