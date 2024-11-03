package connections

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestBaseApiCallSuccess - Test for the ApiCallSetup function
func TestBaseApiCallSuccess(t *testing.T) {
	expectedData := map[string]string{"key": "value"}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode(expectedData)
		assert.Nil(t, err)
	}))
	defer ts.Close()

	var target map[string]string

	err := ApiCallSetup(ts.URL, &target)
	if err != nil {
		return
	}

	assert.Equal(t, expectedData, target)
}

// TestPokemonApiCallSuccess - Test for the PokemonApiCall function
func TestPokemonApiCallSuccess(t *testing.T) {
	expectedPokemon := PokemonJSONStruct{
		Name: "pikachu",
		ID:   25,
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
		assert.Nil(t, err)
	}))
	defer ts.Close()

	pokemon, name, id := PokemonApiCall("/pokemon", "pikachu", ts.URL)

	assert.Equal(t, expectedPokemon, pokemon)
	assert.Equal(t, "pikachu", name)
	assert.Equal(t, 25, id)
}

// TestTypesApiCallSuccess - Test for the TypesApiCall function
func TestTypesApiCallSuccess(t *testing.T) {
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
		assert.Nil(t, err)
	}))
	defer ts.Close()

	types, name, id := TypesApiCall("/type", "electric", ts.URL)

	assert.Equal(t, expectedTypes, types)
	assert.Equal(t, "electric", name)
	assert.Equal(t, 13, id)
}

func TestApiCallSetup_NotFound(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Println(w, `{"error": "not found"}`)
	}))
	defer ts.Close()

	var target map[string]string
	err := ApiCallSetup(ts.URL, &target)
	if err != nil {
		return
	}
	// TODO: Add assertions for the output or error message handling
}

func TestPokemonApiCall_UnmarshalError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Println(w, `{"name": "123", "id": "not_a_number"}`) // Partially malformed JSON
	}))
	defer ts.Close()

	var pokemonStruct PokemonJSONStruct
	err := ApiCallSetup(ts.URL, &pokemonStruct)
	assert.NotNil(t, err, "Expected unmarshalling error due to type mismatch")

	var typesStruct TypesJSONStruct
	err = ApiCallSetup(ts.URL, &typesStruct)
	assert.NotNil(t, err, "Expected unmarshalling error due to type mismatch")
}

func TestTypesApiCall_SuccessWithAllFields(t *testing.T) {
	expectedTypes := TypesJSONStruct{
		Name: "electric",
		ID:   13,
		// TODO: Add fields to test complex struct parsing like `DamageRelations`
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode(expectedTypes)
		assert.Nil(t, err)
	}))
	defer ts.Close()

	types, _, _ := TypesApiCall("/type", "electric", ts.URL)
	assert.Equal(t, expectedTypes, types)
}

func TestApiCallSetup_Handles404(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer ts.Close()

	var target map[string]string
	err := ApiCallSetup(ts.URL, &target)

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "404 error")
}
