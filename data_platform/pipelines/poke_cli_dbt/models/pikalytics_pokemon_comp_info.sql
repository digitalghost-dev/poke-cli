{{ config(
    materialized='incremental',
    pre_hook="truncate table {{ this }}"
) }}

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
