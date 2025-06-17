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
	t.Run("Successful API call returns expected ability", func(t *testing.T) {
		expectedAbility := structs.AbilityJSONStruct{
			Name: "unaware",
		}

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			err := json.NewEncoder(w).Encode(expectedAbility)
			assert.NoError(t, err, "Expected no error for encoding response")
		}))
		defer ts.Close()

		ability, name, err := AbilityApiCall("/ability", "unaware", ts.URL)

		require.NoError(t, err, "Expected no error on successful API call")
		assert.Equal(t, expectedAbility, ability, "Expected ability struct does not match")
		assert.Equal(t, "unaware", name, "Expected ability name does not match")
	})

	t.Run("Failed API call returns styled error", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Simulate API failure (e.g., 404 Not Found)
			http.Error(w, "Not Found", http.StatusNotFound)
		}))
		defer ts.Close()

		ability, _, err := AbilityApiCall("/ability", "non-existent-ability", ts.URL)

		require.Error(t, err, "Expected an error for invalid ability")
		assert.Equal(t, structs.AbilityJSONStruct{}, ability, "Expected empty ability struct on error")

		assert.Contains(t, err.Error(), "Ability not found", "Expected 'Ability not found' in error message")
		assert.Contains(t, err.Error(), "Perhaps a typo?", "Expected helpful suggestion in error message")
	})
}

func TestMoveApiCall(t *testing.T) {
	t.Run("Successful API call returns expected move", func(t *testing.T) {
		expectedMove := structs.MoveJSONStruct{
			Name: "shadow-ball",
		}

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			err := json.NewEncoder(w).Encode(expectedMove)
			assert.NoError(t, err, "Expected no error for encoding response")
		}))
		defer ts.Close()

		move, name, err := MoveApiCall("/move", "shadow-ball", ts.URL)

		require.NoError(t, err, "Expected no error on successful API call")
		assert.Equal(t, expectedMove, move, "Expected move struct does not match")
		assert.Equal(t, "shadow-ball", name, "Expected move name does not match")
	})

	t.Run("Failed API call returns styled error", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Simulate API failure (e.g., 404 Not Found)
			http.Error(w, "Not Found", http.StatusNotFound)
		}))
		defer ts.Close()

		move, _, err := MoveApiCall("/move", "non-existent-move", ts.URL)

		require.Error(t, err, "Expected an error for invalid move")
		assert.Equal(t, structs.MoveJSONStruct{}, move, "Expected empty move struct on error")

		assert.Contains(t, err.Error(), "Move not found", "Expected 'Move not found' in error message")
		assert.Contains(t, err.Error(), "Perhaps a typo?", "Expected helpful suggestion in error message")
	})
}

func TestPokemonApiCall(t *testing.T) {
	t.Run("Successful API call returns expected pokemon", func(t *testing.T) {
		expectedPokemon := structs.PokemonJSONStruct{
			Name: "flareon",
		}

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			err := json.NewEncoder(w).Encode(expectedPokemon)
			assert.NoError(t, err, "Expected no error for encoding response")
		}))
		defer ts.Close()

		pokemon, name, err := PokemonApiCall("/pokemon", "flareon", ts.URL)

		require.NoError(t, err, "Expected no error on successful API call")
		assert.Equal(t, expectedPokemon, pokemon, "Expected pokemon struct does not match")
		assert.Equal(t, "flareon", name, "Expected pokemon name does not match")
	})

	t.Run("Failed API call returns styled error", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Simulate API failure (e.g., 404 Not Found)
			http.Error(w, "Not Found", http.StatusNotFound)
		}))
		defer ts.Close()

		pokemon, _, err := PokemonApiCall("/pokemon", "non-existent-pokemon", ts.URL)

		require.Error(t, err, "Expected an error for invalid pokemon")
		assert.Equal(t, structs.PokemonJSONStruct{}, pokemon, "Expected empty pokemon struct on error")

		assert.Contains(t, err.Error(), "Pokémon not found", "Expected 'Pokémon not found' in error message")
		assert.Contains(t, err.Error(), "Perhaps a typo?", "Expected helpful suggestion in error message")
	})
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
