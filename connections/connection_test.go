package connections

import (
	"encoding/json"
	"github.com/digitalghost-dev/poke-cli/structs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
		assert.NoError(t, err)
	}))
	defer ts.Close()

	var target map[string]string

	// Call ApiCallSetup with skipHTTPSCheck set to true
	err := ApiCallSetup(ts.URL, &target, true)
	require.NoError(t, err, "Expected no error for skipHTTPSCheck")

	assert.Equal(t, expectedData, target, "Expected data does not match the response")

	t.Run("invalid URL", func(t *testing.T) {
		var target map[string]string
		err := ApiCallSetup(":", &target, true)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid URL")
	})

	t.Run("GET request fails", func(t *testing.T) {
		var target map[string]string
		err := ApiCallSetup("https://nonexistent.example.com", &target, true)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "error making GET request")
	})

	t.Run("invalid JSON response", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, err := w.Write([]byte("not-json"))
			assert.NoError(t, err)
		}))
		defer ts.Close()

		var target map[string]string
		err := ApiCallSetup(ts.URL, &target, true)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "error unmarshalling JSON")
	})
}

func TestAbilityApiCall(t *testing.T) {
	expectedAbility := structs.AbilityJSONStruct{
		Name: "Unaware",
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode(expectedAbility)
		assert.NoError(t, err, "Expected no error for skipHTTPSCheck")
	}))
	defer ts.Close()

	_, name, _ := AbilityApiCall("/ability", "unaware", ts.URL)

	assert.Equal(t, "Unaware", name, "Expected name does not match the response")
}

func TestPokemonApiCall(t *testing.T) {
	expectedPokemon := structs.PokemonJSONStruct{
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
		assert.NoError(t, err, "Expected no error for skipHTTPSCheck")
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
	expectedTypes := structs.TypesJSONStruct{
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
		assert.NoError(t, err, "Expected no error for skipHTTPSCheck")
	}))
	defer ts.Close()

	typesStruct, name, id := TypesApiCall("/type", "electric", ts.URL)

	assert.Equal(t, expectedTypes, typesStruct)
	assert.Equal(t, "electric", name)
	assert.Equal(t, 13, id)
}
