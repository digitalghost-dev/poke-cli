use serde::Deserialize;

#[derive(Deserialize, Debug)]
pub struct NamedRef {
    pub name: String,
    pub url: String,
}

#[derive(Deserialize, Debug)]
pub struct RawPokemon {
    pub id: u32,
    pub name: String,
    pub height: u32,
    pub weight: u32,
    pub species: NamedRef,
    pub types: Vec<RawTypingEntry>,
    pub abilities: Vec<RawAbilitiesEntry>,
    pub stats: Vec<RawStatsEntry>,
    pub moves: Vec<RawMovesEntry>,
    pub sprites: RawSprites,
}

#[derive(Deserialize, Debug)]
pub struct RawTypingEntry {
    pub slot: u8,
    // `type` is a Rust keyword, so rename the JSON key onto a legal field name.
    #[serde(rename = "type")]
    pub typing: NamedRef,
}

#[derive(Deserialize, Debug)]
pub struct RawAbilitiesEntry {
    pub ability: NamedRef,
    #[serde(default)]
    pub is_hidden: bool,
    pub slot: u8,
}

#[derive(Deserialize, Debug)]
pub struct RawStatsEntry {
    pub base_stat: u16,
    pub effort: u8,
    pub stat: NamedRef,
}

#[derive(Deserialize, Debug)]
pub struct RawMovesEntry {
    // `move` is a Rust keyword; `r#move` escapes it.
    pub r#move: NamedRef,
    pub version_group_details: Vec<RawVersionGroupDetail>,
}

#[derive(Deserialize, Debug)]
pub struct RawVersionGroupDetail {
    pub level_learned_at: u8,
    pub move_learn_method: NamedRef,
    pub version_group: NamedRef,
}

#[derive(Deserialize, Debug)]
pub struct RawSprites {
    pub front_default: Option<String>,
}

#[derive(Deserialize, Debug)]
pub struct RawPokemonSpecies {
    pub name: String,
    pub egg_groups: Vec<NamedRef>,
    pub gender_rate: i8,
    pub hatch_counter: u8,
    pub evolves_from_species: Option<NamedRef>,
}

#[derive(Deserialize, Debug)]
pub struct RawType {
    pub damage_relations: RawDamageRelations,
}

#[derive(Deserialize, Debug)]
pub struct RawDamageRelations {
    pub double_damage_from: Vec<NamedRef>,
    pub half_damage_from: Vec<NamedRef>,
    pub no_damage_from: Vec<NamedRef>,
}

#[derive(Deserialize, Debug)]
pub struct RawMove {
    pub name: String,
    // `type` is a Rust keyword, so rename the JSON key onto a legal field name.
    #[serde(rename = "type")]
    pub typing: NamedRef,
    pub damage_class: NamedRef, // category lives here
    pub power: Option<u16>,
    pub accuracy: Option<u8>,
    pub pp: Option<u8>,
}

#[cfg(test)]
mod tests {
    use super::*;

    const CHARIZARD: &str = include_str!("../tests/fixtures/charizard.json");
    const CHARIZARD_SPECIES: &str = include_str!("../tests/fixtures/charizard-species.json");

    #[test]
    fn raw_pokemon_deserialize() {
        let pokemon: RawPokemon = serde_json::from_str(CHARIZARD).unwrap();

        assert_eq!(pokemon.id, 6);
        assert_eq!(pokemon.name, "charizard");
        assert_eq!(pokemon.types.len(), 2);
    }

    #[test]
    fn raw_pokemon_species_deserialize() {
        let pokemon_species: RawPokemonSpecies = serde_json::from_str(CHARIZARD_SPECIES).unwrap();

        assert_eq!(pokemon_species.name, "charizard");
        assert!(pokemon_species.evolves_from_species.is_some());
        assert_eq!(pokemon_species.gender_rate, 1);
    }
}
