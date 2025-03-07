package structs

import (
	"encoding/json"
	"testing"
)

func TestPokemonJSONStruct_Unmarshal(t *testing.T) {
	// Sample JSON data for a Pokémon
	jsonData := `{
		"name": "pikachu",
		"id": 25,
		"weight": 60,
		"height": 4,
		"abilities": [
			{
				"ability": {
					"name": "static",
					"url": "https://pokeapi.co/api/v2/ability/9/"
				},
				"hidden": false,
				"slot": 1
			}
		],
		"types": [
			{
				"slot": 1,
				"type": {
					"name": "electric",
					"url": "https://pokeapi.co/api/v2/type/13/"
				}
			}
		],
		"sprites": {
			"front_default": "https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/25.png"
		},
		"stats": [
			{
				"base_stat": 35,
				"stat": {
					"name": "hp"
				}
			}
		]
	}`

	var pokemon PokemonJSONStruct
	err := json.Unmarshal([]byte(jsonData), &pokemon)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	// Assertions
	if pokemon.Name != "pikachu" {
		t.Errorf("Expected name to be 'pikachu', got '%s'", pokemon.Name)
	}
	if pokemon.ID != 25 {
		t.Errorf("Expected ID to be 25, got %d", pokemon.ID)
	}
	if len(pokemon.Abilities) != 1 || pokemon.Abilities[0].Ability.Name != "static" {
		t.Errorf("Expected ability 'static', got '%s'", pokemon.Abilities[0].Ability.Name)
	}
	if len(pokemon.Types) != 1 || pokemon.Types[0].Type.Name != "electric" {
		t.Errorf("Expected type 'electric', got '%s'", pokemon.Types[0].Type.Name)
	}
	if pokemon.Sprites.FrontDefault == "" {
		t.Errorf("Expected a sprite URL but got an empty string")
	}
}
