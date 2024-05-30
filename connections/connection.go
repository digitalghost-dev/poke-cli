package connections

import (
	"encoding/json"
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"io"
	"log"
	"net/http"
)

var httpGet = http.Get
var red = lipgloss.Color("#F2055C")
var errorColor = lipgloss.NewStyle().Foreground(red)

// Helper function to handle API calls and JSON unmarshalling
func baseApiCall(url string, target interface{}) {
	res, err := httpGet(url)
	if err != nil {
		log.Fatalf("Error making GET request: %v", err)
	}
	defer res.Body.Close()

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

func NameApiCall(pokemonName string, baseURL string) string {
	type Pokemon struct {
		Name string `json:"name"`
	}

	url := baseURL + pokemonName
	var pokemon Pokemon

	baseApiCall(url, &pokemon)

	return pokemon.Name
}

func TypeApiCall(pokemonName string, baseURL string) {
	type Pokemon struct {
		Types []struct {
			Slot int `json:"slot"`
			Type struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"type"`
		} `json:"types"`
	}

	url := baseURL + pokemonName
	var pokemon Pokemon

	baseApiCall(url, &pokemon)

	for _, pokeType := range pokemon.Types {
		fmt.Printf("Type %d: %s\n", pokeType.Slot, pokeType.Type.Name)
	}
}
