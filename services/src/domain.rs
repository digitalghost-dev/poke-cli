use serde::Serialize;

#[derive(Serialize, Debug)]
pub struct Pokemon {
    pub id: u32,
    pub name: String,
    pub height: u32,
    pub weight: u32,
    pub species: PokemonSpeciesInfo,
    pub types: Vec<PokemonTyping>,

    // Optional sections are omitted from the JSON entirely when not requested,
    // rather than serialized as `null`.
    #[serde(skip_serializing_if = "Option::is_none")]
    pub abilities: Option<Vec<PokemonAbility>>,

    #[serde(skip_serializing_if = "Option::is_none")]
    pub stats: Option<Vec<PokemonStats>>,

    #[serde(skip_serializing_if = "Option::is_none")]
    pub moves: Option<Vec<LearnableMove>>,

    pub source: ResourceSourceMetadata,

    #[serde(skip_serializing_if = "Option::is_none")]
    pub defenses: Option<TypeDefenseProfile>,
}

#[derive(Serialize, Debug)]
pub struct PokemonTyping {
    pub name: String,
    pub slot: u8,
}

#[derive(Serialize, Debug)]
pub struct PokemonAbility {
    pub name: String,
    pub is_hidden: bool,
}

#[derive(Serialize, Debug)]
pub struct PokemonStats {
    pub name: String,
    pub base_stat: u16,
}

#[derive(Serialize, Debug)]
pub struct PokemonSpeciesInfo {
    pub name: String,
    pub egg_groups: Vec<String>,
    pub gender_rate: i8,
    pub hatch_counter: u8,

    #[serde(skip_serializing_if = "Option::is_none")]
    pub evolves_from: Option<String>,
}

#[derive(Serialize, Debug)]
pub struct ResourceSourceMetadata {
    pub fetched_at: String, // RFC3339
    pub partial_errors: Vec<PartialResourceError>,
}

#[derive(Serialize, Debug)]
pub struct PartialResourceError {
    pub resource: String,
    pub name: String,
    pub error: String,
}

#[derive(Serialize, Debug)]
pub struct TypeDefenseProfile {
    pub weak_to: Vec<TypeEffectiveness>,
    pub resistant_to: Vec<TypeEffectiveness>,
    pub immune_to: Vec<TypeEffectiveness>,
    pub normal_damage: Vec<String>,
}

#[derive(Serialize, Debug)]
pub struct TypeEffectiveness {
    #[serde(rename = "type")]
    pub type_name: String,
    pub multiplier: f32,
}

#[derive(Serialize, Debug)]
pub struct LearnableMove {
    pub name: String,
    pub level: u8,                 // ← from the Pokémon's move list (version_group_details)
    #[serde(rename = "type")]
    pub type_name: String,         // ← from the move detail (RawMove.typing)
    pub category: String,          // ← from the move detail (RawMove.damage_class)
    pub power: Option<u16>,        // ← move detail; null for status moves
    pub accuracy: Option<u8>,      // ← move detail; null for never-miss
    pub pp: Option<u8>,            // ← move detail
    pub learn_method: String,      // ← from the filter (e.g. "level-up")
    pub version_group: String,     // ← from the filter (e.g. "scarlet-violet")
}

#[cfg(test)]
mod tests {
    use super::*;

    fn sample_pokemon() -> Pokemon {
        Pokemon {
            id: 6,
            name: "charizard".to_string(),
            height: 17,
            weight: 905,
            species: PokemonSpeciesInfo {
                name: "charizard".to_string(),
                egg_groups: vec!["monster".to_string(), "dragon".to_string()],
                gender_rate: 1,
                hatch_counter: 20,
                evolves_from: Some("charmeleon".to_string()),
            },
            types: vec![
                PokemonTyping { name: "fire".to_string(), slot: 1 },
                PokemonTyping { name: "flying".to_string(), slot: 2 },
            ],
            abilities: None,
            stats: None,
            moves: None,
            defenses: None,
            source: ResourceSourceMetadata {
                fetched_at: "2026-05-30T00:00:00Z".to_string(),
                partial_errors: vec![],
            },
        }
    }

    #[test]
    fn omits_unrequested_sections() {
        let value = serde_json::to_value(sample_pokemon()).unwrap();

        assert_eq!(value["name"], "charizard");
        assert_eq!(value["types"].as_array().unwrap().len(), 2);

        // None sections are skipped entirely, not serialized as null.
        assert!(value.get("abilities").is_none());
        assert!(value.get("stats").is_none());
    }

    #[test]
    fn includes_requested_sections() {
        let mut pokemon = sample_pokemon();
        pokemon.abilities = Some(vec![PokemonAbility {
            name: "blaze".to_string(),
            is_hidden: false,
        }]);

        let value = serde_json::to_value(pokemon).unwrap();

        assert_eq!(value["abilities"][0]["name"], "blaze");
        assert_eq!(value["abilities"][0]["is_hidden"], false);
    }
}
