package pokemon

import (
	"bytes"
	"strings"
	"testing"

	"github.com/digitalghost-dev/poke-cli/structs"
	"github.com/stretchr/testify/assert"
)

func TestRenderEntry(t *testing.T) {
	tests := []struct {
		name     string
		species  structs.PokemonSpeciesJSONStruct
		contains string
	}{
		{
			name: "first matching english entry returned",
			species: structs.PokemonSpeciesJSONStruct{
				FlavorTextEntries: []struct {
					FlavorText string `json:"flavor_text"`
					Language   struct {
						Name string `json:"name"`
						URL  string `json:"url"`
					} `json:"language"`
					Version struct {
						Name string `json:"name"`
						URL  string `json:"url"`
					} `json:"version"`
				}{
					{FlavorText: "A scarlet entry.", Language: struct {
						Name string `json:"name"`
						URL  string `json:"url"`
					}{Name: "en"}, Version: struct {
						Name string `json:"name"`
						URL  string `json:"url"`
					}{Name: "scarlet"}},
					{FlavorText: "A shield entry.", Language: struct {
						Name string `json:"name"`
						URL  string `json:"url"`
					}{Name: "en"}, Version: struct {
						Name string `json:"name"`
						URL  string `json:"url"`
					}{Name: "shield"}},
				},
			},
			contains: "A scarlet entry.",
		},
		{
			name: "non-english entries are skipped",
			species: structs.PokemonSpeciesJSONStruct{
				FlavorTextEntries: []struct {
					FlavorText string `json:"flavor_text"`
					Language   struct {
						Name string `json:"name"`
						URL  string `json:"url"`
					} `json:"language"`
					Version struct {
						Name string `json:"name"`
						URL  string `json:"url"`
					} `json:"version"`
				}{
					{FlavorText: "Un texto en español.", Language: struct {
						Name string `json:"name"`
						URL  string `json:"url"`
					}{Name: "es"}, Version: struct {
						Name string `json:"name"`
						URL  string `json:"url"`
					}{Name: "scarlet"}},
				},
			},
			contains: "",
		},
		{
			name:     "empty flavor text entries",
			species:  structs.PokemonSpeciesJSONStruct{},
			contains: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			renderEntry(&buf, tt.species)
			if tt.contains == "" {
				assert.Empty(t, buf.String())
			} else {
				assert.Contains(t, buf.String(), tt.contains)
			}
		})
	}
}

func TestRenderEggInformation(t *testing.T) {
	tests := []struct {
		name     string
		species  structs.PokemonSpeciesJSONStruct
		contains []string
	}{
		{
			name: "known legacy egg group names are modernized",
			species: structs.PokemonSpeciesJSONStruct{
				EggGroups: []struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				}{{Name: "indeterminate"}, {Name: "ground"}},
				GenderRate:   4,
				HatchCounter: 20,
			},
			contains: []string{"Amorphous", "Field", "50% F", "20"},
		},
		{
			name: "genderless pokemon",
			species: structs.PokemonSpeciesJSONStruct{
				EggGroups: []struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				}{{Name: "no-eggs"}},
				GenderRate:   -1,
				HatchCounter: 120,
			},
			contains: []string{"Undiscovered", "Genderless", "120"},
		},
		{
			name: "regular egg group title-cased",
			species: structs.PokemonSpeciesJSONStruct{
				EggGroups: []struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				}{{Name: "monster"}},
				GenderRate:   1,
				HatchCounter: 5,
			},
			contains: []string{"Monster", "12.5% F"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			renderEggInformation(&buf, tt.species)
			output := buf.String()
			for _, s := range tt.contains {
				assert.Contains(t, output, s)
			}
		})
	}
}

func TestRenderEffortValues(t *testing.T) {
	tests := []struct {
		name     string
		pokemon  structs.PokemonJSONStruct
		contains string
		absent   string
	}{
		{
			name: "single EV stat",
			pokemon: structs.PokemonJSONStruct{
				Stats: []struct {
					BaseStat int `json:"base_stat"`
					Effort   int `json:"effort"`
					Stat     struct {
						Name string `json:"name"`
					} `json:"stat"`
				}{
					{Effort: 2, Stat: struct{ Name string `json:"name"` }{Name: "speed"}},
					{Effort: 0, Stat: struct{ Name string `json:"name"` }{Name: "attack"}},
				},
			},
			contains: "2 Spd",
			absent:   "Atk",
		},
		{
			name: "unknown stat name falls back",
			pokemon: structs.PokemonJSONStruct{
				Stats: []struct {
					BaseStat int `json:"base_stat"`
					Effort   int `json:"effort"`
					Stat     struct {
						Name string `json:"name"`
					} `json:"stat"`
				}{
					{Effort: 1, Stat: struct{ Name string `json:"name"` }{Name: "mystery-stat"}},
				},
			},
			contains: "Missing from API",
		},
		{
			name: "no EVs produces empty list",
			pokemon: structs.PokemonJSONStruct{
				Stats: []struct {
					BaseStat int `json:"base_stat"`
					Effort   int `json:"effort"`
					Stat     struct {
						Name string `json:"name"`
					} `json:"stat"`
				}{
					{Effort: 0, Stat: struct{ Name string `json:"name"` }{Name: "hp"}},
				},
			},
			contains: "Effort Values:",
			absent:   "HP",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			renderEffortValues(&buf, tt.pokemon)
			output := buf.String()
			assert.Contains(t, output, tt.contains)
			if tt.absent != "" {
				assert.NotContains(t, output, tt.absent)
			}
		})
	}
}

func TestRenderMetrics(t *testing.T) {
	tests := []struct {
		name     string
		pokemon  structs.PokemonJSONStruct
		contains []string
	}{
		{
			name:    "pikachu-like metrics",
			pokemon: structs.PokemonJSONStruct{ID: 25, Weight: 60, Height: 4},
			contains: []string{
				"National Pokédex #: 25",
				"6.0kg",
				"0.4m",
			},
		},
		{
			name:    "heavy pokemon",
			pokemon: structs.PokemonJSONStruct{ID: 131, Weight: 2160, Height: 20},
			contains: []string{
				"National Pokédex #: 131",
				"216.0kg",
				"2.0m",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			renderMetrics(&buf, tt.pokemon)
			output := buf.String()
			for _, s := range tt.contains {
				assert.Contains(t, output, s)
			}
		})
	}
}

func TestRenderSpecies(t *testing.T) {
	tests := []struct {
		name     string
		species  structs.PokemonSpeciesJSONStruct
		contains string
	}{
		{
			name: "evolves from a previous stage",
			species: structs.PokemonSpeciesJSONStruct{
				EvolvesFromSpecies: struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				}{Name: "pikachu"},
			},
			contains: "Evolves from: Pikachu",
		},
		{
			name:     "basic pokemon with no pre-evolution",
			species:  structs.PokemonSpeciesJSONStruct{},
			contains: "Basic Pokémon",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			renderSpecies(&buf, tt.species)
			assert.Contains(t, buf.String(), tt.contains)
		})
	}
}

func TestRenderTyping(t *testing.T) {
	tests := []struct {
		name     string
		pokemon  structs.PokemonJSONStruct
		contains string
	}{
		{
			name: "fire type renders type name",
			pokemon: structs.PokemonJSONStruct{
				Types: []struct {
					Slot int `json:"slot"`
					Type struct {
						Name string `json:"name"`
						URL  string `json:"url"`
					} `json:"type"`
				}{{Type: struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				}{Name: "fire"}}},
			},
			contains: "Fire",
		},
		{
			name: "unknown type is skipped",
			pokemon: structs.PokemonJSONStruct{
				Types: []struct {
					Slot int `json:"slot"`
					Type struct {
						Name string `json:"name"`
						URL  string `json:"url"`
					} `json:"type"`
				}{{Type: struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				}{Name: "faketype"}}},
			},
			contains: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			renderTyping(&buf, tt.pokemon)
			output := buf.String()
			if tt.contains == "" {
				assert.Empty(t, strings.TrimSpace(output))
			} else {
				assert.Contains(t, output, tt.contains)
			}
		})
	}
}
