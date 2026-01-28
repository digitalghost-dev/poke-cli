package connections

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/digitalghost-dev/poke-cli/structs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

	t.Run("non-200 status code returns error", func(t *testing.T) {
		testCases := []struct {
			name       string
			statusCode int
		}{
			{"404 Not Found", http.StatusNotFound},
			{"500 Internal Server Error", http.StatusInternalServerError},
			{"403 Forbidden", http.StatusForbidden},
			{"503 Service Unavailable", http.StatusServiceUnavailable},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(tc.statusCode)
				}))
				defer ts.Close()

				var target map[string]string
				err := ApiCallSetup(ts.URL, &target, true)
				require.Error(t, err)
				assert.Contains(t, err.Error(), "non-200 response")
				assert.Contains(t, err.Error(), strconv.Itoa(tc.statusCode))
			})
		}
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

		ability, name, err := AbilityApiCall("/ability", "non-existent-ability", ts.URL)

		require.Error(t, err, "Expected an error for invalid ability")
		assert.Equal(t, structs.AbilityJSONStruct{}, ability, "Expected empty ability struct on error")
		assert.Empty(t, name, "Expected empty name string on error")

		assert.Contains(t, err.Error(), "Ability not found", "Expected 'Ability not found' in error message")
		assert.Contains(t, err.Error(), "Perhaps a typo?", "Expected helpful suggestion in error message")
	})
}

func TestItemApiCall(t *testing.T) {
	t.Run("Successful API call returns expected item", func(t *testing.T) {
		expectedItem := structs.ItemJSONStruct{
			Name: "choice-band",
		}

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			err := json.NewEncoder(w).Encode(expectedItem)
			assert.NoError(t, err, "Expected no error for encoding response")
		}))
		defer ts.Close()

		item, name, err := ItemApiCall("/item", "choice-band", ts.URL)

		require.NoError(t, err, "Expected no error on successful API call")
		assert.Equal(t, expectedItem, item, "Expected item struct does not match")
		assert.Equal(t, "choice-band", name, "Expected item name does not match")
	})

	t.Run("Failed API call returns styled error", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Simulate API failure (e.g., 404 Not Found)
			http.Error(w, "Not Found", http.StatusNotFound)
		}))
		defer ts.Close()

		item, name, err := ItemApiCall("/item", "non-existent-item", ts.URL)

		require.Error(t, err, "Expected an error for invalid item")
		assert.Equal(t, structs.ItemJSONStruct{}, item, "Expected empty item struct on error")
		assert.Empty(t, name, "Expected empty name string on error")

		assert.Contains(t, err.Error(), "Item not found", "Expected 'Item not found' in error message")
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

		move, name, err := MoveApiCall("/move", "non-existent-move", ts.URL)

		require.Error(t, err, "Expected an error for invalid move")
		assert.Equal(t, structs.MoveJSONStruct{}, move, "Expected empty move struct on error")
		assert.Empty(t, name, "Expected empty name string on error")

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

		pokemon, name, err := PokemonApiCall("/pokemon", "non-existent-pokemon", ts.URL)

		require.Error(t, err, "Expected an error for invalid pokemon")
		assert.Equal(t, structs.PokemonJSONStruct{}, pokemon, "Expected empty pokemon struct on error")
		assert.Empty(t, name, "Expected empty name string on error")

		assert.Contains(t, err.Error(), "Pokémon not found", "Expected 'Pokémon not found' in error message")
		assert.Contains(t, err.Error(), "Perhaps a typo?", "Expected helpful suggestion in error message")
	})
}

// TestTypesApiCall - Test for the TypesApiCall function
func TestTypesApiCall(t *testing.T) {
	t.Run("Successful API call returns expected type", func(t *testing.T) {
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
			assert.NoError(t, err, "Expected no error for encoding response")
		}))
		defer ts.Close()

		typesStruct, name, err := TypesApiCall("/type", "electric", ts.URL)

		require.NoError(t, err, "Expected no error on successful API call")
		assert.Equal(t, expectedTypes, typesStruct, "Expected types struct does not match")
		assert.Equal(t, "electric", name, "Expected type name does not match")
	})

	t.Run("Failed API call returns styled error", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Simulate API failure (e.g., 404 Not Found)
			http.Error(w, "Not Found", http.StatusNotFound)
		}))
		defer ts.Close()

		typesStruct, name, err := TypesApiCall("/type", "non-existent-type", ts.URL)

		require.Error(t, err, "Expected an error for invalid type")
		assert.Equal(t, structs.TypesJSONStruct{}, typesStruct, "Expected empty types struct on error")
		assert.Empty(t, name, "Expected empty name string on error")

		assert.Contains(t, err.Error(), "Type not found", "Expected 'Type not found' in error message")
		assert.Contains(t, err.Error(), "Perhaps a typo?", "Expected helpful suggestion in error message")
	})
}

func TestPokemonSpeciesApiCall(t *testing.T) {
	t.Run("Successful API call returns expected species", func(t *testing.T) {
		expectedSpecies := structs.PokemonSpeciesJSONStruct{
			Name: "flareon",
		}

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			err := json.NewEncoder(w).Encode(expectedSpecies)
			assert.NoError(t, err, "Expected no error for encoding response")
		}))
		defer ts.Close()

		species, name, err := PokemonSpeciesApiCall("/pokemon-species", "flareon", ts.URL)

		require.NoError(t, err, "Expected no error on successful API call")
		assert.Equal(t, expectedSpecies, species, "Expected species struct does not match")
		assert.Equal(t, "flareon", name, "Expected species name does not match")
	})

	t.Run("Failed API call returns styled error", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Simulate API failure (e.g., 404 Not Found)
			http.Error(w, "Not Found", http.StatusNotFound)
		}))
		defer ts.Close()

		species, name, err := PokemonSpeciesApiCall("/pokemon-species", "non-existent-species", ts.URL)

		require.Error(t, err, "Expected an error for invalid species")
		assert.Equal(t, structs.PokemonSpeciesJSONStruct{}, species, "Expected empty species struct on error")
		assert.Empty(t, name, "Expected empty name string on error")

		assert.Contains(t, err.Error(), "PokémonSpecies not found", "Expected 'PokémonSpecies not found' in error message")
		assert.Contains(t, err.Error(), "Perhaps a typo?", "Expected helpful suggestion in error message")
	})
}

// testSupabaseKey is the publishable API key used in tests.
const testSupabaseKey = "sb_publishable_oondaaAIQC-wafhEiNgpSQ_reRiEp7j"

func TestCallTCGData(t *testing.T) {
	t.Run("sends correct headers and returns body", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Validate headers
			if got := r.Header.Get("apikey"); got != testSupabaseKey {
				t.Fatalf("missing or wrong apikey header: %q", got)
			}
			if got := r.Header.Get("Authorization"); got != "Bearer "+testSupabaseKey {
				t.Fatalf("missing or wrong Authorization header: %q", got)
			}
			if got := r.Header.Get("Content-Type"); got != "application/json" {
				t.Fatalf("missing or wrong Content-Type header: %q", got)
			}

			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"ok":true}`))
		}))
		defer srv.Close()

		body, err := CallTCGData(srv.URL)
		require.NoError(t, err)
		assert.Equal(t, `{"ok":true}`, string(body))
	})

	t.Run("returns error for non-200 status", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "boom", http.StatusInternalServerError)
		}))
		defer srv.Close()

		_, err := CallTCGData(srv.URL)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "unexpected status code: 500")
	})

	t.Run("returns error for bad URL", func(t *testing.T) {
		_, err := CallTCGData("http://%41:80/") // invalid URL host
		require.Error(t, err)
	})
}
