package connections

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/digitalghost-dev/poke-cli/structs"
	"io"
	"net/http"
	"net/url"
)

var (
	errorBorder = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#F2055C"))
	errorColor = lipgloss.NewStyle().Foreground(lipgloss.Color("#F2055C"))
)

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

	res, err := http.Get(parsedURL.String())
	if err != nil {
		return fmt.Errorf("error making GET request: %w", err)
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %w", err)
	}

	err = json.Unmarshal(body, target)
	if err != nil {
		return fmt.Errorf("error unmarshalling JSON: %w", err)
	}

	return nil
}

func AbilityApiCall(endpoint string, abilityName string, baseURL string) (structs.AbilityJSONStruct, string, error) {
	fullURL := baseURL + endpoint + "/" + abilityName

	var abilityStruct structs.AbilityJSONStruct
	err := ApiCallSetup(fullURL, &abilityStruct, false)

	if err != nil {
		errMessage := errorBorder.Render(
			errorColor.Render("Error!"),
			"\nAbility not found.\nPerhaps a typo?",
		)
		return structs.AbilityJSONStruct{}, "", fmt.Errorf("%s", errMessage)
	}

	return abilityStruct, abilityStruct.Name, nil
}

func PokemonApiCall(endpoint string, pokemonName string, baseURL string) (structs.PokemonJSONStruct, string, int, int, int, error) {
	fullURL := baseURL + endpoint + "/" + pokemonName

	var pokemonStruct structs.PokemonJSONStruct
	err := ApiCallSetup(fullURL, &pokemonStruct, false)

	if err != nil {
		errMessage := errorBorder.Render(
			errorColor.Render("Error!"),
			"\nPokémon not found.\nPerhaps a typo?",
		)
		return structs.PokemonJSONStruct{}, "", 0, 0, 0, fmt.Errorf("%s", errMessage)
	}

	return pokemonStruct, pokemonStruct.Name, pokemonStruct.ID, pokemonStruct.Weight, pokemonStruct.Height, nil
}

func TypesApiCall(endpoint string, typesName string, baseURL string) (structs.TypesJSONStruct, string, int, error) {
	fullURL := baseURL + endpoint + "/" + typesName
	var typesStruct structs.TypesJSONStruct

	err := ApiCallSetup(fullURL, &typesStruct, false)

	if err != nil {
		errMessage := errorBorder.Render(
			errorColor.Render("Error!"),
			"\nType not found.\nPerhaps a typo?",
		)
		return structs.TypesJSONStruct{}, "", 0, fmt.Errorf("%s", errMessage)
	}

	return typesStruct, typesStruct.Name, typesStruct.ID, nil
}
