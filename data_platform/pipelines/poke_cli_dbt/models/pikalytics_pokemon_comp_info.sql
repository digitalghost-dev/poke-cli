{{ config(
    materialized='table',
    post_hook="{{ enable_rls() }}"
) }}

-- Replaces pikalytics_pokemon_comp_stats. n8n reads the fresh top-50 from
-- staging.pikalytics_usage, scrapes each Pokémon's AI pokedex page, and lands RAW rows
-- (pokemon, web_url + 4 JSONB "common" sections) in staging.pikalytics_pokemon_comp_info
-- (full replace each run). Derivations live here: pokemon_slug, pokemon_id (resolve_pokemon_id),
-- format/source constants. The common_* JSONB sections pass through unchanged.

WITH staged AS (
    SELECT
        pokemon,
        web_url,
        common_moves,
        common_abilities,
        common_items,
        common_teammates,
        TRIM(BOTH '-' FROM LOWER(REGEXP_REPLACE(pokemon, '[^a-zA-Z0-9]+', '-', 'g'))) AS pokemon_slug
    FROM {{ source('staging', 'pikalytics_pokemon_comp_info') }}
)

SELECT
    'gen9championsvgc2026regma'              AS format,
    pokemon,
    pokemon_slug,
    {{ resolve_pokemon_id('pokemon_slug') }} AS pokemon_id,
    web_url,
    common_moves,
    common_abilities,
    common_items,
    common_teammates,
    'pikalytics'                             AS source
FROM staged
ORDER BY pokemon
