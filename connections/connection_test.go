package connections

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestApiCallSetup - Test for the ApiCallSetup function
func TestApiCallSetup(t *testing.T) {
	expectedData := map[string]string{"key": "value"}

	// Create a test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode(expectedData)
		assert.Nil(t, err)
	}))
	defer ts.Close()

	var target map[string]string

	// Call ApiCallSetup with skipHTTPSCheck set to true
	err := ApiCallSetup(ts.URL, &target, true)
	assert.Nil(t, err, "Expected no error for skipHTTPSCheck")

	assert.Equal(t, expectedData, target, "Expected data does not match the response")
}

func TestAbilityApiCall(t *testing.T) {
	expectedAbility := AbilityJSONStruct{
		Name: "Unaware",
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode(expectedAbility)
		assert.Nil(t, err, "Expected no error for skipHTTPSCheck")
	}))
	defer ts.Close()

	_, name, _ := AbilityApiCall("/ability", "unaware", ts.URL)

	assert.Equal(t, "Unaware", name, "Expected name does not match the response")
}

func TestPokemonApiCall(t *testing.T) {
	expectedPokemon := PokemonJSONStruct{
		Name:   "pikachu",
		ID:     25,
		Weight: 60,
		Height: 4,
		Types: []struct {
			Slot int `json:"slot"`
			Type struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"type"`
		}{
			{Slot: 1, Type: struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			}{Name: "electric", URL: "https://pokeapi.co/api/v2/type/13/"}},
		},
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode(expectedPokemon)
		assert.Nil(t, err, "Expected no error for skipHTTPSCheck")
	}))
	defer ts.Close()

	pokemonStruct, name, id, weight, height, _ := PokemonApiCall("/pokemon", "pikachu", ts.URL)

	assert.Equal(t, expectedPokemon, pokemonStruct, "Expected Pokémon struct does not match")
	assert.Equal(t, "pikachu", name, "Expected name does not match")
	assert.Equal(t, 25, id, "Expected ID does not match")
	assert.Equal(t, 60, weight, "Expected weight does not match")
	assert.Equal(t, 4, height, "Expected height does not match")
}

// TestTypesApiCall - Test for the TypesApiCall function
func TestTypesApiCall(t *testing.T) {
	expectedTypes := TypesJSONStruct{
		Name: "electric",
		ID:   13,
		Pokemon: []struct {
			Pokemon struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"pokemon"`
			Slot int `json:"slot"`
		}{
			{Pokemon: struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			}{Name: "pikachu", URL: "https://pokeapi.co/api/v2/pokemon/25/"},
				Slot: 1},
		},
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode(expectedTypes)
		assert.Nil(t, err, "Expected no error for skipHTTPSCheck")
	}))
	defer ts.Close()

	typesStruct, name, id, _ := TypesApiCall("/type", "electric", ts.URL)

	assert.Equal(t, expectedTypes, typesStruct)
	assert.Equal(t, "electric", name)
	assert.Equal(t, 13, id)
}
