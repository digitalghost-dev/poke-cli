use serde::de::DeserializeOwned;
use crate::api::{RawPokemon, RawPokemonSpecies};
use crate::domain::{PokemonTyping, PokemonAbility, PokemonStats, PokemonSpeciesInfo};
use crate::domain::{Pokemon, ResourceSourceMetadata};
use std::time::{SystemTime, UNIX_EPOCH};

pub struct ProfileOptions {
    pub abilities: bool,
    pub defense: bool,
    pub image: Option<String>,
    pub moves: bool,
    pub stats: bool,
}

pub fn run(name: &str, opts: &ProfileOptions) -> anyhow::Result<Pokemon> {
    get_pokemon_profile(name, opts)
}

// functions focused on retrieving data
fn fetch_json<T: DeserializeOwned>(url: &str) -> anyhow::Result<T> {
    let data: T = reqwest::blocking::get(url)?
        .error_for_status()?
        .json::<T>()?;

    Ok(data)
}

fn get_pokemon(name: &str) -> anyhow::Result<RawPokemon> {
    let url: String = format!("https://pokeapi.co/api/v2/pokemon/{name}");

    fetch_json(&url)
    
}

fn get_pokemon_species(name: &str) -> anyhow::Result<RawPokemonSpecies> {
    let url: String = format!("https://pokeapi.co/api/v2/pokemon-species/{name}");

    fetch_json(&url)
}

// functions focused on building data
fn build_stats(pokemon: &RawPokemon) -> Vec<PokemonStats> {
    pokemon.stats
        .iter()
        .map(|s| PokemonStats {
            name: s.stat.name.clone(),
            base_stat: s.base_stat,
        })
        .collect()
}


fn build_abilities(pokemon: &RawPokemon) -> Vec<PokemonAbility> {
    pokemon.abilities
        .iter()
        .map(|a| PokemonAbility {
            name: a.ability.name.clone(),
            is_hidden: a.is_hidden,
        })
        .collect()
}

fn build_types(pokemon: &RawPokemon) -> Vec<PokemonTyping> {
    pokemon.types
        .iter()
        .map(|t| PokemonTyping {
            name: t.typing.name.clone(),
            slot: t.slot,
        })
        .collect()
}

fn build_species(species: &RawPokemonSpecies) -> PokemonSpeciesInfo {
    PokemonSpeciesInfo{
        name: species.name.clone(),
        egg_groups: species.egg_groups.iter().map(|g| g.name.clone()).collect(),
        gender_rate: species.gender_rate,
        hatch_counter: species.hatch_counter,
        evolves_from: species.evolves_from_species.as_ref().map(|n| n.name.clone()),
    }
}

pub fn get_pokemon_profile(name: &str, opts: &ProfileOptions) -> anyhow::Result<Pokemon> {
    let pokemon  = get_pokemon(name)?;
    let pokemon_species = get_pokemon_species(name)?;

    let profile = Pokemon {
        id: pokemon.id,
        name: pokemon.name.clone(),
        height: pokemon.height,
        weight: pokemon.weight,
        species: build_species(&pokemon_species),
        abilities: if opts.abilities { Some(build_abilities(&pokemon)) } else { None },
        types: build_types(&pokemon),
        stats: if opts.stats { Some(build_stats(&pokemon)) } else { None },
        source: ResourceSourceMetadata {
            fetched_at: now_epoch_secs(),
            partial_errors: vec![],
          }
    };

    Ok(profile)

}

fn now_epoch_secs() -> String {
    SystemTime::now()
        .duration_since(UNIX_EPOCH)
        .unwrap()
        .as_secs()
        .to_string()
}

#[cfg(test)]
mod tests {
    use super::*;

    const CHARIZARD: &str = include_str!("../tests/fixtures/charizard.json");
    const CHARIZARD_SPECIES: &str = include_str!("../tests/fixtures/charizard-species.json");

    fn charizard() -> RawPokemon {
        serde_json::from_str(CHARIZARD).unwrap()
    }

    fn charizard_species() -> RawPokemonSpecies {
        serde_json::from_str(CHARIZARD_SPECIES).unwrap()
    }

    #[test]
    fn build_abilities_maps_names_and_hidden() {
        let abilities = build_abilities(&charizard());

        assert_eq!(abilities.len(), 2);
        assert_eq!(abilities[0].name, "blaze");
        // solar-power is charizard's hidden ability
        assert!(abilities.iter().any(|a| a.name == "solar-power" && a.is_hidden));
    }

    #[test]
    fn build_types_maps_names_and_slots() {
        let types = build_types(&charizard());

        assert_eq!(types.len(), 2);
        assert_eq!(types[0].name, "fire");
        assert_eq!(types[0].slot, 1);
        assert_eq!(types[1].name, "flying");
        assert_eq!(types[1].slot, 2);
    }

    #[test]
    fn build_stats_maps_all_six() {
        let stats = build_stats(&charizard());

        assert_eq!(stats.len(), 6);
        assert_eq!(stats[0].name, "hp");
    }

    #[test]
    fn build_species_maps_summary() {
        let species = build_species(&charizard_species());

        assert_eq!(species.evolves_from, Some("charmeleon".to_string()));
        assert_eq!(species.gender_rate, 1);
        assert!(species.egg_groups.contains(&"monster".to_string()));
    }
}