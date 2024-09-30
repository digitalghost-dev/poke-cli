package connections

import (
	"encoding/json"
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"io"
	"log"
	"net/http"
)

type PokemonJSONStruct struct {
	Name  string `json:"name"`
	ID    int    `json:"id"`
	Types []struct {
		Slot int `json:"slot"`
		Type struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"type"`
	} `json:"types"`
	Abilities []struct {
		Ability struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"ability"`
		Hidden bool `json:"hidden"`
		Slot   int  `json:"slot"`
	} `json:"abilities"`
}

type TypesJSONStruct struct {
	Name    string `json:"name"`
	ID      int    `json:"id"`
	Pokemon []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		Slot int `json:"slot"`
	} `json:"pokemon"`
}

var httpGet = http.Get
var red = lipgloss.Color("#F2055C")
var errorColor = lipgloss.NewStyle().Foreground(red)

// ApiCallSetup Helper function to handle API calls and JSON unmarshalling
func ApiCallSetup(url string, target interface{}) {
	res, err := httpGet(url)
	if err != nil {
		log.Fatalf("Error making GET request: %v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("Failed to close body: %v", err)
		}
	}(res.Body)

	if res.StatusCode == http.StatusNotFound {
		fmt.Println(errorColor.Render("Page not found. 404 error."))
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	err = json.Unmarshal(body, target)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)
	}
}

func PokemonApiCall(endpoint string, pokemonName string, baseURL string) (PokemonJSONStruct, string, int) {

	url := baseURL + endpoint + "/" + pokemonName
	var pokemonStruct PokemonJSONStruct

	ApiCallSetup(url, &pokemonStruct)

	return pokemonStruct, pokemonStruct.Name, pokemonStruct.ID
}

func TypesApiCall(endpoint string, typesName string, baseURL string) (TypesJSONStruct, string, int) {

	url := baseURL + endpoint + "/" + typesName
	var typesStruct TypesJSONStruct

	ApiCallSetup(url, &typesStruct)

	return typesStruct, typesStruct.Name, typesStruct.ID
}
