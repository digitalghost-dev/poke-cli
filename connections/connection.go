package connections

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"io"
	"net/http"
	"net/url"
	"os"
)

var errorBorder = lipgloss.NewStyle().
	BorderStyle(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("#F2055C"))

type AbilitiesJSONStruct struct {
	Name          string `json:"name"`
	EffectEntries []struct {
		Effect   string `json:"effect"`
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
		ShortEffect string `json:"short_effect"`
	} `json:"effect_entries"`
	FlavorTextEntries []struct {
		FlavorText string `json:"flavor_text"`
		Language   struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
		VersionGroup struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"version_group"`
	} `json:"flavor_text_entries"`
	Pokemon []struct {
		Hidden      bool `json:"hidden"`
		PokemonName struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
	} `json:"pokemon"`
}

type PokemonJSONStruct struct {
	Name      string `json:"name"`
	ID        int    `json:"id"`
	Weight    int    `json:"weight"`
	Height    int    `json:"height"`
	Abilities []struct {
		Ability struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"ability"`
		Hidden bool `json:"hidden"`
		Slot   int  `json:"slot"`
	} `json:"abilities"`
	Types []struct {
		Slot int `json:"slot"`
		Type struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"type"`
	} `json:"types"`
	Sprites struct {
		FrontDefault string `json:"front_default"`
	} `json:"sprites"`
	Stats []struct {
		BaseStat int `json:"base_stat"`
		Stat     struct {
			Name string `json:"name"`
		} `json:"stat"`
	} `json:"stats"`
}

type TypesJSONStruct struct {
	Name  string `json:"name"`
	ID    int    `json:"id"`
	Moves []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"moves"`
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

var red = lipgloss.Color("#F2055C")
var errorColor = lipgloss.NewStyle().Foreground(red)

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

	if res.StatusCode == http.StatusNotFound {
		errMessage := errorBorder.Render(
			errorColor.Render("Error!"),
			"\nPokémon or ability not found.\nPerhaps a typo in the name?",
		)
		fmt.Println(errMessage)

		if flag.Lookup("test.v") != nil {
			return fmt.Errorf("page not found: 404 error")
		} else {
			os.Exit(1)
		}
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

func AbilityApiCall(endpoint string, abilityName string, baseURL string) (AbilitiesJSONStruct, string, error) {
	fullURL := baseURL + endpoint + "/" + abilityName

	var abilityStruct AbilitiesJSONStruct
	err := ApiCallSetup(fullURL, &abilityStruct, false)
	if err != nil {
		return AbilitiesJSONStruct{}, "", fmt.Errorf("error fetching ability %s: %w", abilityName, err)
	}

	return abilityStruct, abilityStruct.Name, nil
}

func PokemonApiCall(endpoint string, pokemonName string, baseURL string) (PokemonJSONStruct, string, int, int, int, error) {
	fullURL := baseURL + endpoint + "/" + pokemonName

	var pokemonStruct PokemonJSONStruct
	err := ApiCallSetup(fullURL, &pokemonStruct, false)
	if err != nil {
		fmt.Printf("Error fetching Pokémon: %v\n", err) // Debugging
		return PokemonJSONStruct{}, "", 0, 0, 0, fmt.Errorf("error fetching pokemon: %w", err)
	}

	return pokemonStruct, pokemonStruct.Name, pokemonStruct.ID, pokemonStruct.Weight, pokemonStruct.Height, nil
}

func TypesApiCall(endpoint string, typesName string, baseURL string) (TypesJSONStruct, string, int, error) {

	fullURL := baseURL + endpoint + "/" + typesName
	var typesStruct TypesJSONStruct

	err := ApiCallSetup(fullURL, &typesStruct, false)
	if err != nil {
		return TypesJSONStruct{}, "", 0, fmt.Errorf("error fetching types: %w", err)
	}

	return typesStruct, typesStruct.Name, typesStruct.ID, nil
}
