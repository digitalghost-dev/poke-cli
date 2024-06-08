package connections

import (
	"encoding/json"
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"io"
	"log"
	"net/http"
)

type Pokemon struct {
	Name  string `json:"name"`
	ID    int    `json:"id"`
	Types []struct {
		Slot int `json:"slot"`
		Type struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"type"`
	} `json:"types"`
}

var httpGet = http.Get
var red = lipgloss.Color("#F2055C")
var errorColor = lipgloss.NewStyle().Foreground(red)

// Helper function to handle API calls and JSON unmarshalling
func baseApiCall(url string, target interface{}) {
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
		fmt.Println(errorColor.Render("Couldn't find that Pokémon... perhaps its named was misspelled?"))
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

func PokemonNameApiCall(pokemonName string, baseURL string) (string, int) {

	url := baseURL + pokemonName
	var pokemonStruct Pokemon

	baseApiCall(url, &pokemonStruct)

	return pokemonStruct.Name, pokemonStruct.ID
}

func PokemonTypeApiCall(pokemonName string, baseURL string) Pokemon {

	url := baseURL + pokemonName
	var pokemonStruct Pokemon

	baseApiCall(url, &pokemonStruct)

	return pokemonStruct
}
