package connections

import (
	"encoding/json"
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"io"
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
	DamageRelations struct {
		DoubleDamageFrom []struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"double_damage_from"`
		DoubleDamageTo []struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"double_damage_to"`
		HalfDamageFrom []struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"half_damage_from"`
		HalfDamageTo []struct {
			Name string `json:"name"`
			URL  string `json:"ul"`
		} `json:"half_damage_to"`
		NoDamageFrom []struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"no_damage_from"`
		NoDamageTo []struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"no_damage_to"`
	} `json:"damage_relations"`
}

var httpGet = http.Get
var red = lipgloss.Color("#F2055C")
var errorColor = lipgloss.NewStyle().Foreground(red)

// ApiCallSetup Helper function to handle API calls and JSON unmarshalling
func ApiCallSetup(url string, target interface{}) error {
	res, err := httpGet(url)
	if err != nil {
		return fmt.Errorf("error making GET request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
		}
	}(res.Body)

	if res.StatusCode == http.StatusNotFound {
		fmt.Println(errorColor.Render("Page not found. 404 error."))
		return fmt.Errorf("page not found: 404 error")
	}

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

func PokemonApiCall(endpoint string, pokemonName string, baseURL string) (PokemonJSONStruct, string, int) {

	url := baseURL + endpoint + "/" + pokemonName
	var pokemonStruct PokemonJSONStruct

	err := ApiCallSetup(url, &pokemonStruct)
	if err != nil {
		return PokemonJSONStruct{}, "", 0
	}

	return pokemonStruct, pokemonStruct.Name, pokemonStruct.ID
}

func TypesApiCall(endpoint string, typesName string, baseURL string) (TypesJSONStruct, string, int) {

	url := baseURL + endpoint + "/" + typesName
	var typesStruct TypesJSONStruct

	err := ApiCallSetup(url, &typesStruct)
	if err != nil {
		return TypesJSONStruct{}, "", 0
	}

	return typesStruct, typesStruct.Name, typesStruct.ID
}
